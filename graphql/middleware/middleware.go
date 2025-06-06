package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const responseWriterKey = contextKey("responseWriter")

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
