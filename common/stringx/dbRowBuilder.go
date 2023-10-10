package stringx

import (
	"fmt"
	"reflect"
	"strings"
)

// default tagname for sqlx
const dbTag = "db"

// convert golang struct fields into []string.
// according to sqlx tag.
// eg. `db:field_name`
func StructFieldNames(in any) (res []string) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		// struct type, step into recursion.
		if field.Type.Kind() == reflect.Struct &&
			field.Type.Name() != "NullBool" &&
			field.Type.Name() != "NullByte" &&
			field.Type.Name() != "NullFloat64" &&
			field.Type.Name() != "NullInt" &&
			field.Type.Name() != "NullInt16" &&
			field.Type.Name() != "NullInt32" &&
			field.Type.Name() != "NullInt64" &&
			field.Type.Name() != "NullString" &&
			field.Type.Name() != "NullTime" {
			subFields := StructFieldNames(v.FieldByIndex([]int{i}).Interface())
			res = append(res, subFields...)
			continue
		}

		tagv := field.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			res = append(res, fmt.Sprintf("%s", field.Name))
		default:
			res = append(res, fmt.Sprintf("%s", tagv))
		}
	}
	return res
}

func FieldsFromStructToSql(li []string) (res []string) {
	res = make([]string, len(li))
	for i := 0; i < len(li); i++ {
		res[i] = SnakeFormat(li[i])
	}
	return res
}

func RowsWithPlaceHolder(rows []string) string {
	return strings.Join(rows, "=?,") + "=?"
}
