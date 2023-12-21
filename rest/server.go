package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	readTimeout       = 2 * time.Minute
	readHeaderTimeout = time.Minute
	readWriteTimeout  = 30 * time.Minute
	idleTimeout       = time.Minute
	gracefulTimeout   = time.Minute
)

// Server is HTTP Server.
type Server struct {
	*http.Server
}

// NewServer will return HTTP server instance.
func NewServer(port int) *Server {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil,
	})

	slog.SetDefault(slog.New(jsonHandler))

	return &Server{
		&http.Server{
			Addr:                         fmt.Sprintf(":%d", port),
			Handler:                      getRoutes(),
			DisableGeneralOptionsHandler: false,
			TLSConfig:                    nil,
			ReadTimeout:                  readTimeout,
			ReadHeaderTimeout:            readHeaderTimeout,
			WriteTimeout:                 readWriteTimeout,
			IdleTimeout:                  idleTimeout,
			MaxHeaderBytes:               http.DefaultMaxHeaderBytes,
			TLSNextProto:                 nil,
			ConnState:                    nil,
			ErrorLog:                     slog.NewLogLogger(jsonHandler, slog.LevelError),
			BaseContext:                  nil,
			ConnContext:                  nil,
		},
	}
}

// Start will start HTTP Server.
func (srv *Server) Start() {
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.LogAttrs(context.TODO(), slog.LevelError,
				"server stopped", slog.String("error", err.Error()))
		}
	}()

	srv.gracefulShutdown()
}

// gracefulShutdown help to gracefull shutdown HTTP Server.
func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(
		context.Background(), gracefulTimeout,
	)

	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "gracefully shutdown failed",
			slog.String("error", err.Error()))
	}

	slog.LogAttrs(context.Background(), slog.LevelInfo, "server shutdown")
}
