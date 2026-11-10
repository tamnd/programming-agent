// SPDX-License-Identifier: Apache-2.0

package check

import (
	"strings"

	"github.com/tamnd/programming-agent/examples"
	"github.com/tamnd/programming-agent/node"
)

// diagnostics checks that each block with expect=error:<code> names a diagnostic card or node of the language of the block (rule X05), so that the site can link the error to its page.
func (r *run) diagnostics(n *node.Node) {
	type file struct {
		path, body string
		line       int
	}
	files := []file{{n.Rel + "/node.md", n.Body, n.BodyLine}}
	for _, ex := range n.Examples {
		files = append(files, file{ex.Path, ex.Body, ex.Line})
	}
	for _, f := range files {
		for _, b := range examples.Parse(f.body, f.line) {
			kind, code, ok := strings.Cut(b.Attrs["expect"], ":")
			if !ok || kind != "error" || b.Mode() != "fails" {
				continue
			}
			slug := node.FullSlug("diagnostic", b.Lang, code)
			if r.c.BySlug[slug] != nil || r.c.CardBySlug[slug] != nil {
				continue
			}
			r.add("X05", Block, f.path, b.Line, "expect=error:%s names no diagnostic: add a card for %s to cards/diagnostic/%s/", code, slug, b.Lang)
		}
	}
}
