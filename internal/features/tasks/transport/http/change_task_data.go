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

func (h *TasksHandler) ChangeTaskData(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/tasks/update"))

	ctx := r.Context()
	userIdStr, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("change task data failed: userID not found in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		h.log.Warn("change task data failed: error parsing user id", slog.Any("err", err))
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	task := &UpdateTaskRequestDTO{}
	if err := request.DecodeAndValidateRequest(r, task); err != nil {
		h.log.Error("decode and validate failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}
	task.AuthorId = userUUID

	saveTask, err := h.tasksService.UpdateTask(ctx, task.ToServiceUpdateTask())
	if err != nil {
		h.log.Error("update task failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("task updated successfully", slog.Any("task", saveTask))

	resp.RespondWithJSON(w, http.StatusOK, saveTask)

}
