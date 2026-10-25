package middleware

import (
	"net/http"

	"charm.land/log/v2"
)

var (
	_ http.ResponseWriter = &response{}
)

type response struct {
	res http.ResponseWriter
}

func (r *response) Header() http.Header {
	return r.res.Header()
}

func (r *response) Write(data []byte) (int, error) {
	log.Debug("logger", "response-body", string(data))
	return r.res.Write(data)
}

func (r *response) WriteHeader(statusCode int) {
	r.res.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &response{w}
		next.ServeHTTP(res, r)
	}
}
