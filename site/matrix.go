// SPDX-License-Identifier: Apache-2.0

package site

import (
	"strings"

	"github.com/tamnd/programming-agent/node"
)

// VerdictRow is one language of the verdict table of a concept page (spec 04.8 and 16.4).
type VerdictRow struct {
	Lang    string
	Name    string
	Verdict string
	// URL goes to the example tab or the feature page of the language.
	URL  string
	Note string
	Tier int
}

// verdictWords are the labels of the verdicts.
var verdictWords = map[string]string{
	"native":  "Native",
	"stdlib":  "Standard library",
	"library": "Library",
	"idiom":   "Idiom",
	"absent":  "Absent",
	"n/a":     "Does not apply",
}

// verdictShort are the labels of the verdicts in the matrix page.
var verdictShort = map[string]string{
	"native":  "yes",
	"stdlib":  "std",
	"library": "lib",
	"idiom":   "idiom",
	"absent":  "no",
	"n/a":     "n/a",
}

// verdictRows gives a row for each tier-1 language, and for each other language that has a cell.
func (b *builder) verdictRows(n *node.Node, examples []ExampleView) []VerdictRow {
	tab := map[string]string{}
	for _, e := range examples {
		tab[e.Lang] = "#" + e.ID
	}
	var out []VerdictRow
	for _, l := range b.c.Langs {
		c, ok := b.c.Cell(n, l.ID)
		if !ok && l.Tier != 1 {
			continue
		}
		row := VerdictRow{Lang: l.ID, Name: l.Name, Tier: l.Tier}
		if !ok {
			row.Verdict = "open"
			row.Note = "Not decided yet."
			out = append(out, row)
			continue
		}
		row.Verdict = c.Status
		switch {
		case c.Feature != "":
			row.URL = "/" + c.Feature + "/"
		case tab[l.ID] != "":
			row.URL = tab[l.ID]
		}
		var note []string
		for _, p := range [][2]string{{"Library: ", c.Library}, {"", c.Reason}, {"Alternative: ", c.Alternative}} {
			if s := strings.TrimSpace(p[1]); s != "" {
				note = append(note, p[0]+sentence(s))
			}
		}
		row.Note = strings.Join(note, " ")
		out = append(out, row)
	}
	return out
}

// sentence starts s with a capital letter and ends it with a full stop.
func sentence(s string) string {
	if s == "" {
		return s
	}
	s = strings.ToUpper(s[:1]) + s[1:]
	if !strings.HasSuffix(s, ".") {
		s += "."
	}
	return s
}

func verdictWord(v string) string {
	if w, ok := verdictWords[v]; ok {
		return w
	}
	if v == "open" {
		return "Open"
	}
	return v
}
