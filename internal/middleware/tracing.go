package middleware

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// Tracing middleware adds OpenTracing spans to HTTP requests
func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create span for this request
		span := opentracing.StartSpan(r.Method + " " + r.URL.Path)
		defer span.Finish()

		// Add tags
		span.SetTag("http.method", r.Method)
		span.SetTag("http.url", r.URL.String())

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// InitNoopTracer initializes a no-op tracer for development
func InitNoopTracer() {
	opentracing.SetGlobalTracer(opentracing.NoopTracer{})
}