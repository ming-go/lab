package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ming-go/lab/echo-error-handling-and-logging-best-practice/models/model"
	"github.com/ming-go/lab/echo-error-handling-and-logging-best-practice/modules/module"
	errors "golang.org/x/xerrors"
)

func ControllerOK(c echo.Context) error {
	return c.JSON(http.StatusOK, &model.ModelResult{Value: 10.0})
}

type ErrorStruct struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID int64  `json:"requestID"`
}

type ErrorResponse struct {
	ErrorStruct `json:"error"`
}

func ResponseWithRequestID(req *http.Request, resp interface{}) {
	//c.Request().Context()
}

func ControllerError(c echo.Context) error {
	result, err := module.Module()
	if err != nil {
		//c.JSON(http.StatusInternalServerError, &ErrorResponse{
		//	ErrorStruct: ErrorStruct{
		//		Code:    "ContextError",
		//		Message: "Context Canceled",
		//	},
		//})

		return errors.Errorf("controllers.ControllerError: %w", err)
	}

	return c.JSON(http.StatusOK, result)
}
