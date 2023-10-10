package copier

import "github.com/jinzhu/copier"

func Copy(toValue any, fromValue any) error {
	return copier.Copy(toValue, fromValue)
}
