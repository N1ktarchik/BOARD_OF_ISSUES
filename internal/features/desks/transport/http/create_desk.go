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

// CreateDesk           godoc
// @Summary             Create a desk
// @Description         Create a new board for tasks
// @Tags                desks
// @Security            ApiKeyAuth
// @Accept              json
// @Produce             json
// @Param               request body DeskRequestDTO true "Desk Info"
// @Success             201 {object} domain.Desk "Successfully created desk"
// @Failure             400 {object} resp.ErrorResponse "Possible: desk_name_too_short, invalid_data"
// @Failure             401 {object} resp.ErrorResponse "unauthorized"
// @Failure             500 {object} resp.ErrorResponse "internal_server_error"
// @Router              /desks/create [post]
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
