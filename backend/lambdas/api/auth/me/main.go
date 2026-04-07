package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lauralesteves/copa-guru-backend/internal/controllers"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
	"github.com/lauralesteves/copa-guru-backend/internal/services/external/google_oauth"
	"github.com/lauralesteves/copa-guru-backend/internal/shared/config"
)

func main() {
	mongo := config.SetupMongo()
	collection := mongo.Database.Collection("users")

	userRepo := repositories.NewUserRepository(collection)
	jwtSvc := services.NewJWTService(config.GetJWTSecret())
	googleOAuth := google_oauth.NewService(google_oauth.NewAdapter(), config.GetGoogleClientID(), config.GetGoogleClientSecret())
	authSvc := services.NewAuthService(userRepo, jwtSvc, googleOAuth)

	ctrl := controllers.NewAuthController(authSvc)
	lambda.Start(ctrl.Me)
}
