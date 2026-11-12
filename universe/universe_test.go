// SPDX-License-Identifier: Apache-2.0

package universe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tamnd/programming-agent/node"
)

func write(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReadAndProblems(t *testing.T) {
	root := t.TempDir()
	write(t, root, "registers/universe/pldb.jsonl", strings.Join([]string{
		`{"id":"python","scope":"languages","outcome":"node","target":"python"}`,
		`{"id":"pi","scope":"languages","outcome":"out-of-scope","reason":"not-programming"}`,
		`{"id":"x","scope":"languages","outcome":"out-of-scope","reason":"too-small","rule":"libs-100"}`,
		`{"id":"y","scope":"languages","outcome":"out-of-scope","reason":"too-small"}`,
		`{"id":"z","scope":"langs","outcome":"merged"}`,
		`{"id":"w","scope":"api","outcome":"alias","target":"no-such-slug"}`,
		`{"id":"v","scope":"api","outcome":"pending"}`,
		``,
	}, "\n"))
	write(t, root, "registers/universe/bad.jsonl", `{"id":"a","scope":"api","outcome":"pending","extra":1}`+"\n")
	rows, errs := Read(root)
	if len(rows) != 7 || len(errs) != 1 {
		t.Fatalf("got %d rows and %d errors, want 7 and 1: %v", len(rows), len(errs), errs)
	}
	if rows[0].Source != "pldb" || rows[0].Line != 1 || rows[0].File != "registers/universe/pldb.jsonl" {
		t.Fatalf("the first row is %+v", rows[0])
	}
	c := &node.Corpus{BySlug: map[string]*node.Node{"python": {Slug: "python"}}, CardBySlug: map[string]*node.Card{}}
	rules := map[string]ScopeRule{"libs-100": {ID: "libs-100", Rule: "fewer than 100 dependents"}}
	want := [][]string{nil, nil, nil, {"U04"}, {"U01", "U02"}, {"U03"}, nil}
	for i, r := range rows {
		var got []string
		for _, p := range Problems(c, rules, r) {
			got = append(got, p.Rule)
		}
		if strings.Join(got, ",") != strings.Join(want[i], ",") {
			t.Errorf("row %d (%s): got %v, want %v", i, r.ID, got, want[i])
		}
	}
}

func TestReportAndGate(t *testing.T) {
	c := &node.Corpus{
		Nodes: []*node.Node{{Kind: "language"}, {Kind: "concept"}, {Kind: "area"}},
		Cards: []*node.Card{{Kind: "language"}, {Kind: "api"}, {Kind: "api"}, {Kind: "weakness"}},
	}
	rows := []Row{
		{ID: "python", Scope: "languages", Outcome: "node", Target: "python"},
		{ID: "pi", Scope: "languages", Outcome: "out-of-scope", Reason: "not-programming"},
		{ID: "v", Scope: "api", Outcome: "pending"},
	}
	counts := Report(c, rows)
	by := map[string]Count{}
	for _, n := range counts {
		by[n.Scope] = n
	}
	if n := by["languages"]; n.Total != 3 || n.ByOutcome["node"] != 1 || n.ByOutcome["card"] != 1 {
		t.Errorf("languages: %+v", n)
	}
	if n := by["api"]; n.Total != 3 || n.Pending != 1 {
		t.Errorf("api: %+v", n)
	}
	if n := by["all"]; n.Total != 8 || n.Pending != 1 || counts[len(counts)-1].Scope != "all" {
		t.Errorf("all: %+v", n)
	}
	if err := Gate(counts, 0.05, []string{"languages"}); err == nil {
		t.Error("the gate passes with 12.5% pending")
	}
	if err := Gate(counts, 0.2, []string{"languages"}); err != nil {
		t.Errorf("the gate fails: %v", err)
	}
	if err := Gate(counts, 0.2, []string{"api"}); err == nil {
		t.Error("the gate passes with a pending item in a strict scope")
	}
}
