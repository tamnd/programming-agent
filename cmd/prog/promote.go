// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tamnd/programming-agent/check"
	"github.com/tamnd/programming-agent/node"
)

// cmdPromote makes a node from a card (spec 04.1). The node keeps the tag and the facts of the card, and the row leaves its register, so the change is one commit.
func cmdPromote(args []string) error {
	f := newFlags("promote", "[-area <path>] <slug or tag>")
	root := f.String("root", ".", "the root of the content repository")
	area := f.String("area", "", "the area path of the node; the default is the area whose URL is the collection of the card, such as languages/rust/errors")
	level := f.String("level", "core", "the level: intro, core, advanced or research")
	if err := parse(f, args); err != nil {
		return err
	}
	if f.NArg() != 1 {
		return usageError{"give the slug or the tag of one card, for example rust/errors/E0382"}
	}
	c, err := load(*root)
	if err != nil {
		return err
	}
	k := c.CardBySlug[strings.Trim(f.Arg(0), "/")]
	if k == nil {
		for _, x := range c.Cards {
			if x.Tag == f.Arg(0) {
				k = x
			}
		}
	}
	if k == nil {
		return fmt.Errorf("no card has the slug or the tag %q", f.Arg(0))
	}
	if n := c.BySlug[k.Slug()]; n != nil {
		return fmt.Errorf("the node %s has the slug %s already", n.Rel, k.Slug())
	}
	if *area == "" {
		for _, a := range c.AreaByID {
			if a.URL() == k.Collection() {
				*area = a.Path()
			}
		}
	}
	if *area == "" {
		return usageError{fmt.Sprintf("no area has the URL %s; give the area with -area", k.Collection())}
	}
	if _, ok := c.AreaByID[*area]; !ok {
		return usageError{fmt.Sprintf("the area %q is not in topics/areas.yaml", *area)}
	}
	dir := filepath.Join(*root, "nodes", filepath.FromSlash(*area), k.Local())
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s exists", dir)
	}
	md := promoted(k, *area, *level)
	shard := filepath.Join(*root, filepath.FromSlash(k.File))
	rest, err := dropLine(shard, k.Line)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	p := filepath.Join(dir, "node.md")
	if err := os.WriteFile(p, []byte(md), 0o644); err != nil {
		return err
	}
	if len(bytes.TrimSpace(rest)) == 0 {
		err = os.Remove(shard)
	} else {
		err = os.WriteFile(shard, rest, 0o644)
	}
	if err != nil {
		return err
	}
	fmt.Println(p)
	return nil
}

// promoted is the node.md of a promoted card: the facts in the front matter and the empty sections of its kind.
func promoted(k *node.Card, area, level string) string {
	var b strings.Builder
	str := func(key, v string) {
		if v != "" {
			fmt.Fprintf(&b, "%s: %s\n", key, strconv.Quote(v))
		}
	}
	list := func(key string, vs []string) {
		if len(vs) == 0 {
			return
		}
		q := make([]string, len(vs))
		for i, v := range vs {
			q[i] = strconv.Quote(v)
		}
		fmt.Fprintf(&b, "%s: [%s]\n", key, strings.Join(q, ", "))
	}
	tag := k.Tag
	if tag == "" {
		tag = "NEW-" + strings.ToLower(strings.NewReplacer("/", "-", ".", "-", "_", "-", "::", "-").Replace(k.Slug()))
	}
	fmt.Fprintf(&b, "---\ntag: %s\nkind: %s\n", tag, k.Kind)
	if k.Owner != "" && node.IsOwned(k.Kind) {
		fmt.Fprintf(&b, "owner: %s\n", k.Owner)
	}
	fmt.Fprintf(&b, "slug: %s\n", k.Local())
	str("name", k.Heading())
	fmt.Fprintf(&b, "area: %s\nlevel: %s\n", area, level)
	sum := k.Summary
	if sum == "" {
		sum = "Write one sentence that says what this is."
	}
	fmt.Fprintf(&b, "summary: >-\n  %s\n", sum)
	switch k.Kind {
	case "diagnostic":
		str("code", k.Code)
		str("severity", k.Severity)
		list("message_templates", k.MessageTemplates)
	case "proposal":
		str("proposal_status", k.Status)
		str("proposal_type", k.Type)
		str("created", k.Created)
		str("release", k.Release)
		list("authors", k.Authors)
	case "api":
		str("path", k.Path)
		str("entity", k.Entity)
		str("signature", k.Signature)
		str("module", k.Module)
	}
	str("since", k.Since)
	list("concepts", k.Concepts)
	list("features", k.Features)
	var srcs []string
	for _, u := range []string{k.Source.URL, k.URL, k.Doc.URL} {
		if u != "" && !strings.Contains(strings.Join(srcs, "\n"), u) {
			srcs = append(srcs, u)
		}
	}
	if len(srcs) > 0 {
		b.WriteString("sources:\n")
		for _, u := range srcs {
			fmt.Fprintf(&b, "  - url: %s\n    title: %s\n    accessed: %s\n", u, strconv.Quote(k.Heading()), time.Now().UTC().Format(time.DateOnly))
		}
	}
	b.WriteString("---\n\n")
	for _, h := range check.Order[k.Kind] {
		fmt.Fprintf(&b, "## %s\n\n", h)
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// dropLine returns the file without its line n, counted from 1.
func dropLine(path string, n int) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := bytes.SplitAfter(data, []byte("\n"))
	if n < 1 || n > len(lines) {
		return nil, fmt.Errorf("%s has no line %d", path, n)
	}
	return bytes.Join(append(lines[:n-1:n-1], lines[n:]...), nil), nil
}
