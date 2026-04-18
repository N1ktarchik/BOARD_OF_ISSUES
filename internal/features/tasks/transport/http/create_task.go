package http

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	"Board_of_issuses/internal/core/transport/request"
	resp "Board_of_issuses/internal/core/transport/response"
	"log/slog"

	"net/http"

	"github.com/google/uuid"
)

func (h *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	h.log.Info("new request", slog.String("path", "/tasks/create"))

	ctx := r.Context()
	authorIdStr, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("create task failed: authorID not found in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	task := &TaskRequestDTO{}
	if err := request.DecodeAndValidateRequest(r, task); err != nil {
		h.log.Error("decode and validate failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	authorUUID, err := uuid.Parse(authorIdStr)
	if err != nil {
		h.log.Warn("create task failed: error parsing author id", slog.Any("err", err))
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	task.Done = false
	task.AuthorId = authorUUID

	saveTask, err := h.tasksService.CreateTask(ctx, task.ToServiceTask())
	if err != nil {
		h.log.Error("create task failed: service error",
			slog.Any("err", err), slog.Any("author_id", authorUUID))

		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("task created successfully", slog.Any("task_id", saveTask.Id))

	resp.RespondWithJSON(w, http.StatusCreated, saveTask)
}
