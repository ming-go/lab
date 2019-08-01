package main

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func controller(c echo.Context) error {
	result, err := module()
	log.Println(result)
	return err
}

func middleIncome() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceCode, err := sf.Load().(*snowflake.Snowflake).NextId()
			if err != nil {
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), snowflakeKey, traceCode)))
			err = next(c)
			if err == nil {
				return nil
			}

			zapFields := []zap.Field{}
			zapFields = append(zapFields, zap.String("error", err.Error()))
			zap.L().Info("", zapFields...)

			return err
		}
	}
}

var sf atomic.Value

func middleOutcome() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Request().Context().Value(snowflakeKey)
		log.Println("Outcome")
	})
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
	e.Use(middleOutcome())

	e.GET("/", controller)

	e.Logger.Fatal(e.Start(":1323"))
}
