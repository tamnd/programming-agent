// SPDX-License-Identifier: Apache-2.0

package examples

import (
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/tamnd/programming-agent/lang"
)

// Case is one program to run: a block, or the blocks of a group, with its output blocks.
type Case struct {
	File string
	// Block is the block that owns the run: the block itself, or the last runnable block of a group.
	Block *Block
	// Parts are the blocks of the program in the order of the page. A block without file= adds its code to the main file of the language.
	Parts   []*Block
	Outputs []*Block
	// Stdin is the text of the input block that stdin= names.
	Stdin string
	// Dir is the directory of the node. The path of deps= is relative to it.
	Dir string
}

// ID is the id of the case: the id= of the block, or ex<n> from the order of the blocks.
func (c *Case) ID() string { return c.Block.Attrs["id"] }

// Mode is run, fails, wrong or test.
func (c *Case) Mode() string { return c.Block.Mode() }

// Label is the mode of the case as the result shows it. A test block with fails is "test fails".
func (c *Case) Label() string {
	if c.Block.TestFails() {
		return "test fails"
	}
	return c.Mode()
}

// Deps is the lock file that deps= names, as a path relative to the node directory (spec 07.7). It is the deps= of the block that owns the run, or of another block of the group. It is "" when the program uses the standard library only.
func (c *Case) Deps() string {
	if d := c.Block.Deps(); d != "" {
		return d
	}
	for _, b := range c.Parts {
		if d := b.Deps(); d != "" {
			return d
		}
	}
	return ""
}

// LockPath is the path of the lock file of deps= on this machine, or "".
func (c *Case) LockPath() string {
	d := c.Deps()
	if d == "" {
		return ""
	}
	return filepath.Join(c.Dir, filepath.FromSlash(d))
}

// Files gives the files of the program. The parts without file= go to the main file of the language, in order.
func (c *Case) Files(p lang.Program) map[string]string {
	files := map[string]string{}
	for _, b := range c.Parts {
		name := b.Attrs["file"]
		if name == "" {
			name = p.Main
		} else {
			// A file of a multi-file example goes next to the main file. A Cargo project has the main file in src.
			name = path.Join(path.Dir(p.Main), name)
		}
		if files[name] != "" && !strings.HasSuffix(files[name], "\n") {
			files[name] += "\n"
		}
		files[name] += b.Code
	}
	return files
}

// TestFile is the name of the test file of a test block in the work directory: the file= of the block next to the main file, or the main file.
func (c *Case) TestFile(p lang.Program) string {
	if f := c.Block.Attrs["file"]; f != "" {
		return path.Join(path.Dir(p.Main), f)
	}
	return p.Main
}

// Issue is a broken rule of spec 07.2.
type Issue struct {
	File string
	Line int
	Rule string
	Msg  string
}

func (i Issue) String() string { return fmt.Sprintf("%s:%d: %s %s", i.File, i.Line, i.Rule, i.Msg) }

// NorunReasons are the codes of spec 07.8.
var NorunReasons = []string{"network", "gui", "platform", "hardware", "ub", "vulnerable", "fragment", "pseudo", "long", "interactive"}

// Expects are the kinds of failure signature of spec 07.2.
var Expects = []string{"compile", "runtime", "exit", "exception", "error", "panic", "timeout", "sanitizer"}

// TestFailure is the failure signature of a test block with fails: at least one test fails (spec 07.11).
const TestFailure = "test-failure"

