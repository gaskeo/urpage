package handlers

import (
	"net/http"
)

type Answer struct {
	Err string
}

type Data interface{}

type TemplateData map[string]Data

func SendJson(writer http.ResponseWriter, jsonAnswer []byte) {
	if len(jsonAnswer) > 0 {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(jsonAnswer)
	}
}
