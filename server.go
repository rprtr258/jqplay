package main

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/rs/zerolog/log"

	"github.com/owenthereal/jqplay/jq"
	"github.com/owenthereal/jqplay/middleware"
)

//go:embed public/index.tmpl
var index string

var tmpl = template.Must(template.
	New("index.tmpl").
	Delims("#{", "}").
	Parse(index))

func newHTTPServer(cfg *Config) (*http.Server, error) {
	h := &JQHandler{
		JQExec: jq.NewJQExec(),
		Config: cfg,
	}

	router := gin.New()
	router.Use(
		middleware.Timeout(5*time.Second),
		middleware.Secure(cfg.IsProd()),
		middleware.RequestID(),
		middleware.Logger(),
		gin.Recovery(),
	)
	router.SetHTMLTemplate(tmpl)
	router.StaticFS("/assets", http.FS(PublicFS))
	router.GET("/", h.handleIndex)
	router.GET("/jq", h.handleJqGet)
	router.POST("/jq", h.handleJqPost)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}, nil
}

func newServer(ctx context.Context, c *Config) error {
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
