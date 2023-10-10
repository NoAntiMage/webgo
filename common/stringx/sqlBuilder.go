package stringx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	moduleName     = "SqlBuild"
	SqlPlaceHolder = "?"
)

var (
	ErrNoneColumn          = errors.New(moduleName + ": columns is none, at least 1.")
	ErrSelectWithNoFrom    = errors.New(moduleName + ": main action [SELECT] must be use with FROM.")
	ErrActionExist         = errors.New(moduleName + ": main action [SELECT, UPDATE, INSERT, DELETE] has been set.")
	ErrActionNotExist      = errors.New(moduleName + ": main action [SELECT, UPDATE, INSERT, DELETE] should be set property.")
	ErrFromExist           = errors.New(moduleName + ": duplicate FROM.")
	ErrFromNotExist        = errors.New(moduleName + ": FROM not exist. JOIN should be set after FROM.")
	ErrWhereExist          = errors.New(moduleName + ": duplicate WHERE.")
	ErrWhereNotExist       = errors.New(moduleName + ": WHERE not exist. condition should be set after WHERE.")
	ErrBetweenNotMatch     = errors.New(moduleName + ": type of args in BETWEEN is not match.")
	ErrGroupWithNoFunc     = errors.New(moduleName + ": GROUP-BY should be used with Aggregate function.")
	ErrGroupBeforeFrom     = errors.New(moduleName + ": FROM not found. GROUP-BY should be used after FROM.")
	ErrHavingWithNoGroup   = errors.New(moduleName + ": HAVING should be used with GROUP-BY.")
	ErrOrderByExist        = errors.New(moduleName + ": duplicate ORDER-BY.")
	ErrOrderByAfterLimit   = errors.New(moduleName + ": LIMIT exist.ORDER-BY should be used before LIMIT.")
	ErrLimitExist          = errors.New(moduleName + ": duplicate LIMIT.")
	ErrLimitNotExist       = errors.New(moduleName + ": OFFSET should be used with LIMIT.")
	ErrOffsetExist         = errors.New(moduleName + ": duplicate OFFSET.")
	ErrHolderParamNotMatch = errors.New(moduleName + ": PlaceHolder and params not match.")
)

type Piority int

// Piority used as binary in PiorityMask.
const (
	action  Piority = 0
	selectx         = 1
	insert          = 2
	update          = 3
	deletex         = 3

	from = 4
	join = 4

	where = 5
	and   = 5
	or    = 5

	aggregate = 6
	groupBy   = 7
	having    = 7

	orderBy = 8

	limit  = 9
	offset = 10
	end    = 11
)

var (
	actionMask int = 1<<from - 1
	selectMask int = 1<<action + 1<<selectx
	insertMask int = 1<<action + 1<<insert
	updateMask int = 1<<action + 1<<update

	withoutActionMask  int = (1<<end - 1) - (1<<from - 1)
	SelectWithFromMask int = 1<<selectx + 1<<from
)

type (
	actionBuilder interface {
		Select(rows ...string) *SqlBuilder
		Insert(table string, kvs map[string]any) *SqlBuilder
		Update(table string, kvs map[string]any) *SqlBuilder
		DeleteSoft(table string) *SqlBuilder

		From(tables ...string) *SqlBuilder
		Join(table string, match string) *SqlBuilder
		LeftJoin(table string, match string) *SqlBuilder
		RightJoin(table string, match string) *SqlBuilder
		OrderBy(column string, desc bool) *SqlBuilder
	}

	whereBuilder interface {
		Where(condition string) *SqlBuilder
		And(condition string) *SqlBuilder
		Or(condition string) *SqlBuilder

		Between(column string, left any, right any) string
		NotBetween(column string, left any, right any) string
		In(column string, values any) string
		NotIn(column string, values any) string
		Like(column string, value any) string
		NotLike(column string, value any) string

		Eq(column string, value any) string
		Ne(column string, value any) string
		Gt(column string, value any) string
		Ge(column string, value any) string
		Lt(column string, value any) string
		Le(column string, value any) string
	}

	aggregateBuilder interface {
		GroupBy(columns ...string) *SqlBuilder
		Having(condition string) *SqlBuilder

		Count(field string) string
		Sum(field string) string
		Avg(field string) string
		Max(field string) string
		Min(field string) string
	}

	limitBuilder interface {
		Limit(num int) *SqlBuilder
		Offset(num int) *SqlBuilder
	}

	sqlBuilder interface {
		actionBuilder
		whereBuilder
		aggregateBuilder
		limitBuilder
	}
)

