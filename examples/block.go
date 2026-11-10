// SPDX-License-Identifier: Apache-2.0

// Package examples reads the fenced code blocks of the pages, runs them and checks their output (spec 07).
package examples

import (
	"slices"
	"strings"
)

// Block is one fenced code block of a Markdown file.
type Block struct {
	// Lang is the language id after the aliases, such as python or cpp. It is "output" or "input" for those blocks.
	Lang  string
	Info  string
	Attrs map[string]string
	// Words are the attributes without a value, such as run or fails.
	Words []string
	Code  string
	// Line is the line of the opening fence, from 1.
	Line int
	// Start and End are the byte offsets of the code in the file. Start is the first byte after the opening fence line and End is the first byte of the closing fence line.
	Start, End int
	// Closed is false when the file ends before the closing fence.
	Closed bool
	// After is the byte offset just after the closing fence line.
	After int
}

// Has tells if the block has a word attribute.
func (b *Block) Has(w string) bool {
	return slices.Contains(b.Words, w)
}

// Deps is the value of deps=: the path of a lock file, relative to the node directory (spec 07.7). It is "" when the block has no deps=.
func (b *Block) Deps() string { return b.Attrs["deps"] }

// Mode is the run mode of the block: run, fails, wrong, test or norun. It is "" when the block has none, and "many" when it has more than one. A test block with fails is a test block that must fail (spec 07.11), so its mode is test.
func (b *Block) Mode() string {
	mode := ""
	for _, m := range []string{"run", "fails", "wrong", "test", "norun"} {
		if b.Has(m) && (m != "fails" || !b.Has("test")) {
			if mode != "" {
				return "many"
			}
			mode = m
		}
	}
	return mode
}

// TestFails tells if the block is a test block with fails, which shows a test that fails on purpose (spec 07.11).
func (b *Block) TestFails() bool { return b.Mode() == "test" && b.Has("fails") }

// Runs tells if prog run runs the block: its mode is run, fails, wrong or test.
func (b *Block) Runs() bool {
	switch b.Mode() {
	case "run", "fails", "wrong", "test":
		return true
	}
	return false
}

// aliases maps the other names of a language in an info string to its id.
var aliases = map[string]string{
	"py": "python", "python3": "python",
	"js": "javascript", "node": "javascript", "mjs": "javascript",
	"ts":     "typescript",
	"golang": "go",
	"rs":     "rust",
	"c++":    "cpp", "cxx": "cpp",
	"cs": "csharp", "c#": "csharp",
	"kt": "kotlin",
	"sh": "bash", "shell": "bash",
	"ps1": "powershell", "pwsh": "powershell",
	"sqlite": "sql",
}

// LangID gives the language id of a name in an info string.
func LangID(name string) string {
	name = strings.ToLower(name)
	if id, ok := aliases[name]; ok {
		return id
	}
	return name
}

// Parse finds the fenced code blocks of a Markdown file, as CommonMark defines them: a fence of three or more backticks or tildes, indented by at most three spaces, closed by a fence of the same character that is at least as long. line0 is the line of the file where src starts, from 1.
func Parse(src string, line0 int) []Block {
	var out []Block
	pos, line := 0, line0
	var cur *Block
	var fch byte
	var flen, indent int
	for pos < len(src) {
		end := strings.IndexByte(src[pos:], '\n')
		next := len(src)
		if end >= 0 {
			next = pos + end + 1
		}
		text := strings.TrimRight(src[pos:next], "\r\n")
		if cur == nil {
			if ch, n, ind, info, ok := openFence(text); ok {
				cur = &Block{Info: info, Line: line, Start: next}
				fch, flen, indent = ch, n, ind
				fillInfo(cur)
			}
		} else if closeFence(text, fch, flen) {
			cur.End, cur.After, cur.Closed = pos, next, true
			cur.Code = dedent(src[cur.Start:cur.End], indent)
			out = append(out, *cur)
			cur = nil
		}
		pos = next
		line++
	}
	if cur != nil {
		cur.End, cur.After = len(src), len(src)
		cur.Code = dedent(src[cur.Start:], indent)
		out = append(out, *cur)
	}
	return out
}

func openFence(s string) (ch byte, n, indent int, info string, ok bool) {
	for indent < len(s) && indent < 4 && s[indent] == ' ' {
		indent++
	}
	if indent > 3 || indent >= len(s) || (s[indent] != '`' && s[indent] != '~') {
		return 0, 0, 0, "", false
	}
	ch = s[indent]
	i := indent
	for i < len(s) && s[i] == ch {
		i++
	}
	n = i - indent
	if n < 3 {
		return 0, 0, 0, "", false
	}
	info = strings.TrimSpace(s[i:])
	if ch == '`' && strings.Contains(info, "`") {
		return 0, 0, 0, "", false
	}
	return ch, n, indent, info, true
}

func closeFence(s string, ch byte, n int) bool {
	t := strings.TrimLeft(s, " ")
	if len(s)-len(t) > 3 {
		return false
	}
	i := 0
	for i < len(t) && t[i] == ch {
		i++
	}
	return i >= n && strings.TrimSpace(t[i:]) == ""
}

// dedent removes up to indent spaces from each line, as CommonMark does for an indented fence.
func dedent(s string, indent int) string {
	if indent == 0 {
		return s
	}
	lines := strings.SplitAfter(s, "\n")
	for i, l := range lines {
		k := 0
		for k < indent && k < len(l) && l[k] == ' ' {
			k++
		}
		lines[i] = l[k:]
	}
	return strings.Join(lines, "")
}

// fillInfo splits the info string into the language, the words and the key=value attributes. A value can be in double quotes.
func fillInfo(b *Block) {
	b.Attrs = map[string]string{}
	f := splitInfo(b.Info)
	if len(f) == 0 {
		return
	}
	b.Lang = LangID(f[0])
	for _, w := range f[1:] {
		if k, v, ok := strings.Cut(w, "="); ok {
			b.Attrs[k] = strings.Trim(v, `"`)
		} else {
			b.Words = append(b.Words, w)
		}
	}
}

func splitInfo(s string) []string {
	var out []string
	var cur strings.Builder
	quoted := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			quoted = !quoted
			cur.WriteByte(c)
		case (c == ' ' || c == '\t') && !quoted:
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(c)
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}
