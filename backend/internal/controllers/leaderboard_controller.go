package controllers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
)

type LeaderboardController struct {
	leaderboardService services.LeaderboardService
}

func NewLeaderboardController(leaderboardService services.LeaderboardService) *LeaderboardController {
	return &LeaderboardController{leaderboardService: leaderboardService}
}

// GetLeaderboard handles GET /leaderboard requests.
// Returns a paginated ranking of all participants.
//
// Query Parameters:
//   - limit: max results per page (default: 20, max: 100)
//   - offset: number of results to skip (default: 0)
//
// Responses:
//   - 200: { entries: LeaderboardEntryDTO[], total: int }
//   - 500: Internal server error
//
// Authentication: None
func (c *LeaderboardController) GetLeaderboard(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	limit := middlewares.GetQueryStringInt(request, "limit", 20)
	offset := middlewares.GetQueryStringInt(request, "offset", 0)

	entries, total, err := c.leaderboardService.GetLeaderboard(limit, offset)
	if err != nil {
		return middlewares.HandleServiceError(err), nil
	}

	dtos := make([]interface{}, len(entries))
	for i, e := range entries {
		dtos[i] = e.ToDTO()
	}

	return middlewares.SuccessResponse(map[string]interface{}{
		"entries": dtos,
		"total":   total,
	}, http.StatusOK), nil
}

// GetUser handles GET /leaderboard/{userId} requests.
// Returns the leaderboard position for a specific user.
//
// Path Parameters:
//   - userId: the user's ID
//
// Responses:
//   - 200: LeaderboardEntryDTO
//   - 400: Invalid user ID
//   - 404: User not found in leaderboard
//
// Authentication: None
func (c *LeaderboardController) GetUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := middlewares.GetId(request)
	if userId == "" {
		return middlewares.BadRequestResponse([]string{"userId is required"}), nil
	}

	entry, err := c.leaderboardService.GetUserPosition(userId)
	if err != nil {
		return middlewares.HandleServiceError(err), nil
	}

	return middlewares.SuccessResponse(entry.ToDTO(), http.StatusOK), nil
}
