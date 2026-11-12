// Package universe reads the universe register of spec 18.3 and counts the outcomes of each scope.
//
// The universe is the set of all source items of programming knowledge. Each card of cards/ and each node of nodes/ is a universe item with the outcome card or node, so the register does not repeat them. The files registers/universe/<source>.jsonl hold the other rows: the items that are aliases, sections, merged, on another site, out of scope, not an item, or still pending.
package universe

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/tamnd/programming-agent/node"
	"go.yaml.in/yaml/v3"
)

// Row is one row of registers/universe/<source>.jsonl.
type Row struct {
	// Source is the id of the source, such as pldb or cwe. An empty source takes the name of the file.
	Source string `json:"source,omitempty"`
	// ID is the id of the item in the source, and Name is the name that the source gives it.
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	// Scope is one of Scopes.
	Scope string `json:"scope"`
	// Outcome is one of Outcomes.
	Outcome string `json:"outcome"`
	// Target is the slug of the node or card that the row points to, or the tag on the other site for elsewhere outcomes.
	Target string `json:"target,omitempty"`
	// Reason says why the row has its outcome. Out of scope needs a reason code.
	Reason string `json:"reason,omitempty"`
	// Rule is the id of the rule of registers/scope-rules.yaml for the reason too-small.
	Rule string `json:"rule,omitempty"`

	File string `json:"-"`
	Line int    `json:"-"`
}

// Outcomes are the outcomes of spec 18.3. The outcomes node and card come from nodes and cards, but a row may also use them to record a source name that maps to a node or card with another id.
var Outcomes = []string{"node", "card", "alias", "section", "merged", "disambiguation", "elsewhere:algorithm", "elsewhere:math2", "out-of-scope", "not-an-item", "pending"}

// Scopes are the scopes of spec 18.4.
var Scopes = []string{"languages", "releases", "concepts", "features", "api", "diagnostics", "proposals", "tasks", "patterns", "weaknesses", "tools", "libraries", "formats", "licences", "people", "terms"}

// Reasons are the reason codes of the outcome out-of-scope.
var Reasons = []string{"not-programming", "product-marketing", "personal", "too-small"}

// Targeted tells if an outcome needs a target that is a node or a card of this site.
func Targeted(outcome string) bool {
	switch outcome {
	case "node", "card", "alias", "section", "merged":
		return true
	}
	return false
}

// ScopeOfKind returns the scope of the nodes and cards of a kind, or "" for a kind that is not a universe item, such as area.
func ScopeOfKind(kind string) string {
	switch kind {
	case "concept", "definition":
		return "concepts"
	case "language":
		return "languages"
	case "feature":
		return "features"
	case "release":
		return "releases"
	case "task", "exercise":
		return "tasks"
	case "pattern", "technique", "practice":
		return "patterns"
	case "tool":
		return "tools"
	case "library":
		return "libraries"
	case "format", "standard":
		return "formats"
	case "weakness":
		return "weaknesses"
	case "diagnostic":
		return "diagnostics"
	case "proposal":
		return "proposals"
	case "api":
		return "api"
	case "licence":
		return "licences"
	case "work", "person", "organization":
		return "people"
	}
	return ""
}

// Read reads every registers/universe/*.jsonl file under root. A missing directory gives no rows. Errors are the rows that cannot be read.
func Read(root string) ([]Row, []error) {
	files, err := filepath.Glob(filepath.Join(root, "registers", "universe", "*.jsonl"))
	if err != nil {
		return nil, []error{err}
	}
	sort.Strings(files)
	var rows []Row
	var errs []error
	for _, f := range files {
		rel, _ := filepath.Rel(root, f)
		rel = filepath.ToSlash(rel)
		src := strings.TrimSuffix(filepath.Base(f), ".jsonl")
		fh, err := os.Open(f)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		sc := bufio.NewScanner(fh)
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		line := 0
		for sc.Scan() {
			line++
			b := sc.Bytes()
			if len(strings.TrimSpace(string(b))) == 0 {
				continue
			}
			var r Row
			dec := json.NewDecoder(strings.NewReader(string(b)))
			dec.DisallowUnknownFields()
			if err := dec.Decode(&r); err != nil {
				errs = append(errs, fmt.Errorf("%s:%d: %w", rel, line, err))
				continue
			}
			if r.Source == "" {
				r.Source = src
			}
			r.File, r.Line = rel, line
			rows = append(rows, r)
		}
		if err := sc.Err(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", rel, err))
		}
		_ = fh.Close()
	}
	return rows, errs
}

// Problem is one problem of a row, with the check rule that finds it.
type Problem struct {
	Rule string
	Msg  string
}

// ScopeRule is one rule of registers/scope-rules.yaml. A row with the reason too-small names the id of a rule, so that a person approves the rule and not each row.
type ScopeRule struct {
	ID    string `yaml:"id"`
	Scope string `yaml:"scope"`
	Rule  string `yaml:"rule"`
}

