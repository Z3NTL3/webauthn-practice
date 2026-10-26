package middleware

import "net/http"

func LimitReqBody(next http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		req.Body = http.MaxBytesReader(res, req.Body, 1024*1024)
		next.ServeHTTP(res, req)
	}
}
