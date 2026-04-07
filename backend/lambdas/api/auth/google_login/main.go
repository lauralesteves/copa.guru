package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lauralesteves/copa-guru-backend/internal/config"
	"github.com/lauralesteves/copa-guru-backend/internal/controllers"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
	gauth "github.com/lauralesteves/copa-guru-backend/internal/services/external/google_oauth"
)

func main() {
	db := config.SetupMongo()
	userRepo := repositories.NewUserRepository(db)
	googleOAuth := gauth.NewService(gauth.NewAdapter(), config.GetGoogleClientID(), config.GetGoogleClientSecret())
	authSvc := services.NewAuthService(userRepo, config.GetJWTSecret(), googleOAuth)
	ctrl := controllers.NewAuthController(authSvc)

	lambda.Start(ctrl.GoogleLogin)
}