// ReadRules reads registers/scope-rules.yaml. A missing file gives no rules.
func ReadRules(root string) (map[string]ScopeRule, error) {
	b, err := os.ReadFile(filepath.Join(root, "registers", "scope-rules.yaml"))
	if errors.Is(err, os.ErrNotExist) {
		return map[string]ScopeRule{}, nil
	}
	if err != nil {
		return nil, err
	}
	var list []ScopeRule
	if err := yaml.Unmarshal(b, &list); err != nil {
		return nil, fmt.Errorf("registers/scope-rules.yaml: %w", err)
	}
	out := map[string]ScopeRule{}
	for _, r := range list {
		if r.ID == "" || r.Rule == "" {
			return nil, fmt.Errorf("registers/scope-rules.yaml: each rule needs an id and a rule")
		}
		if _, dup := out[r.ID]; dup {
			return nil, fmt.Errorf("registers/scope-rules.yaml: the id %s is used twice", r.ID)
		}
		out[r.ID] = r
	}
	return out, nil
}

// Problems returns the problems of one row: U01 for an unknown scope or outcome, U02 for a missing reason or target, U03 for a target that is not a node or card of c, and U04 for a reason too-small without a rule of rules.
func Problems(c *node.Corpus, rules map[string]ScopeRule, r Row) []Problem {
	var out []Problem
	add := func(rule, f string, a ...any) { out = append(out, Problem{rule, fmt.Sprintf(f, a...)}) }
	if r.ID == "" {
		add("U01", "the row has no id")
	}
	if !slices.Contains(Scopes, r.Scope) {
		add("U01", "the scope %q is not a scope of spec 18.4; use one of %s", r.Scope, strings.Join(Scopes, ", "))
	}
	if !slices.Contains(Outcomes, r.Outcome) {
		add("U01", "the outcome %q is not an outcome of spec 18.3", r.Outcome)
		return out
	}
	if r.Outcome == "pending" {
		return out
	}
	if Targeted(r.Outcome) {
		switch {
		case r.Target == "":
			add("U02", "the outcome %s needs a target", r.Outcome)
		case c != nil && c.BySlug[r.Target] == nil && c.CardBySlug[r.Target] == nil:
			add("U03", "the target %q is not a node or a card", r.Target)
		}
	} else if r.Reason == "" && r.Target == "" {
		add("U02", "the outcome %s needs a reason or a target", r.Outcome)
	}
	if r.Outcome == "out-of-scope" {
		switch {
		case !slices.Contains(Reasons, r.Reason):
			add("U02", "out-of-scope needs one of the reason codes %s", strings.Join(Reasons, ", "))
		case r.Reason == "too-small" && r.Rule == "":
			add("U04", "the reason too-small needs a rule of registers/scope-rules.yaml")
		case r.Reason == "too-small" && rules != nil:
			if _, ok := rules[r.Rule]; !ok {
				add("U04", "the rule %q is not in registers/scope-rules.yaml", r.Rule)
			}
		}
	}
	return out
}

// Count holds the outcomes of one scope.
type Count struct {
	Scope     string         `json:"scope"`
	Total     int            `json:"total"`
	Pending   int            `json:"pending"`
	ByOutcome map[string]int `json:"by_outcome"`
}

// Share is the part of the items of the scope that is pending, from 0 to 1.
func (n Count) Share() float64 {
	if n.Total == 0 {
		return 0
	}
	return float64(n.Pending) / float64(n.Total)
}

// Report counts the items of each scope: the nodes and cards of c by kind, and the rows. A row with the outcome node or card that points to a node or card is not counted again. The last Count has the scope "all".
func Report(c *node.Corpus, rows []Row) []Count {
	by := map[string]*Count{}
	get := func(s string) *Count {
		if by[s] == nil {
			by[s] = &Count{Scope: s, ByOutcome: map[string]int{}}
		}
		return by[s]
	}
	for _, n := range c.Nodes {
		if s := ScopeOfKind(n.Kind); s != "" {
			get(s).ByOutcome["node"]++
		}
	}
	for _, k := range c.Cards {
		if s := ScopeOfKind(k.Kind); s != "" {
			get(s).ByOutcome["card"]++
		}
	}
	for _, r := range rows {
		if (r.Outcome == "node" || r.Outcome == "card") && r.Target != "" {
			continue
		}
		get(r.Scope).ByOutcome[r.Outcome]++
	}
	all := &Count{Scope: "all", ByOutcome: map[string]int{}}
	var out []Count
	for _, s := range Scopes {
		n := by[s]
		if n == nil {
			n = &Count{Scope: s, ByOutcome: map[string]int{}}
		}
		for o, v := range n.ByOutcome {
			n.Total += v
			all.ByOutcome[o] += v
		}
		n.Pending = n.ByOutcome["pending"]
		all.Total += n.Total
		all.Pending += n.Pending
		out = append(out, *n)
	}
	return append(out, *all)
}

// Gate tells if the counts meet the exit criteria of M1 (spec 21): less than max of all items pending, and no pending item in the strict scopes. It returns the reasons when they are not met.
func Gate(counts []Count, max float64, strict []string) error {
	var errs []error
	for _, n := range counts {
		if n.Scope == "all" && n.Share() >= max {
			errs = append(errs, fmt.Errorf("%.2f%% of all %d items are pending, the limit is %.2f%%", 100*n.Share(), n.Total, 100*max))
		}
		if slices.Contains(strict, n.Scope) && n.Pending > 0 {
			errs = append(errs, fmt.Errorf("the scope %s has %d pending items, the limit is 0", n.Scope, n.Pending))
		}
	}
	return errors.Join(errs...)
}
