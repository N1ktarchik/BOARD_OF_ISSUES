package handlers

import (
	er "Board_of_issuses/internal/core"
	dto "Board_of_issuses/internal/features/transport/http/dto"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (h *UserHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "read request body error")
		return
	}

	task := &dto.Task{}

	if err := json.Unmarshal(reqBody, task); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
		return
	}

	if task.DeskId <= 0 {
		RespondWithError(w, http.StatusBadRequest, "task id can not be less than or equal to zero ")
	}

	if len(task.Name) < 3 {
		RespondWithError(w, http.StatusBadRequest, "length name can not be less than 3 ")
		return
	}

	task.Created_at = time.Now()
	task.UserId = userID
	task.Done = false
	if task.Time.IsZero() {
		task.Time = time.Now().Add(30 * 24 * time.Hour)
	}

	if err := h.serv.CreateTask(r.Context(), task.ToServicenTask()); err != nil {
		switch {

		case er.IsError(err, "USER_HAVE_NOT_ACCES"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "add task error")
		}

		return
	}

	RespondWithJSON(w, http.StatusCreated, "task have been created")

}

func (h *UserHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	taskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser task id")
		return
	}

	if taskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	if err := h.serv.DeleteTask(r.Context(), taskID, userId); err != nil {
		switch {

		case er.IsError(err, "USER_IS_NOT_OWNER_OF_TASK"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to delete desk")
		}

		return
	}

	RespondWithJSON(w, http.StatusNoContent, "task had deleted")

}

func (h *UserHandler) HandleComplyteTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	taskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser task id")
		return
	}

	if taskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	if err := h.serv.ComplyteTask(r.Context(), userId, taskID); err != nil {
		switch {

		case er.IsError(err, "USER_HAVE_NOT_ACCES"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to complyte task")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, "you have complyted the task")

}

func (h *UserHandler) HandleAddTimeToTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	taskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser task id")
		return
	}

	if taskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "read request body error")
		return
	}

	taskTime := &dto.UpdateTaskTimeRequest{}
	if err := json.Unmarshal(reqBody, taskTime); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
		return
	}

	if err := h.serv.UpdateTaskTime(r.Context(), userId, taskID, taskTime.Hours); err != nil {
		switch {

		case er.IsError(err, "USER_IS_NOT_OWNER_OF_TASK"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to delete desk")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, "task time have updated")

}

func (h *UserHandler) HandleChangeTaskDescription(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	taskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser task id")
		return
	}

	if taskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "read request body error")
		return
	}

	newDescription := &dto.UpdateTaskDescriptionRequest{}

	if err := json.Unmarshal(reqBody, newDescription); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
		return
	}

	if err := h.serv.ChangeTaskDescription(r.Context(), userId, taskID, newDescription.Description); err != nil {
		switch {

		case er.IsError(err, "USER_IS_NOT_OWNER"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to change desk name")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, "task description had change")

}

func (h *UserHandler) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	deskID, err := strconv.Atoi(mux.Vars(r)["deskId"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser desk id")
		return
	}

	if deskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	tasks, err := h.serv.GetAllTasks(r.Context(), userId, deskID)
	if err != nil {
		switch {

		case er.IsError(err, "USER_HAVE_NOT_ACCES"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to get tasks")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *UserHandler) HandleGetTasksWithParams(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)

	deskID, err := strconv.Atoi(r.URL.Query().Get("desk_id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser desk id")
		return
	}

	done, err := strconv.ParseBool(r.URL.Query().Get("done"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser complyte param")
		return
	}

	if deskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	tasks, err := h.serv.GetTasksWithParams(r.Context(), userId, deskID, done)
	if err != nil {
		switch {

		case er.IsError(err, "USER_HAVE_NOT_ACCES"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to get tasks")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, tasks)
}
