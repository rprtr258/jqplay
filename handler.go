package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/owenthereal/jqplay/jq"
	"github.com/owenthereal/jqplay/middleware"
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

func (h *JQHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, &JQHandlerContext{Config: h.Config})
}

func (h *JQHandler) handleJqPost(w http.ResponseWriter, r *http.Request) {
	var j jq.JQ
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		err = fmt.Errorf("error parsing JSON: %s", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	out := bytes.NewBuffer(nil)
	mw := io.MultiWriter(w, out)
	if err := h.JQExec.Eval(r.Context(), j, mw); err != nil {
		if err == jq.ErrExecAborted || err == jq.ErrExecTimeout {
			middleware.GetLogger(r).Error("jq error", "error", err, "out", out.String(), "in", j)
		}
		fmt.Fprint(w, err.Error())
	}
}

func (h *JQHandler) handleJqGet(w http.ResponseWriter, r *http.Request) {
	jqObj := &jq.JQ{
		Input: r.URL.Query().Get("j"),
		Query: r.URL.Query().Get("q"),
	}

	var jqData string
	if err := jqObj.Validate(); err == nil {
		d, err := json.Marshal(jqObj)
		if err == nil {
			jqData = string(d)
		}
	}

	renderTemplate(w, &JQHandlerContext{Config: h.Config, JQ: jqData})
}
