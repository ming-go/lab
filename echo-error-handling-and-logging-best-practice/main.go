package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/labstack/echo"
	"github.com/ming-go/pkg/snowflake"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type modelResult struct {
	Value float64
}

type key int

const (
	snowflakeKey key = iota
)

func model() (*modelResult, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:8787"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	var result modelResult
	collection := client.Database("testing").Collection("testing")
	err = collection.FindOne(ctx, bson.M{"name": "pi"}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func module() (*modelResult, error) {
	return model()
}

func controllerOK(c echo.Context) error {
	return c.JSON(http.StatusOK, &modelResult{Value: 10.0})
}

func controllerError(c echo.Context) error {
	result, err := module()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func GetRequestURL(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	return scheme + r.Host + r.RequestURI
}

type ErrorStruct struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID int64  `json:"requestID"`
}

type ErrorResponse struct {
	ErrorStruct `json:"error"`
}

func middleIncome() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			// Dump Request Body
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			requestID, err := sf.Load().(*snowflake.Snowflake).NextId()
			if err != nil {
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), snowflakeKey, requestID)))
			err = next(c)
			if err == nil {
				return nil
			}

			// Error Handling Example
			c.JSON(http.StatusInternalServerError, &ErrorResponse{
				ErrorStruct: ErrorStruct{
					Code:      "ContextError",
					Message:   "Context Canceled",
					RequestID: requestID,
				},
			})

			zapFields := []zap.Field{}
			zapFields = append(zapFields, zap.Int64("RequestID", requestID))
			zapFields = append(zapFields, zap.String("Error", err.Error()))
			zapFields = append(zapFields, zap.String("Request Method", c.Request().Method))
			zapFields = append(zapFields, zap.String("Request URL", GetRequestURL(c.Request())))
			zapFields = append(zapFields, zap.String("Request Protocol", c.Request().Proto))
			zapFields = append(zapFields, zap.Any("Request Header", c.Request().Header))
			zapFields = append(zapFields, zap.Any("Remote Address", c.Request().RemoteAddr))

			zapFields = append(zapFields, zap.ByteString("Request Body", reqBody))
			zapFields = append(zapFields, zap.ByteString("Response Body", resBody.Bytes()))

			zap.L().Info("IncomeLog", zapFields...)

			return err
		}
	}
}

var sf atomic.Value

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("zap init failed", err)
	}
	zap.ReplaceGlobals(logger)

	sfp, err := snowflake.New(0, 0)
	if err != nil {
		log.Fatal("snowflake init failed", err)
	}
	sf.Store(sfp)

	e := echo.New()
	e.Use(middleIncome())
	//e.Use(middleOutcome())

	e.GET("/ok", controllerOK)
	e.Any("/error", controllerError)

	e.Logger.Fatal(e.Start(":1323"))
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *bodyDumpResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