//try *SqlBuilder with grammar tree in next implement.

//Builder makes sql statement.
//Most of result is prepared statment.
//Use it with sqlx.
type SqlBuilder struct {
	// blocks for concating sql stmt.
	stmt []string
	// ORDER-BY, LIMIT added at the last of stmt.
	lastStmt []string
	//params are set for placeholder
	params      []any
	lastParams  []any
	KeywordMask int
	Errs        map[error]bool
}

var (
	_ sqlBuilder = (*SqlBuilder)(nil)
)

func newSqlBuilder() *SqlBuilder {
	return &SqlBuilder{
		stmt:       make([]string, 0),
		lastStmt:   make([]string, 0),
		params:     make([]any, 0),
		lastParams: make([]any, 0),
		Errs:       make(map[error]bool),
	}
}

func NewSqlBuilder() *SqlBuilder {
	return newSqlBuilder()
}

//=== main action piority ===

func (s *SqlBuilder) selectx(rows ...string) *SqlBuilder {
	selectBlock := fmt.Sprintf("SELECT %v", strings.Join(rows, ","))
	s.stmt = append(s.stmt, selectBlock)
	return s
}

func (s *SqlBuilder) Select(rows ...string) *SqlBuilder {
	if len(rows) == 0 {
		s.Errs[ErrNoneColumn] = true
	}
	if s.KeywordMask>>action%2 == 1 {
		s.Errs[ErrActionExist] = true
	}
	s.KeywordMask |= 1 << action
	s.KeywordMask |= 1 << selectx

	return s.selectx(rows...)
}

func (s *SqlBuilder) insert(table string, kvs map[string]any) *SqlBuilder {
	keys := make([]string, 0, len(kvs))
	values := make([]any, 0, len(kvs))
	placeHolderArr := make([]string, 0, len(kvs))

	for k, v := range kvs {
		keys = append(keys, SnakeFormat(k))
		values = append(values, v)
		placeHolderArr = append(placeHolderArr, SqlPlaceHolder)
	}

	insertBlock := fmt.Sprintf(fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v)", table, strings.Join(keys, ","), strings.Join(placeHolderArr, ",")))

	s.stmt = append(s.stmt, insertBlock)
	s.params = append(s.params, values...)
	return s
}

func (s *SqlBuilder) Insert(table string, kvs map[string]any) *SqlBuilder {
	if len(kvs) == 0 {
		s.Errs[ErrNoneColumn] = true
	}
	if s.KeywordMask>>action%2 == 1 {
		s.Errs[ErrActionExist] = true
	}
	s.KeywordMask |= 1 << action
	s.KeywordMask |= 1 << insert

	return s.insert(table, kvs)
}

func (s *SqlBuilder) update(table string, kvs map[string]any) *SqlBuilder {
	kvStmt := []string{}
	values := []any{}
	for k, v := range kvs {
		kvStmt = append(kvStmt, fmt.Sprintf("%v = %v", SnakeFormat(k), SqlPlaceHolder))
		values = append(values, v)
	}
	updateBlock := fmt.Sprintf("UPDATE %v SET %v", table, strings.Join(kvStmt, ", "))
	s.stmt = append(s.stmt, updateBlock)
	s.params = append(s.params, values...)
	return s
}

