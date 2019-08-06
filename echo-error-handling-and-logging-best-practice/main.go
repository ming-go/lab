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

	"github.com/labstack/echo"
	"github.com/ming-go/lab/echo-error-handling-and-logging-best-practice/controllers"
	"github.com/ming-go/pkg/snowflake"
	"go.uber.org/zap"
)

type key int

const (
	snowflakeKey key = iota
)

func getRequestURL(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	return scheme + r.Host + r.RequestURI
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

			zapFields := []zap.Field{}
			zapFields = append(zapFields, zap.Int64("RequestID", requestID))
			zapFields = append(zapFields, zap.NamedError("Error", err))
			zapFields = append(zapFields, zap.String("Request Method", c.Request().Method))
			zapFields = append(zapFields, zap.String("Request URL", getRequestURL(c.Request())))
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

func router() {
	e := echo.New()
	e.Use(middleIncome())

	e.GET("/ok", controllers.ControllerOK)
	e.Any("/error", controllers.ControllerError)

	e.Logger.Fatal(e.Start(":1323"))
}

// TODO: graceful shutdown
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

	go router()

	forever := make(chan struct{})
	<-forever
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
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
