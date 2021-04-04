package lambda

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/maxstanley/tango/handler"
)

type gatewayHandler func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func Wrapper(handlerFactory func() handler.Handler) gatewayHandler {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		h := handlerFactory()

		if err := h.UnmarshalBody([]byte(req.Body)); err != nil {
			return response(http.StatusBadRequest, fmt.Sprintf("could not unmarshal body: %s", err.Error()))
		}

		h.UnmarshalPath(req.PathParameters)

		return response(h.Handler())
	}
}

func response(status int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
	}, nil
}
