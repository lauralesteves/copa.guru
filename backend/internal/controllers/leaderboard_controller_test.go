package controllers

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	svcerr "github.com/lauralesteves/copa-guru-backend/internal/middlewares"
	"github.com/lauralesteves/copa-guru-backend/internal/models"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
	"go.uber.org/mock/gomock"
)

func setupLeaderboardController(t *testing.T) (*gomock.Controller, *services.MockLeaderboardService, *LeaderboardController) {
	ctrl := gomock.NewController(t)
	mockSvc := services.NewMockLeaderboardService(ctrl)
	controller := NewLeaderboardController(mockSvc)
	return ctrl, mockSvc, controller
}

// --- GetLeaderboard ---

func TestGetLeaderboard_Success(t *testing.T) {
	ctrl, mockSvc, controller := setupLeaderboardController(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetLeaderboard(int64(20), int64(0)).Return([]*models.LeaderboardEntry{
		{UserName: "Alice", TotalPoints: 30, Rank: 1},
		{UserName: "Bob", TotalPoints: 20, Rank: 2},
	}, int64(2), nil)

	req := events.APIGatewayProxyRequest{}
	resp, err := controller.GetLeaderboard(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	var body map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &body)
	entries := body["entries"].([]interface{})
	if len(entries) != 2 {
		t.Errorf("len(entries) = %d, want 2", len(entries))
	}
	if body["total"].(float64) != 2 {
		t.Errorf("total = %v, want 2", body["total"])
	}
}

func TestGetLeaderboard_WithPagination(t *testing.T) {
	ctrl, mockSvc, controller := setupLeaderboardController(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetLeaderboard(int64(5), int64(10)).Return([]*models.LeaderboardEntry{}, int64(0), nil)

	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"limit": "5", "offset": "10"},
	}
	resp, _ := controller.GetLeaderboard(req)
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
}

func TestGetLeaderboard_ServiceError(t *testing.T) {
	ctrl, mockSvc, controller := setupLeaderboardController(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetLeaderboard(int64(20), int64(0)).Return(nil, int64(0), svcerr.NewInternalError("failed", nil))

	req := events.APIGatewayProxyRequest{}
	resp, _ := controller.GetLeaderboard(req)
	if resp.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", resp.StatusCode)
	}
}

// --- GetUser ---

func TestGetUser_Success(t *testing.T) {
	ctrl, mockSvc, controller := setupLeaderboardController(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetUserPosition("abc123").Return(&models.LeaderboardEntry{
		UserName: "Alice", TotalPoints: 30, Rank: 1,
	}, nil)

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"userId": "abc123"},
	}
	resp, err := controller.GetUser(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
}

func TestGetUser_MissingUserId(t *testing.T) {
	ctrl, _, controller := setupLeaderboardController(t)
	defer ctrl.Finish()

	req := events.APIGatewayProxyRequest{}
	resp, _ := controller.GetUser(req)
	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	ctrl, mockSvc, controller := setupLeaderboardController(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().GetUserPosition("nonexistent").Return(nil, svcerr.NewNotFoundError("not found"))

	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"userId": "nonexistent"},
	}
	resp, _ := controller.GetUser(req)
	if resp.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", resp.StatusCode)
	}
}
