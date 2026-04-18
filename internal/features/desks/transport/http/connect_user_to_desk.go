package http

//перенести в фичю досок потом удалить
//доработать функцию!

import (
	"Board_of_issuses/internal/core/domain"
	"Board_of_issuses/internal/core/errors"
	req "Board_of_issuses/internal/core/transport/request"
	resp "Board_of_issuses/internal/core/transport/response"
	"log/slog"
	"net/http"
)

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
