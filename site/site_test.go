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
