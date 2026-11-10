// SPDX-License-Identifier: Apache-2.0

package site

import (
	"fmt"
	"slices"
	"strings"

	"github.com/tamnd/programming-agent/node"
)

// FactRow is one row of the facts box of a node page (spec 14.7.4). A row has text, links or both.
type FactRow struct {
	Key   string
	Text  string
	Code  bool
	URL   string
	Links []Link
}

// facts gives the rows of the facts box for the kinds that have no summary card. The rows come from the front matter of the kind (spec 04.5), in a fixed order.
func (b *builder) facts(n *node.Node) []FactRow {
	var out []FactRow
	text := func(k, v string) {
		if v != "" {
			out = append(out, FactRow{Key: k, Text: v})
		}
	}
	list := func(k string, v []string) {
		if len(v) > 0 {
			out = append(out, FactRow{Key: k, Text: strings.Join(v, ", ")})
		}
	}
	links := func(k string, slugs []string) {
		if len(slugs) > 0 {
			out = append(out, FactRow{Key: k, Links: b.links(slugs)})
		}
	}
	code := func(k string, v []string) {
		if len(v) > 0 {
			out = append(out, FactRow{Key: k, Text: strings.Join(v, " "), Code: true})
		}
	}
	people := func(k string, slugs []string) {
		if len(slugs) > 0 {
			out = append(out, FactRow{Key: k, Links: b.personLinks(slugs)})
		}
	}
	url := func(k, u string) {
		if u != "" {
			out = append(out, FactRow{Key: k, Text: strings.TrimPrefix(strings.TrimPrefix(u, "https://"), "http://"), URL: u})
		}
	}
	if n.Owner != "" && n.Kind != "language" {
		out = append(out, FactRow{Key: ownerKey(b, n.Owner), Links: []Link{b.link(n.Owner)}})
	}
	switch n.Kind {
	case "concept":
		text("Applies to", n.AppliesTo)
		links("Realized by", b.realizedBy(n.Slug))
	case "language":
		text("Class", n.Class)
		if n.Tier > 0 {
			text("Tier", fmt.Sprint(n.Tier))
		}
		// The box holds who and when. The technical model is in the table of the section "The language at a glance" (glance).
		if n.Appeared > 0 {
			text("Appeared", fmt.Sprint(n.Appeared))
		}
		people("Designed by", n.Designers)
		text("Standard", n.Standard)
		text("Licence of the reference implementation", n.ReferenceImpl)
	case "feature":
		text("Since", n.Since)
		links("Realizes", n.Realizes)
		list("Proposals", n.Proposals)
	case "release":
		text("Version", n.Version)
		text("Released", n.Released)
		text("End of life", n.EOL)
	case "task":
		text("Category", n.Category)
		if len(n.Impls) > 0 {
			text("Solutions", plural(len(n.Impls), "language", "languages"))
		}
	case "pattern", "technique", "practice":
		text("Family", n.Family)
		text("Intent", n.Intent)
		list("Also called", n.AlsoKnownAs)
		links("Related", n.Related)
		links("Refactorings", n.Refactorings)
	case "tool", "library":
		list("Role", n.Role)
		text("Licence", n.Licence)
		if v, ok := n.Latest["version"]; ok {
			text("Latest release", fmt.Sprint(v))
		}
		text("First release", n.FirstRelease)
		url("Repository", n.Repo)
		url("Home page", n.Homepage)
		links("Alternatives", n.AlternativeTo)
	case "format", "standard":
		text("Media type", n.MediaType)
		code("Extensions", n.Extensions)
		list("Specifications", n.Specs)
		text("Publisher", n.Publisher)
	case "weakness":
		if n.CWE > 0 {
			out = append(out, FactRow{Key: "CWE", Text: fmt.Sprintf("CWE-%d", n.CWE), URL: fmt.Sprintf("https://cwe.mitre.org/data/definitions/%d.html", n.CWE)})
		}
		text("Abstraction", n.Abstraction)
		text("OWASP Top 10 (2025)", n.OWASPTop10)
		if n.CWETop25 > 0 {
			text("CWE Top 25 (2025)", fmt.Sprintf("rank %d", n.CWETop25))
		}
		list("Languages affected", n.LanguagesAffected)
		links("Prevented by", n.Prevents)
	case "diagnostic":
		if n.Code != "" {
			out = append(out, FactRow{Key: "Code", Text: n.Code, Code: true})
		}
		text("Severity", n.Severity)
		links("Concepts", n.Concepts)
		links("Features", n.Features)
	case "api":
		text("Module", n.Module)
		text("Kind of entity", n.Entity)
		text("Since", n.Since)
	case "person":
		text("Born", n.Born)
		text("Died", n.Died)
	case "organization":
		text("Founded", n.Founded)
		url("Home page", n.URL)
	case "work":
		list("Authors", n.Authors)
		if n.Year > 0 {
			text("Year", fmt.Sprint(n.Year))
		}
		url("Link", n.URL)
	}
	return out
}

// ownerKey names the owner row by the kind of the owner node.
func ownerKey(b *builder, owner string) string {
	if o, ok := b.c.BySlug[owner]; ok {
		switch o.Kind {
		case "tool":
			return "Tool"
		case "library":
			return "Library"
		case "format", "standard":
			return "Specification"
		}
	}
	if b.c.LangByID[owner] != nil {
		return "Language"
	}
	return "Owner"
}

// realizedBy lists the features that realize a concept (the reverse of the realizes edge, spec 04.6).
func (b *builder) realizedBy(slug string) []string {
	var out []string
	for _, m := range b.c.Nodes {
		if slices.Contains(m.Realizes, slug) {
			out = append(out, m.Slug)
		}
	}
	slices.Sort(out)
	return out
}

