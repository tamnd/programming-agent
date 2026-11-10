// SPDX-License-Identifier: Apache-2.0

package compat

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Prober runs probes in containers (spec 09.4). Each run gets a new container with no network, a read-only root and a small writable /tmp.
type Prober struct {
	// Engine is the container command, such as docker or podman.
	Engine string
	// Timeout limits one run, including the compile step.
	Timeout time.Duration
	// Exec runs a command and returns its combined output and exit code. Tests replace it. The default runs the command with os/exec.
	Exec func(ctx context.Context, name string, args ...string) (string, int, error)

	digests map[string]string
}

// NewProber returns a prober that uses engine.
func NewProber(engine string) *Prober {
	return &Prober{Engine: engine, Timeout: 3 * time.Minute, Exec: execCmd, digests: map[string]string{}}
}

func execCmd(ctx context.Context, name string, args ...string) (string, int, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return out.String(), ee.ExitCode(), nil
	}
	if err != nil {
		return out.String(), -1, err
	}
	return out.String(), 0, nil
}

// Task is one probe to prove: a statement of an entry for one implementation.
type Task struct {
	Feature string
	// Dir is the node directory on the host.
	Dir    string
	Impl   Impl
	Probe  string
	Added  string
	Expect Expect
	// Edition is the edition of the statement, such as 2018 for "2018+". It fills {edition} in the run command.
	Edition string
}

// Edition returns the first edition of a statement, such as 2018 for "2018+", and def when the statement has none.
func Edition(s, def string) string {
	s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "+"))
	if i := strings.IndexAny(s, " -,"); i > 0 {
		s = s[:i]
	}
	if s == "" {
		return def
	}
	return s
}

// Probe runs the task on version_added, and on the release before when the register has one. It returns the result with the hash of the probe.
func (p *Prober) Probe(ctx context.Context, t Task) (Result, error) {
	hash, err := ProbeHash(t.Dir, t.Probe, t.Expect)
	if err != nil {
		return Result{}, err
	}
	r := Result{Feature: t.Feature, Impl: t.Impl.ID, Probe: t.Probe, Added: t.Added, Hash: hash, Date: time.Now().UTC().Format("2006-01-02")}
	if t.Impl.Image == "" || len(t.Impl.Run) == 0 {
		return r, fmt.Errorf("%s has no image or run command in registers/implementations.yaml", t.Impl.ID)
	}
	var sig *regexp.Regexp
	if t.Expect.FailsWith != "" {
		if sig, err = regexp.Compile(t.Expect.FailsWith); err != nil {
			return r, fmt.Errorf("fails_with of %s: %w", path.Base(t.Probe), err)
		}
	}
	r.Pass = p.run(ctx, t, t.Added)
	r.Pass.OK = r.Pass.Error == "" && r.Pass.Exit == 0 && (t.Expect.Output == "" || normalize(r.Pass.Output) == normalize(t.Expect.Output))
	if before, ok := t.Impl.Before(t.Added); ok {
		f := p.run(ctx, t, before)
		f.OK = f.Error == "" && f.Exit != 0 && sig != nil && sig.MatchString(f.Output)
		r.Fail = &f
	}
	return r, nil
}

// run runs the probe on one release.
func (p *Prober) run(ctx context.Context, t Task, version string) Run {
	ref := t.Impl.Ref(version)
	out := Run{Version: version, Image: ref}
	d, err := p.digest(ctx, ref)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	out.Image = d
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()
	dir := filepath.Dir(filepath.Join(t.Dir, filepath.FromSlash(t.Probe)))
	file := "/probe/" + path.Base(t.Probe)
	args := []string{"run", "--rm", "--network", "none", "--read-only", "--tmpfs", "/tmp:exec,size=512m", "--memory", "1g", "--pids-limit", "256", "-v", dir + ":/probe:ro", "-w", "/tmp", out.Image}
	for _, a := range t.Impl.Run {
		a = strings.ReplaceAll(a, "{file}", file)
		a = strings.ReplaceAll(a, "{edition}", t.Edition)
		args = append(args, a)
	}
	o, code, err := p.Exec(ctx, p.Engine, args...)
	out.Output, out.Exit = trimOutput(o), code
	if ctx.Err() != nil {
		out.Error = "the run took longer than " + p.Timeout.String()
	} else if err != nil {
		out.Error = err.Error()
	}
	return out
}

// digest pulls the image once and returns its reference pinned by digest.
func (p *Prober) digest(ctx context.Context, ref string) (string, error) {
	if d, ok := p.digests[ref]; ok {
		return d, nil
	}
	if o, code, err := p.Exec(ctx, p.Engine, "pull", "--quiet", ref); err != nil || code != 0 {
		return "", fmt.Errorf("pull %s: %s", ref, strings.TrimSpace(firstLines(o, 3)))
	}
	o, code, err := p.Exec(ctx, p.Engine, "image", "inspect", "--format", "{{index .RepoDigests 0}}", ref)
	if err != nil || code != 0 {
		return "", fmt.Errorf("inspect %s: %s", ref, strings.TrimSpace(firstLines(o, 3)))
	}
	d := strings.TrimSpace(o)
	if !strings.Contains(d, "@sha256:") {
		return "", fmt.Errorf("inspect %s: no digest", ref)
	}
	p.digests[ref] = d
	return d, nil
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.TrimRight(s, "\n ")
}

// trimOutput keeps the first 4 KB of the output, which holds the error of a failed run.
func trimOutput(s string) string {
	if len(s) > 4096 {
		return s[:4096] + "\n[output cut]"
	}
	return s
}

func firstLines(s string, n int) string {
	ls := strings.SplitN(s, "\n", n+1)
	if len(ls) > n {
		ls = ls[:n]
	}
	return strings.Join(ls, "\n")
}
