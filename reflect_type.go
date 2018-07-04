package util

import (
	"reflect"
)

// InitPointer with value
func InitPointer(src reflect.Type) (dst reflect.Value) {
	if src.Kind() == reflect.Ptr {
		data := InitPointer(src.Elem())
		dst = reflect.Indirect(reflect.New(src))
		if dst.CanSet() {
			if data.CanAddr() {
				dst.Set(data.Addr())
			}
		}
	} else {
		dst = InitWithZeroValue(src)
	}
	return
}

// InitWithZeroValue init type with default value
func InitWithZeroValue(src reflect.Type) (dst reflect.Value) {
	dst = reflect.Indirect(reflect.New(src))
	switch dst.Kind() {
	case reflect.Invalid:
	case reflect.Bool:
		dst.SetBool(false)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dst.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dst.SetUint(0)
	case reflect.Uintptr:

	case reflect.Float32, reflect.Float64:
		dst.SetFloat(0)
	case reflect.Complex64, reflect.Complex128:
		dst.SetComplex(complex(0, 0))
	case reflect.Array:
		dst.Set(reflect.MakeSlice(dst.Type(), 0, 0))
	case reflect.Chan:
		dst.Set(reflect.MakeChan(dst.Type(), 0))
	case reflect.Func:
		dst.Set(reflect.MakeFunc(dst.Type(), func(args []reflect.Value) (results []reflect.Value) {
			return
		}))
	case reflect.Interface:

	case reflect.Map:
		dst.Set(reflect.MakeMap(dst.Type()))
	case reflect.Ptr:

	case reflect.Slice:
		dst.Set(reflect.MakeSlice(dst.Type(), 0, 0))
	case reflect.String:
		dst.SetString("")
	case reflect.Struct:
		numField := dst.NumField()
		for i := 0; i < numField; i++ {
			fieldValue := dst.Field(i)
			fieldValue.Set(InitWithZeroValue(fieldValue.Type()))
		}
	case reflect.UnsafePointer:

	}
	return
}

// IsExportedOrBuiltinType Is this type exported or a builtin?
func IsExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return IsExported(t.Name()) || t.PkgPath() == ""
}
