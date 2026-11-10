// SPDX-License-Identifier: Apache-2.0

package markup

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// The inline marks of spec 06.2 and 06.6 are [^n], a citation of the source n of the front matter, and {rule:...}, {std:...} and {since:...}. The renderer changes each mark to a placeholder of two private-use characters around an index before goldmark runs, and changes the placeholder to HTML after it. So a mark is never parsed as a link or as emphasis, and a mark inside code stays as it is.
const (
	markOpen  = "\ue000"
	markClose = "\ue001"
)

var (
	markRE        = regexp.MustCompile(`\[\^([0-9]+)\]|\{(rule|std|since):([^{}\s]+)\}`)
	placeholderRE = regexp.MustCompile(markOpen + `([0-9]+)` + markClose)
	footDefRE     = regexp.MustCompile(`(?m)^\[\^([^\]\s]+)\]:`)
)

// marks holds the HTML and the plain text of each placeholder of one render.
type marks struct {
	html, text []string
}

func (m *marks) add(h, t string) string {
	m.html = append(m.html, h)
	m.text = append(m.text, t)
	return markOpen + strconv.Itoa(len(m.html)-1) + markClose
}

// pre changes the marks outside code to placeholders.
func (m *marks) pre(src string, opt Options, res *Result) string {
	defined := map[string]bool{}
	for _, d := range footDefRE.FindAllStringSubmatch(src, -1) {
		defined[d[1]] = true
	}
	lines := strings.Split(src, "\n")
	fence := ""
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			if fence == "" {
				fence = t[:3]
			} else if strings.HasPrefix(t, fence) {
				fence = ""
			}
			continue
		}
		if fence != "" || strings.HasPrefix(l, "    ") {
			continue
		}
		lines[i] = outsideCode(l, func(s string) string {
			return markRE.ReplaceAllStringFunc(s, func(x string) string {
				sm := markRE.FindStringSubmatch(x)
				if sm[1] != "" {
					if defined[sm[1]] {
						return x
					}
					return m.cite(sm[1], opt, res)
				}
				return m.ref(sm[2], sm[3], opt, res)
			})
		})
	}
	return strings.Join(lines, "\n")
}

func (m *marks) cite(n string, opt Options, res *Result) string {
	k, _ := strconv.Atoi(n)
	res.Cites = append(res.Cites, k)
	if k < 1 || k > opt.Sources {
		res.BadCites = append(res.BadCites, k)
		return m.add(fmt.Sprintf(`<sup class="cite missing" title="The front matter has no source %d">%d</sup>`, k, k), "")
	}
	return m.add(fmt.Sprintf(`<sup class="cite"><a href="#ref-%d" aria-label="Source %d">%d</a></sup>`, k, k, k), "")
}

func (m *marks) ref(kind, arg string, opt Options, res *Result) string {
	switch kind {
	case "rule":
		owner, id, ok := strings.Cut(arg, "/")
		if !ok {
			owner, id = "", arg
		}
		title := "Rule " + id
		if owner != "" {
			title = "Rule " + id + " of the " + ownerName(owner) + " reference"
		}
		inner := `<code>[` + html.EscapeString(id) + `]</code>`
		if opt.RuleURL != nil {
			if u := opt.RuleURL(owner, id); u != "" {
				inner = `<a href="` + html.EscapeString(u) + `">` + inner + `</a>`
			}
		}
		return m.add(`<span class="rule-ref" title="`+html.EscapeString(title)+`">`+inner+`</span>`, "["+id+"]")
	case "std":
		std, clause, _ := strings.Cut(arg, "/")
		label := strings.ToUpper(std)
		if clause != "" {
			label += " " + clause
		}
		return m.add(`<span class="std-ref"><code>`+html.EscapeString(label)+`</code></span>`, label)
	case "since":
		v, ok := "", false
		if opt.Since != nil {
			v, ok = opt.Since(arg)
		}
		if !ok || v == "" {
			res.Missing = append(res.Missing, "since:"+arg)
			return m.add(`<span class="since missing" title="No version is known for `+html.EscapeString(arg)+`">?</span>`, "?")
		}
		return m.add(`<span class="since">`+html.EscapeString(v)+`</span>`, v)
	}
	return ""
}

// post changes the placeholders in the HTML to their HTML.
func (m *marks) post(s string) string {
	if len(m.html) == 0 {
		return s
	}
	return placeholderRE.ReplaceAllStringFunc(s, func(x string) string {
		i, _ := strconv.Atoi(placeholderRE.FindStringSubmatch(x)[1])
		return m.html[i]
	})
}

