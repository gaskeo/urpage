package handlers

import "net/http"

func SendJson(writer http.ResponseWriter, jsonAnswer []byte) {
	if len(jsonAnswer) > 0 {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(jsonAnswer)
	}
}
