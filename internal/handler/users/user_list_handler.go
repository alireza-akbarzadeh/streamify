package users

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/utils"
)

// UserList returns a paginated list of users with optional search.
// @Summary      List users
// @Description  Get a paginated list of users. Supports search by username and email.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        limit     query     int     false  "Max results per page (default 20, max 100)"
// @Param        offset    query     int     false  "Offset for pagination (default 0)"
// @Param        username  query     string  false  "Search by username (partial match)"
// @Param        email     query     string  false  "Search by email (partial match)"
// @Success      200       {object}  users.UserListResponse
// @Failure      401       {object}  utils.ErrorResponse
// @Failure      500       {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users [get]
func (h *UserHandler) UserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Get Params using your helper
	limit := h.parseInt(r.URL.Query().Get("limit"), 20)
	if limit > 100 {
		limit = 100
	}

	offset := h.parseInt(r.URL.Query().Get("offset"), 0)

	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")

	// 2. Prepare Database Arguments
	arg := database.GetUsersParams{
		Limit:          int32(limit),
		Offset:         int32(offset),
		SearchUsername: sql.NullString{String: username, Valid: username != ""},
		SearchEmail:    sql.NullString{String: email, Valid: email != ""},
	}

	// 3. Fetch Data
	dbUsers, err := h.app.DB.GetUsers(ctx, arg)
	if err != nil {
		logger.Error(ctx, "UserList: failed to fetch users", err, "args", arg)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}

	total, err := h.app.DB.CountUsers(ctx, database.CountUsersParams{
		SearchUsername: arg.SearchUsername,
		SearchEmail:    arg.SearchEmail,
	})
	if err != nil {
		logger.Error(ctx, "UserList: failed to count users", err, "args", arg)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to count users", err)
		return
	}

	// 4. Map and Respond
	// MapUserListToResponse returns []UserResponse, which now matches our struct!
	utils.RespondWithJSON(w, http.StatusOK, UserListResponse{
		Users:   MapUserListToResponse(dbUsers),
		Total:   total,
		Limit:   int32(limit),
		Offset:  int32(offset),
		HasMore: int64(offset+limit) < total,
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