// plain changes the placeholders in a text, such as a heading of the table of contents, to plain text.
func (m *marks) plain(s string) string {
	if len(m.text) == 0 {
		return s
	}
	return placeholderRE.ReplaceAllStringFunc(s, func(x string) string {
		i, _ := strconv.Atoi(placeholderRE.FindStringSubmatch(x)[1])
		return m.text[i]
	})
}

func ownerName(o string) string {
	switch o {
	case "rust":
		return "Rust"
	case "fls":
		return "Ferrocene Language Specification"
	}
	return o
}

// outsideCode calls f on each part of the line that is not in a code span.
func outsideCode(l string, f func(string) string) string {
	if !strings.Contains(l, "`") {
		return f(l)
	}
	var b strings.Builder
	rest := l
	for rest != "" {
		i := strings.IndexByte(rest, '`')
		if i < 0 {
			b.WriteString(f(rest))
			break
		}
		b.WriteString(f(rest[:i]))
		rest = rest[i:]
		n := 0
		for n < len(rest) && rest[n] == '`' {
			n++
		}
		run := rest[:n]
		end := closingRun(rest[n:], n)
		if end < 0 {
			b.WriteString(run)
			rest = rest[n:]
			continue
		}
		b.WriteString(rest[:n+end+n])
		rest = rest[n+end+n:]
	}
	return b.String()
}

// closingRun finds a run of exactly n backticks in s and returns its start, or -1.
func closingRun(s string, n int) int {
	for i := 0; i < len(s); {
		if s[i] != '`' {
			i++
			continue
		}
		j := i
		for j < len(s) && s[j] == '`' {
			j++
		}
		if j-i == n {
			return i
		}
		i = j
	}
	return -1
}

// Callouts are block quotes whose first line is [!note], [!warning], [!pitfall], [!history] or [!since v="3.10"] (spec 06.2). The renderer gives the block quote the class callout and a title.
var calloutRE = regexp.MustCompile(`^\[!([a-z]+)((?:\s+[a-z]+="[^"]*")*)\]\s*$`)
var calloutAttrRE = regexp.MustCompile(`([a-z]+)="([^"]*)"`)

// CalloutKinds are the kinds of callout and their titles.
var CalloutKinds = map[string]string{"note": "Note", "warning": "Warning", "pitfall": "Pitfall", "history": "History", "since": "Since", "tip": "Tip"}

type calloutExt struct{ res *Result }

func (e calloutExt) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(util.Prioritized(calloutTransformer(e), 100)))
}

type calloutTransformer struct{ res *Result }

func (t calloutTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	src := reader.Source()
	var quotes []*ast.Blockquote
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if q, ok := n.(*ast.Blockquote); ok && entering {
			quotes = append(quotes, q)
		}
		return ast.WalkContinue, nil
	})
	for _, q := range quotes {
		p, ok := q.FirstChild().(*ast.Paragraph)
		if !ok || p.Lines().Len() == 0 {
			continue
		}
		first := p.Lines().At(0)
		sm := calloutRE.FindStringSubmatch(strings.TrimSpace(string(first.Value(src))))
		if sm == nil {
			continue
		}
		kind := sm[1]
		title, known := CalloutKinds[kind]
		if !known {
			t.res.BadCallouts = append(t.res.BadCallouts, kind)
			continue
		}
		attrs := map[string]string{}
		for _, a := range calloutAttrRE.FindAllStringSubmatch(sm[2], -1) {
			attrs[a[1]] = a[2]
		}
		if kind == "since" && attrs["v"] != "" {
			title = "Since " + attrs["v"]
		}
		// Remove the inline nodes of the first line, and the line break after them.
		for c := p.FirstChild(); c != nil; {
			next := c.NextSibling()
			tx, isText := c.(*ast.Text)
			if isText && tx.Segment.Start >= first.Stop {
				break
			}
			if isText && tx.Segment.Stop > first.Stop {
				break
			}
			p.RemoveChild(p, c)
			c = next
		}
		if p.ChildCount() == 0 {
			q.RemoveChild(q, p)
		}
		q.SetAttributeString("class", []byte("callout callout-"+kind))
		head := ast.NewParagraph()
		head.SetAttributeString("class", []byte("callout-title"))
		head.AppendChild(head, ast.NewString([]byte(title)))
		q.InsertBefore(q, q.FirstChild(), head)
	}
}
