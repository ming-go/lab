package ratelimit

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/ming-go/pkg/ratelimiting"
)

const (
	HeaderXRateLimitLimit     = "X-RateLimit-Limit"
	HeaderXRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderXRateLimitReset     = "X-RateLimit-Reset"
)

type resp struct {
	Message string
}

type RateLimitConfig struct {
	ExtractorFunc func(interface{}) (string, error)
	Impl          ratelimiting.RateLimitingInf
}

func ExtractorFuncByURLPath(httpRequest interface{}) (string, error) {
	r, ok := httpRequest.(*http.Request)
	if !ok {
		return "", errors.New("Params must be a *http.Request")
	}

	return r.URL.Path, nil
}

func RateLimit(cfg *RateLimitConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isExceedLimit := false
			token, err := cfg.ExtractorFunc(c.Request())
			if err != nil {
				log.Println("API rate limit extractorFunc returns an error: ", err)
				return c.JSON(http.StatusInternalServerError, resp{Message: "API rate limit extractorFunc returns an error"})
			}

			result, err := cfg.Impl.Take(token)
			if err != nil {
				isExceedLimit = true
			}

			remaining := result.Remaining

			reset := strconv.FormatInt(result.Reset, 10)

			c.Response().Header().Set(HeaderXRateLimitLimit, strconv.Itoa(cfg.Impl.GetLimit()))
			c.Response().Header().Set(HeaderXRateLimitRemaining, strconv.Itoa(remaining))
			c.Response().Header().Set(HeaderXRateLimitReset, reset)

			if isExceedLimit {
				return c.JSON(http.StatusForbidden, resp{Message: "API rate limit exceeded for " + token})
			}

			return next(c)
		}
	}
}
