package jsonx

import (
	"encoding/json"
	"reflect"
)

func Json2Map(jsonStr string) (map[string]any, error) {
	var m map[string]any
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Map2Json(m map[string]any) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonByte), nil
}

// C style: return structPtr
func Json2Struct(structPtr any, jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), structPtr)
	if err != nil {
		return err
	}
	return nil
}

func Struct2Json(structPtr any) (string, error) {
	jsonBytes, err := json.Marshal(structPtr)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), err
}

// C style: return structPtr
func Map2Struct(structPtr any, m map[string]any) error {
	s, err := Map2Json(m)
	if err != nil {
		return err
	}

	Json2Struct(structPtr, s)
	return nil
}

//convert struct to map, embedded struct will be ignored.
func Struct2Map(structPtr any) (map[string]any, error) {
	s, err := Struct2Json(structPtr)
	if err != nil {
		return nil, err
	}

	m, err := Json2Map(s)
	if err != nil {
		return nil, err
	}
	return m, nil
}

//convert struct to map, all fields will be converted.
func Struct2MapV2(structPtr any) (_ map[string]any, err error) {
	m := make(map[string]any)

	t := reflect.TypeOf(structPtr)
	v := reflect.ValueOf(structPtr)

	for i := 0; i < t.Elem().NumField(); i++ {
		if t.Elem().Field(i).Type.Kind() == reflect.Struct {
			subMap, _ := Struct2MapV2(v.Elem().Field(i).Addr().Interface())
			for k, v := range subMap {
				m[k] = v
			}
			continue
		}
		m[t.Elem().Field(i).Name] = v.Elem().Field(i).Interface()
	}
	return m, err
}

//handle struct's tag or null.
//TODO PLAN: handle by err or boolean?
func Struct2MapWithRule(structPtr any, fns ...ruleFn) (map[string]any, error) {
	m := make(map[string]any)
	errKeys := make([]string, 0)

	t := reflect.TypeOf(structPtr)
	v := reflect.ValueOf(structPtr)

	for i := 0; i < t.Elem().NumField(); i++ {
		if t.Elem().Field(i).Type.Kind() == reflect.Struct {
			subMap, _ := Struct2MapWithRule(v.Elem().Field(i).Addr().Interface(), fns...)
			for k, v := range subMap {
				m[k] = v
			}
			continue
		}
		for _, fn := range fns {
			err := fn(t.Elem().Field(i), v.Elem().Field(i))
			if err == nil {
				m[t.Elem().Field(i).Name] = v.Elem().Field(i).Interface()
			} else {
				errKeys = append(errKeys, t.Elem().Field(i).Name)
			}
		}
	}

	for _, k := range errKeys {
		delete(m, k)
	}

	return m, nil
}

//C style, return structPtr
func Struct2StructWithRule(structPtr any, fns ...ruleFn) error {
	t := reflect.TypeOf(structPtr)
	v := reflect.ValueOf(structPtr)

	for i := 0; i < t.Elem().NumField(); i++ {
		if t.Elem().Field(i).Type.Kind() == reflect.Struct {
			if err := Struct2StructWithRule(v.Elem().Field(i).Addr().Interface(), fns...); err != nil {
				return err
			}
			continue
		}
		for _, fn := range fns {
			if err := fn(t.Elem().Field(i), v.Elem().Field(i)); err != nil {
				return err
			}
		}
	}
	return nil
}
