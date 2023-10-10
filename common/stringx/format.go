package stringx

import (
	"fmt"
	"strings"
)

//e.g: HelloWorldWujimaster -> hello_world_wujimaster
func SnakeFormat(fieldName string) (sqlName string) {
	sqlBuilder := []string{}
	bs := []byte(fieldName)
	for i, b := range bs {
		if i == 0 {
			sqlBuilder = append(sqlBuilder, strings.ToLower(string(b)))
			continue
		}
		if string(b) == strings.ToUpper(string(b)) {
			tmp := fmt.Sprintf("_%v", strings.ToLower(string(b)))
			sqlBuilder = append(sqlBuilder, tmp)
			continue
		}
		sqlBuilder = append(sqlBuilder, string(b))
	}
	return strings.Join(sqlBuilder, "")
}
