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
		list("Paradigms", n.Paradigms)
		text("Typing", typing(n.Typing))
		text("Memory", n.Memory)
		list("Execution", n.Execution)
		if n.Appeared > 0 {
			text("Appeared", fmt.Sprint(n.Appeared))
		}
		list("Designed by", n.Designers)
		links("Influenced by", n.InfluencedBy)
		text("Standard", n.Standard)
		code("File extensions", n.FileExtensions)
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

// typing writes the typing map of a language as a short phrase, such as "dynamic, strong, checked at runtime and optional-static".
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
			c = append(c, fmt.Sprint(x))
		}
		if len(c) > 0 {
			parts = append(parts, "checked at "+strings.Join(c, " and "))
		}
	case string:
		if v != "" {
			parts = append(parts, "checked at "+v)
		}
	}
	if v, ok := m["inference"].(string); ok && v != "" && v != "n/a" && v != "none" {
		parts = append(parts, v+" inference")
	}
	return strings.Join(parts, ", ")
}
