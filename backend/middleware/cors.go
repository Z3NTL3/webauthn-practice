package middleware

import "net/http"

func CORS(origin string, next http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.Header().Set("Access-Control-Allow-Methods", "POST")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if req.Method == http.MethodOptions {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(res, req)
	}
}
