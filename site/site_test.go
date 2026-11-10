package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tamnd/programming-agent/node"
)

func buildTest(t *testing.T) string {
	t.Helper()
	c, err := node.Load("../testdata/corpus")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "site")
	r, err := Build(c, Options{Out: out, Commit: "0123456789abcdef", Version: "test", Now: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.MathErrors) != 0 {
		t.Errorf("math errors: %v", r.MathErrors)
	}
	if r.Pages == 0 || r.Files < r.Pages {
		t.Errorf("the report is wrong: %d pages, %d files", r.Pages, r.Files)
	}
	return out
}

func TestBuild(t *testing.T) {
	out := buildTest(t)
	for _, f := range []string{
		"index.html", "404.html", "search-a-sorted-list/index.html", "closures/index.html",
		"python/index.html", "python/closures-and-nonlocal/index.html", "python/3.14/index.html",
		"python/functions/index.html", "data-and-values/index.html", "tasks/index.html", "spec/index.html", "search.json", "sitemap.xml", "robots.txt",
		"opensearch.xml", "assets/site.css", "assets/player.js",
	} {
		if _, err := os.Stat(filepath.Join(out, f)); err != nil {
			t.Errorf("%s is missing", f)
		}
	}
	page, err := os.ReadFile(filepath.Join(out, "search-a-sorted-list/index.html"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(page)
	for _, want := range []string{
		"<title>Search a sorted list",
		`id="see-it-run"`,
		`data-scene="search-a-sorted-list"`,
		`id="references"`,
		"https://programming.tamnd.com/search-a-sorted-list/",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("the task page does not contain %q", want)
		}
	}
	if strings.Count(s, `id="references"`) != 1 {
		t.Error("the task page has more than one References section")
	}
}

func TestBuildReplacesOutput(t *testing.T) {
	c, err := node.Load("../testdata/corpus")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "site")
	if err := os.MkdirAll(out, 0o755); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(out, "stale.html")
	if err := os.WriteFile(stale, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Build(c, Options{Out: out}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Error("the build kept a file from the last build")
	}
	if _, err := os.Stat(out + ".next"); !os.IsNotExist(err) {
		t.Error("the build left its work directory")
	}
}

func TestKindPages(t *testing.T) {
	out := buildTest(t)
	read := func(p string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(out, p))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	checks := map[string][]string{
		// A concept shows its example files as tabs, and the features that realize it.
		"closures/index.html": {`class="code examples"`, `id="code-python"`, `id="code-go"`, `data-lang="python"`, `<figcaption>Output</figcaption>`, `href="/python/closures-and-nonlocal/"`},
		// A feature names its language and its first version.
		"python/closures-and-nonlocal/index.html": {`<dt>Language</dt>`, `<dd><a href="/python/">Python</a></dd>`, `<dt>Since</dt><dd>3.0</dd>`},
		// A language page has the facts of the language and the parts of its guide.
		// The table of the section "The language at a glance" comes from the front matter.
		"python/index.html": {`<table class="glance">`, `<th scope="row">Typing</th><td>dynamic, strong, checked at run time and by optional static checkers</td>`, `<dt>Appeared</dt><dd>1991</dd>`, `<code>.py .pyi</code>`, `href="/python/functions/"`, `href="/python/3.14/"`},
	}
	for p, wants := range checks {
		s := read(p)
		for _, w := range wants {
			if !strings.Contains(s, w) {
				t.Errorf("%s does not contain %q", p, w)
			}
		}
	}
	if s := read("matrix/index.html"); !strings.Contains(s, `href="/closures/#code-python"`) {
		t.Error("the matrix has no cell for the Python example of closures")
	}
	for _, p := range []string{"matrix/index.html", "languages/index.html", "llms.txt"} {
		if s := read(p); strings.Contains(s, "algorithm") || strings.Contains(s, "Algorithm") {
			t.Errorf("%s still has text of the algorithm site", p)
		}
	}
}

func TestCardPages(t *testing.T) {
	out := buildTest(t)
	read := func(p string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(out, p))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	cases := map[string][]string{
		"rust/errors/index.html":         {`id="E0382"`, "use of moved value", `href="/rust/errors/E0382.json"`, "Help write this page"},
		"python/proposals/index.html":    {`id="pep-634"`, "Structural Pattern Matching: Specification", "Brandt Bucher, Guido van Rossum", "peps.python.org"},
		"python/api/index.html":          {`href="/python/api/str/"`, "1 card"},
		"python/api/str/index.html":      {`id="str.split"`, "str.split(sep=None, maxsplit=-1)"},
		"javascript/api/JSON/index.html": {`id="JSON.parse"`},
		"python/index.html":              {`href="/python/api/"`, `href="/python/proposals/"`},
		"api/moved.json":                 {`"/rust/errors/E0382":"/rust/errors/#E0382"`, `"/t/9STRSP":"/python/api/str/#str.split"`},
		"python/api/str.split.json":      {`"url":"https://programming.tamnd.com/python/api/str/#str.split"`},
	}
	for p, wants := range cases {
		s := read(p)
		for _, w := range wants {
			if !strings.Contains(s, w) {
				t.Errorf("%s has no %q", p, w)
			}
		}
	}
}

func TestNaturalLess(t *testing.T) {
	for _, c := range [][2]string{{"e0382", "e0499"}, {"pep-8", "pep-634"}, {"e99", "e100"}, {"str.join", "str.split"}} {
		if !naturalLess(c[0], c[1]) || naturalLess(c[1], c[0]) {
			t.Errorf("%s is not before %s", c[0], c[1])
		}
	}
}

func TestTitleWords(t *testing.T) {
	for in, want := range map[string]string{
		"guido-van-rossum": "Guido van Rossum",
		"modula-3":         "Modula-3",
		"brendan-eich":     "Brendan Eich",
		"icon":             "Icon",
	} {
		if got := titleWords(in); got != want {
			t.Errorf("titleWords(%q) = %q, want %q", in, got, want)
		}
	}
}
