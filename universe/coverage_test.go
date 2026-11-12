// SPDX-License-Identifier: Apache-2.0

package universe

import (
	"testing"

	"github.com/tamnd/programming-agent/node"
)

func TestCoverage(t *testing.T) {
	broken := &node.Card{Kind: "api", File: "cards/api/go/st.jsonl", Line: 3}
	c := &node.Corpus{
		Nodes: []*node.Node{{Slug: "python", Kind: "language"}, {Slug: "rust", Kind: "language"}, {Slug: "x", Kind: "area"}},
		Cards: []*node.Card{{Kind: "language"}, {Kind: "api"}, broken, {Kind: "person"}},
	}
	w := map[string]int{"python": 3, "rust": 1}
	levels := Coverage(c, func(k *node.Card) bool { return k == broken }, func(n *node.Node) int { return w[n.Slug] })
	by := map[string]Level{}
	for _, l := range levels {
		by[l.Scope] = l
	}
	if l := by["languages"]; l.Cards != 1 || l.Nodes != 2 || l.ByLevel["C0"] != 1 || l.ByLevel["C1"] != 1 || l.ByLevel["C2"] != 1 {
		t.Errorf("languages: %+v", l)
	}
	if l := by["api"]; l.Cards != 2 || l.ByLevel["C0"] != 1 || l.ByLevel["C1"] != 1 || len(l.Failing) != 1 || l.Failing[0] != "cards/api/go/st.jsonl:3" {
		t.Errorf("api: %+v", l)
	}
	if l := by["people"]; l.ByLevel["C1"] != 1 || l.AtLeast("C1") != 1 {
		t.Errorf("people: %+v", l)
	}
	if l := levels[len(levels)-1]; l.Scope != "all" || l.Cards != 4 || l.Nodes != 2 || l.AtLeast("C0") != 6 || l.AtLeast("C2") != 1 {
		t.Errorf("all: %+v", l)
	}
	// The languages scope has a node below C2, but all its cards pass, so only the api scope fails the gate.
	err := CardGate(levels)
	if err == nil {
		t.Fatal("the gate passes with a failing api card")
	}
	if got := err.Error(); got != "not every card scope is at C1: api has 1 of 2 cards below C1 (cards/api/go/st.jsonl:3)" {
		t.Errorf("the gate error is %q", got)
	}
	if err := CardGate(Coverage(c, func(*node.Card) bool { return false }, func(*node.Node) int { return 0 })); err != nil {
		t.Errorf("the gate fails with no failing card: %v", err)
	}
}