// typing writes the typing map of a language as a short phrase, such as "dynamic, strong, checked at run time and by optional static checkers".
// glance gives the rows of the table of the section "The language at a glance" (spec 06.3): the paradigms, the type system, the memory model and the execution model, from the front matter.
func (b *builder) glance(n *node.Node) []FactRow {
	var out []FactRow
	add := func(k, v string) {
		if v != "" {
			out = append(out, FactRow{Key: k, Text: v})
		}
	}
	add("Paradigms", words(n.Paradigms))
	add("Typing", typing(n.Typing))
	add("Memory", phrase(n.Memory))
	add("Execution", words(n.Execution))
	if len(n.InfluencedBy) > 0 {
		out = append(out, FactRow{Key: "Influenced by", Links: b.nameLinks(n.InfluencedBy)})
	}
	if len(n.FileExtensions) > 0 {
		ext := make([]string, len(n.FileExtensions))
		for i, e := range n.FileExtensions {
			ext[i] = "." + strings.TrimPrefix(e, ".")
		}
		out = append(out, FactRow{Key: "File extensions", Text: strings.Join(ext, " "), Code: true})
	}
	if u, ok := n.Reference["url"].(string); ok && u != "" {
		out = append(out, FactRow{Key: "Definition", Text: strings.TrimPrefix(u, "https://"), URL: u})
	}
	return out
}

// phrases are the words for the ids of the language front matter (spec 04.5) that humanize does not give well.
var phrases = map[string]string{
	"interpreted-bytecode": "compiled to bytecode and interpreted",
	"jit":                  "just-in-time compiled",
	"aot":                  "compiled ahead of time",
	"native":               "compiled to machine code",
	"vm":                   "runs on a virtual machine",
	"transpiled":           "compiled to another language",
	"garbage-collected":    "garbage collected",
	"ownership":            "ownership and borrowing, checked at compile time",
	"manual":               "managed by hand",
	"reference-counted":    "reference counted",
}

// checks are the words for the values of typing.checking.
var checks = map[string]string{
	"runtime":         "at run time",
	"run-time":        "at run time",
	"compile-time":    "at compile time",
	"static":          "at compile time",
	"optional-static": "by optional static checkers",
	"gradual":         "by gradual typing",
}

func checkPhrase(id string) string {
	if c, ok := checks[id]; ok {
		return c
	}
	return "at " + strings.ReplaceAll(id, "-", " ")
}

func phrase(id string) string {
	if p, ok := phrases[id]; ok {
		return p
	}
	return strings.ReplaceAll(id, "-", " ")
}

func words(ids []string) string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = phrase(id)
	}
	return strings.Join(out, ", ")
}

// personLinks links the person nodes. A person with no node shows the name from the slug, with each word capitalized.
func (b *builder) personLinks(slugs []string) []Link {
	out := b.links(slugs)
	for i := range out {
		if out[i].Missing {
			out[i].Name = titleWords(slugs[i])
		}
	}
	return out
}

// nameLinks links the language nodes. A language with no node shows its register name, or the slug.
func (b *builder) nameLinks(slugs []string) []Link {
	out := b.links(slugs)
	for i := range out {
		if !out[i].Missing {
			continue
		}
		if l := b.c.LangByID[slugs[i]]; l != nil {
			out[i].Name = l.Name
		} else {
			out[i].Name = titleWords(strings.TrimSuffix(slugs[i], "-language"))
		}
	}
	return out
}

// particles stay in lower case inside a name, as in "Guido van Rossum".
var particles = map[string]bool{"van": true, "von": true, "de": true, "der": true, "den": true, "da": true, "di": true, "du": true, "del": true, "la": true, "le": true, "bin": true, "ibn": true}

// acronyms are the words that a name writes in capitals, as in "Standard ML" and "ABC". A language that has a node or a card gets its name from there, so this list is only for the names that have neither.
var acronyms = map[string]bool{"abc": true, "ml": true, "sml": true, "apl": true, "awk": true, "clu": true, "sql": true, "cpl": true, "bcpl": true, "pl": true, "jvm": true, "clr": true, "gcc": true, "llvm": true}

// titleWords makes a name from a slug: "guido-van-rossum" gives "Guido van Rossum", and "modula-3" gives "Modula-3".
func titleWords(slug string) string {
	parts := strings.Split(slug, "-")
	var b strings.Builder
	for i, w := range parts {
		if w == "" {
			continue
		}
		if i > 0 {
			if w[0] >= '0' && w[0] <= '9' {
				b.WriteByte('-')
			} else {
				b.WriteByte(' ')
			}
		}
		if i > 0 && particles[w] {
			b.WriteString(w)
			continue
		}
		if acronyms[w] {
			b.WriteString(strings.ToUpper(w))
			continue
		}
		b.WriteString(strings.ToUpper(w[:1]) + w[1:])
	}
	return b.String()
}

func typing(m map[string]any) string {
	if len(m) == 0 {
		return ""
	}
	var parts []string
	for _, k := range []string{"discipline", "strength"} {
		if v, ok := m[k].(string); ok && v != "" && v != "n/a" {
			parts = append(parts, v)
		}
	}
	switch v := m["checking"].(type) {
	case []any:
		var c []string
		for _, x := range v {
			c = append(c, checkPhrase(fmt.Sprint(x)))
		}
		if len(c) > 0 {
			parts = append(parts, "checked "+strings.Join(c, " and "))
		}
	case string:
		if v != "" {
			parts = append(parts, "checked "+checkPhrase(v))
		}
	}
	if v, ok := m["inference"].(string); ok && v != "" && v != "n/a" && v != "none" {
		parts = append(parts, v+" inference")
	}
	return strings.Join(parts, ", ")
}