func (s *SqlBuilder) Update(table string, kvs map[string]any) *SqlBuilder {
	if len(kvs) == 0 {
		s.Errs[ErrNoneColumn] = true
	}
	if s.KeywordMask>>action%2 == 1 {
		s.Errs[ErrActionExist] = true
	}
	s.KeywordMask |= 1 << action
	s.KeywordMask |= 1 << update

	return s.update(table, kvs)
}

func (s *SqlBuilder) delete(table string) *SqlBuilder {
	deleteBlock := fmt.Sprintf("DELETE FROM %v", table)
	s.stmt = append(s.stmt, deleteBlock)
	return s
}

func (s *SqlBuilder) Delete(table string) *SqlBuilder {
	if s.KeywordMask>>action%2 == 1 {
		s.Errs[ErrActionExist] = true
	}
	s.KeywordMask |= 1 << action
	s.KeywordMask |= 1 << deletex

	return s.delete(table)
}

// BaseModel.isDeleted is the flag to define delete status in biz.
func (s *SqlBuilder) deleteSoft(table string) *SqlBuilder {
	deleteBlock := fmt.Sprintf("UPDATE %v set is_deleted=1", table)
	s.stmt = append(s.stmt, deleteBlock)
	return s
}

// delete logically.
func (s *SqlBuilder) DeleteSoft(table string) *SqlBuilder {
	if s.KeywordMask>>action%2 == 1 {
		s.Errs[ErrActionExist] = true
	}
	s.KeywordMask |= 1 << action
	s.KeywordMask |= 1 << update

	return s.deleteSoft(table)
}

//=== from piority ===

func (s *SqlBuilder) from(tables ...string) *SqlBuilder {
	fromBlock := fmt.Sprintf(" FROM %v", strings.Join(tables, ","))
	s.stmt = append(s.stmt, fromBlock)
	return s
}

// From in sql statement.
//alias style is legal.
// e.g. "user u", "person p"
func (s *SqlBuilder) From(tables ...string) *SqlBuilder {
	if s.KeywordMask>>action%2 == 0 {
		s.Errs[ErrActionNotExist] = true
	}
	if s.KeywordMask>>from%2 == 1 {
		s.Errs[ErrFromExist] = true
	}
	s.KeywordMask |= 1 << from
	return s.from(tables...)
}

func (s *SqlBuilder) join(table string, match string) *SqlBuilder {
	joinBlock := fmt.Sprintf(" JOIN %v ON %v", table, match)
	s.stmt = append(s.stmt, joinBlock)
	return s
}

func (s *SqlBuilder) Join(table string, match string) *SqlBuilder {
	if s.KeywordMask>>from%2 == 0 {
		s.Errs[ErrFromNotExist] = true
	}
	return s.join(table, match)
}

func (s *SqlBuilder) leftJoin(table string, match string) *SqlBuilder {
	joinBlock := fmt.Sprintf(" LEFT JOIN %v ON %v ", table, match)
	s.stmt = append(s.stmt, joinBlock)
	return s
}

func (s *SqlBuilder) LeftJoin(table string, match string) *SqlBuilder {
	if s.KeywordMask>>from%2 == 0 {
		s.Errs[ErrFromNotExist] = true
	}
	return s.leftJoin(table, match)
}

func (s *SqlBuilder) rightJoin(table string, match string) *SqlBuilder {
	joinBlock := fmt.Sprintf(" RIGHT JOIN %v ON %v ", table, match)
	s.stmt = append(s.stmt, joinBlock)
	return s
}

//RightJoin() not for PROD!
//sqlite not support...
func (s *SqlBuilder) RightJoin(table string, match string) *SqlBuilder {
	if s.KeywordMask>>from%2 == 0 {
		s.Errs[ErrFromNotExist] = true
	}
	return s.rightJoin(table, match)

}

//=== where piority ===

func (s *SqlBuilder) where(condition string) *SqlBuilder {
	var whereBlock string
	if len(condition) == 0 {
		whereBlock = fmt.Sprintf(" WHERE 1 = 1")
	} else {
		whereBlock = fmt.Sprintf(" WHERE %v", condition)
	}
	s.stmt = append(s.stmt, whereBlock)
	return s
}

