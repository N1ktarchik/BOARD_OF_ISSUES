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

func (h *DesksHandler) ChangeDeskDate(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/desks/update"))

	ctx := r.Context()
	userIdStr, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("change desk data failed: userID not faund in context")

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
		h.log.Warn("change desk data failed: error parsing user id",
			slog.Any("deskID", deskDTO.Id), slog.Any("err", err))

		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	saveDesk, err := h.desksService.ChangeDesksData(ctx, deskDTO.ToServiceDesk(), userUUID)
	if err != nil {
		h.log.Error("service change desk data failed", slog.Any("ownerID", userUUID),
			slog.Any("deskID", deskDTO.Id), slog.Any("err", err))

		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("desk data changed successfully", slog.Any("deskID", saveDesk.Id), slog.Any("ownerID", saveDesk.OwnerId))
	resp.RespondWithJSON(w, http.StatusOK, saveDesk)

}
