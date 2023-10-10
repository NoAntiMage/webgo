package jsonx

import (
	"errors"
	"reflect"
	"strconv"
)

//convert is under controll of rules.
//implement map,filter,reduce rules by Tag, Null, Default.
//C style, the params, t & v, will be return.
type ruleFn func(t reflect.StructField, v reflect.Value) error

//tag: default. set struct default value if zero.
func RuleDefault(t reflect.StructField, v reflect.Value) error {
	if v.IsZero() {
		defaultx := t.Tag.Get(DefaultTag)
		if defaultx != "" && defaultx != "-" {
			switch t.Type.String() {
			case "int", "int8", "int32", "int64":
				result, err := strconv.ParseInt(defaultx, 10, 64)
				if err != nil {
					return err
				}
				v.SetInt(int64(result))
			case "uint", "uint8", "uint32", "uint64":
				result, err := strconv.ParseUint(defaultx, 10, 64)
				if err != nil {
					return err
				}
				v.SetUint(uint64(result))
			case "string":
				v.SetString(defaultx)
			default:
			}
		}
	}
	return nil
}

//ignore the field if null.
func RuleWithoutZero(t reflect.StructField, v reflect.Value) error {
	if v.IsZero() {
		return errors.New("jsonx: field is null.")
	}
	return nil
}