// Where block for sqlBuilder.
// Where(""),as empty condition string, is legal which default is `Where 1 = 1`
func (s *SqlBuilder) Where(condition string) *SqlBuilder {
	if s.KeywordMask>>where%2 == 1 {
		s.Errs[ErrWhereExist] = true
	}
	s.KeywordMask |= 1 << where
	return s.where(condition)
}
func (s *SqlBuilder) or(condition string) *SqlBuilder {
	orBlock := fmt.Sprintf(" OR %v", condition)
	s.stmt = append(s.stmt, orBlock)
	return s
}

func (s *SqlBuilder) Or(condition string) *SqlBuilder {
	if s.KeywordMask>>where%2 == 0 {
		s.Errs[ErrWhereNotExist] = true
	}
	return s.or(condition)
}

func (s *SqlBuilder) and(condition string) *SqlBuilder {
	orBlock := fmt.Sprintf(" AND %v", condition)
	s.stmt = append(s.stmt, orBlock)
	return s
}

func (s *SqlBuilder) And(condition string) *SqlBuilder {
	if s.KeywordMask>>where%2 == 0 {
		s.Errs[ErrWhereNotExist] = true
	}
	return s.and(condition)
}

//=== group piority ===

func (s *SqlBuilder) groupBy(columns ...string) *SqlBuilder {
	groupBlock := fmt.Sprintf(" GROUP BY %v", strings.Join(columns, ","))
	s.stmt = append(s.stmt, groupBlock)
	return s
}

func (s *SqlBuilder) GroupBy(columns ...string) *SqlBuilder {
	if s.KeywordMask>>aggregate%2 == 0 {
		s.Errs[ErrGroupWithNoFunc] = true
	}
	if s.KeywordMask>>from%2 == 0 {
		s.Errs[ErrGroupBeforeFrom] = true
	}
	s.KeywordMask |= 1 << groupBy
	return s.groupBy(columns...)
}

func (s *SqlBuilder) having(condition string) *SqlBuilder {
	havingBlock := fmt.Sprintf(" HAVING %v", condition)
	s.stmt = append(s.stmt, havingBlock)
	return s
}

func (s *SqlBuilder) Having(condition string) *SqlBuilder {
	if s.KeywordMask>>groupBy%2 == 0 {
		s.Errs[ErrHavingWithNoGroup] = true
	}
	return s.having(condition)
}

//=== orderBy piority ===

func (s *SqlBuilder) orderBy(column string, desc bool) *SqlBuilder {
	var orderBlock string
	if desc {
		orderBlock = fmt.Sprintf(" ORDER BY %v DESC", column)
	} else {
		orderBlock = fmt.Sprintf(" ORDER BY %v ASC", column)
	}
	s.lastStmt = append(s.lastStmt, orderBlock)
	return s
}

func (s *SqlBuilder) OrderBy(column string, desc bool) *SqlBuilder {
	if s.KeywordMask>>action%2 == 0 {
		s.Errs[ErrActionNotExist] = true
	}
	if s.KeywordMask>>orderBy%2 == 1 {
		s.Errs[ErrOrderByExist] = true
	}
	if s.KeywordMask>>limit%2 == 1 {
		s.Errs[ErrOrderByAfterLimit] = true
	}
	s.KeywordMask |= 1 << orderBy
	return s.orderBy(column, desc)
}

//=== Limit piority ===

func (s *SqlBuilder) limit(num int) *SqlBuilder {
	s.lastParams = append(s.lastParams, num)
	limitBLock := fmt.Sprintf(" LIMIT %v", SqlPlaceHolder)
	s.lastStmt = append(s.lastStmt, limitBLock)
	return s
}

func (s *SqlBuilder) Limit(num int) *SqlBuilder {
	if s.KeywordMask>>limit%2 == 1 {
		s.Errs[ErrLimitExist] = true
	}
	s.KeywordMask |= 1 << limit
	return s.limit(num)
}

