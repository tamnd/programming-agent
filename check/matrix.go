// SPDX-License-Identifier: Apache-2.0

package check

import (
	"slices"
	"strings"

	"github.com/tamnd/programming-agent/node"
)

// matrix runs the rules of the concept and language matrix (spec 04.8).
func (r *run) matrix(n *node.Node) {
	for key, c := range n.Cells {
		p := c.From
		if !slices.Contains(node.Verdicts, c.Status) {
			r.add("E12", Block, p, c.Line, "%s: the status %q is not one of %s", key, c.Status, strings.Join(node.Verdicts, ", "))
			continue
		}
		switch n.Kind {
		case "concept":
			if r.c.LangByID[key] == nil {
				r.add("E12", Block, p, c.Line, "%s is not a language of registers/languages.yaml", key)
			}
			if n.Example(key) != nil {
				r.add("E12", Warn, p, c.Line, "%s has an example file too; the cell here wins over the verdict of the example", key)
			}
		case "language":
			if t, ok := r.c.BySlug[key]; ok && t.Kind != "concept" {
				r.add("E12", Block, p, c.Line, "%s is a %s, not a concept", key, t.Kind)
			}
		}
		switch c.Status {
		case "absent", "n/a", "idiom":
			if strings.TrimSpace(c.Reason) == "" {
				r.add("E12", Block, p, c.Line, "%s: a cell with the status %s needs a reason", key, c.Status)
			}
		case "library":
			if c.Library == "" && c.Reason == "" {
				r.add("E12", Block, p, c.Line, "%s: a library cell names the library", key)
			}
		}
	}
	for _, ex := range n.Examples {
		if ex.Verdict != "" && !slices.Contains(node.Verdicts, ex.Verdict) {
			r.add("E12", Warn, ex.Path, 0, "the verdict %q is not one of %s", ex.Verdict, strings.Join(node.Verdicts, ", "))
		}
	}
	if n.Kind != "concept" || n.AppliesTo != "core-languages" {
		return
	}
	var miss []string
	for _, l := range r.c.Langs {
		if l.Tier != 1 {
			continue
		}
		if _, ok := r.c.Cell(n, l.ID); !ok {
			miss = append(miss, l.ID)
		}
	}
	if len(miss) > 0 {
		r.add("E11", Level, n.Rel+"/node.md", 0, "no matrix cell for the tier-1 languages %s; add an example file or a cell in matrix.yaml", strings.Join(miss, ", "))
	}
}
