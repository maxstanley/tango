package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maxstanley/tango/handler"
)

func Wrapper(path string, handlerFactory func() handler.Handler) (string, http.HandlerFunc) {
	pathParams := map[int]string{}
	pathSplit := strings.Split(path, "/")[1:]

	// Create a mapping between :params and the path index they appear in.
	for index, value := range pathSplit {
		if strings.HasPrefix(value, "{") {
			// slice removes the "{}" from the start and end of the string.
			pathParams[index] = value[1 : len(value)-1]
		}
	}

	return path, func(w http.ResponseWriter, r *http.Request) {
		h := handlerFactory()

		pathMap := parsePath(pathParams, r.URL.Path)
		h.UnmarshalPath(pathMap)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not read body: %s", err.Error())
			return
		}

		if err := h.UnmarshalBody(body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not unmarshal body: %s", err.Error())
			return
		}

		status, stringResponse := h.Handler()

		w.WriteHeader(status)
		fmt.Fprintf(w, stringResponse)
	}
}

func parsePath(pathParams map[int]string, requestPath string) map[string]string {
	requestPathSplit := strings.Split(requestPath, "/")[1:]
	pathParamsPair := map[string]string{}

	// Iterate over pathParams, get the value of the path at each index.``
	for k, v := range pathParams {
		pathParamsPair[v] = requestPathSplit[k]
	}
	return pathParamsPair
}
