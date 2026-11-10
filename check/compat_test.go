// SPDX-License-Identifier: Apache-2.0

package check

import (
	"strings"
	"testing"
)

const featureDir = "nodes/languages/python/functions/closures-and-nonlocal/"

const implRegister = `python:
  - id: cpython
    name: CPython
    main: true
    releases: ["2.7", "3.0", "3.1"]
  - id: pypy
    name: PyPy
`

func compatFindings(t *testing.T, files map[string]string) string {
	t.Helper()
	files["registers/implementations.yaml"] = implRegister
	c := fixture(t, files)
	fs := Run(c, Options{Main: true})
	var out []string
	for _, rule := range []string{"K01", "K02", "K03", "K04", "K05"} {
		for _, f := range fired(fs, rule) {
			out = append(out, rule+" "+strings.TrimPrefix(f, featureDir))
		}
	}
	return strings.Join(out, "\n")
}

func TestKRulesClean(t *testing.T) {
	got := compatFindings(t, map[string]string{
		featureDir + "compat.yaml": `closures-and-nonlocal:
  description: "The nonlocal statement"
  support:
    cpython: {version_added: "3.0", probe: probes/nonlocal.py}
    pypy: {version_added: false}
`,
		featureDir + "probes/nonlocal.py": "def f():\n    n = 0\n    def g():\n        nonlocal n\n",
		featureDir + "probes/probes.yaml": "nonlocal.py: {fails_with: SyntaxError}\n",
	})
	if got != "" {
		t.Fatalf("findings:\n%s", got)
	}
}

func TestKRules(t *testing.T) {
	got := compatFindings(t, map[string]string{
		featureDir + "compat.yaml": `nonlocal:
  description: "The nonlocal statement"
  support:
    cpython: {version_added: "3.2", probe: probes/nonlocal.py}
    jython: {version_added: "2.7"}
    pypy: {version_added: null}
other.x:
  colour: red
  support:
    cpython: {version_added: "3.0", probe: probes/missing.py}
    pypy: {partial_implementation: true, version_added: "7.3.0"}
other.y:
  description: "A probe with no failure signature"
  support:
    cpython: {version_added: "3.1", probe: probes/nonlocal.py}
`,
		featureDir + "probes/nonlocal.py": "nonlocal n\n",
		featureDir + "probes/probes.yaml": "{}\n",
	})
	for _, want := range []string{
		`K01 compat.yaml: block: the first entry is "nonlocal"; it must be the slug "closures-and-nonlocal" of the node`,
		"K01 compat.yaml: block: nonlocal: jython is not a column of python in registers/implementations.yaml",
		"version_added is null",
		`the entry "other.x" does not start with "nonlocal."`,
		`the entry has the key "colour"`,
		"other.x: the entry has no description",
		"the probe probes/missing.py does not exist",
		"a partial implementation needs a note",
		"K02 compat.yaml: block: nonlocal: cpython: 3.2 is not a release of cpython",
		`K03 node.md: block: since is "3.0", but the compat entry nonlocal gives 3.2 for cpython`,
		"K05 compat.yaml: level: other.y: cpython: the probe nonlocal.py has no failure signature",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no finding %q in:\n%s", want, got)
		}
	}
}

func TestK01BadFile(t *testing.T) {
	got := compatFindings(t, map[string]string{featureDir + "compat.yaml": "- a list\n"})
	if !strings.Contains(got, "K01 compat.yaml: block") || !strings.Contains(got, "not a map of entries") {
		t.Fatalf("findings:\n%s", got)
	}
}

func TestK04(t *testing.T) {
	c := fixture(t, nil)
	n := c.BySlug["python/closures-and-nonlocal"]
	if n == nil {
		t.Fatal("no feature python/closures-and-nonlocal in the test corpus")
	}
	n.Body = "## Semantics\n\nThe statement was added in Python 3.0 for this.\n\nIt is {since:python/closures-and-nonlocal} and `since 3.0` in code.\n\n```python\n# new in 3.0\n```\n\n## History\n\nPEP 3104 came in Python 3.0, added in 3.0.\n"
	got := fired(Run(c, Options{Main: true, Only: map[string]bool{n.Slug: true}}), "K04")
	if len(got) != 1 || !strings.Contains(got[0], `"added in Python 3.0"`) {
		t.Fatalf("K04 = %q, want one finding in the section Semantics", got)
	}
}
