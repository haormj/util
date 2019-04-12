package util

import (
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cast"
)

// StructField struct filed
type StructField struct {
	Field      reflect.StructField
	FieldName  string
	FieldValue reflect.Value
	IsExist    bool
}

// Bool convert StructField to bool
func (sf *StructField) Bool() (bool, bool) {
	if sf.IsExist {
		return cast.ToBool(sf.FieldValue.Interface()), sf.IsExist
	}
	return false, sf.IsExist
}

// BoolE convert StructField to bool with error
func (sf *StructField) BoolE() (bool, bool, error) {
	if sf.IsExist {
		b, err := cast.ToBoolE(sf.FieldValue.Interface())
		return b, sf.IsExist, err
	}
	return false, sf.IsExist, nil
}

// Int convert StructField to int
func (sf *StructField) Int() (int, bool) {
	if sf.IsExist {
		return cast.ToInt(sf.FieldValue.Interface()), sf.IsExist
	}
	return 0, sf.IsExist
}

// IntE convert StructField to int with error
func (sf *StructField) IntE() (int, bool, error) {
	if sf.IsExist {
		i, err := cast.ToIntE(sf.FieldValue.Interface())
		return i, sf.IsExist, err
	}
	return 0, sf.IsExist, nil
}

// Float64 convert StructField to float64
func (sf *StructField) Float64() (float64, bool) {
	if sf.IsExist {
		return cast.ToFloat64(sf.FieldValue.Interface()), sf.IsExist
	}
	return 0, sf.IsExist
}

// Float64E convert StructField to float64 with error
func (sf *StructField) Float64E() (float64, bool, error) {
	if sf.IsExist {
		f, err := cast.ToFloat64E(sf.FieldValue.Interface())
		return f, sf.IsExist, err
	}
	return 0, sf.IsExist, nil
}

// String convert StructField to string
func (sf *StructField) String() (string, bool) {
	if sf.IsExist {
		return cast.ToString(sf.FieldValue.Interface()), sf.IsExist
	}
	return "", sf.IsExist
}

// StringE convert StructField to string with error
func (sf *StructField) StringE() (string, bool, error) {
	if sf.IsExist {
		str, err := cast.ToStringE(sf.FieldValue.Interface())
		return str, sf.IsExist, err
	}
	return "", sf.IsExist, nil
}

// Interface return StructField value
func (sf *StructField) Interface() (interface{}, bool) {
	return sf.FieldValue.Interface(), sf.IsExist
}

// SetStructField v is ptr to struct
func SetStructField(v interface{}, fieldValue interface{}, fieldNames ...string) bool {
	sf := GetStructField(v, fieldNames...)
	if !sf.IsExist {
		return false
	}
	if !sf.FieldValue.CanSet() {
		return false
	}
	fieldValueValue := reflect.ValueOf(fieldValue)
	fieldValueType := reflect.TypeOf(fieldValue)
	switch sf.Field.Type.Kind() {
	case fieldValueType.Kind():
		// Field type == Field Value Type, use fieldValueValue
	case reflect.Bool:
		b, err := cast.ToBoolE(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(b)
	case reflect.Int:
		i, err := cast.ToIntE(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(i)
	case reflect.Int8:
		i, err := cast.ToInt8E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(i)
	case reflect.Int16:
		i, err := cast.ToInt16E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(i)
	case reflect.Int32:
		i, err := cast.ToInt32E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(i)
	case reflect.Int64:
		i, err := cast.ToInt64E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(i)
	case reflect.Uint:
		u, err := cast.ToUintE(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(u)
	case reflect.Uint8:
		u, err := cast.ToUint8E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(u)
	case reflect.Uint16:
		u, err := cast.ToUint16E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(u)
	case reflect.Uint32:
		u, err := cast.ToUint32E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(u)
	case reflect.Uint64:
		u, err := cast.ToUint64E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(u)
	case reflect.Float32:
		f, err := cast.ToFloat32E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(f)
	case reflect.Float64:
		f, err := cast.ToFloat64E(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(f)
	case reflect.String:
		s, err := cast.ToStringE(fieldValue)
		if err != nil {
			return false
		}
		fieldValueValue = reflect.ValueOf(s)

	case reflect.Ptr:
		// Field type *int, Field Value type int, support
		if sf.Field.Type.Elem() == fieldValueType {
			v := reflect.New(fieldValueType)
			v.Elem().Set(fieldValueValue)
			fieldValueValue = v
		} else {
			return false
		}
	case reflect.Complex64:
		return false
	case reflect.Complex128:
		return false
	case reflect.Array:
		return false
	case reflect.Chan:
		return false
	case reflect.Func:
		return false
	case reflect.Interface:
		return false
	case reflect.Map:
		return false
	case reflect.Slice:
		return false
	case reflect.Struct:
		return false
	case reflect.Uintptr:
		return false
	case reflect.UnsafePointer:
		return false
	case reflect.Invalid:
		return false
	}
	sf.FieldValue.Set(fieldValueValue)
	return true
}

// GetStructField get struct field
func GetStructField(v interface{}, fieldNames ...string) *StructField {
	var sf *StructField
	fieldValue := reflect.ValueOf(v)
	for _, fieldName := range fieldNames {
		sf = getStructField(fieldValue, fieldName)
		if !sf.IsExist {
			break
		}
		fieldValue = sf.FieldValue
	}
	return sf
}

func getStructField(v reflect.Value, fieldName string) *StructField {
	sf := &StructField{}
	vValue := v
	for {
		if vValue.Kind() != reflect.Ptr {
			break
		}
		vValue = reflect.Indirect(vValue)
	}
	if vValue.Kind() != reflect.Struct {
		return sf
	}
	vType := vValue.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		if field.Name == fieldName {
			fValue := vValue.Field(i)
			sf := &StructField{
				Field:      field,
				FieldName:  fieldName,
				FieldValue: fValue,
				IsExist:    true,
			}
			return sf
		}
	}
	return sf
}

// IsExported Is this an exported - upper case - name?
func IsExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}
