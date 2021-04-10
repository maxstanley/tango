package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxstanley/tango/handler"
	"github.com/maxstanley/tango/wrapper"
)

func parseContext(c *gin.Context) map[string]string {
	pathMap := map[string]string{}

	for _, v := range c.Params {
		pathMap[v.Key] = v.Value
	}

	return pathMap
}

// Wrapper wraps the gin.Context so that a custom context can be
// passed to the requests.
func Wrapper(path string, handlerFactory func() handler.Handler) (string, func(c *gin.Context)) {
	return path, func(c *gin.Context) {
		h := handlerFactory()
		pathMap := parseContext(c)
		body, _ := c.GetRawData()

		contentType := c.ContentType()
		if err := wrapper.UnmarshalRequestBody(h, body, contentType); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		h.UnmarshalPath(pathMap)
		status, protoMessage := h.Handler()
		acceptType := c.Request.Header["Accept"][0]
		c.Header("Content-Type", acceptType)

		responseBytes, errStatus, err := wrapper.MarshalResponseBody(protoMessage, acceptType)
		if err != nil {
			c.String(errStatus, err.Error())
			return
		}

		c.String(status, string(responseBytes))
		return
	}
}
