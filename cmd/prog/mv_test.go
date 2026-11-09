// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMv(t *testing.T) {
	root := copyCorpus(t)
	if got := run([]string{"mv", "search-a-sorted-list", "find-in-a-sorted-list", "-root", root}); got != 0 {
		t.Fatalf("prog mv exits with %d", got)
	}
	newDir := filepath.Join(root, "nodes/data-and-values/collections/find-in-a-sorted-list")
	if _, err := os.Stat(filepath.Join(newDir, "impl")); err != nil {
		t.Fatalf("the node did not move with its files: %v", err)
	}
	if exists(filepath.Join(root, "nodes/data-and-values/collections/search-a-sorted-list")) {
		t.Fatal("the old directory is still there")
	}
	md, _ := os.ReadFile(filepath.Join(newDir, "node.md"))
	if !strings.Contains(string(md), "\nslug: find-in-a-sorted-list\n") || !strings.Contains(string(md), "\ntag: 2B7QDX\n") {
		t.Errorf("node.md has the wrong slug or tag:\n%s", md[:200])
	}
	li, _ := os.ReadFile(filepath.Join(root, "nodes/program-analysis-and-verification/assertions-and-contracts/loop-invariant/node.md"))
	if !strings.Contains(string(li), "](/find-in-a-sorted-list/)") || strings.Contains(string(li), "](/search-a-sorted-list/)") {
		t.Error("the link in loop-invariant still names the old slug")
	}
	reg, _ := os.ReadFile(filepath.Join(root, "registers/tags.reg"))
	if !strings.Contains(string(reg), "2B7QDX find-in-a-sorted-list\n") {
		t.Errorf("tags.reg is\n%s", reg)
	}
	al, _ := os.ReadFile(filepath.Join(root, "registers/aliases.reg"))
	if !strings.HasSuffix(string(al), "\nsearch-a-sorted-list find-in-a-sorted-list\n") {
		t.Errorf("aliases.reg is\n%s", al)
	}
	if got := run([]string{"check", "-root", root}); got != 0 {
		t.Fatalf("prog check exits with %d after the move", got)
	}
	out := t.TempDir()
	if got := run([]string{"site", "build", "-root", root, "-out", out}); got != 0 {
		t.Fatalf("prog site build exits with %d", got)
	}
	red, _ := os.ReadFile(filepath.Join(out, "_redirects"))
	for _, want := range []string{"/search-a-sorted-list/* /find-in-a-sorted-list/:splat 301\n", "/search-a-sorted-list /find-in-a-sorted-list/ 301\n"} {
		if !strings.Contains(string(red), want) {
			t.Errorf("_redirects has no line %q", want)
		}
	}
	if moved, _ := os.ReadFile(filepath.Join(out, "api", "moved.json")); !strings.Contains(string(moved), `"/t/2B7QDX":"/find-in-a-sorted-list/"`) {
		t.Errorf("moved.json has no tag of the moved node:\n%s", moved)
	}

	// A second move points the first old slug to the newest slug, so no redirect is a chain.
	if got := run([]string{"mv", "find-in-a-sorted-list", "sorted-search", "-root", root}); got != 0 {
		t.Fatalf("the second prog mv exits with %d", got)
	}
	al, _ = os.ReadFile(filepath.Join(root, "registers/aliases.reg"))
	if !strings.Contains(string(al), "\nsearch-a-sorted-list sorted-search\nfind-in-a-sorted-list sorted-search\n") {
		t.Errorf("aliases.reg after two moves is\n%s", al)
	}
	if got := run([]string{"mv", "sorted-search", "search-a-sorted-list", "-root", root}); got != 2 {
		t.Errorf("a move back to an old slug exits with %d, not 2", got)
	}
}

func TestMvArea(t *testing.T) {
	root := copyCorpus(t)
	if got := run([]string{"mv", "loop-invariant", "performance/complexity-in-practice", "-root", root}); got != 0 {
		t.Fatalf("prog mv exits with %d", got)
	}
	md, err := os.ReadFile(filepath.Join(root, "nodes/performance/complexity-in-practice/loop-invariant/node.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(md), "\narea: performance/complexity-in-practice\n") {
		t.Errorf("the area line did not change:\n%s", md[:200])
	}
	if exists(filepath.Join(root, "registers/aliases.reg")) {
		t.Error("a move to another area wrote aliases.reg, but the slug did not change")
	}
	if exists(filepath.Join(root, "nodes/program-analysis-and-verification/assertions-and-contracts")) {
		t.Error("the empty area directory is still there")
	}
	if got := run([]string{"check", "-root", root}); got != 0 {
		t.Fatalf("prog check exits with %d after the move", got)
	}
	for _, args := range [][]string{
		{"mv", "loop-invariant", "no/such-area"},
		{"mv", "loop-invariant", "Bad_Slug"},
		{"mv", "loop-invariant", "search-a-sorted-list"},
		{"mv", "no-such-node", "x"},
	} {
		if got := run(append(args, "-root", root)); got != 2 {
			t.Errorf("prog %s exits with %d, not 2", strings.Join(args, " "), got)
		}
	}
}
