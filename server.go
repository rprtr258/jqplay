package main

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"

	"github.com/rprtr258/jqplay/jq"
	"github.com/rprtr258/jqplay/middleware"
)

//go:embed public/index.tmpl
var index string

var tmpl = template.Must(template.
	New("index.tmpl").
	Delims("#{", "}").
	Parse(index))

func renderTemplate(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newHTTPServer(cfg Config) (*http.Server, error) {
	h := &JQHandler{
		JQExec: jq.NewJQExec(),
		Config: cfg,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.handleIndex)
	mux.HandleFunc("GET /jq", h.handleJqGet)
	mux.HandleFunc("POST /jq", h.handleJqPost)
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong")) //nolint:errcheck
	})

	// Serve assets from embedded FS
	assetsFS, _ := http.FS(PublicFS).Open("/public")
	if assetsFS != nil {
		assetsFS.Close() //nolint:errcheck
	}
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(PublicFS))))

	// Chain middleware
	handler := middleware.Timeout(5 * time.Second)(middleware.Secure(cfg.IsProd())(middleware.Logger(mux)))

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}, nil
}

func newServer(ctx context.Context, c Config) error {
	srv, err := newHTTPServer(c)
	if err != nil {
		return err
	}

	var g run.Group
	g.Add(run.SignalHandler(ctx, syscall.SIGTERM))
	g.Add(func() error {
		return srv.ListenAndServe()
	}, func(error) {
		ctx, cancel := context.WithTimeout(context.Background(), 28*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("error shutting down server")
		}
	})
	return g.Run()
}
