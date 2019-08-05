package module

import (
	"github.com/ming-go/lab/echo-error-handling-and-logging-best-practice/models/model"
	errors "golang.org/x/xerrors"
)

func Module() (*model.ModelResult, error) {
	r, err := model.Model()
	if err != nil {
		return nil, errors.Errorf("module.Module: %w", err)
	}

	return r, nil
}
