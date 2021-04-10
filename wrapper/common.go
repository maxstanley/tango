package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/maxstanley/tango/handler"
	"google.golang.org/protobuf/proto"
)

func IsBinaryType(contentType string) bool {
	return map[string]bool{
		"application/protobuf":   true,
		"application/x-protobuf": true,
	}[contentType]
}

func UnmarshalRequestBody(h handler.Handler, body []byte, contentType string) error {
	switch contentType {
	case "application/json":
		if err := h.UnmarshalJSONBody(body); err != nil {
			return fmt.Errorf("could not unmarshal body: %s", err.Error())
		}
	case "application/x-protobuf":
		if err := h.UnmarshalProtoBody(body); err != nil {
			return fmt.Errorf("could not unmarshal protobuf body: %s", err.Error())
		}
	default:
		return fmt.Errorf("%s is not a supported Content-Type", contentType)
	}

	return nil
}

func MarshalResponseBody(protoMessage proto.Message, acceptType string) (protoBytes []byte, errStatus int, err error) {
	switch acceptType {
	case "application/json":
		protoBytes, err = json.Marshal(protoMessage)
	case "application/x-protobuf":
		protoBytes, err = proto.Marshal(protoMessage)
	default:
		return nil, http.StatusUnsupportedMediaType, fmt.Errorf("%s in not a supported Accept Type", acceptType)
	}

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not marshal response")
	}

	return
}
