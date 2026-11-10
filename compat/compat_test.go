// SPDX-License-Identifier: Apache-2.0

package compat

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const sample = `match-statement:
  description: "The match statement"
  support:
    cpython: {version_added: "3.10", probe: probes/m.py}
    micropython: {version_added: false, note: "No."}
match-statement.x:
  description: "A sub-feature"
  colour: red
  support:
    rustc:
      - version_added: "1.31.0"
        edition: "2018+"
      - version_added: null
        weight: 3
`

func write(t *testing.T, dir, name, s string) string {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRead(t *testing.T) {
	es, err := Read(write(t, t.TempDir(), "compat.yaml", sample))
	if err != nil {
		t.Fatal(err)
	}
	if len(es) != 2 || es[0].Key != "match-statement" || es[1].Key != "match-statement.x" {
		t.Fatalf("entries = %+v", es)
	}
	if got := strings.Join(es[0].Columns, ","); got != "cpython,micropython" {
		t.Errorf("columns = %s, want the order of the file", got)
	}
	cp := es[0].Support["cpython"][0]
	if cp.Added != "3.10" || cp.Probe != "probes/m.py" || !cp.IsVersion() || cp.Line != 4 {
		t.Errorf("cpython = %+v", cp)
	}
	if mp := es[0].Support["micropython"][0]; mp.Added != "false" || mp.Supported() || mp.Note != "No." {
		t.Errorf("micropython = %+v", mp)
	}
	rs := es[1].Support["rustc"]
	if len(rs) != 2 || rs[0].Edition != "2018+" || rs[1].Added != "null" || strings.Join(rs[1].Unknown, ",") != "weight" {
		t.Errorf("rustc = %+v", rs)
	}
	if strings.Join(es[1].Unknown, ",") != "colour" {
		t.Errorf("unknown keys = %v", es[1].Unknown)
	}
}

func TestImpl(t *testing.T) {
	im := Impl{ID: "cpython", Image: "python:{v}-slim", Releases: []string{"3.8", "3.9", "3.10"}}
	if b, ok := im.Before("3.10"); !ok || b != "3.9" {
		t.Errorf("Before(3.10) = %s, %v", b, ok)
	}
	if _, ok := im.Before("3.8"); ok {
		t.Error("the first release has a release before")
	}
	if im.Known("3.11") || !im.Known("3.9") {
		t.Error("Known is wrong")
	}
	if im.Ref("3.9") != "python:3.9-slim" {
		t.Errorf("Ref = %s", im.Ref("3.9"))
	}
	if Edition("2018+", "2015") != "2018" || Edition("", "2015") != "2015" || Edition("2015", "x") != "2015" {
		t.Error("Edition is wrong")
	}
}

func TestResultsRoundTrip(t *testing.T) {
	root := t.TempDir()
	rs := Results{}
	r := Result{Feature: "f", Impl: "cpython", Probe: "probes/m.py", Added: "3.10", Hash: "h", Pass: Run{OK: true}, Fail: &Run{OK: true}}
	rs[r.Key()] = r
	if err := rs.Write(root); err != nil {
		t.Fatal(err)
	}
	got, err := ReadResults(root)
	if err != nil {
		t.Fatal(err)
	}
	if g := got[r.Key()]; !g.Verified("h") || g.Verified("other") {
		t.Errorf("result = %+v", g)
	}
	if empty, err := ReadResults(t.TempDir()); err != nil || len(empty) != 0 {
		t.Errorf("a missing file gives %v, %v", empty, err)
	}
}

func TestProbe(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "probes/m.py", "print('north')\n")
	var runs []string
	p := NewProber("docker")
	p.Exec = func(_ context.Context, _ string, args ...string) (string, int, error) {
		switch args[0] {
		case "pull":
			return "", 0, nil
		case "image":
			return "python@sha256:" + strings.Repeat("a", 64) + "\n", 0, nil
		}
		runs = append(runs, strings.Join(args, " "))
		// The first run of each probe is on version_added, and it passes.
		if len(runs) == 1 {
			return "north\n", 0, nil
		}
		return "  File \"/probe/m.py\"\nSyntaxError: invalid syntax\n", 1, nil
	}
	p.Timeout = time.Minute
	im := Impl{ID: "cpython", Image: "python:{v}-slim", Run: []string{"python", "{file}"}, Releases: []string{"3.9", "3.10"}}
	r, err := p.Probe(context.Background(), Task{Feature: "f", Dir: dir, Impl: im, Probe: "probes/m.py", Added: "3.10", Expect: Expect{Output: "north", FailsWith: "SyntaxError"}})
	if err != nil {
		t.Fatal(err)
	}
	if !r.Pass.OK || r.Fail == nil || !r.Fail.OK || r.Fail.Version != "3.9" || !r.Verified(r.Hash) {
		t.Fatalf("result = %+v, fail = %+v", r, r.Fail)
	}
	if !strings.Contains(runs[0], "--network none") || !strings.Contains(runs[0], "/probe/m.py") || !strings.Contains(runs[0], dir+"/probes:/probe:ro") {
		t.Errorf("run = %s", runs[0])
	}
	// With no failure signature, the run on the release before cannot count.
	runs = nil
	r, _ = p.Probe(context.Background(), Task{Feature: "f", Dir: dir, Impl: im, Probe: "probes/m.py", Added: "3.10"})
	if r.Fail == nil || r.Fail.OK || r.Verified(r.Hash) {
		t.Errorf("a probe with no signature is verified: %+v", r.Fail)
	}
}
