// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/tamnd/programming-agent/check"
	"github.com/tamnd/programming-agent/node"
	"github.com/tamnd/programming-agent/status"
	"github.com/tamnd/programming-agent/universe"
)

// cardRules are the rules of a card that C1 needs (spec 18.2): the R rules and the tag rule T02.
var cardRules = map[string]bool{"R01": true, "R02": true, "R04": true, "R05": true, "T02": true}

func cmdCoverage(args []string) error {
	f := newFlags("coverage", "[flags]")
	root := f.String("root", ".", "the root of the content repository")
	asJSON := f.Bool("json", false, "print the levels as JSON")
	gate := f.Bool("gate", false, "fail when a card of any scope is below C1, the exit criterion of M1")
	if err := parse(f, args); err != nil {
		return err
	}
	c, err := load(*root)
	if err != nil {
		return err
	}
	bad := map[string]bool{}
	for _, fd := range check.Run(c, check.Options{}) {
		if fd.Severity == check.Block && cardRules[fd.Rule] && fd.Line > 0 {
			bad[fmt.Sprintf("%s:%d", fd.Path, fd.Line)] = true
		}
	}
	env := status.Env{Corpus: c}
	levels := universe.Coverage(c,
		func(k *node.Card) bool { return bad[fmt.Sprintf("%s:%d", k.File, k.Line)] },
		func(n *node.Node) int { return status.Of(n, env).W.Level })
	if *asJSON {
		if err := printJSON(struct {
			Levels     []universe.Level `json:"levels"`
			Unreadable int              `json:"unreadable_card_rows"`
		}{levels, len(c.CardErrors)}); err != nil {
			return err
		}
	} else {
		printLevels(levels, len(c.CardErrors))
	}
	if *gate {
		var errs []error
		if n := len(c.CardErrors); n > 0 {
			errs = append(errs, fmt.Errorf("%d card rows cannot be read; prog check lists them as R01", n))
		}
		errs = append(errs, universe.CardGate(levels))
		return errors.Join(errs...)
	}
	return nil
}

// printLevels prints one line for each scope: the cards, the nodes, the count at each level, and the share of the cards at C1.
func printLevels(levels []universe.Level, unreadable int) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "scope\tcards\tnodes\t%s\tcards at C1\t\n", strings.Join(universe.Levels, "\t"))
	for _, l := range levels {
		var cols []string
		for _, v := range universe.Levels {
			cols = append(cols, fmt.Sprint(l.ByLevel[v]))
		}
		share := "-"
		if l.Cards > 0 {
			share = fmt.Sprintf("%.2f%%", 100*float64(l.Cards-l.CardsBelow)/float64(l.Cards))
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\t\n", l.Scope, l.Cards, l.Nodes, strings.Join(cols, "\t"), share)
	}
	_ = w.Flush()
	if unreadable > 0 {
		fmt.Printf("\n%d card rows cannot be read; prog check lists them as R01.\n", unreadable)
	}
	fmt.Println("\nC0 is a node or card, C1 is a card that passes the R rules, and C2 is a node at W2 (spec 18.2). The G ladder is not built yet, so C2 does not check G1, and C3 to C6 are not measured yet.")
}
