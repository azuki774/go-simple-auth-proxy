package telemetry

import (
	"context"
	"net/http"
)

type contextKey string

var RequestIDKey contextKey

const requestIDHeaderName = "X-Request-ID"

// X-Request-ID: bdedbd7f90385bf9966288a5f18a9d6a
func GetRequestID(r *http.Request) contextKey {
	v := r.Header.Get(requestIDHeaderName)
	return contextKey(v)
}

func GetRequestIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(RequestIDKey).(contextKey)
	return string(v)
}
