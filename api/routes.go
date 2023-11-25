package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirijagadeesh/simplerestserver/middlewares"
)

func getRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.Heartbeat("/ping"),
		middleware.Recoverer,
		middlewares.Logger, // custom logger
		middleware.RequestID,
		middleware.RealIP,
		middleware.CleanPath,
	)

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, `{"message": "Hello Ramya Ravali"}`) //nolint: errcheck
	})

	router.Get("/{name}", func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "name")
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, fmt.Sprintf(`{"message": "Hello %s"}`, name)) //nolint: errcheck
	})

	if err := chi.Walk(router, printRountes()); err != nil {
		slog.LogAttrs(context.TODO(), slog.LevelError, "Failed to walk routes:",
			slog.String("err", err.Error()))
	}

	return router
}

func printRountes() chi.WalkFunc {
	return func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		slog.LogAttrs(context.Background(), slog.LevelInfo,
			"Registered routes",
			slog.String("method", method),
			slog.String("route", route),
		)
		return nil
	}
}
