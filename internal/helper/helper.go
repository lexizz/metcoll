package helper

import (
	"errors"
	"fmt"
	"strings"
)

func GetType(value interface{}) (string, error) {
	asd := fmt.Sprintf("%T", value)
	valueTypeData := strings.Split(asd, ".")

	if len(valueTypeData) > 0 {
		resultType := valueTypeData[len(valueTypeData)-1]

		return strings.ToLower(resultType), nil
	}

	return "", errors.New("the type is not defined")
}
