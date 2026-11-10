// SPDX-License-Identifier: Apache-2.0

package site

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/tamnd/programming-agent/compat"
	"github.com/tamnd/programming-agent/node"
)

// CompatView is the compat table of a feature page (spec 09.2 and 16.4).
type CompatView struct {
	Columns []compat.Impl
	Rows    []CompatRow
	Notes   []string
	// Verified and Documented count the cells that a probe proves and the cells that come from a source only.
	Verified, Documented int
}

// CompatRow is one compat entry.
type CompatRow struct {
	Key         string
	Description string
	Cells       []CompatCell
}

// CompatCell is the support of one implementation. A cell has one line for each statement.
type CompatCell struct {
	Lines []CompatLine
}

// CompatLine is one statement.
type CompatLine struct {
	Text string
	// Class is yes, no, preview or partial.
	Class string
	// Verified tells if a current probe result proves the statement.
	Verified bool
	Detail   string
	// Notes are the numbers of the notes of the statement.
	Notes []int
}

// compatView builds the table, or nil when the node has no compat entries.
func (b *builder) compatView(n *node.Node) *CompatView {
	if len(n.CompatEntries) == 0 {
		return nil
	}
	lang := n.Owner
	v := &CompatView{}
	used := map[string]bool{}
	for _, e := range n.CompatEntries {
		for _, c := range e.Columns {
			used[c] = true
		}
	}
	for _, im := range b.c.Impls[lang] {
		if used[im.ID] {
			v.Columns = append(v.Columns, im)
			delete(used, im.ID)
		}
	}
	for _, e := range n.CompatEntries {
		for _, c := range e.Columns {
			if used[c] {
				v.Columns = append(v.Columns, compat.Impl{ID: c, Name: c})
				delete(used, c)
			}
		}
	}
	noteNum := map[string]int{}
	for _, e := range n.CompatEntries {
		row := CompatRow{Key: e.Key, Description: e.Description}
		for _, col := range v.Columns {
			var cell CompatCell
			for _, s := range e.Support[col.ID] {
				l := b.compatLine(n, col, s)
				if s.Note != "" {
					if _, ok := noteNum[s.Note]; !ok {
						v.Notes = append(v.Notes, s.Note)
						noteNum[s.Note] = len(v.Notes)
					}
					l.Notes = append(l.Notes, noteNum[s.Note])
				}
				if s.IsVersion() {
					if l.Verified {
						v.Verified++
					} else {
						v.Documented++
					}
				}
				cell.Lines = append(cell.Lines, l)
			}
			row.Cells = append(row.Cells, cell)
		}
		v.Rows = append(v.Rows, row)
	}
	return v
}

func (b *builder) compatLine(n *node.Node, col compat.Impl, s compat.Statement) CompatLine {
	l := CompatLine{}
	switch {
	case s.Added == "false":
		l.Text, l.Class = "No", "no"
	case s.Added == "preview":
		l.Text, l.Class = "Preview", "preview"
	case s.IsVersion():
		l.Text, l.Class = s.Added, "yes"
		if s.Removed != "" {
			l.Text += " to " + s.Removed
			l.Class = "no"
		}
	default:
		l.Text, l.Class = "Unknown", "no"
	}
	if s.Partial {
		l.Class = "partial"
	}
	var d []string
	if s.Edition != "" {
		d = append(d, "edition "+s.Edition)
	}
	if s.Std != "" {
		d = append(d, s.Std)
	}
	if len(s.Flags) > 0 {
		d = append(d, "with "+strings.Join(s.Flags, " "))
	}
	if s.Deprecated != "" {
		d = append(d, "deprecated in "+s.Deprecated)
	}
	if s.AltName != "" {
		d = append(d, "as "+s.AltName)
	}
	if s.Partial {
		d = append(d, "partial")
	}
	l.Detail = strings.Join(d, ", ")
	if s.IsVersion() && s.Probe != "" {
		probe := path.Clean(s.Probe)
		if r, ok := b.c.Probes[compat.Key(n.Slug, col.ID, probe, s.Added)]; ok {
			if h, err := compat.ProbeHash(n.Dir, probe, n.Probes[path.Base(probe)]); err == nil {
				l.Verified = r.Verified(h)
			}
		}
	}
	return l
}

// compatSummary is the line under the table that says how many cells a probe proves.
func compatSummary(v *CompatView) string {
	if v.Verified+v.Documented == 0 {
		return ""
	}
	if v.Documented == 0 {
		return fmt.Sprintf("A probe proves each of the %d versions in this table (spec 09.4).", v.Verified)
	}
	return fmt.Sprintf("A probe proves %d of the %d versions in this table. The others come from the documentation (spec 09.4).", v.Verified, v.Verified+v.Documented)
}

// compatData writes /data/compat/<lang>.json: the compat entries of each language in the shape of the browser-compat-data of MDN, with a __compat object at each level (spec 09.9). The data is CC0.
func (b *builder) compatData() error {
	byLang := map[string]map[string]any{}
	for _, n := range b.c.Nodes {
		if len(n.CompatEntries) == 0 || n.Owner == "" {
			continue
		}
		tree := byLang[n.Owner]
		if tree == nil {
			tree = map[string]any{}
			byLang[n.Owner] = tree
		}
		for _, e := range n.CompatEntries {
			at := tree
			for _, part := range strings.Split(e.Key, ".") {
				sub, _ := at[part].(map[string]any)
				if sub == nil {
					sub = map[string]any{}
					at[part] = sub
				}
				at = sub
			}
			support := map[string]any{}
			for _, col := range e.Columns {
				var list []map[string]any
				for _, s := range e.Support[col] {
					m := map[string]any{"version_added": bcdVersion(s.Added)}
					for k, v := range map[string]string{"version_removed": s.Removed, "version_deprecated": s.Deprecated, "alternative_name": s.AltName, "edition": s.Edition, "std": s.Std, "notes": s.Note, "source": s.Source} {
						if v != "" {
							m[k] = v
						}
					}
					if len(s.Flags) > 0 {
						m["flags"] = s.Flags
					}
					if s.Partial {
						m["partial_implementation"] = true
					}
					if s.IsVersion() && s.Probe != "" {
						probe := path.Clean(s.Probe)
						r, ok := b.c.Probes[compat.Key(n.Slug, col, probe, s.Added)]
						h, err := compat.ProbeHash(n.Dir, probe, n.Probes[path.Base(probe)])
						m["verified"] = ok && err == nil && r.Verified(h)
					}
					list = append(list, m)
				}
				if len(list) == 1 {
					support[col] = list[0]
				} else {
					support[col] = list
				}
			}
			c := map[string]any{"description": e.Description, "support": support, "url": b.opt.BaseURL + "/" + n.Slug + "/#compat-" + e.Key}
			if len(e.Spec) > 0 {
				c["spec"] = e.Spec
			}
			if len(e.Status) > 0 {
				c["status"] = e.Status
			}
			at["__compat"] = c
		}
	}
	for lang, tree := range byLang {
		doc := map[string]any{"__meta": map[string]any{"licence": "CC0-1.0", "source": b.opt.BaseURL + "/", "format": "the browser-compat-data shape of MDN, with implementations as columns"}, lang: tree}
		bs, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		if err := b.write("/data/compat/"+lang+".json", append(bs, '\n')); err != nil {
			return err
		}
	}
	return nil
}

// bcdVersion writes false as a JSON false, as browser-compat-data does.
func bcdVersion(v string) any {
	if v == "false" {
		return false
	}
	return v
}
