package http

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	resp "Board_of_issuses/internal/core/transport/response"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *TasksHandler) GetTasksFromOneDesk(w http.ResponseWriter, r *http.Request) {

	h.log.Info("new request", slog.String("path", "/tasks/all/{deskId}"))

	tasksFilter := &domain.TaskFilter{}

	userIdStr, ok := domain.GetUserID(r.Context())
	if !ok {
		h.log.Error("get all tasks failed: userID not found in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		h.log.Warn("get all tasks failed: error parsing user id", slog.Any("err", err))
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}
	tasksFilter.UserId = userUUID

	deskIdStr, ok := mux.Vars(r)["deskId"]
	if !ok {
		h.log.Warn("get all tasks failed: desk id not found in URL")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	deskUUID, err := uuid.Parse(deskIdStr)
	if err != nil {
		h.log.Warn("get all tasks failed: error parsing desk id", slog.Any("err", err))
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}
	tasksFilter.DeskId = deskUUID

	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 || pageInt > 100 {
		tasksFilter.Offset = 1
	} else {
		tasksFilter.Offset = pageInt
	}

	limit := r.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 || limitInt > 50 {
		tasksFilter.Limit = 20
	} else {
		tasksFilter.Limit = limitInt
	}

	done := r.URL.Query().Get("done")
	switch done {
	case "true":
		doneBool := true
		tasksFilter.Done = &doneBool
	case "false":
		doneBool := false
		tasksFilter.Done = &doneBool
	default:
		tasksFilter.Done = nil
	}

	tasks, err := h.tasksService.GetTasksFromOneDesk(r.Context(), tasksFilter)
	if err != nil {
		h.log.Error("get all tasks failed: service error",
			slog.Any("err", err), slog.Any("desk_id", deskIdStr),
			slog.Any("user_id", userIdStr))

		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("tasks retrieved successfully", slog.Any("desk_id", deskIdStr), slog.Any("user_id", userIdStr))

	resp.RespondWithArray(w, http.StatusOK, "task", tasks)
}

func (h *TasksHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {

	h.log.Info("new request", slog.String("path", "/tasks/{taskId}"))

	taskIdStr, ok := mux.Vars(r)["taskId"]
	if !ok {
		h.log.Warn("get task by id failed: task id not found in URL")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	userIdStr, ok := domain.GetUserID(r.Context())
	if !ok {
		h.log.Error("get task by id failed: userID not found in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	task, err := h.tasksService.GetTaskByID(r.Context(), taskIdStr, userIdStr)
	if err != nil {
		h.log.Error("get task by id failed: service error",
			slog.Any("err", err), slog.Any("task_id", taskIdStr),
			slog.Any("user_id", userIdStr))

		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("task retrieved successfully", slog.Any("task_id", taskIdStr), slog.Any("user_id", userIdStr))

	resp.RespondWithJSON(w, http.StatusOK, task)
}
