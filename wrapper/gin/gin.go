package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxstanley/tango/handler"
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
		if err := h.UnmarshalBody(body); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("could not unmarshal body: %s", err.Error()))
			return
		}

		h.UnmarshalPath(pathMap)

		c.String(h.Handler())
	}
}
