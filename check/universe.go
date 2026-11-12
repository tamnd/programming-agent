// SPDX-License-Identifier: Apache-2.0

package check

import (
	"fmt"

	"github.com/tamnd/programming-agent/universe"
)

func init() {
	Rules["U01"] = "each row of registers/universe/*.jsonl is valid JSON with an id, a scope of spec 18.4 and an outcome of spec 18.3"
	Rules["U02"] = "an outcome other than pending has a reason or a target; node, card, alias, section and merged have a target; out-of-scope has a reason code"
	Rules["U03"] = "each target of a universe row is a node or a card"
	Rules["U04"] = "out-of-scope with too-small names a rule of registers/scope-rules.yaml"
	Rules["U05"] = "each source id is in the universe register once"
}

// universeRows runs the rules U01 to U05 on the universe register.
func (r *run) universeRows() {
	rows, errs := universe.Read(r.root)
	for _, err := range errs {
		r.add("U01", Block, "registers/universe", 0, "%v", err)
	}
	rules, err := universe.ReadRules(r.root)
	if err != nil {
		r.add("U04", Block, "registers/scope-rules.yaml", 0, "%v", err)
	}
	seen := map[string]string{}
	for _, row := range rows {
		for _, p := range universe.Problems(r.c, rules, row) {
			r.add(p.Rule, Block, row.File, row.Line, "%s", p.Msg)
		}
		key := row.Source + "\x00" + row.ID
		if at, dup := seen[key]; dup {
			r.add("U05", Block, row.File, row.Line, "the id %q of the source %s is also at %s", row.ID, row.Source, at)
			continue
		}
		seen[key] = fmt.Sprintf("%s:%d", row.File, row.Line)
	}
}
