// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/tamnd/programming-agent/compat"
	"github.com/tamnd/programming-agent/node"
)

func cmdCompat(args []string) error {
	if len(args) == 0 || args[0] != "probe" {
		return usageError{"use prog compat probe [slug ...]"}
	}
	f := newFlags("compat probe", "[flags] [slug ...]")
	root := f.String("root", ".", "the root of the content repository")
	engine := f.String("engine", "docker", "the container command")
	all := f.Bool("all", false, "run each probe again, also when its result is current")
	timeout := f.Duration("timeout", 3*time.Minute, "the limit of one run")
	if err := parse(f, args[1:]); err != nil {
		return err
	}
	c, err := load(*root)
	if err != nil {
		return err
	}
	only := map[string]bool{}
	for _, s := range f.Args() {
		if c.BySlug[s] == nil {
			return usageError{fmt.Sprintf("no node has the slug %q", s)}
		}
		only[s] = true
	}
	tasks, skipped := probeTasks(c, only)
	p := compat.NewProber(*engine)
	p.Timeout = *timeout
	ran, failed := 0, 0
	for _, t := range tasks {
		key := compat.Key(t.Feature, t.Impl.ID, t.Probe, t.Added)
		hash, err := compat.ProbeHash(t.Dir, t.Probe, t.Expect)
		if err != nil {
			return err
		}
		if old, ok := c.Probes[key]; ok && !*all && old.Verified(hash) {
			continue
		}
		r, err := p.Probe(context.Background(), t)
		if err != nil {
			fmt.Printf("%s %s %s: %v\n", t.Feature, t.Impl.ID, path.Base(t.Probe), err)
			failed++
			continue
		}
		c.Probes[key] = r
		ran++
		fmt.Println(probeLine(r))
		if !r.Verified(r.Hash) {
			failed++
		}
	}
	for _, s := range skipped {
		fmt.Println(s)
	}
	if err := c.Probes.Write(c.Root); err != nil {
		return err
	}
	fmt.Printf("%d probes ran, %d did not prove their entry, %d skipped; the results are in %s\n", ran, failed, len(skipped), compat.ResultsFile)
	if failed > 0 {
		return errFindings
	}
	return nil
}

// probeTasks lists a task for each support statement with a version and a probe, once for each feature, implementation, probe and version. A statement whose implementation has no image is skipped.
func probeTasks(c *node.Corpus, only map[string]bool) ([]compat.Task, []string) {
	seen := map[string]bool{}
	var tasks []compat.Task
	var skipped []string
	for _, n := range c.Nodes {
		if len(n.CompatEntries) == 0 || (len(only) > 0 && !only[n.Slug]) {
			continue
		}
		lang := n.Owner
		for _, e := range n.CompatEntries {
			for _, col := range e.Columns {
				for _, s := range e.Support[col] {
					if s.Probe == "" || !s.IsVersion() {
						continue
					}
					probe := path.Clean(s.Probe)
					key := compat.Key(n.Slug, col, probe, s.Added)
					if seen[key] {
						continue
					}
					seen[key] = true
					im, ok := c.Impls.Impl(lang, col)
					if !ok || im.Image == "" {
						skipped = append(skipped, fmt.Sprintf("%s %s %s: skipped, %s has no image in registers/implementations.yaml", n.Slug, col, path.Base(probe), col))
						continue
					}
					tasks = append(tasks, compat.Task{
						Feature: n.Slug,
						Dir:     filepath.Join(c.Root, filepath.FromSlash(n.Rel)),
						Impl:    im,
						Probe:   probe,
						Added:   s.Added,
						Expect:  n.Probes[path.Base(probe)],
						Edition: compat.Edition(s.Edition, "2015"),
					})
				}
			}
		}
	}
	sort.SliceStable(tasks, func(i, j int) bool { return tasks[i].Feature < tasks[j].Feature })
	sort.Strings(skipped)
	return tasks, skipped
}

func probeLine(r compat.Result) string {
	s := fmt.Sprintf("%s %s %s: %s %s", r.Feature, r.Impl, path.Base(r.Probe), r.Pass.Version, okWord(r.Pass.OK, "passes", "does not pass"))
	if r.Pass.Error != "" {
		s += " (" + r.Pass.Error + ")"
	}
	if r.Fail != nil {
		s += fmt.Sprintf(", %s %s", r.Fail.Version, okWord(r.Fail.OK, "fails with the signature", "does not fail with the signature"))
		if r.Fail.Error != "" {
			s += " (" + r.Fail.Error + ")"
		}
	}
	return s
}

func okWord(ok bool, yes, no string) string {
	if ok {
		return yes
	}
	return no
}
