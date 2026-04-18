package http

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	resp "Board_of_issuses/internal/core/transport/response"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *TasksHandler)DeleteTask(w http.ResponseWriter, r *http.Request) {
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
