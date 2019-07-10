package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/ming-go/lab/echo-custom-middlware/pkg/ratelimit"
	"github.com/ming-go/pkg/mredigo"
	"github.com/ming-go/pkg/ratelimiting"
)

type resp struct {
	Message string
}

func ctest1(c echo.Context) error {
	return c.JSON(http.StatusOK, resp{Message: "Hello, ctest1"})
}

func ctest2(c echo.Context) error {
	return c.JSON(http.StatusOK, resp{Message: "Hello, ctest2"})
}

func ctest3(c echo.Context) error {
	return c.JSON(http.StatusOK, resp{Message: "Hello, ctest3"})
}

func main() {
	e := echo.New()
	e.Debug = true

	cfg := mredigo.NewConfig()
	cfg.Host = "172.77.0.99:6379"
	cfg.Database = "7"

	pool, err := mredigo.CreatePool("key", true, cfg)
	if err != nil {
		log.Fatal(err)
	}

	tGroup := e.Group("/test", ratelimit.RateLimit(
		&ratelimit.RateLimitConfig{
			Impl: ratelimiting.NewRedisLeakyBucketImpl(
				&ratelimiting.RedisLeakyBucketImplConfig{
					RPool:  pool,
					Limit:  2,
					Period: 60 * time.Second,
				},
			),
			ExtractorFunc: ratelimit.ExtractorFuncByURLPath,
		},
	))

	tGroup.GET("/c1", ctest1)
	tGroup.GET("/c2", ctest2)
	tGroup.GET("/c3", ctest3)

	e.Logger.Fatal(e.Start(":1323"))
}
