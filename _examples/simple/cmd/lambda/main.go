package main

import (
	"simple/handlers"

	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/maxstanley/tango/wrapper/lambda"
)

func main() {
	aws_lambda.Start(lambda.Wrapper(handlers.NewPostCanDrink))
}
