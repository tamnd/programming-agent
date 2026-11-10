// SPDX-License-Identifier: Apache-2.0

package markup

import (
	"strings"
	"testing"
)

func TestCitations(t *testing.T) {
	src := "A move ends the old owner.[^1] See `x[^2]` and [^3].\n\n```rust\nlet a = b[^1];\n```\n\nA local note.[^n]\n\n[^n]: A note in the body.\n"
	r, err := RenderWith(src, Options{Sources: 2})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<sup class="cite"><a href="#ref-1" aria-label="Source 1">1</a></sup>`,
		"<code>x[^2]</code>",
		`<sup class="cite missing"`,
		`<span class="o">^</span>`,
		`class="footnotes"`,
	} {
		if !strings.Contains(r.HTML, want) {
			t.Errorf("no %q in:\n%s", want, r.HTML)
		}
	}
	if n := strings.Count(r.HTML, `href="#ref-1"`); n != 1 {
		t.Errorf("%d links to source 1, want 1 (the one in the code block must stay as text)", n)
	}
	if len(r.Cites) != 2 || r.Cites[0] != 1 || r.Cites[1] != 3 {
		t.Errorf("cites = %v, want [1 3]", r.Cites)
	}
	if len(r.BadCites) != 1 || r.BadCites[0] != 3 {
		t.Errorf("bad cites = %v, want [3]", r.BadCites)
	}
	if strings.ContainsAny(r.HTML, markOpen+markClose) {
		t.Errorf("a placeholder is left in:\n%s", r.HTML)
	}
}

func TestInlineRefs(t *testing.T) {
	opt := Options{
		Since: func(slug string) (string, bool) {
			if slug == "python/match-statement" {
				return "3.10", true
			}
			return "", false
		},
		RuleURL: func(owner, id string) string {
			if owner == "rust" {
				return "https://doc.rust-lang.org/reference/" + id
			}
			return ""
		},
	}
	r, err := RenderWith("## The rule {rule:rust/expr.move.intro}\n\nIn {since:python/match-statement} and {since:python/nope}, see {std:c23/6.5.7p4}.\n", opt)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<a href="https://doc.rust-lang.org/reference/expr.move.intro"><code>[expr.move.intro]</code></a>`,
		`<span class="since">3.10</span>`,
		`<span class="since missing"`,
		`<code>C23 6.5.7p4</code>`,
	} {
		if !strings.Contains(r.HTML, want) {
			t.Errorf("no %q in:\n%s", want, r.HTML)
		}
	}
	if len(r.TOC) != 1 || r.TOC[0].Text != "The rule [expr.move.intro]" {
		t.Errorf("toc = %+v", r.TOC)
	}
	if len(r.Missing) != 1 || r.Missing[0] != "since:python/nope" {
		t.Errorf("missing = %v", r.Missing)
	}
}

func TestCallouts(t *testing.T) {
	src := "> [!since v=\"3.10\"]\n> The `match` statement is in 3.10.\n\n> [!pitfall]\n> A *name* in a pattern binds.\n\n> [!wrong]\n> Not a kind.\n\n> A plain quote.\n"
	r, err := RenderWith(src, Options{})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<blockquote class="callout callout-since">`,
		`<p class="callout-title">Since 3.10</p>`,
		"<p>The <code>match</code> statement is in 3.10.</p>",
		`<blockquote class="callout callout-pitfall">`,
		`<p class="callout-title">Pitfall</p>`,
		"<p>A <em>name</em> in a pattern binds.</p>",
		"<blockquote>\n<p>A plain quote.</p>",
	} {
		if !strings.Contains(r.HTML, want) {
			t.Errorf("no %q in:\n%s", want, r.HTML)
		}
	}
	if strings.Contains(r.HTML, "[!since") || strings.Contains(r.HTML, "[!pitfall") {
		t.Errorf("a marker is left in:\n%s", r.HTML)
	}
	if len(r.BadCallouts) != 1 || r.BadCallouts[0] != "wrong" {
		t.Errorf("bad callouts = %v", r.BadCallouts)
	}
}

func TestOwnedWikiLinks(t *testing.T) {
	res := func(slug string) (string, string, bool) {
		switch slug {
		case "python/api/str.split":
			return "str.split", "/python/api/str/#str.split", true
		case "rust/errors/E0382":
			return "E0382", "/rust/errors/#E0382", true
		case "closures":
			return "Closures", "/closures/", true
		}
		return "", "", false
	}
	r, err := Render("Use [[python/api/str.split]], read [[rust/errors/E0382]], [[closures#history]] and [[al:merge-sort]].\n", res)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<a href="/python/api/str/#str.split"><code>str.split</code></a>`,
		`<a href="/rust/errors/#E0382">E0382</a>`,
		`<a href="/closures/#history">Closures</a>`,
		`<a href="https://algorithm.tamnd.com/merge-sort/">Merge sort</a>`,
	} {
		if !strings.Contains(r.HTML, want) {
			t.Errorf("no %q in:\n%s", want, r.HTML)
		}
	}
}
