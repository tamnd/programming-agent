package node

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/tamnd/programming-agent/compat"
)

// Area is one entry of topics/areas.yaml (spec 05.2).
type Area struct {
	ID       string `yaml:"id" json:"id"`
	Name     string `yaml:"name" json:"name"`
	Summary  string `yaml:"summary,omitempty" json:"summary,omitempty"`
	Children []Area `yaml:"children,omitempty" json:"children,omitempty"`
	Parent   string `yaml:"-" json:"parent,omitempty"`
}

// Path is "top/sub" for a sub-area and "top" for a top area. A sub-area of a language tree has three parts, such as languages/python/control-flow.
func (a Area) Path() string {
	if a.Parent == "" {
		return a.ID
	}
	return a.Parent + "/" + a.ID
}

// URL is the path of the area page. A topic area is at /<id>/. A language tree is at /<lang>/ and its sub-areas are at /<lang>/<sub>/ (spec 16.2).
func (a Area) URL() string {
	if p := a.Path(); strings.HasPrefix(p, "languages/") {
		return "/" + strings.TrimPrefix(p, "languages/") + "/"
	}
	return "/" + a.ID + "/"
}

// InLanguage tells if the area is a language tree or a part of one.
func (a Area) InLanguage() bool { return strings.HasPrefix(a.Path(), "languages/") }

// Lang is one entry of registers/languages.yaml (spec 08.2).
type Lang struct {
	ID        string   `yaml:"id" json:"id"`
	Name      string   `yaml:"name" json:"name"`
	Tier      int      `yaml:"tier" json:"tier"`
	Ext       []string `yaml:"ext" json:"ext"`
	FileCase  string   `yaml:"file_case" json:"file_case,omitempty"`
	Test      LangTool `yaml:"test" json:"test"`
	Format    LangTool `yaml:"format" json:"format"`
	Lint      LangTool `yaml:"lint" json:"lint"`
	Toolchain string   `yaml:"toolchain" json:"toolchain,omitempty"` // the pinned version, for the language page
	Run       string   `yaml:"run" json:"run,omitempty"`             // how a reader runs the tests without prog
	Comment   string   `yaml:"comment" json:"comment,omitempty"`
	Highlight string   `yaml:"highlight" json:"highlight"` // the Chroma lexer
	Int       string   `yaml:"int" json:"int,omitempty"`
	List      string   `yaml:"list" json:"list,omitempty"`
	Since     string   `yaml:"since" json:"since,omitempty"`
}

// LangTool is a tool of a language: the test framework, the formatter or the linter.
type LangTool struct {
	Tool      string   `yaml:"tool,omitempty" json:"tool,omitempty"`
	Framework string   `yaml:"framework,omitempty" json:"framework,omitempty"`
	File      string   `yaml:"file,omitempty" json:"file,omitempty"`
	Args      []string `yaml:"args,omitempty" json:"args,omitempty"`
}

// Tag is one line of registers/tags.reg.
type Tag struct {
	Tag  string
	Slug string
	Line int
}

// Corpus is the content tree.
type Corpus struct {
	Root     string
	Nodes    []*Node
	BySlug   map[string]*Node
	ByTag    map[string]*Node
	Areas    []Area
	AreaByID map[string]*Area // keyed by Area.Path()
	// Template is the list of sub-areas of a language tree, from topics/languages.yaml (spec 05.3).
	Template []Area
	Langs    []Lang
	LangByID map[string]*Lang
	Reserved map[string]bool
	Tags     []Tag
	// OldSlugs maps each old slug of registers/aliases.reg to the slug it moved to.
	OldSlugs map[string]string
	// Errors are the node directories that could not be read. Package check reports them as T01.
	Errors []error
	// Cards are the rows of the card registers under cards/, sorted by slug (spec 04.1).
	Cards []*Card
	// CardBySlug holds the first card of each slug.
	CardBySlug map[string]*Card
	// CardErrors are the card rows that could not be read. Package check reports them as C01.
	CardErrors []error
	// Impls is registers/implementations.yaml: the compat columns of each language (spec 09.2). Probes are the results of prog compat probe, from manifests/probes.jsonl.
	Impls  compat.Register
	Probes compat.Results
	// realizers maps a concept slug to the features that realize it.
	realizers map[string][]*Node
}

