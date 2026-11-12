// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/tamnd/programming-agent/universe"
)

// strictScopes are the scopes that may have no pending item at the end of M1 (spec 21).
var strictScopes = []string{"languages", "releases", "weaknesses"}

func cmdUniverse(args []string) error {
	f := newFlags("universe", "[flags]")
	root := f.String("root", ".", "the root of the content repository")
	asJSON := f.Bool("json", false, "print the counts as JSON")
	gate := f.Bool("gate", false, "fail when 5% or more of the items are pending, or when a languages, releases or weaknesses item is pending")
	max := f.Float64("max", 0.05, "with -gate, the largest share of pending items that fails")
	if err := parse(f, args); err != nil {
		return err
	}
	c, err := load(*root)
	if err != nil {
		return err
	}
	rows, errs := universe.Read(c.Root)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	counts := universe.Report(c, rows)
	if *asJSON {
		if err := printJSON(counts); err != nil {
			return err
		}
	} else {
		printCounts(counts)
	}
	if *gate {
		if err := universe.Gate(counts, *max, strictScopes); err != nil {
			return err
		}
	}
	return nil
}

// printCounts prints one line for each scope: the total, the pending share, and the count of each outcome.
func printCounts(counts []universe.Count) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "scope\titems\tpending\tshare\toutcomes\t")
	for _, n := range counts {
		var keys []string
		for k := range n.ByOutcome {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var parts []string
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s %d", k, n.ByOutcome[k]))
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%.2f%%\t%s\t\n", n.Scope, n.Total, n.Pending, 100*n.Share(), strings.Join(parts, ", "))
	}
	w.Flush()
}
