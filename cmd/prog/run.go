// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/tamnd/programming-agent/examples"
	"github.com/tamnd/programming-agent/lang"
)

// cmdRun runs the code blocks of the pages and checks their output (spec 07). With -update, it writes the output of each run into the output blocks.
func cmdRun(args []string) error {
	f := newFlags("run", "[flags] [slug ...]")
	root := f.String("root", ".", "the root of the content repository")
	only := f.String("lang", "", "a comma list of languages to run; the default is all")
	update := f.Bool("update", false, "write the output of each run into the output blocks, and add the output blocks that are missing")
	out := f.String("out", "build/runs/runs.json", "the file for the results")
	unsafe := f.Bool("unsafe-local", false, "run the blocks on this machine, not in the sandbox; CI sets CI=true instead")
	sandbox := f.String("sandbox", "auto", "the container engine of the sandbox: docker, podman, or none; auto uses docker when the image manifest is there")
	images := f.String("images", "manifests/images.yaml", "the manifest of the runner images, pinned by digest")
	keep := f.Bool("keep", false, "keep the temporary directories")
	jobs := f.Int("j", max(1, runtime.NumCPU()/2), "the number of blocks that run at the same time")
	strict := f.Bool("strict", false, "fail when a block is skipped because its toolchain is missing")
	if err := parse(f, args); err != nil {
		return err
	}
	c, err := load(*root)
	if err != nil {
		return err
	}
	nodes, err := selectNodes(c, f.Args())
	if err != nil {
		return err
	}
	want := map[string]bool{}
	for _, l := range strings.Split(*only, ",") {
		if l = strings.TrimSpace(l); l != "" {
			if _, ok := lang.Programs[l]; !ok {
				return usageError{fmt.Sprintf("prog run has no runner for the language %q", l)}
			}
			want[l] = true
		}
	}
	box, err := openSandbox(*sandbox, under(c.Root, *images), *unsafe)
	if err != nil {
		return err
	}
	if box == nil && !*unsafe && os.Getenv("CI") != "true" {
		return usageError{"the code blocks are programs from the repository; use the sandbox or CI, or give -unsafe-local"}
	}
	r := &lang.Runner{Keep: *keep, Box: box, Log: func(format string, a ...any) { fmt.Fprintf(os.Stderr, format+"\n", a...) }}
	if box != nil {
		abs, err := filepath.Abs(filepath.Join(c.Root, "build", "sandbox"))
		if err != nil {
			return err
		}
		r.TempDir = abs
		if err := os.MkdirAll(r.TempDir, 0o755); err != nil {
			return err
		}
	}
	diag := map[string]bool{}
	for _, n := range c.Nodes {
		if n.Kind == "diagnostic" && n.Code != "" {
			diag[n.Code] = true
		}
	}

	// Read the files: node.md and examples/*.md of each node.
	type page struct {
		path, src string
		cases     []examples.Case
	}
	var pages []*page
	var issues []examples.Issue
	var warnings []string
	for _, n := range nodes {
		files := []string{filepath.Join(n.Dir, "node.md")}
		for _, e := range n.Examples {
			files = append(files, filepath.Join(c.Root, filepath.FromSlash(e.Path)))
		}
		for _, file := range files {
			b, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(c.Root, file)
			rel = filepath.ToSlash(rel)
			cases, is := examples.Collect(rel, string(b), 1)
			issues = append(issues, is...)
			pg := &page{path: file, src: string(b)}
			for _, k := range cases {
				if len(want) > 0 && !want[k.Block.Lang] {
					continue
				}
				if code, ok := strings.CutPrefix(k.Block.Attrs["expect"], "error:"); ok && !diag[code] {
					warnings = append(warnings, fmt.Sprintf("%s:%d: X05 no diagnostic card has the code %s", rel, k.Block.Line, code))
				}
				pg.cases = append(pg.cases, k)
			}
			pages = append(pages, pg)
		}
	}

	// Run the cases.
	x := &examples.Runner{R: r}
	type job struct {
		pg *page
		i  int
	}
	var all []job
	for _, pg := range pages {
		for i := range pg.cases {
			all = append(all, job{pg, i})
		}
	}
	results := make([]examples.Result, len(all))
	var wg sync.WaitGroup
	sem := make(chan struct{}, max(1, *jobs))
	for i, j := range all {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = x.Run(&j.pg.cases[j.i])
		}()
	}
	wg.Wait()

	// Write the output blocks.
	if *update {
		byPage := map[*page][]examples.Result{}
		for i, j := range all {
			byPage[j.pg] = append(byPage[j.pg], results[i])
		}
		changed := 0
		fixed := map[*examples.Case]bool{}
		for _, pg := range pages {
			rs := byPage[pg]
			src, n := examples.Update(pg.src, rs)
			for _, res := range rs {
				if res.Updated {
					fixed[res.Case] = true
				}
			}
			if n == 0 {
				continue
			}
			if err := os.WriteFile(pg.path, []byte(src), 0o644); err != nil {
				return err
			}
			changed += n
			fmt.Printf("updated %d output blocks in %s\n", n, pg.path)
		}
		// The new output blocks answer X02 and the runs that differed now match, so read the files again.
		issues = issues[:0]
		for _, pg := range pages {
			b, err := os.ReadFile(pg.path)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(c.Root, pg.path)
			_, is := examples.Collect(filepath.ToSlash(rel), string(b), 1)
			issues = append(issues, is...)
		}
		for i := range results {
			if results[i].Status == "fail" && fixed[results[i].Case] && strings.HasPrefix(results[i].Why, "the output at line") {
				results[i].Status, results[i].Why = "pass", "updated"
			}
		}
		if changed == 0 {
			fmt.Println("no output block changed")
		}
	}

	// Report.
	var pass, fail, skip int
	skipped := map[string]int{}
	for _, res := range results {
		switch res.Status {
		case "pass":
			pass++
		case "fail":
			fail++
			fmt.Fprintf(os.Stderr, "FAIL %s:%d %s %s (%s)\n%s\n", res.File, res.Line, res.Lang, res.ID, res.Mode, indent(res.Why))
		case "skip":
			skip++
			skipped[res.Lang+": "+res.Why]++
		}
	}
	for _, is := range issues {
		fmt.Fprintln(os.Stderr, is.String())
	}
	for _, w := range warnings {
		fmt.Fprintln(os.Stderr, w)
	}
	var keys []string
	for k := range skipped {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stderr, "skip %d blocks of %s\n", skipped[k], k)
	}
	fmt.Printf("%d blocks pass, %d fail, %d skipped; %d issues\n", pass, fail, skip, len(issues))
	if *out != "" {
		p := under(*root, *out)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		b, err := json.MarshalIndent(results, "", " ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(p, append(b, '\n'), 0o644); err != nil {
			return err
		}
	}
	if fail > 0 || len(issues) > 0 || (*strict && skip > 0) {
		return errFindings
	}
	return nil
}
