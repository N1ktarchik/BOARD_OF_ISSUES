package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"N1ktarchik/Board_of_issues/internal/core/transport/request"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// ChangeTaskData           godoc
// @Summary                 Update task
// @Description             Update task name, description or deadline
// @Tags                    tasks
// @Security                ApiKeyAuth
// @Accept                  json
// @Produce                 json
// @Param                   request body UpdateTaskRequestDTO true "New task data"
// @Success                 200 {object} domain.Task "Updated task"
// @Failure                 400 {object} resp.ErrorResponse "Possible: invalid_deadline_format, task_name_too_short"
// @Failure                 401 {object} resp.ErrorResponse "unauthorized"
// @Failure                 403 {object} resp.ErrorResponse "not_a_desk_member"
// @Failure                 404 {object} resp.ErrorResponse "task_not_found"
// @Failure                 500 {object} resp.ErrorResponse "internal_server_error"
// @Router                  /tasks/update [patch]
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
