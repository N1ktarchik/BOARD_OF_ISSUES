package http

import (
	req "N1ktarchik/Board_of_issues/internal/core/transport/request"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"log/slog"
	"net/http"
)

// Login            godoc
// @Summary         Login user
// @Description     Authenticate user and return JWT token
// @Tags            users
// @Accept          json
// @Produce         json
// @Param           request body UsersRequestDTO true "Login credentials"
// @Success         201 {object} resp.JWTResponse "Successfully logged in"
// @Failure         400 {object} resp.ErrorResponse "Possible: invalid_credentials, bad_request"
// @Failure         401 {object} resp.ErrorResponse "invalid_password"
// @Failure         404 {object} resp.ErrorResponse "user_not_found"
// @Failure         500 {object} resp.ErrorResponse "internal_server_error"
// @Router          /login [post]
func (h *UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/login"))

	reqData := &UsersRequestDTO{}
	ctx := r.Context()

	if err := req.DecodeAndValidateRequest(r, reqData); err != nil {
		h.log.Error("decode and validate failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	JWTtoken, err := h.usersService.LoginUser(ctx, reqData.ToServiceUser())
	if err != nil {
		h.log.Error("service login failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	resp.RespondWithJWT(w, http.StatusCreated, JWTtoken)

}
