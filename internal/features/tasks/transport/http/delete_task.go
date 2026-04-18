package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// DeleteTask           godoc
// @Summary             Delete task
// @Description         Remove a task by its ID
// @Tags                tasks
// @Security            ApiKeyAuth
// @Accept              json
// @Produce             json
// @Param               id path string true "TASK ID" format(uuid) example("550e8400-e29b-41d4-a716-446655440000")
// @Success             200 {object} map[string]string "message: task with ID ... has been deleted"
// @Failure             400 {object} resp.ErrorResponse "Possible: invalid_uuid, bad_request"
// @Failure             401 {object} resp.ErrorResponse "unauthorized"
// @Failure             403 {object} resp.ErrorResponse "not_an_owner"
// @Failure             404 {object} resp.ErrorResponse "task_not_found"
// @Failure             500 {object} resp.ErrorResponse "internal_server_error"
// @Router              /tasks/{id} [delete]
func (h *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/tasks/{id}"),
		slog.String("method", http.MethodDelete))

	ctx := r.Context()
	userId, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("delete task failed: userID not found in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	taskID, ok := mux.Vars(r)["id"]
	if !ok {
		h.log.Warn("delete task failed: task id not found in URL")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID, userId); err != nil {
		h.log.Error("delete task failed: service error",
			slog.Any("err", err), slog.Any("task_id", taskID), slog.Any("user_id", userId))

		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("task deleted successfully", slog.Any("task_id", taskID))

	resp.RespondWithJSON(w, http.StatusOK,
		map[string]string{"message": fmt.Sprintf("task with ID %s has been deleted", taskID)})

}
