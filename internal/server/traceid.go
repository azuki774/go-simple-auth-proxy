package server

import "github.com/google/uuid"

type contextKey string

const traceIdKey contextKey = "trace.id" // Like traceID, not real traceID

func GenerateTraceID() contextKey {
	return contextKey(uuid.New().String())
}
