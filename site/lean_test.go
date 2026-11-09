// SPDX-License-Identifier: Apache-2.0

package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tamnd/programming-agent/node"
)

func TestLeanPage(t *testing.T) {
	root := filepath.Join(t.TempDir(), "corpus")
	if err := os.CopyFS(root, os.DirFS("../testdata/corpus")); err != nil {
		t.Fatal(err)
	}
	mod := filepath.Join(root, "lean", "Algo", "Searching")
	if err := os.MkdirAll(mod, 0o755); err != nil {
		t.Fatal(err)
	}
	src := "module\n\nnamespace Algo.Searching.BinarySearch\n\ndef search (xs : List Nat) (x : Nat) : Bool := xs.contains x\n\ntheorem search_correct : True := trivial\n\nend Algo.Searching.BinarySearch\n"
	if err := os.WriteFile(filepath.Join(mod, "BinarySearch.lean"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	decls := filepath.Join(t.TempDir(), "decls.jsonl")
	lines := `{"decl":"Algo.Searching.BinarySearch.search","kind":"def","module":"Algo.Searching.BinarySearch","type":"List ℕ → ℕ → Bool","axioms":[],"computable":true,"pos":{"line":5,"end_line":5}}
{"decl":"Algo.Searching.BinarySearch.search_correct","kind":"theorem","module":"Algo.Searching.BinarySearch","type":"True","axioms":["propext"],"pos":{"line":7,"end_line":7}}
`
	if err := os.WriteFile(decls, []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := node.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	n := c.BySlug["search-a-sorted-list"]
	if n == nil {
		t.Fatal("no search-a-sorted-list node")
	}
	n.Lean = map[string]string{"module": "Algo.Searching.BinarySearch", "model": "Algo.Searching.BinarySearch.search", "correct": "Algo.Searching.BinarySearch.search_correct"}
	out := filepath.Join(t.TempDir(), "site")
	if _, err := Build(c, Options{Out: out, Decls: decls, Commit: "0123456789abcdef", Version: "test", Now: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	read := func(p string) string {
		b, err := os.ReadFile(filepath.Join(out, p))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	page := read("search-a-sorted-list/index.html")
	for _, want := range []string{`id="formal-proof"`, `href="/search-a-sorted-list/lean/"`, "<code>search_correct</code>", "only the standard axioms"} {
		if !strings.Contains(page, want) {
			t.Errorf("the node page does not contain %q", want)
		}
	}
	lean := read("search-a-sorted-list/lean/index.html")
	for _, want := range []string{"<title>Lean: Search a sorted list", "lean/Algo/Searching/BinarySearch.lean", "#L7-L7", "search_correct", "propext"} {
		if !strings.Contains(lean, want) {
			t.Errorf("the Lean page does not contain %q", want)
		}
	}
}
