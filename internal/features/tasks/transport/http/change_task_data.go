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
	h.log.Info("new request", slog.String("path", "/tasks/change"))

	ctx := r.Context()
	userIdStr, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("change task data failed: userID not found in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		h.log.Error("change task data failed: error parsing user id", slog.Any("err", err))
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	//To service
	// if taskID <= 0 {
	// 	RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
	// 	return
	// }

	task := &TaskRequestDTO{}
	if err := request.DecodeAndValidateRequest(r, task); err != nil {
		h.log.Error("decode and validate failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}
	task.AuthorId = userUUID

	saveTask, err := h.tasksService.UpdateTask(ctx, task.ToServiceTask())
	if err != nil {
		h.log.Error("update task failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("task updated successfully", slog.Any("task", saveTask))

	resp.RespondWithJSON(w, http.StatusOK, saveTask)

}

// func (h *UserHandler) HandleChangeTaskDescription(w http.ResponseWriter, r *http.Request) {
// 	userId := getUserIDFromContext(r)
// 	taskID, err := strconv.Atoi(mux.Vars(r)["id"])
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, "error to parser task id")
// 		return
// 	}

// 	if taskID <= 0 {
// 		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
// 		return
// 	}

// 	reqBody, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, "read request body error")
// 		return
// 	}

// 	newDescription := &dto.UpdateTaskDescriptionRequest{}

// 	if err := json.Unmarshal(reqBody, newDescription); err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
// 		return
// 	}

// 	if err := h.serv.ChangeTaskDescription(r.Context(), userId, taskID, newDescription.Description); err != nil {
// 		switch {

// 		case er.IsError(err, "USER_IS_NOT_OWNER"):
// 			var appErr *er.ErrorApp = err.(*er.ErrorApp)
// 			RespondWithError(w, http.StatusForbidden, appErr.Message)

// 		default:
// 			RespondWithError(w, http.StatusInternalServerError, "error to change desk name")
// 		}

// 		return
// 	}

// 	RespondWithJSON(w, http.StatusOK, "task description had change")

// }
