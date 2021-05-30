package handlers

import (
	"fmt"
	"net/http"
)

func errorHandler(writer http.ResponseWriter, request *http.Request, status int) {
	writer.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(writer, "my 404")
	}
}
