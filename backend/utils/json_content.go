package utils

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	enc, err := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"message":"internal server error"}`))
		return
	}
	w.WriteHeader(status)
	w.Write(enc)
}

func OK[T any](w http.ResponseWriter, status int, data T) {
	writeJSON(w, status, Response[T]{Success: true, Data: data})
}

func Fail(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, Response[any]{Success: false, Message: msg})
}
