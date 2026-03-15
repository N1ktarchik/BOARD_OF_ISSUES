package store

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DeskUser struct {
	deskid int
	userid int
}

func TestConnectUserToDeskAndDeleteUserDesk(t *testing.T) {
	const (
		undefindUserId int = 789654123
		undefindDeskId int = 789654123
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []DeskUser{
		{1, 11},
		{2, 22},
		{3, 33},
		{4, 44},
		{5, 55},
		{6, 66},
		{7, 77},
		{8, 88},
		{9, 99},
		{10, 110},
	}

	for _, v := range tests {
		err := db.ConnectUserToDesk(ctx, v.userid, v.deskid)
		require.NoError(t, err)

		var actualDeskUser = &DeskUser{}
		query := `SELECT userid,deskid FROM desksusers where userid=$1 `
		err = db.db.QueryRow(ctx, query, v.userid).Scan(&actualDeskUser.userid, &actualDeskUser.deskid)
		require.NoError(t, err)

		assert.Equal(t, v.deskid, actualDeskUser.deskid)
		assert.Equal(t, v.userid, actualDeskUser.userid)

		err = db.DeleteUserDesk(ctx, v.userid, v.deskid)
		assert.NoError(t, err)

		var userID int
		query = "SELECT userid FROM desksusers where userid=$1 "
		err = db.db.QueryRow(ctx, query, v.userid).Scan(&userID)
		require.Error(t, err)

	}

	err = db.DeleteUserDesk(ctx, undefindUserId, undefindDeskId)
	assert.Equal(t, sql.ErrNoRows, err)

}

func TestGetUserDesks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	expected := [][]int{
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
		{11, 22, 33, 44, 55, 66, 77, 88, 99},
		{},
	}

	tests := []DeskUser{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
		{5, 1},
		{6, 1},
		{7, 1},
		{8, 1},
		{9, 1},

		{11, 2},
		{22, 2},
		{33, 2},
		{44, 2},
		{55, 2},
		{66, 2},
		{77, 2},
		{88, 2},
		{99, 2},
	}

	for _, v := range tests {
		err := db.ConnectUserToDesk(ctx, v.userid, v.deskid)
		require.NoError(t, err)
	}

	for i := 0; i < len(expected); i++ {
		mas, err := db.GetUserDesks(ctx, i+1)
		assert.NoError(t, err)

		assert.Equal(t, expected[i], mas)

	}

	for _, v := range tests {
		err := db.DeleteUserDesk(ctx, v.userid, v.deskid)
		assert.NoError(t, err)
	}

}

func TestCheckUserDesk(t *testing.T) {

	const (
		undefindUserId int = 789654123
		undefindDeskId int = 789654123
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []DeskUser{
		{1, 11},
		{2, 22},
		{3, 33},
		{4, 44},
		{5, 55},
		{6, 66},
		{7, 77},
		{8, 88},
		{9, 99},
		{10, 110},
	}

	for _, v := range tests {
		err := db.ConnectUserToDesk(ctx, v.userid, v.deskid)
		require.NoError(t, err)
	}

	for _, v := range tests {
		result, err := db.CheckUserDesk(ctx, v.userid, v.deskid)
		assert.NoError(t, err)

		assert.True(t, result)
	}

	result, err := db.CheckUserDesk(ctx, undefindUserId, undefindDeskId)
	assert.NoError(t, err)
	assert.False(t, result)

	for _, v := range tests {
		err := db.DeleteUserDesk(ctx, v.userid, v.deskid)
		assert.NoError(t, err)
	}

}
