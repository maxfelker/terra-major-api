package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ReturnError(writer http.ResponseWriter, message string, codes ...int) {
	code := http.StatusBadRequest
	if len(codes) > 0 {
		code = codes[0]
	}

	resp := ErrorResponse{Error: message}
	jsonResponse, _ := json.Marshal(resp)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(jsonResponse)
}
