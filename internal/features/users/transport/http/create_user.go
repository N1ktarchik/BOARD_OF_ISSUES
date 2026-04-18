package http

import (
	req "N1ktarchik/Board_of_issues/internal/core/transport/request"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"log/slog"
	"net/http"
)

// Register         godoc
// @Summary         Register a new user
// @Description     Create a new user account with login, password and email
// @Tags            users
// @Accept          json
// @Produce         json
// @Param           request body UsersRequestDTO true "Register request body"
// @Success         201 {object} resp.JWTResponse "Successfully registered user"
// @Failure         400 {object} resp.ErrorResponse "Possible: invalid_email, password_too_short"
// @Failure         409 {object} resp.ErrorResponse "Possible: user_already_exists, email_already_taken"
// @Failure         500 {object} resp.ErrorResponse "internal_server_error"
// @Router          /register [post]
func (h *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/register"))

	reqData := &UsersRequestDTO{}
	ctx := r.Context()

	if err := req.DecodeAndValidateRequest(r, reqData); err != nil {
		h.log.Error("decode and validate failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	JWTtoken, err := h.usersService.RegisterUser(ctx, reqData.ToServiceUser())
	if err != nil {
		h.log.Error("service register failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	resp.RespondWithJWT(w, http.StatusCreated, JWTtoken)

}
