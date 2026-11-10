// SPDX-License-Identifier: Apache-2.0

package check

import (
	"strings"
	"testing"
)

const closuresDir = "nodes/functions-and-abstraction/closures-and-scope/closures/"

func TestE11(t *testing.T) {
	c := fixture(t, nil)
	got := fired(Run(c, Options{Main: true}), "E11")
	if len(got) != 1 || !strings.Contains(got[0], "javascript, rust") || !strings.Contains(got[0], ": level: ") {
		t.Fatalf("E11 = %q, want one level finding for javascript and rust", got)
	}

	c = fixture(t, map[string]string{
		closuresDir + "matrix.yaml": "javascript: native\nrust:\n  status: native\n",
	})
	if got := fired(Run(c, Options{Main: true}), "E11"); len(got) != 0 {
		t.Fatalf("E11 = %q, want none when each tier-1 language has a cell", got)
	}
}

func TestE12(t *testing.T) {
	c := fixture(t, map[string]string{
		closuresDir + "matrix.yaml": strings.Join([]string{
			"javascript: maybe",
			"cobol: native",
			"go: native",
			"rust:",
			"  status: absent",
			"",
		}, "\n"),
	})
	got := strings.Join(fired(Run(c, Options{Main: true}), "E12"), "\n")
	for _, want := range []string{
		`javascript: the status "maybe" is not one of`,
		"cobol is not a language of registers/languages.yaml",
		"warn: go has an example file too",
		"rust: a cell with the status absent needs a reason",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("E12 has no finding %q in:\n%s", want, got)
		}
	}
	if !strings.Contains(got, closuresDir+"matrix.yaml") {
		t.Errorf("E12 does not name matrix.yaml: %s", got)
	}

	c = fixture(t, map[string]string{
		closuresDir + "matrix.yaml": "javascript:\n  status: library\n  library: lodash\nrust:\n  status: idiom\n  reason: a struct that holds the state and a method\n",
	})
	if got := fired(Run(c, Options{Main: true}), "E12"); len(got) != 0 {
		t.Fatalf("E12 = %q, want none", got)
	}
}

func TestE12LanguageCells(t *testing.T) {
	c := fixture(t, map[string]string{
		"nodes/languages/python/concepts.yaml": "closures: native\nsearch-a-sorted-list: native\n",
	})
	got := fired(Run(c, Options{Main: true}), "E12")
	if len(got) != 1 || !strings.Contains(got[0], "search-a-sorted-list is a task") || !strings.Contains(got[0], "not a concept") {
		t.Fatalf("E12 = %q, want one finding for the task search-a-sorted-list", got)
	}
}

func TestE12ExampleVerdict(t *testing.T) {
	c := fixture(t, map[string]string{
		closuresDir + "examples/rust.md": "Rust closures capture by reference or by move.\n\n- Verdict: sometimes.\n",
	})
	got := fired(Run(c, Options{Main: true}), "E12")
	if len(got) != 1 || !strings.Contains(got[0], `the verdict "sometimes"`) {
		t.Fatalf("E12 = %q, want one finding for the verdict of the Rust example", got)
	}
}
