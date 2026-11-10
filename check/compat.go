// SPDX-License-Identifier: Apache-2.0

package check

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/tamnd/programming-agent/compat"
	"github.com/tamnd/programming-agent/node"
)

// compat runs the rules of the compat entries of a node (spec 09.2 and 09.4).
func (r *run) compat(n *node.Node) {
	r.versionClaims(n)
	if n.CompatErr != nil {
		r.add("K01", Block, n.Rel+"/"+compatFile(n), 0, "%s", strings.ReplaceAll(n.CompatErr.Error(), n.Rel+"/", ""))
		return
	}
	if len(n.CompatEntries) == 0 {
		if n.Compat != "" {
			r.add("K01", Block, n.Rel+"/node.md", 0, "the field compat names %s, but the file has no entries", n.Compat)
		}
		return
	}
	p := n.Rel + "/" + compatFile(n)
	lang := n.Owner
	if n.Kind == "language" {
		lang = n.Local
	}
	cols := r.c.Impls[lang]
	first := n.CompatEntries[0].Key
	if n.Local != "" && first != n.Local {
		r.add("K01", Block, p, n.CompatEntries[0].Line, "the first entry is %q; it must be the slug %q of the node", first, n.Local)
	}
	used := map[string]bool{}
	for i, e := range n.CompatEntries {
		if i > 0 && !strings.HasPrefix(e.Key, first+".") {
			r.add("K01", Block, p, e.Line, "the entry %q does not start with %q", e.Key, first+".")
		}
		for _, k := range e.Unknown {
			r.add("K01", Block, p, e.Line, "%s: the entry has the key %q, which is not in the schema", e.Key, k)
		}
		if strings.TrimSpace(e.Description) == "" {
			r.add("K01", Block, p, e.Line, "%s: the entry has no description", e.Key)
		}
		if len(e.Support) == 0 {
			r.add("K01", Block, p, e.Line, "%s: the entry has no support statements", e.Key)
		}
		for _, col := range e.Columns {
			im, known := compat.Impl{}, false
			if len(cols) > 0 {
				if im, known = r.c.Impls.Impl(lang, col); !known {
					r.add("K01", Block, p, e.Line, "%s: %s is not a column of %s in registers/implementations.yaml", e.Key, col, lang)
				}
			}
			for _, s := range e.Support[col] {
				r.statement(n, p, e.Key, col, s, im, known, used)
			}
		}
	}
	for f := range n.Probes {
		if !used[f] {
			r.add("K01", Warn, n.Rel+"/probes/probes.yaml", 0, "%s: no support statement uses this probe", f)
		}
	}
	r.since(n, p, lang)
}

// statement checks one support statement.
func (r *run) statement(n *node.Node, p, key, col string, s compat.Statement, im compat.Impl, known bool, used map[string]bool) {
	for _, k := range s.Unknown {
		r.add("K01", Block, p, s.Line, "%s: %s: the statement has the key %q, which is not in the schema", key, col, k)
	}
	switch s.Added {
	case "":
		r.add("K01", Block, p, s.Line, "%s: %s: the statement has no version_added", key, col)
	case "true", "null":
		r.add("K01", Block, p, s.Line, "%s: %s: version_added is %s; write the first version, false or preview", key, col, s.Added)
	}
	if s.Partial && strings.TrimSpace(s.Note) == "" {
		r.add("K01", Block, p, s.Line, "%s: %s: a partial implementation needs a note", key, col)
	}
	if known && len(im.Releases) > 0 {
		for _, v := range []string{s.Added, s.Removed, s.Deprecated} {
			if v != "" && v != "false" && v != "preview" && v != "true" && v != "null" && !im.Known(v) {
				r.add("K02", Block, p, s.Line, "%s: %s: %s is not a release of %s in registers/implementations.yaml", key, col, v, im.ID)
			}
		}
	}
	if s.Probe == "" {
		return
	}
	file := path.Clean(s.Probe)
	if strings.HasPrefix(file, "../") || path.IsAbs(file) {
		r.add("K01", Block, p, s.Line, "%s: %s: the probe %s is outside the node directory", key, col, s.Probe)
		return
	}
	if !exists(r.c.Root, n.Rel+"/"+file) {
		r.add("K01", Block, p, s.Line, "%s: %s: the probe %s does not exist", key, col, s.Probe)
		return
	}
	used[path.Base(file)] = true
	if !s.IsVersion() {
		return
	}
	if _, has := im.Before(s.Added); known && has && n.Probes[path.Base(file)].FailsWith == "" {
		r.add("K05", Level, p, s.Line, "%s: %s: the probe %s has no failure signature (fails_with) in probes/probes.yaml, so the run on the release before %s cannot count", key, col, path.Base(file), s.Added)
	}
}

