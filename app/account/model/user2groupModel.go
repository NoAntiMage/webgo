package model

import (
	"context"
	"database/sql"
	"fmt"
	ctx "goweb/common/context"
	"goweb/common/sqlc"
	"goweb/common/stringx"
)

type User2GroupModel interface {
	GroupsListByUser(c context.Context, userId int64) (groupIds []int64, _ error)
	GroupsAddToUser(c context.Context, userId int64, groupIds []int64) (rowsAffected int64, _ error)
	GroupsDeleteFromUser(c context.Context, userId int64, groupIds []int64) (sql.Result, error)

	UserListByGroup(c context.Context, groupId int64) (userIds []int64, _ error)
	UsersAddToGroup(c context.Context, groupId int64, userIds []int64) (rowsAffected int64, _ error)
	UsersDeleteFromGroup(c context.Context, groupId int64, userIds []int64) (sql.Result, error)
}

type user2GroupModel struct {
	ctx.ModelContext
	table string
}

type User2Group struct {
	UserId   int64
	GroupxId int64
}

func NewUser2GroupModel(mctx ctx.ModelContext) User2GroupModel {
	return &user2GroupModel{
		ModelContext: mctx,
		table:        "user_groupx_m2m",
	}
}

func (ug *user2GroupModel) GroupsListByUser(c context.Context, userId int64) (groupIds []int64, _ error) {
	b := new(stringx.SqlBuilder)
	stmt, params, err :=
		b.Select(GroupTable + "_id").
			From(ug.table).
			Where(b.Eq(UserTable+"_id", userId)).
			Result()
	if err != nil {
		return nil, err
	}

	if err := ug.ModelContext.SelectContext(c, &groupIds, stmt, params...); err != nil {
		return nil, err
	}

	return groupIds, nil
}

func (ug *user2GroupModel) GroupsAddToUser(c context.Context, userId int64, groupIds []int64) (rowsAffected int64, _ error) {
	err := ug.Trans(c, func(tx *sqlc.Tx) error {
		sqlQuery := fmt.Sprintf("INSERT INTO %v (user_id, groupx_id) VALUES (?, ?)", ug.table)
		stmt, err := tx.PrepareContextWithLog(c, sqlQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, groupId := range groupIds {
			result, err := stmt.ExecContextWithLog(c, userId, groupId)
			if err != nil {
				return err
			}
			num, err := result.RowsAffected()
			if err != nil {
				return err
			}
			rowsAffected += num
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
func (ug *user2GroupModel) GroupsDeleteFromUser(c context.Context, userId int64, groupIds []int64) (sql.Result, error) {
	b := new(stringx.SqlBuilder)
	stmt, params, err :=
		b.Delete(ug.table).
			Where(b.Eq("user_id", userId)).
			And(b.In("groupx_id", groupIds)).
			Result()
	if err != nil {
		return nil, err
	}

	return ug.ModelContext.ExecContextWithLog(c, stmt, params...)
}

func (ug *user2GroupModel) UserListByGroup(c context.Context, groupId int64) (userIds []int64, _ error) {
	b := new(stringx.SqlBuilder)
	stmt, params, err :=
		b.Select(UserTable + "_id").
			From(ug.table).
			Where(b.Eq(GroupTable+"_id", groupId)).
			Result()
	if err != nil {
		return nil, err
	}

	if err := ug.ModelContext.SelectContext(c, &userIds, stmt, params...); err != nil {
		return nil, err
	}
	return userIds, nil
}

func (ug *user2GroupModel) UsersAddToGroup(c context.Context, groupId int64, userIds []int64) (rowsAffected int64, _ error) {
	err := ug.Trans(c, func(tx *sqlc.Tx) error {
		sqlQuery := fmt.Sprintf("INSERT INTO %v (user_id, groupx_id) VALUES (?, ?)", ug.table)
		stmt, err := tx.PrepareContextWithLog(c, sqlQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, userId := range userIds {
			result, err := stmt.ExecContextWithLog(c, userId, groupId)
			if err != nil {
				return err
			}
			num, err := result.RowsAffected()
			if err != nil {
				return err
			}
			rowsAffected += num
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (ug *user2GroupModel) UsersDeleteFromGroup(c context.Context, groupId int64, userIds []int64) (sql.Result, error) {
	b := new(stringx.SqlBuilder)
	stmt, params, err :=
		b.Delete(ug.table).
			Where(b.Eq("groupx_id", groupId)).
			And(b.In("user_id", userIds)).
			Result()
	if err != nil {
		return nil, err
	}

	return ug.ModelContext.ExecContextWithLog(c, stmt, params...)
}
