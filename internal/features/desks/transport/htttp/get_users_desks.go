package http

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	resp "Board_of_issuses/internal/core/transport/response"
	"log/slog"
	"net/http"
)

func (h *DesksHandler) HandleGetAllDesksId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdStr, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("get users desks failed: userID not faund in context")
		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	desks, err := h.desksService.GetAllUsersDesks(ctx, userIdStr)
	if err != nil {
		h.log.Error("service get users desks failed", slog.Any("userID", userIdStr), slog.Any("err", err))
		resp.RespondWithError(w, err)
		return
	}

	resp.RespondWithArray(w, http.StatusOK, "desks", desks)
}