func (s *SqlBuilder) offset(num int) *SqlBuilder {
	s.lastParams = append(s.lastParams, num)
	offsetBLock := fmt.Sprintf(" OFFSET %v", SqlPlaceHolder)
	s.lastStmt = append(s.lastStmt, offsetBLock)
	return s
}

func (s *SqlBuilder) Offset(num int) *SqlBuilder {
	if s.KeywordMask>>offset%2 == 1 {
		s.Errs[ErrOffsetExist] = true
	}
	if s.KeywordMask>>limit%2 == 0 {
		s.Errs[ErrLimitNotExist] = true
	}
	s.KeywordMask |= 1 << offset
	return s.offset(num)
}

// === Aggregate func ===
// For Select Column.

func (s *SqlBuilder) min(field string) string {
	return "MIN(" + field + ")"
}

func (s *SqlBuilder) Min(field string) string {
	s.KeywordMask |= 1 << aggregate
	return s.min(field)
}

func (s *SqlBuilder) max(field string) string {
	return "MAX(" + field + ")"
}

func (s *SqlBuilder) Max(field string) string {
	s.KeywordMask |= 1 << aggregate
	return s.max(field)
}

func (s *SqlBuilder) avg(field string) string {
	return "AVG(" + field + ")"
}

func (s *SqlBuilder) Avg(field string) string {
	s.KeywordMask |= 1 << aggregate
	return s.avg(field)
}

func (s *SqlBuilder) count(field string) string {
	return "COUNT(" + field + ")"
}

func (s *SqlBuilder) Count(field string) string {
	s.KeywordMask |= 1 << aggregate
	return s.count(field)
}

func (s *SqlBuilder) sum(field string) string {
	return "SUM(" + field + ")"
}

func (s *SqlBuilder) Sum(field string) string {
	s.KeywordMask |= 1 << aggregate
	return s.sum(field)
}

//=== Where func ===
// For Where condition.
// According to shell Conditional Expressions

func (s *SqlBuilder) Eq(column string, value any) string {
	operator := "="
	return s.conditionTypeSwitch(column, value, operator)
}

func (s *SqlBuilder) Ne(column string, value any) string {
	operator := "!="
	return s.conditionTypeSwitch(column, value, operator)
}

func (s *SqlBuilder) Lt(column string, value any) string {
	operator := "<"
	return s.conditionTypeSwitch(column, value, operator)
}

func (s *SqlBuilder) Le(column string, value any) string {
	operator := "<="
	return s.conditionTypeSwitch(column, value, operator)
}

func (s *SqlBuilder) Gt(column string, value any) string {
	operator := ">"
	return s.conditionTypeSwitch(column, value, operator)
}

func (s *SqlBuilder) Ge(column string, value any) string {
	operator := ">="
	return s.conditionTypeSwitch(column, value, operator)
}

func (s SqlBuilder) conditionTypeSwitch(column string, value any, operator string) string {
	switch value.(type) {
	case string:
		return fmt.Sprintf("%v %v \"%v\"", column, operator, value)
	default:
		return fmt.Sprintf("%v %v %v", column, operator, value)
	}
}

func (s *SqlBuilder) Between(column string, left any, right any) string {
	if reflect.TypeOf(left).String() != reflect.TypeOf(right).String() {
		s.Errs[ErrBetweenNotMatch] = true
	}

	switch left.(type) {
	case int:
		return fmt.Sprintf("%v BETWEEN %v AND %v", column, left, right)
	default: // string
		return fmt.Sprintf("%v BETWEEN \"%v\" AND \"%v\"", column, left, right)
	}
}

func (s *SqlBuilder) NotBetween(column string, left any, right any) string {
	switch left.(type) {
	case int:
		return fmt.Sprintf("%v NOT BETWEEN %v AND %v", column, left, right)
	default: // string
		return fmt.Sprintf("%v NOT BETWEEN \"%v\" AND \"%v\"", column, left, right)
	}
}

