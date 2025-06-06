package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const AuthTokenKey = contextKey("auth_token")

const responseWriterKey = contextKey("response_writer")

func InjectResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseWriterKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	w, _ := ctx.Value(responseWriterKey).(http.ResponseWriter)
	return w
}

func AuthFromCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			r = r.WithContext(context.WithValue(r.Context(), AuthTokenKey, cookie.Value))
		}
		next.ServeHTTP(w, r)
	})
}
