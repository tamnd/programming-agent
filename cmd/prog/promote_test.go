// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tamnd/programming-agent/node"
)

func TestPromote(t *testing.T) {
	root := filepath.Join(t.TempDir(), "corpus")
	if err := os.CopyFS(root, os.DirFS("../../testdata/corpus")); err != nil {
		t.Fatal(err)
	}
	if err := cmdPromote([]string{"-root", root, "-area", "languages/python/versions", "python/proposals/pep-634"}); err != nil {
		t.Fatal(err)
	}
	md, err := os.ReadFile(filepath.Join(root, "nodes/languages/python/versions/pep-634/node.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"tag: 9PEP63\n", "kind: proposal\n", "owner: python\n", "slug: pep-634\n", `proposal_status: "final"`, `authors: ["Brandt Bucher", "Guido van Rossum"]`, "url: https://peps.python.org/pep-0634/"} {
		if !strings.Contains(string(md), s) {
			t.Errorf("node.md has no %q:\n%s", s, md)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "cards/proposal/python/pe.jsonl")); !os.IsNotExist(err) {
		t.Errorf("the shard with one row is still there: %v", err)
	}
	c, err := node.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	n := c.BySlug["python/proposals/pep-634"]
	if n == nil || n.Tag != "9PEP63" {
		t.Fatalf("the promoted node does not load with its slug and tag: %+v %v", n, c.Errors)
	}
	if c.CardBySlug["python/proposals/pep-634"] != nil {
		t.Error("the card is still in the register")
	}
	if err := cmdPromote([]string{"-root", root, "python/proposals/pep-634"}); err == nil {
		t.Error("a second promote of the same slug did not fail")
	}
}
