package handlers

import (
	er "Board_of_issuses/internal/core"
	dto "Board_of_issuses/internal/features/transport/http/dto"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *UserHandler) HandleCreateDesk(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "read request body error")
		return
	}

	desk := &dto.Desk{}

	if err := json.Unmarshal(reqBody, desk); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
		return
	}

	if len(desk.Name) < 3 {
		RespondWithError(w, http.StatusBadRequest, "length name can not be less than 3 ")
		return
	}

	desk.OwnerId = userID

	if err = h.serv.CreateDesk(r.Context(), desk.ToServiceDeskr()); err != nil {
		switch {

		case er.IsError(err, "PASSWORD_IS_SHORT"), er.IsError(err, "PASSWORD_IS_lONG"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusBadRequest, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to create desk")
		}

		return
	}

	RespondWithJSON(w, http.StatusCreated, "desk have been created")

}

func (h *UserHandler) HandleChangeDeskName(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	deskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser desk id")
		return
	}

	if deskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "read request body error")
		return
	}

	newName := &dto.UpdateDeskNameRequest{}

	if err := json.Unmarshal(reqBody, newName); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
		return
	}

	if len(newName.Name) < 3 {
		RespondWithError(w, http.StatusBadRequest, "length name can not be less than 3 ")
		return
	}

	if err := h.serv.ChangeDeskName(r.Context(), newName.Name, deskID, userId); err != nil {
		switch {

		case er.IsError(err, "USER_IS_NOT_OWNER"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to change desk name")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, "desk name had change")

}

func (h *UserHandler) HandleChangeDeskPassword(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	deskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser desk id")
		return
	}

	if deskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "read request body error")
		return
	}

	newPassword := &dto.UpdateDeskPasswordRequest{}

	if err := json.Unmarshal(reqBody, newPassword); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
		return
	}

	if err := h.serv.ChangeDeskPassword(r.Context(), newPassword.Password, deskID, userId); err != nil {
		switch {

		case er.IsError(err, "USER_IS_NOT_OWNER"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		case er.IsError(err, "PASSWORD_IS_SHORT"), er.IsError(err, "PASSWORD_IS_lONG"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusBadRequest, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to change desk name")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, "desk password had change")

}

func (h *UserHandler) HandleChangeDeskOwner(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	deskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parser desk id")
		return
	}

	if deskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "read request body error")
		return
	}

	newOwner := &dto.UpdateDeskOwnerRequest{}

	if err := json.Unmarshal(reqBody, newOwner); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "parse request data error")
		return
	}

	if newOwner.ID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "owner id can not be less than or equal to zero")
		return
	}

	if err := h.serv.ChangeDeskOwner(r.Context(), deskID, userId, newOwner.ID); err != nil {
		switch {

		case er.IsError(err, "USER_IS_NOT_OWNER"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)
		default:
			RespondWithError(w, http.StatusInternalServerError, "error to change desk owner")
		}

		return
	}

	RespondWithJSON(w, http.StatusOK, "desk owner had change")

}

func (h *UserHandler) HandleDeleteDesk(w http.ResponseWriter, r *http.Request) {
	userId := getUserIDFromContext(r)
	deskIdStr := mux.Vars(r)["id"]
	deskID, err := strconv.Atoi(deskIdStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error to parse desk id")
		return
	}

	if deskID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
		return
	}

	if err := h.serv.DeleteDesk(r.Context(), deskID, userId); err != nil {
		switch {

		case er.IsError(err, "USER_IS_NOT_OWNER"):
			var appErr *er.ErrorApp = err.(*er.ErrorApp)
			RespondWithError(w, http.StatusForbidden, appErr.Message)

		default:
			RespondWithError(w, http.StatusInternalServerError, "error to delete desk")
		}

		return
	}

	RespondWithJSON(w, http.StatusNoContent, "desk had deleted")

}

func (h *UserHandler) HandleGetAllDesksId(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	desksID, err := h.serv.GetAllDesks(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error to get desks id")
		return
	}

	RespondWithJSON(w, http.StatusOK, desksID)
}
