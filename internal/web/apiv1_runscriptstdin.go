package web

import (
	"fmt"
	"net/http"
)

// APIV1RunscriptstdinGetHandler creates a http response for the API /runscript http get requests
func APIV1RunscriptstdinGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(`{"endpoints": ["executes a script by piping it to the standard input (stdin) of the specified command"]}`))
}

// APIV1RunscriptstdinPostHandler creates executes a script by piping it to the standard input (stdin) of the specified command
func APIV1RunscriptstdinPostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	script, error := JsonDecodeScript(responseWriter, request)
	if error != nil {
		return
	}

	if script.StdIn == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write(processResult(responseWriter, 3, fmt.Sprintf("%d Bad Request - This endpoint requires stdin", http.StatusBadRequest)))
		return
	}

	responseWriter.Write(runScript(responseWriter, script))
}
