// SPDX-License-Identifier: Apache-2.0

package node

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadCells(t *testing.T) {
	p := filepath.Join(t.TempDir(), "matrix.yaml")
	src := "python: native\nc:\n  status: absent\n  reason: no nested functions\n  alternative: a function pointer and a context argument\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	cells, err := readCells(p)
	if err != nil {
		t.Fatal(err)
	}
	if c := cells["python"]; c == nil || c.Status != "native" || c.Line != 1 {
		t.Errorf("python = %+v, want native at line 1", c)
	}
	if c := cells["c"]; c == nil || c.Status != "absent" || c.Reason != "no nested functions" || c.Alternative == "" || c.Line != 2 {
		t.Errorf("c = %+v, want absent with a reason and an alternative at line 2", c)
	}
	if cells, err := readCells(filepath.Join(t.TempDir(), "none.yaml")); err != nil || len(cells) != 0 {
		t.Errorf("a missing file gives %v, %v; want no cells and no error", cells, err)
	}
}

func TestExampleVerdict(t *testing.T) {
	for body, want := range map[string]string{
		"Text.\n\n- Verdict: idiom.\n":    "idiom",
		"- **Verdict:** library\n":        "library",
		"- Verdict: n/a. Markup.\n":       "n/a",
		"The verdict is native.\n":        "",
		"Text.\n  - Verdict: absent.\n":   "",
		"- verdict: Absent, because...\n": "absent",
	} {
		if got := exampleVerdict(body); got != want {
			t.Errorf("exampleVerdict(%q) = %q, want %q", body, got, want)
		}
	}
}

func TestCellOrder(t *testing.T) {
	c, err := Load("../testdata/corpus")
	if err != nil {
		t.Fatal(err)
	}
	n := c.BySlug["closures"]
	if n == nil {
		t.Fatal("the test corpus has no closures concept")
	}
	if cell, ok := c.Cell(n, "go"); !ok || cell.From != "example" || cell.Status != "native" {
		t.Errorf("go = %+v, %v; want a native cell from the example", cell, ok)
	}
	if _, ok := c.Cell(n, "rust"); ok {
		t.Error("rust has a cell, but the corpus decides nothing for it")
	}
	n.Cells = map[string]*Cell{"go": {Status: "idiom", Reason: "r", From: "matrix.yaml"}}
	if cell, _ := c.Cell(n, "go"); cell.From != "matrix.yaml" {
		t.Errorf("go = %+v; the cell of matrix.yaml must win over the example", cell)
	}
	c.BySlug["python"].Cells = map[string]*Cell{"closures": {Status: "native", From: "concepts.yaml"}}
	n.Cells = nil
	if cell, _ := c.Cell(n, "python"); cell.From != "concepts.yaml" {
		t.Errorf("python = %+v; the cell of the language must win over the example", cell)
	}
}
