// SPDX-License-Identifier: Apache-2.0

package check

import (
	"strings"

	"github.com/tamnd/programming-agent/examples"
	"github.com/tamnd/programming-agent/node"
)

// norunAreas are the areas where the norun blocks can be 40% of the blocks of a node (rule X08). In the other areas the limit is 10%.
var norunAreas = []string{"embedded-and-hardware", "graphics-and-games-programming", "mobile-and-desktop", "distributed-and-cloud"}

// norunLimit is the largest share of norun blocks of the node, in percent.
func norunLimit(n *node.Node) int {
	for _, a := range append([]string{n.Area}, n.Areas...) {
		top, _, _ := strings.Cut(a, "/")
		for _, x := range norunAreas {
			if top == x {
				return 40
			}
		}
	}
	return 10
}

// norun checks rule X08 (spec 07.8): the norun blocks of node.md and of the example files are at most 10% of the code blocks of the node, or 40% in the areas of norunAreas. A norun block of a group that runs is a part of a program that runs, so it does not count as norun.
func (r *run) norun(n *node.Node) {
	if !r.wants("X08") {
		return
	}
	norun, total := examples.NorunShare(n.Body)
	for _, ex := range n.Examples {
		a, b := examples.NorunShare(ex.Body)
		norun, total = norun+a, total+b
	}
	limit := norunLimit(n)
	if total == 0 || norun*100 <= limit*total {
		return
	}
	r.add("X08", Level, n.Rel+"/node.md", 0, "%d of the %d code blocks of the node are norun, which is more than %d%%; run them as run, fails, wrong or test blocks, or put a fragment in a group that runs", norun, total, limit)
}
