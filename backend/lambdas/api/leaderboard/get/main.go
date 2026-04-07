package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/controllers"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
)

func main() {
	db := config.SetupMongo()
	leaderboardRepo := repositories.NewLeaderboardRepository(db)
	leaderboardSvc := services.NewLeaderboardService(leaderboardRepo, nil)

	ctrl := controllers.NewLeaderboardController(leaderboardSvc)
	lambda.Start(ctrl.GetLeaderboard)
}
