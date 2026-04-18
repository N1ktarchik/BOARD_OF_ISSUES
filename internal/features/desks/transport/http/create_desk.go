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

func (h *DesksHandler) CreateDesk(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/desks/create"))

	ctx := r.Context()
	userIdStr, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("create desk failed: userID not faund in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	deskDTO := &DeskRequestDTO{}
	if err := request.DecodeAndValidateRequest(r, deskDTO); err != nil {
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

	deskDTO.OwnerId = userUUID

	saveDesk, err := h.desksService.CreateDesk(ctx, deskDTO.ToServiceDesk())
	if err != nil {
		h.log.Error("service create desk failed", slog.Any("userID", userUUID), slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("desk created successfully", slog.Any("deskID", saveDesk.Id), slog.Any("ownerID", saveDesk.OwnerId))

	resp.RespondWithJSON(w, http.StatusCreated, saveDesk)

}
