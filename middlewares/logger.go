package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		logRes := middleware.NewWrapResponseWriter(res, req.ProtoMajor)

		start := time.Now()
		defer func() {
			slog.LogAttrs(context.Background(), slog.LevelInfo,
				"request details",
				slog.String("request_id", middleware.GetReqID(req.Context())),
				slog.Time("request_date", start),
				slog.String("proto", req.Proto),
				slog.String("path", req.URL.Path),
				slog.Duration("response_time", time.Since(start)),
				slog.Int("status", logRes.Status()),
				slog.Int("size", logRes.BytesWritten()),
			)
		}()

		next.ServeHTTP(logRes, req)
	})
}
