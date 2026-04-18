package http

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	resp "N1ktarchik/Board_of_issues/internal/core/transport/response"
	"log/slog"
	"net/http"
)

// GetUsersDesks            godoc
// @Summary                 Get my desks
// @Description             Get all desks where current user is a member or owner
// @Tags                    desks
// @Security                ApiKeyAuth
// @Accept                  json
// @Produce                 json
// @Success                 200 {array} domain.Desk "Successfully retrieved list of desks"
// @Failure                 400 {object} resp.ErrorResponse "Possible: invalid_user_id, bad_request"
// @Failure                 401 {object} resp.ErrorResponse "unauthorized"
// @Failure                 500 {object} resp.ErrorResponse "internal_server_error"
// @Router                  /desks/my [get]
func (h *DesksHandler) GetUsersDesks(w http.ResponseWriter, r *http.Request) {
	h.log.Info("new request", slog.String("path", "/desks/my"))

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
