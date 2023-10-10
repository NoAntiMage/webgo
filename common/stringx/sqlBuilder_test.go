package stringx

import (
	"testing"
)

// Sample table
/*
CREATE TABLE "user" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "is_deleted" integer NOT NULL DEFAULT 0,
  "username" TEXT NOT NULL,
  "realname" TEXT NOT NULL,
  "age" INTEGER,
  "gender" INTEGER,
  "email" TEXT
);

CREATE TABLE "groupx" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT NOT NULL
);

CREATE TABLE "user_groupx_m2m" (
  "user_id" INTEGER NOT NULL,
  "groupx_id" INTEGER NOT NULL,
  CONSTRAINT "userId_groupxId" UNIQUE ("user_id" COLLATE BINARY ASC, "groupx_id" COLLATE BINARY ASC) ON CONFLICT FAIL
);
*/

// INSERT user(username, realname, gender) VALUES("wujimaster", "wuji", 1)
func TestInsertInto(t *testing.T) {
	b := NewSqlBuilder()
	kvs := map[string]any{
		"username": "wujimaster",
		"realname": "wuji",
		"gender":   1,
	}

	handleResult(t, b.Insert("user", kvs).Result)

	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")

}

// SELECT id, username, realname, email FROM user
func TestSelectFrom(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("id", "username", "realname", "email").
		From("user").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// UPDATE user u SET username = "wuji"
func TestUpdate(t *testing.T) {
	b := NewSqlBuilder()
	kvs := map[string]any{
		"username": "wuji",
	}

	handleResult(t, b.Update("user", kvs).Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// DeleteSoft
func TestDeleteSoft(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.DeleteSoft("user").Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname, age, gender FROM user u JOIN user_groupx_m2m ug on u.id = ug.user_id JOIN groupx g on g.id = ug.groupx_id
func TestJoin(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname", "age", "gender").
		From("user u").
		Join("user_groupx_m2m ug", b.Eq("u.id", "ug.user_id")).
		Join("groupx g", b.Eq("g.id", "ug.groupx_id")).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname, age, gender FROM user u LEFT JOIN user_groupx_m2m ug on u.id = ug.user_id LEFT JOIN groupx g on g.id = ug.groupx_id
func TestLeftJoin(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname", "age", "gender").
		From("user u").
		LeftJoin("user_groupx_m2m ug", b.Eq("u.id", "ug.user_id")).
		LeftJoin("groupx g", b.Eq("g.id", "ug.groupx_id")).
		Result)

	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname, age, gender FROM user u RIGHT JOIN user_groupx_m2m ug on u.id = ug.user_id RIGHT JOIN groupx g on g.id = ug.groupx_id
func TestRightJoin(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname", "age", "gender").
		From("user u").
		RightJoin("user_groupx_m2m ug", b.Eq("u.id", "ug.user_id")).
		RightJoin("groupx g", b.Eq("g.id", "ug.groupx_id")).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE id = 1
func TestWhereInt(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where("id = 1").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE username = "wuji"
func TestWhereString(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where(b.Eq("username", "wuji")).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE gender = 1 AND age < 40
func TestWhereAnd(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where("gender = 1").
		And("age < 40").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE 1 = 1 AND age < 40
func TestWhereNoneAnd(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where("").
		And("gender = 1").
		And("age < 40").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE gender = 0 OR age > 40
func TestWhereOr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where("gender = 0").
		Or("age > 40").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE username IN ("foo", "bar")
func TestWhereInString(t *testing.T) {
	b := NewSqlBuilder()

	names := []string{"first", "second", "third"}

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where(b.In("username", names)).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE age IN (12,20)
func TestWhereInInt(t *testing.T) {
	b := NewSqlBuilder()

	ages := []int64{12, 20}

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where(b.In("age", ages)).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE age BETWEEN(30, 40)
func TestWhereBetween(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where(b.Between("age", 30, 40)).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname FROM user WHERE username LIKE "%wu%"
func TestWhereLike(t *testing.T) {
	b := NewSqlBuilder()
	likePattern := "wu"

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where(b.Like("username", likePattern)).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname, MIN(age) FROM user
func TestMin(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname", b.Min("age")).
		From("user").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT username, realname, MAX(age) FROM user
func TestMax(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname", b.Max("age")).
		From("user").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT gender, AVG(age) FROM user GROUP BY gender
func TestAvg(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("gender", b.Avg("age")).
		From("user").
		GroupBy("gender").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT gender, COUNT(age) FROM user GROUP BY gender
func TestCount(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("gender", b.Count("age")).
		From("user").
		GroupBy("gender").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT SUM(age) FROM user
func TestSum(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select(b.Sum("age")).
		From("user").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT gender, COUNT(id) FROM user GROUP BY gender HAVING COUNT(id) > 2
func TestGroupBy(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("gender", b.Count("id")).
		From("user").
		GroupBy("gender").
		Having(b.Count(b.Gt("id", 2))).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT age, username, realname FROM user ORDER BY age DESC
func TestOrderBy(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("age", "username", "realname").
		From("user").
		OrderBy("age", true).
		Limit(2).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT age, username, realname FROM user ORDER BY age ASC LIMIT 2
func TestLimit(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("age", "username", "realname").
		From("user").
		OrderBy("age", false).
		Limit(2).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// SELECT age, username, realname FROM user ORDER BY age ASC LIMIT 2 OFFSET 1
func TestOffset(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("age", "username", "realname").
		From("user").
		OrderBy("age", false).
		Limit(2).
		Offset(1).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

// testcase of ToRowString()
func TestToRowString(t *testing.T) {}

//TODO=== Error Process Test ===

func TestSelectNoneColumnErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select().
		From("user").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")

}

func TestInsertNoneColumnErr(t *testing.T) {
	b := NewSqlBuilder()
	kvs := make(map[string]any, 0)

	handleResult(t, b.
		Insert("user", kvs).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}
func TestUpdateNoneColumnErr(t *testing.T) {
	b := NewSqlBuilder()
	kvs := make(map[string]any, 0)

	handleResult(t, b.
		Update("user", kvs).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestSelectWithNoFromErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestActionExistErr(t *testing.T) {
	b := NewSqlBuilder()
	kvs := make(map[string]any, 0)

	handleResult(t, b.
		Select("username").
		Insert("username", kvs).
		From("user").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestActionNotExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		From("user").
		Where("").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestFromNotExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		Join("user_groupx_m2m ug", "u.id = ug.user_id").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestFromExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where("").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestWhereExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where("").
		Where(b.Eq("gender", 0)).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestWhereNotExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		And(b.Ge("age", 30)).
		Or(b.Eq("gender", 1)).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestBetweenNotMatchErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Where(b.Between("age", 31, "hello")).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestGroupWithNoAggregateErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		GroupBy("gender").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestGroupBeforeFromErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("gender", b.Avg("age")).
		GroupBy("gender").
		From("user").
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestHavingWithoutGroupErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname", b.Count("gender")).
		From("user").
		Having(b.Gt(b.Count("gender"), 1)).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestOrderByExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		OrderBy("gender", true).
		OrderBy("gender", true).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestOrderByAfterLimitErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Limit(10).
		OrderBy("gender", true).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestLimitExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Limit(10).
		Limit(10).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestLimitNotExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Offset(2).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

func TestOffsetExistErr(t *testing.T) {
	b := NewSqlBuilder()

	handleResult(t, b.
		Select("username", "realname").
		From("user").
		Limit(10).
		Offset(2).
		Offset(2).
		Result)
	t.Logf("key mask: %b", b.KeywordMask)
	t.Error("END")
}

//=== sqlCondFn ===

//=== Util func ===

// func handleResult(fn func() (stmt string, params []any, err error)) {
func handleResult(t *testing.T, fn func() (stmt string, params []any, err error)) {
	var stmtStr string
	var params []any
	var err error

	stmtStr, params, err = fn()

	t.Log("stmt: ", stmtStr)
	if len(params) != 0 {
		t.Log("params: ", params)
	}
	if err != nil {
		t.Log("err: ", err.Error())
	}
}
