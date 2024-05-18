package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owenthereal/jqplay/jq"
)

type JQHandlerContext struct {
	*Config
	JQ string
}

func (c *JQHandlerContext) Asset(path string) string {
	return fmt.Sprintf("%s/assets/public/%s", c.AssetHost, path)
}

func (c *JQHandlerContext) ShouldInitJQ() bool {
	return c.JQ != ""
}

type JQHandler struct {
	JQExec *jq.JQExec
	Config *Config
}

func (h *JQHandler) handleIndex(c *gin.Context) {
	c.HTML(200, "index.tmpl", &JQHandlerContext{Config: h.Config})
}

func (h *JQHandler) handleJqPost(c *gin.Context) {
	var j jq.JQ
	if err := c.BindJSON(&j); err != nil {
		err = fmt.Errorf("error parsing JSON: %s", err)
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")

	// Evaling into ResponseWriter sets the status code to 200
	// appending error message in the end if there's any
	out := bytes.NewBuffer(nil)
	if err := h.JQExec.Eval(c.Request.Context(), j, io.MultiWriter(c.Writer, out)); err != nil {
		if err == jq.ErrExecAborted || err == jq.ErrExecTimeout {
			h.logger(c).Error("jq error", "error", err, "out", out.String(), "in", j)
		}

		fmt.Fprint(c.Writer, err.Error())
	}
}

func (h *JQHandler) handleJqGet(c *gin.Context) {
	jq := &jq.JQ{
		Input: c.Query("j"),
		Query: c.Query("q"),
	}

	var jqData string
	if err := jq.Validate(); err == nil {
		d, err := json.Marshal(jq)
		if err == nil {
			jqData = string(d)
		}
	}

	c.HTML(http.StatusOK, "index.tmpl", &JQHandlerContext{Config: h.Config, JQ: jqData})
}

func (h *JQHandler) logger(c *gin.Context) *slog.Logger {
	l, _ := c.Get("logger")
	return l.(*slog.Logger)
}
