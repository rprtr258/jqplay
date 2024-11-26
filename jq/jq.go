package jq

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	ErrExecTimeout   = errors.New("jq execution was timeout")
	ErrExecCancelled = errors.New("jq execution was cancelled")
	ErrExecAborted   = errors.New("jq execution was aborted")
	allowedOpts      = map[string]struct{}{
		"slurp":          {},
		"null-input":     {},
		"compact-output": {},
		"raw-input":      {},
		"raw-output":     {},
		"sort-keys":      {},
	}
)

type JQ struct {
	Input   string  `json:"j"`
	Query   string  `json:"q"`
	Options []JQOpt `json:"o"`
}

func (j *JQ) optIsEnabled(name string) bool {
	for _, o := range j.Options {
		if o.Name == name {
			return o.Enabled
		}
	}
	return false
}

type JQOpt struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (o *JQOpt) String() string {
	return fmt.Sprintf("%s (%t)", o.Name, o.Enabled)
}

func (j *JQ) Opts() []string {
	opts := []string{}
	for _, opt := range j.Options {
		if opt.Enabled {
			opts = append(opts, fmt.Sprintf("--%s", opt.Name))
		}
	}
	return opts
}

func (j *JQ) Validate() error {
	errMsgs := []string{}
	if j.Query == "" {
		errMsgs = append(errMsgs, "missing filter")
	}
	if j.Input == "" && !j.optIsEnabled("null-input") {
		errMsgs = append(errMsgs, "missing JSON")
	}
	for _, opt := range j.Options {
		if _, allowed := allowedOpts[opt.Name]; !allowed {
			errMsgs = append(errMsgs, fmt.Sprintf("disallow option %q", opt.Name))
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf("invalid input: %s", strings.Join(errMsgs, ", "))
	}

	return nil
}

func (j JQ) String() string {
	return fmt.Sprintf("j=%s, q=%s, o=%v", j.Input, j.Query, j.Opts())
}

type JQExec struct {
	LimitResourcesFunc func(*os.Process) error
}

func NewJQExec() *JQExec {
	return &JQExec{
		LimitResourcesFunc: func(p *os.Process) error {
			const limitMemory = 48 * 1024 * 1024 // 48 MiB
			const limitCPUTime = 20              // 20 percentage
			return limitResources(p, limitMemory, limitCPUTime)
		},
	}
}

func (e *JQExec) Eval(ctx context.Context, jq JQ, w io.Writer) error {
	if err := jq.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, Path, append(jq.Opts(), jq.Query)...)
	cmd.Stdin = bytes.NewBufferString(jq.Input)
	cmd.Env = make([]string, 0)
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := e.LimitResourcesFunc(cmd.Process); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		if ctxErr := ctx.Err(); ctxErr == context.DeadlineExceeded {
			return ErrExecTimeout
		} else if ctxErr == context.Canceled {
			return ErrExecCancelled
		}

		if strings.Contains(err.Error(), "signal: segmentation fault") ||
			strings.Contains(err.Error(), "signal: aborted") {
			return ErrExecAborted
		}

		return err
	}

	return nil
}
