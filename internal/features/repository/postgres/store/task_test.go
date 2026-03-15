package store

import (
	"Board_of_issuses/internal/features/repository"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDeleteTask(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.Task{
		{UserId: 1, DeskId: 1, Name: "name1", Description: "descr1", Time: time.Now().UTC().Add(1 * time.Hour)},
		{UserId: 2, DeskId: 1, Name: "name2", Description: "descr2", Time: time.Now().UTC().Add(2 * time.Hour)},
		{UserId: 3, DeskId: 2, Name: "name3", Description: "descr3", Time: time.Now().UTC().Add(3 * time.Hour)},
		{UserId: 4, DeskId: 2, Name: "name4", Description: "descr4", Time: time.Now().UTC().Add(4 * time.Hour)},
		{UserId: 5, DeskId: 3, Name: "name5", Description: "descr5", Time: time.Now().UTC().Add(5 * time.Hour)},
		{UserId: 6, DeskId: 3, Name: "name6", Description: "descr6", Time: time.Now().UTC().Add(6 * time.Hour)},
		{UserId: 7, DeskId: 4, Name: "name7", Description: "descr7", Time: time.Now().UTC().Add(7 * time.Hour)},
		{UserId: 8, DeskId: 4, Name: "name8", Description: "descr8", Time: time.Now().UTC().Add(8 * time.Hour)},
		{UserId: 9, DeskId: 5, Name: "name9", Description: "descr9", Time: time.Now().UTC().Add(9 * time.Hour)},
		{UserId: 10, DeskId: 5, Name: "name10", Description: "descr10", Time: time.Now().UTC().Add(10 * time.Hour)},
	}

	ids := make([]int, len(tests))

	for i, v := range tests {
		err := db.CreateTask(ctx, &v)
		assert.NoError(t, err)

		var actualTask = &repository.Task{}

		query := `SELECT id,userid,deskid,name,description,done,time,created_at FROM tasks WHERE deskid = $1 and userid = $2`
		err = db.db.QueryRow(ctx, query, v.DeskId, v.UserId).Scan(
			&actualTask.Id,
			&actualTask.UserId,
			&actualTask.DeskId,
			&actualTask.Name,
			&actualTask.Description,
			&actualTask.Done,
			&actualTask.Time,
			&actualTask.Created_at)

		require.NoError(t, err)

		assert.Equal(t, v.UserId, actualTask.UserId)
		assert.Equal(t, v.DeskId, actualTask.DeskId)
		assert.Equal(t, v.Name, actualTask.Name)
		assert.Equal(t, v.Description, actualTask.Description)
		assert.Equal(t, v.Time.UTC().Truncate(time.Millisecond), actualTask.Time.UTC().Truncate(time.Millisecond))
		assert.True(t, v.Created_at.Before(time.Now().Add(5*time.Minute)))

		ids[i] = actualTask.Id
		err = db.DeleteTask(ctx, actualTask.Id)
		assert.NoError(t, err)

		query = `SELECT userid FROM tasks WHERE id=$1`

		var undefindUserID int
		err = db.db.QueryRow(ctx, query, actualTask.Id).Scan(&undefindUserID)
		assert.Contains(t, err.Error(), "no rows")

	}
}

func TestUpdateDescriptionTimeDoneTask(t *testing.T) {

	const (
		newDescr string = "newDescr"
	)
	var (
		newTime time.Time = time.Now().UTC().Add(123 * time.Hour)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.Task{
		{UserId: 1, DeskId: 1, Name: "name1", Description: "descr1", Time: time.Now().UTC().Add(1 * time.Hour)},
		{UserId: 2, DeskId: 1, Name: "name2", Description: "descr2", Time: time.Now().UTC().Add(2 * time.Hour)},
		{UserId: 3, DeskId: 2, Name: "name3", Description: "descr3", Time: time.Now().UTC().Add(3 * time.Hour)},
		{UserId: 4, DeskId: 2, Name: "name4", Description: "descr4", Time: time.Now().UTC().Add(4 * time.Hour)},
		{UserId: 5, DeskId: 3, Name: "name5", Description: "descr5", Time: time.Now().UTC().Add(5 * time.Hour)},
		{UserId: 6, DeskId: 3, Name: "name6", Description: "descr6", Time: time.Now().UTC().Add(6 * time.Hour)},
		{UserId: 7, DeskId: 4, Name: "name7", Description: "descr7", Time: time.Now().UTC().Add(7 * time.Hour)},
		{UserId: 8, DeskId: 4, Name: "name8", Description: "descr8", Time: time.Now().UTC().Add(8 * time.Hour)},
		{UserId: 9, DeskId: 5, Name: "name9", Description: "descr9", Time: time.Now().UTC().Add(9 * time.Hour)},
		{UserId: 10, DeskId: 5, Name: "name10", Description: "descr10", Time: time.Now().UTC().Add(10 * time.Hour)},
	}

	ids := make([]int, len(tests))

	for i, v := range tests {
		err := db.CreateTask(ctx, &v)
		assert.NoError(t, err)

		var actualTask = &repository.Task{}

		query := `SELECT id,userid,deskid,name,description,done,time,created_at FROM tasks WHERE deskid = $1 and userid = $2`
		err = db.db.QueryRow(ctx, query, v.DeskId, v.UserId).Scan(
			&actualTask.Id,
			&actualTask.UserId,
			&actualTask.DeskId,
			&actualTask.Name,
			&actualTask.Description,
			&actualTask.Done,
			&actualTask.Time,
			&actualTask.Created_at)

		require.NoError(t, err)

		ids[i] = actualTask.Id
	}

	for i, taskId := range ids {
		err := db.UpdateTaskDecription(ctx, taskId, newDescr)
		assert.NoError(t, err)
		err = db.UpdateTaskDone(ctx, taskId)
		assert.NoError(t, err)
		err = db.UpdateTaskTime(ctx, taskId, newTime)
		assert.NoError(t, err)

		var actualTask = &repository.Task{}
		query := `SELECT id,userid,deskid,name,description,done,time,created_at FROM tasks WHERE id=$1`
		err = db.db.QueryRow(ctx, query, taskId).Scan(
			&actualTask.Id,
			&actualTask.UserId,
			&actualTask.DeskId,
			&actualTask.Name,
			&actualTask.Description,
			&actualTask.Done,
			&actualTask.Time,
			&actualTask.Created_at)

		assert.NoError(t, err)

		assert.Equal(t, tests[i].UserId, actualTask.UserId)
		assert.Equal(t, tests[i].DeskId, actualTask.DeskId)
		assert.Equal(t, tests[i].Name, actualTask.Name)
		assert.Equal(t, newDescr, actualTask.Description)
		assert.WithinDuration(t, newTime, actualTask.Time, time.Millisecond)
		assert.True(t, tests[i].Created_at.Before(time.Now().Add(5*time.Minute)))

		err = db.DeleteTask(ctx, taskId)
		assert.NoError(t, err)

	}
}

func TestGetAllTasksFromOneDesk(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.Task{
		{UserId: 1, DeskId: 1, Name: "name1", Description: "descr1", Time: time.Now().UTC().Add(1 * time.Hour)},
		{UserId: 2, DeskId: 1, Name: "name2", Description: "descr2", Time: time.Now().UTC().Add(2 * time.Hour)},
		{UserId: 3, DeskId: 2, Name: "name3", Description: "descr3", Time: time.Now().UTC().Add(3 * time.Hour)},
		{UserId: 4, DeskId: 2, Name: "name4", Description: "descr4", Time: time.Now().UTC().Add(4 * time.Hour)},
		{UserId: 5, DeskId: 3, Name: "name5", Description: "descr5", Time: time.Now().UTC().Add(5 * time.Hour)},
		{UserId: 6, DeskId: 3, Name: "name6", Description: "descr6", Time: time.Now().UTC().Add(6 * time.Hour)},
		{UserId: 7, DeskId: 4, Name: "name7", Description: "descr7", Time: time.Now().UTC().Add(7 * time.Hour)},
		{UserId: 8, DeskId: 4, Name: "name8", Description: "descr8", Time: time.Now().UTC().Add(8 * time.Hour)},
		{UserId: 9, DeskId: 5, Name: "name9", Description: "descr9", Time: time.Now().UTC().Add(9 * time.Hour)},
		{UserId: 10, DeskId: 5, Name: "name10", Description: "descr10", Time: time.Now().UTC().Add(10 * time.Hour)},
	}

	result := [][]repository.Task{}
	desks_id := []int{1, 2, 3, 4, 5}

	for _, v := range tests {
		err := db.CreateTask(ctx, &v)
		assert.NoError(t, err)

	}

	for _, deskID := range desks_id {
		actualTask, err := db.GetAllTasksFromOneDesk(ctx, deskID)
		assert.NoError(t, err)

		result = append(result, actualTask)
	}

	test_number := 0
	for i := 0; i < len(result); i++ {
		tasks := result[i]

		for _, actualTask := range tasks {
			assert.Equal(t, tests[test_number].UserId, actualTask.UserId)
			assert.Equal(t, tests[test_number].DeskId, actualTask.DeskId)
			assert.Equal(t, tests[test_number].Name, actualTask.Name)
			assert.Equal(t, tests[test_number].Description, actualTask.Description)
			assert.WithinDuration(t, tests[test_number].Time, actualTask.Time, time.Millisecond)

			test_number++

			err = db.DeleteTask(ctx, actualTask.Id)
			assert.NoError(t, err)
		}

	}
}

func TestGetTaskOwnerDeskId(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	tests := []repository.Task{
		{UserId: 1, DeskId: 1, Name: "name1", Description: "descr1", Time: time.Now().UTC().Add(1 * time.Hour)},
		{UserId: 2, DeskId: 1, Name: "name2", Description: "descr2", Time: time.Now().UTC().Add(2 * time.Hour)},
		{UserId: 3, DeskId: 2, Name: "name3", Description: "descr3", Time: time.Now().UTC().Add(3 * time.Hour)},
		{UserId: 4, DeskId: 2, Name: "name4", Description: "descr4", Time: time.Now().UTC().Add(4 * time.Hour)},
		{UserId: 5, DeskId: 3, Name: "name5", Description: "descr5", Time: time.Now().UTC().Add(5 * time.Hour)},
		{UserId: 6, DeskId: 3, Name: "name6", Description: "descr6", Time: time.Now().UTC().Add(6 * time.Hour)},
		{UserId: 7, DeskId: 4, Name: "name7", Description: "descr7", Time: time.Now().UTC().Add(7 * time.Hour)},
		{UserId: 8, DeskId: 4, Name: "name8", Description: "descr8", Time: time.Now().UTC().Add(8 * time.Hour)},
		{UserId: 9, DeskId: 5, Name: "name9", Description: "descr9", Time: time.Now().UTC().Add(9 * time.Hour)},
		{UserId: 10, DeskId: 5, Name: "name10", Description: "descr10", Time: time.Now().UTC().Add(10 * time.Hour)},
	}

	ids := make([]int, len(tests))

	for i, v := range tests {
		err := db.CreateTask(ctx, &v)
		assert.NoError(t, err)

		var actualTask = &repository.Task{}

		query := `SELECT id,userid,deskid,name,description,done,time,created_at FROM tasks WHERE deskid = $1 and userid = $2`
		err = db.db.QueryRow(ctx, query, v.DeskId, v.UserId).Scan(
			&actualTask.Id,
			&actualTask.UserId,
			&actualTask.DeskId,
			&actualTask.Name,
			&actualTask.Description,
			&actualTask.Done,
			&actualTask.Time,
			&actualTask.Created_at)

		require.NoError(t, err)

		ids[i] = actualTask.Id
	}

	for i, taskId := range ids {
		actualOwner, err := db.GetTaskOwner(ctx, taskId)
		assert.NoError(t, err)
		actualDeskId, err := db.GetDeskIDByTask(ctx, taskId)
		assert.NoError(t, err)

		assert.Equal(t, tests[i].UserId, actualOwner)
		assert.Equal(t, tests[i].DeskId, actualDeskId)

		err = db.DeleteTask(ctx, taskId)
		assert.NoError(t, err)

	}
}

func TestGetTasksWithParams(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := createTestPool(ctx)
	require.NoError(t, err)

	db := CreateConnectToDB(pool)

	baseTime := time.Now().UTC().Truncate(time.Millisecond)

	tasks := []repository.Task{
		{UserId: 1, DeskId: 1, Name: "task1", Description: "desc1", Done: false, Time: baseTime.Add(1 * time.Hour)},
		{UserId: 2, DeskId: 1, Name: "task2", Description: "desc2", Done: true, Time: baseTime.Add(2 * time.Hour)},
		{UserId: 3, DeskId: 1, Name: "task3", Description: "desc3", Done: false, Time: baseTime.Add(3 * time.Hour)},
		{UserId: 4, DeskId: 2, Name: "task4", Description: "desc4", Done: true, Time: baseTime.Add(4 * time.Hour)},
		{UserId: 5, DeskId: 2, Name: "task5", Description: "desc5", Done: false, Time: baseTime.Add(5 * time.Hour)},
		{UserId: 6, DeskId: 3, Name: "task6", Description: "desc6", Done: true, Time: baseTime.Add(6 * time.Hour)},
		{UserId: 7, DeskId: 3, Name: "task7", Description: "desc7", Done: false, Time: baseTime.Add(7 * time.Hour)},
	}

	ids := make([]int, len(tasks))

	for i, v := range tasks {
		err := db.CreateTask(ctx, &v)
		require.NoError(t, err)

		var actualTask = &repository.Task{}

		query := `SELECT id,userid,deskid,name,description,done,time,created_at FROM tasks WHERE deskid = $1 and userid = $2`
		err = db.db.QueryRow(ctx, query, v.DeskId, v.UserId).Scan(
			&actualTask.Id,
			&actualTask.UserId,
			&actualTask.DeskId,
			&actualTask.Name,
			&actualTask.Description,
			&actualTask.Done,
			&actualTask.Time,
			&actualTask.Created_at)

		require.NoError(t, err)

		ids[i] = actualTask.Id

		if v.Done {
			err = db.UpdateTaskDone(ctx, actualTask.Id)
			require.NoError(t, err)
		}
	}

	result1, err := db.GetTasksWithParams(ctx, 1, true)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result1))
	if len(result1) > 0 {
		assert.Equal(t, 1, result1[0].DeskId)
		assert.Equal(t, true, result1[0].Done)
	}

	result2, err := db.GetTasksWithParams(ctx, 1, false)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result2))
	for _, task := range result2 {
		assert.Equal(t, 1, task.DeskId)
		assert.Equal(t, false, task.Done)
	}

	result3, err := db.GetTasksWithParams(ctx, 2, true)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result3))
	if len(result3) > 0 {
		assert.Equal(t, 2, result3[0].DeskId)
		assert.Equal(t, true, result3[0].Done)
	}

	result4, err := db.GetTasksWithParams(ctx, 2, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result4))
	if len(result4) > 0 {
		assert.Equal(t, 2, result4[0].DeskId)
		assert.Equal(t, false, result4[0].Done)
	}

	result5, err := db.GetTasksWithParams(ctx, 3, true)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result5))
	if len(result5) > 0 {
		assert.Equal(t, 3, result5[0].DeskId)
		assert.Equal(t, true, result5[0].Done)
	}

	result6, err := db.GetTasksWithParams(ctx, 3, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result6))
	if len(result6) > 0 {
		assert.Equal(t, 3, result6[0].DeskId)
		assert.Equal(t, false, result6[0].Done)
	}

	result7, err := db.GetTasksWithParams(ctx, 999, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result7))

	result8, err := db.GetTasksWithParams(ctx, 999, false)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result8))

	for _, id := range ids {
		err := db.DeleteTask(ctx, id)
		assert.NoError(t, err)
	}
}
