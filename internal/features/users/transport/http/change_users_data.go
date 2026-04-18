package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	req "N1ktarchik/Board_of_issues/internal/core/transport/request"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// ChangesUserData          godoc
// @Summary                 Update user profile
// @Description             Update current user's name or email or password
// @Tags                    users
// @Security                ApiKeyAuth
// @Accept                  json
// @Produce                 json
// @Param                   request body UsersRequestDTO true "New User Data"
// @Success                 200 {object} domain.User "Successfully updated user information"
// @Failure                 400 {object} resp.ErrorResponse "Possible: invalid_email, bad_request"
// @Failure                 401 {object} resp.ErrorResponse "unauthorized"
// @Failure                 404 {object} resp.ErrorResponse "user_not_found"
// @Failure                 409 {object} resp.ErrorResponse "email_already_exists"
// @Failure                 500 {object} resp.ErrorResponse "internal_server_error"
// @Router                  /users/update [patch]
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
