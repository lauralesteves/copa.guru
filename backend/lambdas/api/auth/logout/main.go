package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lauralesteves/copa-guru-backend/internal/controllers"
	"github.com/lauralesteves/copa-guru-backend/internal/repositories"
	"github.com/lauralesteves/copa-guru-backend/internal/services"
	"github.com/lauralesteves/copa-guru-backend/internal/shared/config"
)

func main() {
	mongo := config.SetupMongo()
	collection := mongo.Database.Collection("users")

	userRepo := repositories.NewUserRepository(collection)
	jwtSvc := services.NewJWTService(config.GetJWTSecret())
	googleOAuth := services.NewGoogleOAuthClient(config.GetGoogleClientID(), config.GetGoogleClientSecret())
	authSvc := services.NewAuthService(userRepo, jwtSvc, googleOAuth)

	ctrl := controllers.NewAuthController(authSvc)
	lambda.Start(ctrl.Logout)
}
