package node

import "testing"

func TestLoad(t *testing.T) {
	c, err := Load("../testdata/corpus")
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Errors) != 0 {
		t.Fatalf("errors: %v", c.Errors)
	}
	if got := len(c.Nodes); got != 7 {
		t.Fatalf("got %d nodes, want 7", got)
	}
	n := c.BySlug["search-a-sorted-list"]
	if n == nil {
		t.Fatal("search-a-sorted-list is missing")
	}
	if n.Kind != "task" {
		t.Errorf("kind is %q, want task", n.Kind)
	}
	if c.ByTag["2B7QDX"] != n {
		t.Error("the tag 2B7QDX does not give search-a-sorted-list")
	}
	if c.LangByID["python"] == nil || c.LangByID["go"] == nil {
		t.Error("the language register is not loaded")
	}
	r := c.BySlug["python/3.14"]
	if r == nil || r.Kind != "release" || r.Local != "3.14" || r.Owner != "python" {
		t.Fatalf("python/3.14 = %+v", r)
	}
	if f := c.BySlug["python/closures-and-nonlocal"]; f == nil || f.Local != "closures-and-nonlocal" {
		t.Error("the feature python/closures-and-nonlocal is missing")
	}
	a := c.AreaByID["languages/python/functions"]
	if a == nil || a.URL() != "/python/functions/" {
		t.Errorf("the language area languages/python/functions = %+v", a)
	}
	cl := c.BySlug["closures"]
	if cl == nil || len(cl.Examples) != 2 || cl.Example("python") == nil {
		t.Fatalf("the concept closures does not have its two example files")
	}
}

func TestIsKind(t *testing.T) {
	for _, k := range []string{"concept", "definition", "language", "feature", "release", "task", "pattern", "weakness", "diagnostic", "api", "exercise"} {
		if !IsKind(k) {
			t.Errorf("IsKind(%q) is false", k)
		}
	}
	for _, k := range []string{"recipe", "algorithm", "problem"} {
		if IsKind(k) {
			t.Errorf("IsKind(%q) is true", k)
		}
	}
	if len(Kinds) != 22 {
		t.Errorf("there are %d kinds, want 22", len(Kinds))
	}
}

func TestFullSlug(t *testing.T) {
	for _, tc := range []struct{ kind, owner, local, want string }{
		{"concept", "", "closures", "closures"},
		{"feature", "python", "match-statement", "python/match-statement"},
		{"release", "python", "3.14", "python/3.14"},
		{"diagnostic", "rust", "E0382", "rust/errors/E0382"},
		{"proposal", "python", "pep-634", "python/proposals/pep-634"},
		{"api", "go", "net/http.Client", "go/api/net/http.Client"},
		{"exercise", "leetcode", "two-sum", "exercises/leetcode/two-sum"},
	} {
		if got := FullSlug(tc.kind, tc.owner, tc.local); got != tc.want {
			t.Errorf("FullSlug(%q, %q, %q) = %q, want %q", tc.kind, tc.owner, tc.local, got, tc.want)
		}
	}
	for _, tc := range []struct {
		kind, local string
		ok          bool
	}{
		{"concept", "closures", true},
		{"concept", "Closures", false},
		{"release", "3.14", true},
		{"release", "c++23", false},
		{"diagnostic", "E0382", true},
		{"api", "net/http.Client", true},
		{"api", "../x", false},
	} {
		if got := LocalOK(tc.kind, tc.local); got != tc.ok {
			t.Errorf("LocalOK(%q, %q) = %v, want %v", tc.kind, tc.local, got, tc.ok)
		}
	}
}

func TestSplitFront(t *testing.T) {
	front, body, line, err := SplitFront([]byte("---\ntag: X\n---\n\n## A\n"))
	if err != nil {
		t.Fatal(err)
	}
	if string(front) != "tag: X\n" {
		t.Errorf("front is %q", front)
	}
	if string(body) != "\n## A\n" {
		t.Errorf("body is %q", body)
	}
	if line != 4 {
		t.Errorf("body line is %d, want 4", line)
	}
	if _, _, _, err := SplitFront([]byte("no front matter\n")); err == nil {
		t.Error("a file without front matter gave no error")
	}
}
