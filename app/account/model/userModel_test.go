package model

import (
	"context"
	mctx "goweb/common/context"
	"goweb/common/db"
	"goweb/common/logx"
	"testing"
)

var (
	textCtx = context.Background()
)

func TestInsert(t *testing.T) {
	db.DbInit()
	defer db.DbClose()
	logx.LoggerInit()

	u := User{
		Username: "tester0914",
		Realname: "testerReal0914",
	}
	t.Logf("insert user: %+v\n", u)

	ModelCtx := mctx.NewModelCtx()
	um := NewUserModel(ModelCtx)
	result, err := um.Insert(textCtx, &u)
	if err != nil {
		t.Log("err: ", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Log("err: ", err)
	}
	t.Log("LastInsertId: ", id)
	t.Error("END")
}

func TestSelect(t *testing.T) {
	db.DbInit()
	defer db.DbClose()
	logx.LoggerInit()

	ModelCtx := mctx.NewModelCtx()
	um := NewUserModel(ModelCtx)
	u, err := um.Select(textCtx, 7)
	if err != nil {
		t.Logf("%+v\n", err)
	}

	t.Logf("Select user: %+v\n", u)

	t.Error("END")

}

func TestUpdate(t *testing.T) {
	db.DbInit()
	defer db.DbClose()
	logx.LoggerInit()

	ModelCtx := mctx.NewModelCtx()
	um := NewUserModel(ModelCtx)

	u, err := um.Select(textCtx, 7)
	if err != nil {
		t.Logf("%+v\n", err)
		t.Fatalf("END")
	}
	u.Age = 21
	t.Logf("user: %+v\n", u)

	result, err := um.Update(textCtx, u)
	if err != nil {
		t.Log("err: ", err)
	}
	num, err := result.RowsAffected()
	if err != nil {
		t.Log("err: ", err)
	}
	t.Log("RowsAffected: ", num)
	t.Error("END")
}

func TestDelete(t *testing.T) {
	db.DbInit()
	defer db.DbClose()
	logx.LoggerInit()

	ModelCtx := mctx.NewModelCtx()
	um := NewUserModel(ModelCtx)

	result, err := um.Delete(textCtx, 3)
	if err != nil {
		t.Log("err: ", err)
	}
	num, err := result.RowsAffected()
	if err != nil {
		t.Log("err: ", err)
	}
	t.Log("RowsAffected: ", num)
	t.Error("END")
}
