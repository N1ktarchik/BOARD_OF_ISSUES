package http

//перенести в фичю досок потом удалить
//доработать функцию!

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	"N1ktarchik/Board_of_issues/internal/core/errors"
	req "N1ktarchik/Board_of_issues/internal/core/transport/request"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"log/slog"
	"net/http"
)

// ConnectUserToDesk        godoc
// @Summary                 Join a desk
// @Description             Connect current user to a desk using desk password
// @Tags                    desks
// @Security                ApiKeyAuth
// @Accept                  json
// @Produce                 json
// @Param                   request body DeskRequestDTO true "Desk ID and Password"
// @Success                 201 {object} map[string]string "message: you have connected to desk"
// @Failure                 400 {object} resp.ErrorResponse "Possible: invalid_data, bad_request"
// @Failure                 401 {object} resp.ErrorResponse "unauthorized"
// @Failure                 404 {object} resp.ErrorResponse "desk_not_found"
// @Failure                 409 {object} resp.ErrorResponse "already_a_member"
// @Failure                 500 {object} resp.ErrorResponse "internal_server_error"
// @Router                  /desks/connect [post]
func (h *DesksHandler) ConnectUserToDesk(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/desks/connect"))

	ctx := r.Context()
	userId, ok := domain.GetUserID(ctx)

	if !ok {
		h.log.Error("connect user to desk failed: userID not faund in context")
		resp.RespondWithError(w, errors.BadRequest())
		return
	}

	deskDTO := &DeskRequestDTO{}
	if err := req.DecodeAndValidateRequest(r, deskDTO); err != nil {
		h.log.Error("decode and validate failed", slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	if err := h.desksService.ConnectUserToDesk(ctx, deskDTO.OwnerId, deskDTO.Id); err != nil {
		h.log.Error("service connect user to desk failed",
			slog.Any("userID", userId), slog.Any("deskID", deskDTO.Id), slog.Any("err", err))

		resp.RespondWithError(w, err)
		return
	}

	resp.RespondWithJSON(w, http.StatusCreated, "you have connected to desk")

}
