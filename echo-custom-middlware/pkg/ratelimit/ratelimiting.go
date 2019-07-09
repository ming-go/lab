package ratelimit

import (
	"errors"
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
				// Log
				return c.JSON(http.StatusInternalServerError, resp{Message: "API rate limit extractorFunc exist a error"})
			}

			if err := cfg.Impl.Take(token); err == nil {
				isExceedLimit = true
			}

			c.Response().Header().Set(HeaderXRateLimitLimit, strconv.Itoa(cfg.Impl.GetLimit()))
			//c.Response().Header().Set(HeaderXRateLimitRemaining, "1")
			//c.Response().Header().Set(HeaderXRateLimitReset, strconv.Itoa(int(time.Now().Unix())))

			if isExceedLimit {
				return c.JSON(http.StatusForbidden, resp{Message: "API rate limit exceeded for " + token})
			}

			return next(c)
		}
	}
}
