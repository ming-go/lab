package callFunc

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

func maxInt(ints ...int) int {
	max := math.MinInt32
	for _, v := range ints {
		if v > max {
			max = v
		}
	}

	return max
}

func TestCallFuncNumInCheckIsTrue(t *testing.T) {
	cf := &CallFunc{
		NumInCheck: true,
	}

	results, err := cf.Call(min, 10, 20)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))

	resultValue := results[0]

	assert.Equal(t, reflect.Int, resultValue.Kind())
	resultInt := resultValue.Interface().(int)
	assert.Equal(t, 10, resultInt)

	results, err = cf.Call(min, 10, 20, 30)
	assert.NotNil(t, err)
	assert.Equal(t, "The number of param is not adapted", err.Error())

	results, err = cf.Call(maxInt, 10, 20, 30, 40, 50, 60, 70)
	assert.NotNil(t, err)
	assert.Equal(t, "The number of param is not adapted", err.Error())
}

func TestCallFuncNumInCheckIsFalse(t *testing.T) {
	cf := &CallFunc{
		NumInCheck: false,
	}

	results, err := cf.Call(maxInt, 10, 20, 30, 40, 50, 60, 70)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))

	resultValue := results[0]

	assert.Equal(t, reflect.Int, resultValue.Kind())

	resultInt := resultValue.Interface().(int)

	assert.Equal(t, 70, resultInt)
}
