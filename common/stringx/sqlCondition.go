package stringx

import (
	"errors"
	"goweb/common/condition"
	"reflect"
)

var (
	ErrComparison = errors.New("cond: ErrComparison")
	ErrOperator   = errors.New("cond: ErrOperator")
)

func WithCond(b *SqlBuilder, cond condition.QueryCondition) *SqlBuilder {
	var comparison string

	switch cond.Comparison {
	case "eq":
		comparison = b.Eq(cond.Fieldname, cond.Value)
	case "ne":
		comparison = b.Ne(cond.Fieldname, cond.Value)
	case "gt":
		comparison = b.Gt(cond.Fieldname, cond.Value)
	case "ge":
		comparison = b.Ge(cond.Fieldname, cond.Value)
	case "lt":
		comparison = b.Lt(cond.Fieldname, cond.Value)
	case "le":
		comparison = b.Le(cond.Fieldname, cond.Value)
	case "in":
		v := reflect.ValueOf(cond.Value)
		if v.Kind() != reflect.Slice {
			comparison = b.Eq(cond.Fieldname, cond.Value)
		} else {
			comparison = b.In(cond.Fieldname, cond.Value)
		}
	case "like":
		comparison = b.Like(cond.Fieldname, cond.Value)
	default:
		b.Errs[ErrComparison] = true
	}

	switch cond.Operator {
	case "and":
		b = b.And(comparison)
	case "or":
		b = b.Or(comparison)
	default:
		b.Errs[ErrOperator] = true
	}

	return b
}