// Collect gives the cases of one Markdown file, and the issues of the rules X01 and X02. The blocks of the languages that prog run can run need a mode; the other blocks, such as text or json, do not.
func Collect(file, src string, line0 int) ([]Case, []Issue) {
	blocks := Parse(src, line0)
	var issues []Issue
	bad := func(b *Block, rule, format string, a ...any) {
		issues = append(issues, Issue{file, b.Line, rule, fmt.Sprintf(format, a...)})
	}
	// Give the code blocks their ids, and check the modes.
	ids := map[string]*Block{}
	inputs := map[string]*Block{}
	n := 0
	for i := range blocks {
		b := &blocks[i]
		if !b.Closed {
			bad(b, "X01", "the block has no closing fence")
		}
		switch b.Lang {
		case "":
			bad(b, "X01", "the block has no language; use text for a block that is not code")
			continue
		case "output":
			continue
		case "input":
			if id := b.Attrs["id"]; id != "" {
				inputs[id] = b
			}
			continue
		}
		if _, ok := lang.Programs[b.Lang]; !ok {
			continue
		}
		n++
		if b.Attrs["id"] == "" {
			b.Attrs["id"] = "ex" + strconv.Itoa(n)
		} else if ids[b.Attrs["id"]] != nil {
			bad(b, "X01", "the id %q is used twice on the page", b.Attrs["id"])
		}
		ids[b.Attrs["id"]] = b
		switch m := b.Mode(); m {
		case "":
			bad(b, "X01", "the %s block needs one of run, fails, wrong, test or norun", b.Lang)
		case "many":
			bad(b, "X01", "the %s block has more than one of run, fails, wrong, test and norun", b.Lang)
		case "norun":
			if r := b.Attrs["reason"]; !slices.Contains(NorunReasons, r) {
				bad(b, "X01", "norun needs reason= with one of %s, not %q", strings.Join(NorunReasons, ", "), r)
			}
		case "fails":
			sig := b.Attrs["expect"]
			kind, _, _ := strings.Cut(sig, ":")
			switch {
			case sig == TestFailure:
				bad(b, "X04", "expect=%s is for a test block; write test fails", TestFailure)
			case !slices.Contains(Expects, kind):
				bad(b, "X04", "fails needs expect= with one of %s, not %q", strings.Join(Expects, ", "), sig)
			}
		case "test":
			sig, ok := b.Attrs["expect"]
			switch {
			case b.Has("fails") && sig != TestFailure:
				bad(b, "X04", "test fails needs expect=%s, not %q", TestFailure, sig)
			case !b.Has("fails") && ok:
				bad(b, "X04", "a test block with expect= needs fails; write test fails expect=%s", TestFailure)
			}
		}
		if d, ok := b.Attrs["deps"]; ok && !validDeps(d) {
			bad(b, "X01", "deps= is the path of a lock file in the node directory, such as deps/rust/Cargo.lock, not %q", d)
		}
		if t := b.Attrs["timeout"]; t != "" {
			if s, err := strconv.Atoi(t); err != nil || s < 1 || s > 30 {
				bad(b, "X01", "timeout is a number of seconds from 1 to 30, not %q", t)
			}
		}
	}
	// Pair the output blocks with their blocks: by for=, or the code block before it.
	outs := map[*Block][]*Block{}
	var last *Block
	for i := range blocks {
		b := &blocks[i]
		switch b.Lang {
		case "output":
			owner := last
			if f := b.Attrs["for"]; f != "" {
				owner = ids[f]
				if owner == nil {
					bad(b, "X02", "no block on the page has the id %q", f)
					continue
				}
			}
			if owner == nil {
				bad(b, "X02", "the output block has no code block before it")
				continue
			}
			outs[owner] = append(outs[owner], b)
		case "input":
		default:
			if _, ok := lang.Programs[b.Lang]; ok {
				last = b
			} else {
				last = nil
			}
		}
	}
	// Join the groups and make the cases.
	var cases []Case
	seen := map[string]bool{}
	for i := range blocks {
		b := &blocks[i]
		if _, ok := lang.Programs[b.Lang]; !ok {
			continue
		}
		g := b.Attrs["group"]
		var parts []*Block
		if g == "" {
			parts = []*Block{b}
		} else {
			key := b.Lang + "\x00" + g
			if seen[key] {
				continue
			}
			seen[key] = true
			for j := range blocks {
				if blocks[j].Lang == b.Lang && blocks[j].Attrs["group"] == g {
					parts = append(parts, &blocks[j])
				}
			}
		}
		var owner *Block
		var outputs []*Block
		for _, p := range parts {
			if p.Runs() {
				owner = p
			}
			outputs = append(outputs, outs[p]...)
		}
		if owner == nil {
			continue
		}
		c := Case{File: file, Block: owner, Parts: parts, Outputs: outputs}
		if id := owner.Attrs["stdin"]; id != "" {
			in := inputs[id]
			if in == nil {
				bad(owner, "X01", "no input block has the id %q", id)
				continue
			}
			c.Stdin = in.Code
		}
		if (owner.Mode() == "run" || owner.Mode() == "wrong" || (owner.Mode() == "test" && !owner.TestFails())) && len(outputs) == 0 && owner.Attrs["match"] != "any" {
			bad(owner, "X02", "the %s block %s needs an output block; prog run -update writes it", owner.Mode(), owner.Attrs["id"])
		}
		cases = append(cases, c)
	}
	slices.SortStableFunc(issues, func(a, b Issue) int { return a.Line - b.Line })
	return cases, issues
}

// validDeps tells if the value of deps= is a path that stays in the node directory. The path uses slashes on every system.
func validDeps(d string) bool {
	return d != "" && !strings.Contains(d, "\\") && filepath.IsLocal(filepath.FromSlash(d)) && path.Clean(d) == d
}
