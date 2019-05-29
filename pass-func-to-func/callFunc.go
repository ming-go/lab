package callFunc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

type CallFunc struct {
	NumInCheck bool
}

func (cf *CallFunc) Call(mFunc interface{}, mFuncParams ...interface{}) ([]reflect.Value, error) {
	mFuncTypeOf := reflect.TypeOf(mFunc)
	if mFuncTypeOf.Kind() != reflect.Func {
		return nil, errors.New("Can only pass func")
	}

	mFuncName := runtime.FuncForPC(reflect.ValueOf((mFunc)).Pointer()).Name()

	if cf.NumInCheck {
		if len(mFuncParams) != mFuncTypeOf.NumIn() {
			return nil, errors.New("The number of param is not adapted")
		}
	}

	inputArgs := make([]reflect.Value, len(mFuncParams))

	for index, mFuncParam := range mFuncParams {
		inputArgs[index] = reflect.ValueOf(mFuncParam)
	}

	fmt.Println("Call", mFuncName)

	mFuncValueOf := reflect.ValueOf(mFunc)
	callMFuncResult := mFuncValueOf.Call(inputArgs)

	return callMFuncResult, nil
}
