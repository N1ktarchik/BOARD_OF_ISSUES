package http

import (
	"Board_of_issuses/internal/core/domain"
	core_errors "Board_of_issuses/internal/core/errors"
	resp "Board_of_issuses/internal/core/transport/response"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *DesksHandler) DeleteDesk(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/desks/{id}"),
		slog.String("method", http.MethodDelete))

	ctx := r.Context()
	userIdStr, ok := domain.GetUserID(ctx)
	if !ok {
		h.log.Error("delete desk failed: userID not faund in context")

		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	deskIdStr, ok := mux.Vars(r)["id"]
	if !ok {
		h.log.Warn("delete desk failed: error get desk id in path", slog.Any("userID", userIdStr))

		resp.RespondWithError(w, core_errors.BadRequest())
		return
	}

	if err := h.desksService.DeleteDesk(r.Context(), deskIdStr, userIdStr); err != nil {
		h.log.Error("service delete desk failed",
			slog.Any("userID", userIdStr), slog.Any("deskID", deskIdStr), slog.Any("err", err))

		resp.RespondWithError(w, err)
		return
	}

	h.log.Info("desk deleted successfully", slog.Any("deskID", deskIdStr), slog.Any("userID", userIdStr))

	resp.RespondWithJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("desk with ID %s has been deleted", deskIdStr)})

}