func (s *SqlBuilder) In(column string, values any) string {
	v := reflect.ValueOf(values)
	InBlock := make([]string, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		value := v.Index(i).Interface()
		switch value.(type) {
		case int, int64:
			InBlock = append(InBlock, fmt.Sprintf("%v", value))
		case string:
			InBlock = append(InBlock, fmt.Sprintf("\"%v\"", value))
		}
	}
	return fmt.Sprintf("%v IN (%v)", column, strings.Join(InBlock, ","))
}

func (s *SqlBuilder) NotIn(column string, values any) string {
	v := reflect.ValueOf(values)
	InBlock := make([]string, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		value := v.Index(i).Interface()
		switch value.(type) {
		case int, int64:
			InBlock = append(InBlock, fmt.Sprintf("%v", value))
		case string:
			InBlock = append(InBlock, fmt.Sprintf("\"%v\"", value))
		}
	}
	return fmt.Sprintf("%v NOT IN (%v)", column, strings.Join(InBlock, ","))
}

func (s *SqlBuilder) Like(column string, value any) string {
	return fmt.Sprintf("%v LIKE \"%%%v%%\"", column, value)
}

func (s *SqlBuilder) NotLike(column string, value any) string {
	return fmt.Sprintf("%v NOT LIKE \"%%%v%%\"", column, value)
}

//=== Final api ===

func (s *SqlBuilder) prepareString() string {
	wholeStmt := append(s.stmt, s.lastStmt...)
	return strings.Join(wholeStmt, "")
}

// PrepareString() return a prepared statement.
//WARNING!!! input map is un-order.
// Use it with Params()
func (s *SqlBuilder) PrepareString() (string, error) {
	Err := Valid(s)
	return s.prepareString(), Err
}

// ToParams() return the parameters to be used in prepared statement.
// Use it with String()
func (s *SqlBuilder) toParams() []any {
	res := append([]any(nil), s.params...)
	res = append(res, s.lastParams...)
	return res
}

func (s *SqlBuilder) Params() []any {
	return s.toParams()
}

//Result() return prepare statement and params, use it in sqlx.
func (s *SqlBuilder) Result() (stmt string, params []any, err error) {
	stmt, err = s.PrepareString()
	params = s.Params()
	return
}

// RawString() return an executable statement.
func (s *SqlBuilder) RawString() (string, error) {
	err := Valid(s)
	if err != nil {
		return "", err
	}
	params := s.toParams()
	prepareStmt, err := s.PrepareString()
	if err != nil {
		return "", err
	}
	blocks := strings.Split(prepareStmt, SqlPlaceHolder)
	if len(blocks) != len(params)+1 {
		return "", ErrHolderParamNotMatch
	}

	for i, param := range params {
		switch param.(type) {
		case string:
			blocks[i] = fmt.Sprintf("%v\"%v\"", blocks[i], param)
		default:
			blocks[i] = fmt.Sprintf("%v%v", blocks[i], param)
		}
	}

	return strings.Join(blocks, ""), nil

}

//=== Common api ===

//Valid() validate whether stmt is legal or not.
// add dependencies of combine keyWords here.
// eg. SELECT must be used with FROM.
// most of validation are implemented in the standalone keyWord func.
func Valid(s *SqlBuilder) error {
	maskx := s.KeywordMask & actionMask
	func(maskx int) {
		switch maskx {
		case selectMask:
			if s.KeywordMask>>selectx%2 == 1 && s.KeywordMask>>from%2 == 0 {
				s.Errs[ErrSelectWithNoFrom] = true
			}
		case insertMask:
		case updateMask:
		default:
		}
	}(maskx)

	if len(s.Errs) == 0 {
		return nil
	}

	ErrTemplate := []string{}
	for k, _ := range s.Errs {
		ErrTemplate = append(ErrTemplate, k.Error())
	}
	return errors.New(strings.Join(ErrTemplate, "\n"))
}
