package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

func (r *TasksRepository) UpdateTask(ctx context.Context, updateData *domain.UpdateTask) (*domain.Task, error) {
	r.log.Info("updating task data in repository",
		slog.Any("taskID", updateData.Id),
		slog.Any("authorID", updateData.AuthorId))

	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if updateData.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argID))
		args = append(args, updateData.Name)
		argID++
	}

	if updateData.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argID))
		args = append(args, updateData.Description)
		argID++
	}

	if !updateData.Deadline.IsZero() {
		setValues = append(setValues, fmt.Sprintf("deadline = $%d", argID))
		args = append(args, updateData.Deadline)
		argID++
	}

	if len(setValues) == 0 {
		r.log.Warn("update skipped: no fields provided", slog.Any("taskID", updateData.Id))
		return nil, core_errors.BadRequest()
	}

	query := fmt.Sprintf(`
		UPDATE tasks 
		SET %s 
		WHERE id = $%d 
		AND desk_id IN (
			SELECT desk_id FROM desk_members WHERE user_id = $%d
		) 
		RETURNING id, author_id, desk_id, name, description, done, deadline, created_at;`,
		strings.Join(setValues, ", "),
		argID,
		argID+1,
	)

	args = append(args, updateData.Id, updateData.AuthorId)

	var t taskModel
	row := r.pool.QueryRow(ctx, query, args...)

	if err := t.scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn("update task failed: access denied or task not found",
				slog.Any("taskID", updateData.Id),
				slog.Any("userID", updateData.AuthorId))

			return nil, core_errors.UserNotOwnerOfTask(updateData.AuthorId.String(), updateData.Id.String())
		}

		r.log.Error("repository error during task update",
			slog.Any("taskID", updateData.Id),
			slog.Any("err", err))

		return nil, core_errors.ServerError()
	}

	r.log.Info("task updated successfully in database", slog.Any("taskID", t.Id))

	domainTask := modelToDomain(t)
	return &domainTask, nil
}
