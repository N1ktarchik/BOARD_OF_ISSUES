package http

//перенести в фичю досок потом удалить
//доработать функцию!



// import (
// 	er "Board_of_issuses/internal/core"
// 	"Board_of_issuses/internal/core/domain"
// 	"Board_of_issuses/internal/core/errors"
// 	req "Board_of_issuses/internal/core/transport/request"
// 	resp "Board_of_issuses/internal/core/transport/response"
// 	"log/slog"
// 	"net/http"
// )

// func (h *UsersHandler) HandleConnectUserToDesk(w http.ResponseWriter, r *http.Request) {
// 	h.log.Info("new request", slog.String("path", "/users/connect"))

// 	ctx := r.Context()
// 	userId, ok := domain.GetUserID(ctx)

// 	if !ok {
// 		h.log.Error("connect user to desk failed: userID not faund in context")
// 		resp.RespondWithError(w, errors.BadRequest())
// 		return
// 	}

// 	reqData := &ConnectUserToDeskRequestDTO{}
// 	if err := req.DecodeAndValidateRequest(r, reqData); err != nil {
// 		h.log.Error("decode and validate failed", slog.Any("err", err))
// 		resp.RespondWithError(w, err)
// 		return
// 	}

// 	//!!!!!!!To server

// 	// if desk.ID <= 0 {
// 	// 	RespondWithError(w, http.StatusBadRequest, "desk_id can not be less than or equal to zero")
// 	// 	return
// 	// }

	
// 	if err := h.serv.ConnectUserToDesk(r.Context(), userId, desk.ID, desk.Password); err != nil {
// 		switch {

// 		case er.IsError(err, "PASSWORD_IS_SHORT"), er.IsError(err, "PASSWORD_IS_lONG"):
// 			var appErr *er.ErrorApp = err.(*er.ErrorApp)
// 			RespondWithError(w, http.StatusBadRequest, appErr.Message)

// 		default:
// 			RespondWithError(w, http.StatusInternalServerError, "error connect to desk ")
// 		}

// 		return
// 	}

// 	RespondWithJSON(w, http.StatusCreated, "you have connected to desk")

// }
