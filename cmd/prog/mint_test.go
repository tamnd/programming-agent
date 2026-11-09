// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func copyCorpus(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS(corpus)); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestMint(t *testing.T) {
	tagSource = bytes.NewReader(bytes.Repeat([]byte{7, 1, 200, 33, 90, 14, 250, 3}, 64))
	t.Cleanup(func() { tagSource = cryptoSource })
	root := copyCorpus(t)
	for _, s := range []string{"currying", "partial-application"} {
		if got := run([]string{"new", "concept", s, "-area", "functions-and-abstraction/closures-and-scope", "-root", root}); got != 0 {
			t.Fatalf("prog new %s exits with %d", s, got)
		}
	}
	reg := filepath.Join(root, "registers", "tags.reg")
	before, _ := os.ReadFile(reg)
	if got := run([]string{"mint", "-dry-run", "-root", root}); got != 0 {
		t.Fatalf("prog mint -dry-run exits with %d", got)
	}
	if after, _ := os.ReadFile(reg); !bytes.Equal(before, after) {
		t.Fatal("prog mint -dry-run changed tags.reg")
	}
	if got := run([]string{"mint", "-root", root, "-author", "tester"}); got != 0 {
		t.Fatalf("prog mint exits with %d", got)
	}
	tagRE := regexp.MustCompile(`(?m)^tag: ([0-9A-HJKMNP-TV-Z]{6})$`)
	for _, s := range []string{"currying", "partial-application"} {
		b, err := os.ReadFile(filepath.Join(root, "nodes/functions-and-abstraction/closures-and-scope", s, "node.md"))
		if err != nil {
			t.Fatal(err)
		}
		m := tagRE.FindSubmatch(b)
		if m == nil {
			t.Fatalf("%s has no permanent tag:\n%s", s, b)
		}
		tag := string(m[1])
		regText, _ := os.ReadFile(reg)
		if !strings.Contains(string(regText), tag+" "+s+"\n") {
			t.Errorf("tags.reg has no line %q", tag+" "+s)
		}
		prov, err := os.ReadFile(filepath.Join(root, "prov", tag+".jsonl"))
		if err != nil || !strings.Contains(string(prov), `"author":"tester"`) {
			t.Errorf("prov/%s.jsonl is %q, %v", tag, prov, err)
		}
	}
	if got := run([]string{"check", "-root", root}); got != 0 {
		t.Fatalf("prog check after mint exits with %d", got)
	}
	again, _ := os.ReadFile(reg)
	if got := run([]string{"mint", "-root", root}); got != 0 {
		t.Fatalf("a second prog mint exits with %d", got)
	}
	if now, _ := os.ReadFile(reg); !bytes.Equal(again, now) {
		t.Fatal("a second prog mint changed tags.reg")
	}
}

func TestMintSkipsUsedTags(t *testing.T) {
	used := map[string]bool{}
	for range 200 {
		tag, err := newTag(used)
		if err != nil {
			t.Fatal(err)
		}
		if len(tag) != 6 || strings.ContainsAny(tag, "ILOU") {
			t.Fatalf("bad tag %q", tag)
		}
	}
	if len(used) != 200 {
		t.Fatalf("got %d different tags, want 200", len(used))
	}
}

func TestFmt(t *testing.T) {
	root := copyCorpus(t)
	p := filepath.Join(root, "nodes/data-and-values/collections/search-a-sorted-list/node.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, append(bytes.ReplaceAll(b, []byte("\n## "), []byte("  \n\n\n## ")), '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := run([]string{"fmt", "-check", "-root", root}); got != 1 {
		t.Fatalf("prog fmt -check on a changed file exits with %d, want 1", got)
	}
	if got := run([]string{"fmt", "-root", root}); got != 0 {
		t.Fatalf("prog fmt exits with %d", got)
	}
	if after, _ := os.ReadFile(p); !bytes.Equal(after, b) {
		t.Fatalf("prog fmt does not give back the canonical file")
	}
	if got := run([]string{"fmt", "-check", "-root", root}); got != 0 {
		t.Fatalf("prog fmt -check after prog fmt exits with %d", got)
	}
}
