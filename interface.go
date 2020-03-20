package util

import (
	"reflect"
)

// InitInterface init interface to zero value
func InitInterface(i interface{}) interface{} {
	return InitWithZeroValueExcludePtr(reflect.TypeOf(i)).Interface()
}
