// SPDX-License-Identifier: Apache-2.0

package node

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Verdicts are the values of a cell of the concept and language matrix (spec 04.8), in order from the strongest support to none.
var Verdicts = []string{"native", "stdlib", "library", "idiom", "absent", "n/a"}

// Cell is one cell of the matrix: how a language does a concept. Written cells come from the matrix.yaml file of a concept or the concepts.yaml file of a language. Other cells come from the example files and the features.
type Cell struct {
	Status      string `yaml:"status" json:"status"`
	Reason      string `yaml:"reason,omitempty" json:"reason,omitempty"`
	Alternative string `yaml:"alternative,omitempty" json:"alternative,omitempty"`
	Library     string `yaml:"library,omitempty" json:"library,omitempty"`
	Source      string `yaml:"source,omitempty" json:"source,omitempty"`
	// From tells where the cell comes from: a file path, "example" or "feature".
	From string `yaml:"-" json:"from,omitempty"`
	// Feature is the slug of the feature node that realizes the concept, for a cell that comes from a feature.
	Feature string `yaml:"-" json:"feature,omitempty"`
	// Line is the line of the cell in its file, for the checks.
	Line int `yaml:"-" json:"-"`
}

// readCells reads a file that maps a key (a language id or a concept slug) to a cell. A missing file gives no cells.
func readCells(path string) (map[string]*Cell, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	out := map[string]*Cell{}
	if len(doc.Content) == 0 {
		return out, nil
	}
	m := doc.Content[0]
	if m.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s: the file must be a map from a key to a cell", path)
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		k, v := m.Content[i], m.Content[i+1]
		c := &Cell{Line: k.Line}
		if v.Kind == yaml.ScalarNode {
			c.Status = v.Value
		} else if err := v.Decode(c); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, k.Line, err)
		}
		out[k.Value] = c
	}
	return out, nil
}

var verdictRE = regexp.MustCompile(`(?mi)^- (?:\*\*)?Verdict:?(?:\*\*)?:? *([a-z/]+)`)

// exampleVerdict reads the verdict that an example file states in a line such as "- Verdict: native.".
func exampleVerdict(body string) string {
	m := verdictRE.FindStringSubmatch(body)
	if m == nil {
		return ""
	}
	return strings.ToLower(m[1])
}

// readMatrix reads matrix.yaml of a concept and concepts.yaml of a language.
func readMatrix(root string, n *Node) error {
	name := ""
	switch n.Kind {
	case "concept":
		name = "matrix.yaml"
	case "language":
		name = "concepts.yaml"
	default:
		return nil
	}
	p := filepath.Join(n.Dir, name)
	cells, err := readCells(p)
	if err != nil {
		return err
	}
	rel, _ := filepath.Rel(root, p)
	for _, c := range cells {
		c.From = filepath.ToSlash(rel)
	}
	n.Cells = cells
	return nil
}

// Cell gives the cell of a concept for a language, and false when nothing decides it. The order is: the matrix.yaml of the concept, the concepts.yaml of the language, the example file of the language, then a feature of the language that realizes the concept.
func (c *Corpus) Cell(concept *Node, lang string) (*Cell, bool) {
	if x, ok := concept.Cells[lang]; ok {
		return x, true
	}
	if l, ok := c.BySlug[lang]; ok && l.Kind == "language" {
		if x, ok := l.Cells[concept.Slug]; ok {
			return x, true
		}
	}
	if ex := concept.Example(lang); ex != nil {
		v := ex.Verdict
		if v == "" {
			v = "native"
		}
		return &Cell{Status: v, From: "example"}, true
	}
	for _, f := range c.realizers[concept.Slug] {
		if f.Owner == lang {
			return &Cell{Status: "native", From: "feature", Feature: f.Slug}, true
		}
	}
	return nil, false
}
