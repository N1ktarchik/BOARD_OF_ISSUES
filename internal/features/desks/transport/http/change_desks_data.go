package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"N1ktarchik/Board_of_issues/internal/core/transport/request"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// ChangeDeskData           godoc
// @Summary                 Update desk info
// @Description             Update desk name or password (must be owner)
// @Tags                    desks
// @Security                ApiKeyAuth
// @Accept                  json
// @Produce                 json
// @Param                   request body DeskRequestDTO true "New desk data"
// @Success                 200 {object} domain.Desk "Successfully updated desk information"
// @Failure                 400 {object} resp.ErrorResponse "Possible: invalid_uuid, bad_request"
// @Failure                 401 {object} resp.ErrorResponse "unauthorized"
// @Failure                 403 {object} resp.ErrorResponse "not_an_owner"
// @Failure                 404 {object} resp.ErrorResponse "desk_not_found"
// @Failure                 500 {object} resp.ErrorResponse "internal_server_error"
// @Router                  /desks/update [patch]
func (h *DesksHandler) ChangeDeskData(w http.ResponseWriter, r *http.Request) {
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
