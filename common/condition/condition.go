package condition

// where condition
type QueryCondition struct {
	Operator   string `validate:"required,oneof=and or"`
	Fieldname  string `validate:"required"`
	Comparison string `validate:"required,oneof=eq ne gt ge lt le in like"`
	Value      any    `validate:"required"`
}
