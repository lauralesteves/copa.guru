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
	userRepo := repositories.NewUserRepository(db)
	authSvc := services.NewAuthService(userRepo, config.GetJWTSecret(), nil)

	ctrl := controllers.NewAuthController(authSvc)
	lambda.Start(ctrl.Me)
}
