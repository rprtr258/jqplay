package jq

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var Path = func() string {
	path, err := exec.LookPath("jq")

	var binDir string
	if err == nil {
		binDir = filepath.Dir(path)
	} else {
		dir, err := os.MkdirTemp("", "jqplay")
		if err != nil {
			panic(err.Error())
		}

		if err := copyJqBin(dir); err != nil {
			panic(err.Error())
		}

		binDir = dir
	}

	return filepath.Join(binDir, "jq")
}()

var Version = func() string {
	// get version from `jq --help`
	// since `jq --version` diffs between versions
	// e.g., 1.3 & 1.4
	var b bytes.Buffer
	cmd := exec.Command(Path, "--help")
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, b.String())
		panic(err.Error())
	}

	out := bytes.TrimSpace(b.Bytes())
	r := regexp.MustCompile(`\[version (.+)\]`)
	if !r.Match(out) {
		panic(fmt.Errorf("can't find jq version: %s", out).Error())
	}

	m := r.FindSubmatch(out)[1]
	return string(m)
}()
