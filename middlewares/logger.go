package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type logWritter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
}

func (lwr *logWritter) WriteHeader(code int) {
	if !lwr.wroteHeader {
		lwr.code = code
		lwr.wroteHeader = true
		lwr.ResponseWriter.WriteHeader(code)
	}
}

func (lwr *logWritter) Write(buf []byte) (int, error) {
	lwr.maybeWriteHeader()
	n, err := lwr.ResponseWriter.Write(buf)
	lwr.bytes += n
	return n, err
}

func (lwr *logWritter) maybeWriteHeader() {
	if !lwr.wroteHeader {
		lwr.WriteHeader(http.StatusOK)
	}
}

// Logger middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		logRes := logWritter{
			ResponseWriter: res,
			wroteHeader:    false,
			code:           0,
			bytes:          0,
		}

		start := time.Now()
		defer func() {
			slog.LogAttrs(context.Background(), slog.LevelInfo,
				"request details",
				slog.String("request_id", middleware.GetReqID(req.Context())),
				slog.Time("request_date", start),
				slog.String("proto", req.Proto),
				slog.String("path", req.URL.Path),
				slog.Duration("response_time", time.Since(start)),
				slog.Int("status", logRes.code),
				slog.Int("size", logRes.bytes),
			)
		}()

		next.ServeHTTP(&logRes, req)
	})
}
