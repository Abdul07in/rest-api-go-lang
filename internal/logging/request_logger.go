package logging

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type RequestLogger struct {
	logger *log.Logger
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewRequestLogger() *RequestLogger {
	return &RequestLogger{
		logger: log.New(os.Stdout, "[API] ", log.Ldate|log.Ltime|log.LUTC),
	}
}

func (l *RequestLogger) LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		traceID := uuid.New().String()

		// Add trace ID to request context
		ctx := r.Context()
		ctx = AddTraceIDToContext(ctx, traceID)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status code
		rw := NewResponseWriter(w)

		// Add trace ID to response headers
		rw.Header().Set("X-Trace-ID", traceID)

		// Log incoming request
		l.logger.Printf("Request: [TraceID: %s] %s %s from %s",
			traceID,
			r.Method,
			r.URL.Path,
			getClientIP(r),
		)

		// Handle the request
		handler.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(startTime)

		// Log response
		l.logger.Printf("Response: [TraceID: %s] %s %s from %s - Status: %d - Duration: %v",
			traceID,
			r.Method,
			r.URL.Path,
			getClientIP(r),
			rw.statusCode,
			duration,
		)
	})
}

func (l *RequestLogger) LogOperation(traceID, operation, details string) {
	l.logger.Printf("Operation: [TraceID: %s] %s - %s", traceID, operation, details)
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}
