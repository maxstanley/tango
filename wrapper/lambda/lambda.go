package lambda

import (
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/maxstanley/tango/handler"
	"github.com/maxstanley/tango/wrapper"
)

type gatewayHandler func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func Wrapper(handlerFactory func() handler.Handler) gatewayHandler {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		h := handlerFactory()

		contentType := req.Headers["Content-Type"]
		body := []byte(req.Body)
		bainaryResponseType := wrapper.IsBinaryType(contentType)
		if bainaryResponseType {
			var err error
			body, err = base64.StdEncoding.DecodeString(req.Body)
			if err != nil {
				return response(
					http.StatusInternalServerError,
					"failed to decode request body",
					"text/plain",
				)
			}
		}

		if err := wrapper.UnmarshalRequestBody(h, body, contentType); err != nil {
			return response(http.StatusBadRequest, err.Error(), "text/plain")
		}

		h.UnmarshalPath(req.PathParameters)
		status, protoMessage := h.Handler()

		acceptType := req.Headers["Accept"]
		responseBytes, errStatus, err := wrapper.MarshalResponseBody(protoMessage, acceptType)
		if err != nil {
			return response(errStatus, err.Error(), "text/plain")
		}

		return response(status, string(responseBytes), acceptType)
	}
}

func response(status int, body string, contentType string) (events.APIGatewayProxyResponse, error) {
	bainaryResponseType := wrapper.IsBinaryType(contentType)
	if bainaryResponseType {
		body = base64.StdEncoding.EncodeToString([]byte(body))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		IsBase64Encoded: bainaryResponseType,
	}, nil
}
