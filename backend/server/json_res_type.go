package server

type JSON_API_SPEC[T any] struct {
	Success bool `json:"success"`
	Message T    `json:"message"`
}
