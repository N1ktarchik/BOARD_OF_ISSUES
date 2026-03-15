package store

import (
	"Board_of_issuses/internal/features/repository"
	"context"
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	//введите для подключения к ДБ
	CONSTR string = "your url"
)

func createPool(ctx context.Context) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(CONSTR)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
func TestCreateGetCheckDeleteUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.User{
		{Login: "login1", Password: "pass1", Email: "email1", Name: "name1"},
		{Login: "login2", Password: "pass2", Email: "email2", Name: "name2"},
		{Login: "login3", Password: "pass3", Email: "email3", Name: "name3"},
		{Login: "login4", Password: "pass4", Email: "email4", Name: "name4"},
		{Login: "login5", Password: "pass5", Email: "email5", Name: "name5"},
		{Login: "login6", Password: "pass6", Email: "email6", Name: "name6"},
		{Login: "login7", Password: "pass7", Email: "email7", Name: "name7"},
		{Login: "login8", Password: "pass8", Email: "email8", Name: "name8"},
		{Login: "login9", Password: "pass9", Email: "email9", Name: "name9"},
		{Login: "login10", Password: "pass10", Email: "email10", Name: "name10"},
	}

	ids := make([]int, len(tests))

	for idx, v := range tests {
		id, err := db.CreateUser(ctx, &v)

		require.NoError(t, err)
		assert.NotEqual(t, 0, id)

		ids[idx] = id

		if idx > 0 {
			assert.Equal(t, ids[idx-1]+1, id)
		}

	}

	for idx := range ids {
		usr, err := db.GetUserByID(ctx, ids[idx])
		assert.NoError(t, err)

		assert.Equal(t, tests[idx].Login, usr.Login)
		assert.Equal(t, tests[idx].Password, usr.Password)
		assert.Equal(t, tests[idx].Email, usr.Email)
		assert.Equal(t, tests[idx].Name, usr.Name)

		usr, err = db.GetUserByLoginOrEmail(ctx, tests[idx].Login, "")
		assert.NoError(t, err)

		assert.Equal(t, tests[idx].Login, usr.Login)
		assert.Equal(t, tests[idx].Password, usr.Password)
		assert.Equal(t, tests[idx].Email, usr.Email)
		assert.Equal(t, tests[idx].Name, usr.Name)

		usr, err = db.GetUserByLoginOrEmail(ctx, "", tests[idx].Email)
		assert.NoError(t, err)

		assert.Equal(t, tests[idx].Login, usr.Login)
		assert.Equal(t, tests[idx].Password, usr.Password)
		assert.Equal(t, tests[idx].Email, usr.Email)
		assert.Equal(t, tests[idx].Name, usr.Name)
	}

	_, err = db.GetUserByID(ctx, ids[len(ids)-1]+100)
	assert.Equal(t, pgx.ErrNoRows, err)

	_, err = db.GetUserByLoginOrEmail(ctx, "", "")
	assert.Equal(t, pgx.ErrNoRows, err)
	_, err = db.GetUserByLoginOrEmail(ctx, "f[skpodpdopdopdopo]", "")
	assert.Equal(t, pgx.ErrNoRows, err)
	_, err = db.GetUserByLoginOrEmail(ctx, "f[skpodpdopdopdopo]", "dklasnjndsaj;dla;sdlajs;ld")
	assert.Equal(t, pgx.ErrNoRows, err)
	_, err = db.GetUserByLoginOrEmail(ctx, "", "dklasnjndsaj;dla;sdlajs;ld")
	assert.Equal(t, pgx.ErrNoRows, err)

	for idx := range ids {
		value, err := db.CheckUserByEmailOrLogin(ctx, tests[idx].Login, "")
		assert.NoError(t, err)

		assert.True(t, value)

		value, err = db.CheckUserByEmailOrLogin(ctx, "", tests[idx].Email)
		assert.NoError(t, err)

		assert.True(t, value)
	}

	value, err := db.CheckUserByEmailOrLogin(ctx, "", "")
	assert.False(t, value)
	value, err = db.CheckUserByEmailOrLogin(ctx, "dsaasfaifjewqpo[FJKWEQ]PO[FK]WE", "")
	assert.False(t, value)
	value, err = db.CheckUserByEmailOrLogin(ctx, "", "e[poafk][PEFKEPWQklfwe[pfw=]]")
	assert.False(t, value)
	value, err = db.CheckUserByEmailOrLogin(ctx, "oiqafd[0PJFWEOo]fqweqw", "qdo]w[pqkf[eqwopfk]ope[jfvwepi]]")
	assert.False(t, value)

	for idx := range ids {
		err := db.DeleteUser(ctx, strconv.Itoa(ids[idx]))
		assert.NoError(t, err)

		_, err = db.GetUserByID(ctx, ids[idx])
		assert.Equal(t, pgx.ErrNoRows, err)
	}

	err = db.DeleteUser(ctx, strconv.Itoa(ids[len(ids)-1]+100))
	assert.Equal(t, sql.ErrNoRows, err)

}

func TestUpdateNameEmailPasswordUser(t *testing.T) {

	const (
		newPass  string = "newPass"
		newEmail string = "newEmail"
		newName  string = "newName"
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.User{
		{Login: "login1", Password: "pass1", Email: "email1", Name: "name1"},
		{Login: "login2", Password: "pass2", Email: "email2", Name: "name2"},
		{Login: "login3", Password: "pass3", Email: "email3", Name: "name3"},
		{Login: "login4", Password: "pass4", Email: "email4", Name: "name4"},
		{Login: "login5", Password: "pass5", Email: "email5", Name: "name5"},
		{Login: "login6", Password: "pass6", Email: "email6", Name: "name6"},
		{Login: "login7", Password: "pass7", Email: "email7", Name: "name7"},
		{Login: "login8", Password: "pass8", Email: "email8", Name: "name8"},
		{Login: "login9", Password: "pass9", Email: "email9", Name: "name9"},
		{Login: "login10", Password: "pass10", Email: "email10", Name: "name10"},
	}

	ids := make([]int, len(tests))

	for idx, v := range tests {
		id, err := db.CreateUser(ctx, &v)
		require.NoError(t, err)

		ids[idx] = id
	}

	for idx := range ids {
		err := db.UpdateUserEmail(ctx, newEmail, ids[idx])
		assert.NoError(t, err)
		err = db.UpdateUserName(ctx, newName, ids[idx])
		assert.NoError(t, err)
		err = db.UpdateUserPassword(ctx, newPass, ids[idx])
		assert.NoError(t, err)

		usr, err := db.GetUserByID(ctx, ids[idx])
		assert.NoError(t, err)

		assert.Equal(t, tests[idx].Login, usr.Login)
		assert.Equal(t, newEmail, usr.Email)
		assert.Equal(t, newPass, usr.Password)
		assert.Equal(t, newName, usr.Name)

		err = db.DeleteUser(ctx, strconv.Itoa(ids[idx]))
		assert.NoError(t, err)

	}
}
