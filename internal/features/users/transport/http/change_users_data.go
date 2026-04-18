package http

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	req "Board_of_issuses/internal/core/transport/request"
	resp "Board_of_issuses/internal/core/transport/response"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *UsersHandler) ChangesUserData(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/users/update"))
	ctx := r.Context()
	userIdStr, ok := domain.GetUserID(ctx)

	if !ok {
		h.log.Error("connect user to desk failed: userID not faund in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	reqData := &UsersRequestDTO{}
	if err := req.DecodeAndValidateRequest(r, reqData); err != nil {
		h.log.Error("decode and validate failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		h.log.Warn("parse userID failed", slog.Any("err", err))
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	domainUser := reqData.ToServiceUser()
	domainUser.ID = userUUID

	saveUser, err := h.usersService.ChangeUsersData(ctx, domainUser)
	if err != nil {
		h.log.Error("service change user data failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	resp.RespondWithJSON(w, http.StatusOK, saveUser)

}
