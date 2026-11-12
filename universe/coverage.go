// SPDX-License-Identifier: Apache-2.0

package universe

import (
	"fmt"
	"slices"
	"strings"

	"github.com/tamnd/programming-agent/node"
)

// Levels are the coverage levels of spec 18.2 that prog coverage measures now. C0 is an item with a node or a card. C1 is a card that passes the R rules. C2 is a node at W2. The G ladder is not built yet, so C2 does not check G1, and the levels C3 to C6 are not measured yet.
var Levels = []string{"C0", "C1", "C2"}

// Level holds the coverage levels of the nodes and cards of one scope.
type Level struct {
	Scope string `json:"scope"`
	// Cards and Nodes are the counts of cards and nodes of the scope.
	Cards int `json:"cards"`
	Nodes int `json:"nodes"`
	// CardsBelow is the number of cards that do not pass the R rules, so they are at C0.
	CardsBelow int `json:"cards_below_c1"`
	// ByLevel counts the items at each level. An item counts at its highest level only.
	ByLevel map[string]int `json:"by_level"`
	// Failing lists the first cards of the scope that do not pass the R rules, as file:line.
	Failing []string `json:"failing,omitempty"`
}

// AtLeast is the number of items of the scope at the level or above.
func (l Level) AtLeast(level string) int {
	i := slices.Index(Levels, level)
	n := 0
	for _, v := range Levels[max(i, 0):] {
		n += l.ByLevel[v]
	}
	return n
}

// maxFailing is the number of failing cards that each scope lists.
const maxFailing = 5

// Coverage gives the coverage level of each card and node of c, by scope. bad tells if a card has a block finding of the R rules. nodeW gives the level of the writing ladder of a node. The last Level has the scope "all".
func Coverage(c *node.Corpus, bad func(*node.Card) bool, nodeW func(*node.Node) int) []Level {
	by := map[string]*Level{}
	get := func(s string) *Level {
		if by[s] == nil {
			by[s] = &Level{Scope: s, ByLevel: map[string]int{}}
		}
		return by[s]
	}
	for _, k := range c.Cards {
		s := ScopeOfKind(k.Kind)
		if s == "" {
			continue
		}
		l := get(s)
		l.Cards++
		if bad(k) {
			l.CardsBelow++
			l.ByLevel["C0"]++
			if len(l.Failing) < maxFailing {
				l.Failing = append(l.Failing, fmt.Sprintf("%s:%d", k.File, k.Line))
			}
			continue
		}
		l.ByLevel["C1"]++
	}
	for _, n := range c.Nodes {
		s := ScopeOfKind(n.Kind)
		if s == "" {
			continue
		}
		l := get(s)
		l.Nodes++
		if nodeW(n) >= 2 {
			l.ByLevel["C2"]++
		} else {
			l.ByLevel["C0"]++
		}
	}
	all := Level{Scope: "all", ByLevel: map[string]int{}}
	var out []Level
	for _, s := range Scopes {
		l := by[s]
		if l == nil {
			l = &Level{Scope: s, ByLevel: map[string]int{}}
		}
		all.Cards += l.Cards
		all.Nodes += l.Nodes
		all.CardsBelow += l.CardsBelow
		for k, v := range l.ByLevel {
			all.ByLevel[k] += v
		}
		out = append(out, *l)
	}
	return append(out, all)
}

// CardGate tells if every card of each scope is at C1 or above, the exit criterion of M1 (spec 21). It returns the scopes that fail.
func CardGate(levels []Level) error {
	var bad []string
	for _, l := range levels {
		if l.Scope == "all" {
			continue
		}
		if l.CardsBelow > 0 {
			bad = append(bad, fmt.Sprintf("%s has %d of %d cards below C1 (%s)", l.Scope, l.CardsBelow, l.Cards, strings.Join(l.Failing, ", ")))
		}
	}
	if len(bad) > 0 {
		return fmt.Errorf("not every card scope is at C1: %s", strings.Join(bad, "; "))
	}
	return nil
}
