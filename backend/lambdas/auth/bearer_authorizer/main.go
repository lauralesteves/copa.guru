package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lauralesteves/copa-guru-backend/internal/middlewares"
)

func main() {
	lambda.Start(middlewares.AuthorizeBearer)
}
