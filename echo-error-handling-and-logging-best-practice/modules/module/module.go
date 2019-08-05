package module

import "github.com/ming-go/lab/echo-error-handling-and-logging-best-practice/models/model"

func Module() (*model.ModelResult, error) {
	return model.Model()
}
