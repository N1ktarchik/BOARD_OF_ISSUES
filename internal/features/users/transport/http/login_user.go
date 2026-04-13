package http

import (
	req "Board_of_issuses/internal/core/transport/request"
	resp "Board_of_issuses/internal/core/transport/response"
	"log/slog"
	"net/http"
)

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