// since checks that the field since of the front matter agrees with the main column of the first entry (rule K03).
func (r *run) since(n *node.Node, p, lang string) {
	if n.Since == "" {
		return
	}
	e := n.CompatEntries[0]
	col := mainColumn(r.c.Impls[lang], e)
	for _, s := range e.Support[col] {
		if s.Added == n.Since {
			return
		}
	}
	var got []string
	for _, s := range e.Support[col] {
		got = append(got, s.Added)
	}
	r.add("K03", Block, n.Rel+"/node.md", 0, "since is %q, but the compat entry %s gives %s for %s", n.Since, e.Key, strings.Join(got, " and "), col)
}

// mainColumn is the first main implementation of the register that the entry has, or else the first column of the entry.
func mainColumn(cols []compat.Impl, e compat.Entry) string {
	for _, im := range cols {
		if _, ok := e.Support[im.ID]; ok && (im.Main || im.Standard) {
			return im.ID
		}
	}
	if len(e.Columns) > 0 {
		return e.Columns[0]
	}
	return ""
}

func compatFile(n *node.Node) string {
	if n.Compat != "" {
		return n.Compat
	}
	return "compat.yaml"
}

// versionClaim finds a version claim in prose, such as "added in Python 3.10", "since Java 16" or "deprecated in C++17".
var versionClaim = regexp.MustCompile(`(?i)\b(?:since|added in|introduced in|new in|available in|available from|deprecated in|removed in|starting with|starting in|as of)\s+(?:(?:[A-Z][\w#+]*(?:\s[A-Z][\w#+]*)?)\s)?(?:version\s)?v?(?:\d+\.\d+(?:\.\d+)?|(?:C\+\+|C|ES)\d{2,4})\b`)

// noClaimSections are the sections where a version in prose is history, not a claim.
var noClaimSections = []string{"History", "References", "Further reading", "Changes", "Release", "Releases"}

// versionClaims finds version claims in the prose of a node (rule K04). A claim belongs in a compat entry, or in a {since:} mark that renders it from compat.
func (r *run) versionClaims(n *node.Node) {
	if n.Kind == "release" || n.Kind == "person" || n.Kind == "work" || n.Kind == "organization" {
		return
	}
	section := ""
	fence := ""
	for i, line := range strings.Split(n.Body, "\n") {
		t := strings.TrimSpace(line)
		if fence != "" {
			if strings.HasPrefix(t, fence) {
				fence = ""
			}
			continue
		}
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			fence = t[:3]
			continue
		}
		if strings.HasPrefix(t, "## ") {
			section = strings.TrimSpace(t[3:])
			continue
		}
		if slices.Contains(noClaimSections, section) || strings.HasPrefix(t, "> [!history") || strings.HasPrefix(t, "> [!since") {
			continue
		}
		prose := inlineCode.ReplaceAllString(line, "")
		prose = sinceMark.ReplaceAllString(prose, "")
		if m := versionClaim.FindString(prose); m != "" {
			r.add("K04", Warn, n.Rel+"/node.md", n.BodyLine+i, "%q is a version claim in prose; put it in compat.yaml and write {since:...}", m)
		}
	}
}

var (
	inlineCode = regexp.MustCompile("`[^`]*`")
	sinceMark  = regexp.MustCompile(`\{since:[^}]*\}`)
)

// exists tells if the file rel under root exists.
func exists(root, rel string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil
}
