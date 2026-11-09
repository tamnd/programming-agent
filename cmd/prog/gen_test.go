// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGen(t *testing.T) {
	out := t.TempDir()
	if got := run([]string{"gen", "-root", corpus, "-out", out, "-lang", "go,python"}); got != 0 {
		t.Fatalf("prog gen exits with %d", got)
	}
	dir := filepath.Join(out, "go", "data-and-values", "collections", "search-a-sorted-list")
	for _, f := range []string{"go.mod", "binary_search.go", "binary_search_test.go", "prog-test.sh"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("the tree has no %s: %v", f, err)
		}
	}
	b, err := os.ReadFile(filepath.Join(out, "go", "SOURCES"))
	if err != nil {
		t.Fatal(err)
	}
	want := "data-and-values/collections/search-a-sorted-list/binary_search.go\tnodes/data-and-values/collections/search-a-sorted-list/impl/go/binary_search.go\n"
	if !strings.Contains(string(b), want) {
		t.Errorf("SOURCES has no line %q:\n%s", want, b)
	}
	if _, err := os.Stat(filepath.Join(corpus, "nodes", "data-and-values", "collections", "search-a-sorted-list", "impl", "go", "go.mod")); err == nil {
		t.Error("prog gen wrote go.mod into the content tree")
	}

	// A second run for one node keeps the SOURCES lines of the other nodes.
	if got := run([]string{"gen", "-root", corpus, "-out", out, "-lang", "go", "search-a-sorted-list"}); got != 0 {
		t.Fatalf("prog gen search-a-sorted-list exits with %d", got)
	}
	b2, err := os.ReadFile(filepath.Join(out, "go", "SOURCES"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b2) != string(b) {
		t.Errorf("SOURCES changed after a run for one node:\n%s\nwas\n%s", b2, b)
	}
}
