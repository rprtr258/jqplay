package jq

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJQValidate(t *testing.T) {
	cases := []struct {
		J      string
		Q      string
		O      []JQOpt
		ErrStr string
	}{
		{
			ErrStr: "missing filter, missing JSON",
		},
		{
			J: "{}",
			Q: ".",
			O: []JQOpt{
				{
					Name: "from-file",
				},
			},
			ErrStr: `disallow option "from-file"`,
		},
		{
			Q: ".",
			O: []JQOpt{
				{
					Name:    "null-input",
					Enabled: true,
				},
			},
		},
		{
			O: []JQOpt{
				{
					Name:    "null-input",
					Enabled: true,
				},
			},
			ErrStr: "missing filter",
		},
		{
			J: `{"foo": "bar"}`,
			O: []JQOpt{
				{
					Name:    "null-input",
					Enabled: true,
				},
			},
			ErrStr: "missing filter",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(fmt.Sprintf("j=%q q=%q o=%v", c.J, c.Q, c.O), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)
			jq := &JQ{
				Input:   c.J,
				Query:   c.Q,
				Options: c.O,
			}
			err := jq.Validate()

			if c.ErrStr != "" {
				assert.ErrorContains(err, c.ErrStr)
			}

			if c.ErrStr == "" {
				assert.NoError(err)
			}
		})
	}
}

func TestJQEvalTimeout(t *testing.T) {
	jq := JQ{
		Input: `{"dependencies":{"capnp":{"version":"0.1.4","dependencies":{"es6-promise":{"version":"1.0.0","dependencies":{"es6-promise":{"version":"1.0.0"}}}}}}}`,
		Query: `.dependencies | recurse(to_entries | map(.values.dependencies))`,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	err := newNoLimitJQExec().Eval(ctx, jq, io.Discard)
	cancel()
	assert.ErrorIs(t, err, ErrExecTimeout)
}

func TestJQEvalCancelled(t *testing.T) {
	jq := JQ{
		Input: `{"dependencies":{"capnp":{"version":"0.1.4","dependencies":{"es6-promise":{"version":"1.0.0","dependencies":{"es6-promise":{"version":"1.0.0"}}}}}}}`,
		Query: `.dependencies | recurse(to_entries | map(.values.dependencies))`,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := newNoLimitJQExec().Eval(ctx, jq, io.Discard)
	assert.ErrorIs(t, err, ErrExecCancelled)
}

func TestJQEvalAborted(t *testing.T) {
	jq := JQ{
		Input: `{"dependencies":{"capnp":{"version":"0.1.4","dependencies":{"es6-promise":{"version":"1.0.0","dependencies":{"es6-promise":{"version":"1.0.0"}}}}}}}`,
		Query: `.dependencies | recurse(to_entries | map(.values.dependencies))`,
	}
	err := newLimitJQExec().Eval(context.Background(), jq, io.Discard)
	assert.ErrorIs(t, err, ErrExecAborted)
}

func TestJQEvalRaceCondition(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			jq := JQ{
				Input: `{ "foo": { "bar": { "baz": 123 } } }`,
				Query: ".",
			}
			assert.NoError(t, newNoLimitJQExec().Eval(ctx, jq, io.Discard))
		}()
	}
	wg.Wait()
}

func newNoLimitJQExec() *JQExec {
	return &JQExec{
		LimitResourcesFunc: func(p *os.Process) error { return nil },
	}
}

func newLimitJQExec() *JQExec {
	return &JQExec{
		LimitResourcesFunc: func(p *os.Process) error {
			const MemoryLimit = 1 * 1024 * 1024 // 1 MiB
			const CPUTimeLimit = 1              // 1 percentage
			return limitResources(p, MemoryLimit, CPUTimeLimit)
		},
	}
}
