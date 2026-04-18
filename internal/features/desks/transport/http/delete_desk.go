package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// DeleteDesk           godoc
// @Summary             Delete a desk
// @Description         Remove a desk by ID (must be owner)
// @Tags                desks
// @Security            ApiKeyAuth
// @Accept              json
// @Produce             json
// @Param               id path string true "DESK ID" format(uuid) example("550e8400-e29b-41d4-a716-446655440000")
// @Success             200 {object} map[string]string "message: desk with ID ... has been deleted"
// @Failure             400 {object} resp.ErrorResponse "Possible: invalid_uuid, bad_request"
// @Failure             401 {object} resp.ErrorResponse "unauthorized"
// @Failure             403 {object} resp.ErrorResponse "not_an_owner"
// @Failure             404 {object} resp.ErrorResponse "desk_not_found"
// @Failure             500 {object} resp.ErrorResponse "internal_server_error"
// @Router              /desks/{id} [delete]
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
