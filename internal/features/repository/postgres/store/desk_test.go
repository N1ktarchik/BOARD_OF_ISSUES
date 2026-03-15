package store

import (
	"Board_of_issuses/internal/features/repository"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDeleteDesk(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.Desk{
		{OwnerId: 1, Password: "pass1", Name: "name1"},
		{OwnerId: 2, Password: "pass2", Name: "name2"},
		{OwnerId: 3, Password: "pass3", Name: "name3"},
		{OwnerId: 4, Password: "pass4", Name: "name4"},
		{OwnerId: 5, Password: "pass5", Name: "name5"},
		{OwnerId: 6, Password: "pass6", Name: "name6"},
		{OwnerId: 7, Password: "pass7", Name: "name7"},
		{OwnerId: 8, Password: "pass8", Name: "name8"},
		{OwnerId: 9, Password: "pass9", Name: "name9"},
		{OwnerId: 10, Password: "pass10", Name: "name10"},
	}

	desks_ids := make([]int, len(tests))

	for i, v := range tests {
		err := db.CreateDesk(ctx, &v)
		require.NoError(t, err)

		query := `SELECT id,name,password,ownerid,created_at FROM desks WHERE ownerid=$1`
		var actualDesk = &repository.Desk{}

		err = db.db.QueryRow(ctx, query, v.OwnerId).Scan(
			&actualDesk.Id,
			&actualDesk.Name,
			&actualDesk.Password,
			&actualDesk.OwnerId,
			&actualDesk.Created_at)

		assert.NoError(t, err)

		assert.Equal(t, v.OwnerId, actualDesk.OwnerId)
		assert.Equal(t, v.Password, actualDesk.Password)
		assert.Equal(t, v.Name, actualDesk.Name)
		assert.Greater(t, actualDesk.Id, 0)

		assert.True(t, v.Created_at.Before(time.Now().Add(5*time.Minute)))

		desks_ids[i] = actualDesk.Id

		query = `SELECT userid,deskid FROM desksusers WHERE deskid=$1`
		var actualDeskUser = &DeskUser{}

		err = db.db.QueryRow(ctx, query, desks_ids[i]).Scan(&actualDeskUser.userid, &actualDeskUser.deskid)
		assert.NoError(t, err)

		assert.Equal(t, v.OwnerId, actualDeskUser.userid)
		assert.Equal(t, desks_ids[i], actualDeskUser.deskid)
	}

	for _, v := range desks_ids {
		err := db.DeleteDesk(ctx, v)
		assert.NoError(t, err)

		query := `SELECT id,name,password,ownerid,created_at FROM desks WHERE id=$1`
		var actualDesk = &repository.Desk{}

		err = db.db.QueryRow(ctx, query, v).Scan(
			&actualDesk.Id,
			&actualDesk.Name,
			&actualDesk.Password,
			&actualDesk.OwnerId,
			&actualDesk.Created_at)

		assert.Equal(t, pgx.ErrNoRows, err)

		query = `SELECT userid,deskid FROM desksusers WHERE deskid=$1`
		var actualDeskUser = &DeskUser{}

		err = db.db.QueryRow(ctx, query, v).Scan(&actualDeskUser.userid, &actualDeskUser.deskid)
		assert.Equal(t, pgx.ErrNoRows, err)

	}
}

func TestUpdateDeskNameOwnerPassword(t *testing.T) {
	const (
		newPass  string = "newPass"
		newOwner int    = 6666
		newName  string = "newName"
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.Desk{
		{OwnerId: 1, Password: "pass1", Name: "name1"},
		{OwnerId: 2, Password: "pass2", Name: "name2"},
		{OwnerId: 3, Password: "pass3", Name: "name3"},
		{OwnerId: 4, Password: "pass4", Name: "name4"},
		{OwnerId: 5, Password: "pass5", Name: "name5"},
		{OwnerId: 6, Password: "pass6", Name: "name6"},
		{OwnerId: 7, Password: "pass7", Name: "name7"},
		{OwnerId: 8, Password: "pass8", Name: "name8"},
		{OwnerId: 9, Password: "pass9", Name: "name9"},
		{OwnerId: 10, Password: "pass10", Name: "name10"},
	}

	desks_ids := make([]int, len(tests))

	for idx, v := range tests {
		err := db.CreateDesk(ctx, &v)
		require.NoError(t, err)

		var id int
		query := `SELECT id FROM desks WHERE ownerid=$1 AND password=$2 AND name=$3`
		err = db.db.QueryRow(ctx, query, v.OwnerId, v.Password, v.Name).Scan(&id)

		require.NoError(t, err)
		require.Greater(t, id, 0)

		desks_ids[idx] = id

	}

	for i := range desks_ids {
		err := db.UpdateDeskName(ctx, desks_ids[i], newName)
		assert.NoError(t, err)
		err = db.UpdateDeskOwner(ctx, newOwner, desks_ids[i])
		assert.NoError(t, err)
		err = db.UpdateDesksPassword(ctx, desks_ids[i], newPass)
		assert.NoError(t, err)

		query := `SELECT id,name,password,ownerid,created_at FROM desks WHERE id=$1`
		var actualDesk = &repository.Desk{}

		err = db.db.QueryRow(ctx, query, desks_ids[i]).Scan(
			&actualDesk.Id,
			&actualDesk.Name,
			&actualDesk.Password,
			&actualDesk.OwnerId,
			&actualDesk.Created_at)

		assert.NoError(t, err)

		assert.Equal(t, newOwner, actualDesk.OwnerId)
		assert.Equal(t, newPass, actualDesk.Password)
		assert.Equal(t, newName, actualDesk.Name)

		err = db.DeleteDesk(ctx, desks_ids[i])
		assert.NoError(t, err)

	}
}

func TestChekDeskOwnerPassword(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.Desk{
		{OwnerId: 1, Password: "pass1", Name: "name1"},
		{OwnerId: 2, Password: "pass2", Name: "name2"},
		{OwnerId: 3, Password: "pass3", Name: "name3"},
		{OwnerId: 4, Password: "pass4", Name: "name4"},
		{OwnerId: 5, Password: "pass5", Name: "name5"},
		{OwnerId: 6, Password: "pass6", Name: "name6"},
		{OwnerId: 7, Password: "pass7", Name: "name7"},
		{OwnerId: 8, Password: "pass8", Name: "name8"},
		{OwnerId: 9, Password: "pass9", Name: "name9"},
		{OwnerId: 10, Password: "pass10", Name: "name10"},
	}

	desks_ids := make([]int, len(tests))

	for idx, v := range tests {
		err := db.CreateDesk(ctx, &v)
		require.NoError(t, err)

		var id int
		query := `SELECT id FROM desks WHERE ownerid=$1 AND password=$2 AND name=$3`
		err = db.db.QueryRow(ctx, query, v.OwnerId, v.Password, v.Name).Scan(&id)

		require.NoError(t, err)
		require.Greater(t, id, 0)

		desks_ids[idx] = id

	}

	for i, deskID := range desks_ids {
		ownerId, err := db.CheckDeskOwner(ctx, deskID)
		assert.NoError(t, err)

		password, err := db.CheckDeskPassword(ctx, deskID)
		assert.NoError(t, err)

		assert.Equal(t, tests[i].OwnerId, ownerId)
		assert.Equal(t, tests[i].Password, password)

		err = db.DeleteDesk(ctx, deskID)
		assert.NoError(t, err)

	}
}