// Load reads the whole content tree at root.
func Load(root string) (*Corpus, error) {
	c := &Corpus{
		Root:     root,
		BySlug:   map[string]*Node{},
		ByTag:    map[string]*Node{},
		AreaByID: map[string]*Area{},
		LangByID: map[string]*Lang{},
		Reserved: map[string]bool{},
	}
	if err := readYAML(filepath.Join(root, "topics", "areas.yaml"), &c.Areas); err != nil {
		return nil, err
	}
	if err := readYAML(filepath.Join(root, "registers", "languages.yaml"), &c.Langs); err != nil {
		return nil, err
	}
	sort.SliceStable(c.Langs, func(i, j int) bool { return c.Langs[i].Tier < c.Langs[j].Tier })
	for i := range c.Langs {
		c.LangByID[c.Langs[i].ID] = &c.Langs[i]
	}
	if err := readYAML(filepath.Join(root, "topics", "languages.yaml"), &c.Template); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	c.languageTrees()
	for i := range c.Areas {
		top := &c.Areas[i]
		c.AreaByID[top.Path()] = top
		for j := range top.Children {
			sub := &top.Children[j]
			sub.Parent = top.ID
			c.AreaByID[sub.Path()] = sub
			for k := range sub.Children {
				leaf := &sub.Children[k]
				leaf.Parent = sub.Path()
				c.AreaByID[leaf.Path()] = leaf
			}
		}
	}
	var reserved []string
	if err := readYAML(filepath.Join(root, "registers", "reserved.yaml"), &reserved); err != nil {
		return nil, err
	}
	for _, w := range reserved {
		c.Reserved[w] = true
	}
	tags, err := readTags(filepath.Join(root, "registers", "tags.reg"))
	if err != nil {
		return nil, err
	}
	c.Tags = tags
	if c.OldSlugs, err = ReadAliases(filepath.Join(root, "registers", "aliases.reg")); err != nil {
		return nil, err
	}

	if c.Impls, err = compat.ReadRegister(root); err != nil {
		return nil, err
	}
	if c.Probes, err = compat.ReadResults(root); err != nil {
		return nil, err
	}

	nodes := filepath.Join(root, "nodes")
	err = filepath.WalkDir(nodes, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "node.md" {
			return nil
		}
		n, err := Read(root, filepath.Dir(p))
		if err != nil {
			c.Errors = append(c.Errors, err)
			return nil
		}
		c.Nodes = append(c.Nodes, n)
		return nil
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	sort.Slice(c.Nodes, func(i, j int) bool { return c.Nodes[i].Slug < c.Nodes[j].Slug })
	for _, n := range c.Nodes {
		if n.Slug != "" {
			if _, dup := c.BySlug[n.Slug]; !dup {
				c.BySlug[n.Slug] = n
			}
		}
		if n.Tag != "" {
			if _, dup := c.ByTag[n.Tag]; !dup {
				c.ByTag[n.Tag] = n
			}
		}
	}
	c.Cards, c.CardErrors = ReadCards(root)
	c.CardBySlug = map[string]*Card{}
	for _, k := range c.Cards {
		if _, dup := c.CardBySlug[k.Slug()]; !dup {
			c.CardBySlug[k.Slug()] = k
		}
	}
	c.realizers = map[string][]*Node{}
	for _, n := range c.Nodes {
		if n.Kind == "feature" {
			for _, r := range n.Realizes {
				c.realizers[r] = append(c.realizers[r], n)
			}
		}
	}
	return c, nil
}

// languageTrees adds a sub-area languages/<lang> for each language of the register, with the sub-areas of the template under it (spec 05.3). It does nothing when the topic tree has no area "languages".
func (c *Corpus) languageTrees() {
	var top *Area
	for i := range c.Areas {
		if c.Areas[i].ID == "languages" {
			top = &c.Areas[i]
		}
	}
	if top == nil {
		return
	}
	have := map[string]bool{}
	for _, a := range top.Children {
		have[a.ID] = true
	}
	for _, l := range c.Langs {
		if have[l.ID] {
			continue
		}
		tree := Area{ID: l.ID, Name: l.Name, Summary: "The guide to " + l.Name + "."}
		for _, t := range c.Template {
			tree.Children = append(tree.Children, Area{ID: t.ID, Name: t.Name, Summary: t.Summary})
		}
		top.Children = append(top.Children, tree)
	}
}

// InArea returns the nodes whose primary area is path or is under path.
func (c *Corpus) InArea(path string) []*Node {
	var out []*Node
	for _, n := range c.Nodes {
		if n.Area == path || strings.HasPrefix(n.Area, path+"/") {
			out = append(out, n)
		}
	}
	return out
}

// Kind returns the nodes of one kind.
func (c *Corpus) Kind(kind string) []*Node {
	var out []*Node
	for _, n := range c.Nodes {
		if n.Kind == kind {
			out = append(out, n)
		}
	}
	return out
}

// Solvers returns the nodes that solve a task. A task holds its own solutions (spec 10.1), so the task is first. The nodes with a `solves` edge to the task come after it.
func (c *Corpus) Solvers(problem string) []*Node {
	var out []*Node
	if t := c.BySlug[problem]; t != nil && t.Kind == "task" {
		out = append(out, t)
	}
	for _, n := range c.Nodes {
		for _, s := range n.Solves {
			if s.Problem == problem {
				out = append(out, n)
				break
			}
		}
	}
	return out
}

// Leads returns the nodes that list slug in prereqs or uses.
func (c *Corpus) Leads(slug string) []*Node {
	var out []*Node
	for _, n := range c.Nodes {
		for _, p := range append(append([]string{}, n.Prereqs...), n.Uses...) {
			if p == slug {
				out = append(out, n)
				break
			}
		}
	}
	return out
}

func readYAML(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, v); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

// readTags reads registers/tags.reg. Each line is "TAG slug". A line that starts with # is a comment.
// ReadAliases reads registers/aliases.reg. Each line is an old slug and the slug it moved to.
func ReadAliases(path string) (map[string]string, error) {
	out := map[string]string{}
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	line := 0
	for s.Scan() {
		line++
		t := strings.TrimSpace(s.Text())
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		fs := strings.Fields(t)
		if len(fs) != 2 {
			return nil, fmt.Errorf("%s:%d: a line needs an old slug and a new slug", path, line)
		}
		out[fs[0]] = fs[1]
	}
	return out, s.Err()
}

func readTags(path string) ([]Tag, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []Tag
	s := bufio.NewScanner(f)
	line := 0
	for s.Scan() {
		line++
		t := strings.TrimSpace(s.Text())
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		fs := strings.Fields(t)
		if len(fs) < 2 {
			return nil, fmt.Errorf("%s:%d: a line needs a tag and a slug", path, line)
		}
		out = append(out, Tag{Tag: fs[0], Slug: fs[1], Line: line})
	}
	return out, s.Err()
}
