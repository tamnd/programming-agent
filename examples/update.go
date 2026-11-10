// SPDX-License-Identifier: Apache-2.0

package examples

import (
	"slices"
	"strings"
)

// edit replaces src[Start:End] with Text.
type edit struct {
	Start, End int
	Text       string
}

// Update writes the output of the runs into the output blocks of a file (prog run -update). It replaces the text of an output block that differs from the run, and adds an output block after a run, wrong or test block that has none. The output of a test block is the report of the test runner with its normalizers (spec 07.11). It does not change an output block with match=regex or match=any, because a person wrote the pattern. The results must be of the cases of src. It returns the new text and the number of blocks that changed, and it sets Updated on each result whose blocks now match the run.
func Update(src string, results []Result) (string, int) {
	var edits []edit
	for i := range results {
		r := &results[i]
		c := r.Case
		if c == nil || r.Status == "skip" || r.Got.TimedOut {
			continue
		}
		// The report of a test block that does not pass, or a test block with fails that has no failed test, is not an output to keep.
		if c.Mode() == "test" && checkTests(c, r.Got) != "" {
			continue
		}
		if len(c.Outputs) == 0 {
			if !writesOutput(r) || c.Block.Attrs["match"] == "any" {
				continue
			}
			stream := reportStream(c, &Block{})
			text := Clean(Stream(r.Got, stream))
			info := "output"
			if stream != "" {
				// A test runner such as unittest writes its report to stderr. The output block says so.
				info += " stream=" + stream
			}
			fence := fenceFor(text)
			edits = append(edits, edit{c.Block.After, c.Block.After, "\n" + fence + info + "\n" + text + fence + "\n"})
			r.Updated = true
			continue
		}
		fixed := true
		for _, o := range c.Outputs {
			if m := o.Attrs["match"]; m == "regex" || m == "any" {
				fixed = false
				continue
			}
			if o.Attrs["version"] != "" && o.Attrs["version"] != versionOf(r.Tool) && !sameVersion(versionOf(r.Tool), o.Attrs["version"]) {
				continue
			}
			text := Clean(Stream(r.Got, reportStream(c, o)))
			norm := strings.Split(o.Attrs["normalize"], ",")
			if Normalize(o.Code, norm) == Normalize(text, norm) {
				continue
			}
			edits = append(edits, edit{o.Start, o.End, text})
		}
		r.Updated = fixed
	}
	slices.SortFunc(edits, func(a, b edit) int { return b.Start - a.Start })
	for _, e := range edits {
		src = src[:e.Start] + e.Text + src[e.End:]
	}
	return src, len(edits)
}

// writesOutput tells if -update adds an output block for the result of a case that has none. A run, wrong or test block needs a run that exits 0. A test block needs a run that checkTests accepts, so a test block with fails needs a test that failed.
func writesOutput(r *Result) bool {
	c := r.Case
	switch {
	case c.Mode() == "test":
		return checkTests(c, r.Got) == ""
	case c.Mode() == "run", c.Mode() == "wrong":
		return r.Got.Exit == 0
	}
	return false
}

// fenceFor gives a fence that is longer than each run of backticks in the text.
func fenceFor(text string) string {
	n, run := 3, 0
	for _, c := range text {
		if c == '`' {
			run++
			if run >= n {
				n = run + 1
			}
		} else {
			run = 0
		}
	}
	return strings.Repeat("`", n)
}
