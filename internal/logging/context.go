package logging

import (
	"context"
)

type contextKey string

const traceIDKey contextKey = "traceID"

func AddTraceIDToContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return "unknown"
}
