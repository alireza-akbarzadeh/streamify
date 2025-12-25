package users

import (
	"net/http"
	"strconv"

	"github.com/techies/streamify/internal/service"
	"github.com/techies/streamify/internal/utils"
)

// UserList returns a paginated list of users with optional search.
// @Summary      List users
// @Description  Get a paginated list of users. Supports search by username, email, and phone number.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        limit         query     int     false  "Max results per page (default 20, max 100)"
// @Param        offset        query     int     false  "Offset for pagination (default 0)"
// @Param        username      query     string  false  "Search by username (partial match)"
// @Param        email         query     string  false  "Search by email (partial match)"
// @Param        phone_number  query     string  false  "Search by phone number (partial match)"
// @Success      200       {object}  UserListResponse
// @Failure      401       {object}  utils.ErrorResponse
// @Failure      500       {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users [get]
func (h *UserHandler) UserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := h.parseInt(r.URL.Query().Get("limit"), 20)
	if limit > 100 {
		limit = 100
	}
	offset := h.parseInt(r.URL.Query().Get("offset"), 0)

	result, appErr := h.Service.ListUsers(ctx, service.ListUsersParams{
		Limit:       int32(limit),
		Offset:      int32(offset),
		Username:    r.URL.Query().Get("username"),
		Email:       r.URL.Query().Get("email"),
		PhoneNumber: r.URL.Query().Get("phone_number"),
	})

	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, UserListResponse{
		Users:   MapUserListToResponse(result.Users),
		Total:   result.Total,
		Limit:   int32(limit),
		Offset:  int32(offset),
		HasMore: int64(offset+limit) < result.Total,
	})
}

// Helper to keep the main handler logic clean
func (h *UserHandler) parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	res, err := strconv.Atoi(value)
	if err != nil || res < 0 {
		return fallback
	}
	return res
}
