// Package node reads the node directories of tamnd/programming (spec 04.3 and 04.4).
package node

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Levels are the editorial levels of spec 05.5.
var Levels = []string{"intro", "core", "advanced", "research"}

// Solve is one `solves` edge. The short form is a problem slug. The long form has conditions (spec 04.4).
type Solve struct {
	Problem  string         `yaml:"problem" json:"problem"`
	Requires map[string]any `yaml:"requires,omitempty" json:"requires,omitempty"`
	Time     string         `yaml:"time,omitempty" json:"time,omitempty"`
	Note     string         `yaml:"note,omitempty" json:"note,omitempty"`
}

// UnmarshalYAML accepts a slug or a mapping.
func (s *Solve) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		s.Problem = n.Value
		return nil
	}
	type plain Solve
	var p plain
	if err := n.Decode(&p); err != nil {
		return err
	}
	*s = Solve(p)
	return nil
}

// Complexity holds the claims of a node (spec 04.4 and 06.6).
type Complexity struct {
	Vars  map[string]string            `yaml:"vars,omitempty" json:"vars,omitempty"`
	Time  map[string]string            `yaml:"time,omitempty" json:"time,omitempty"`
	Space map[string]string            `yaml:"space,omitempty" json:"space,omitempty"`
	Other map[string]map[string]string `yaml:",inline" json:"other,omitempty"`
}

// History says where a node comes from.
type History struct {
	Year   int            `yaml:"year,omitempty" json:"year,omitempty"`
	People []string       `yaml:"people,omitempty" json:"people,omitempty"`
	First  map[string]any `yaml:"first,omitempty" json:"first,omitempty"`
	Note   string         `yaml:"note,omitempty" json:"note,omitempty"`
}

// Source is one entry of `sources`. A source names a work of registers/works.yaml, a book, a paper or a URL, and a locator inside it (spec 04.5).
type Source struct {
	Work     string `yaml:"work,omitempty" json:"work,omitempty"`
	Book     string `yaml:"book,omitempty" json:"book,omitempty"`
	URL      string `yaml:"url,omitempty" json:"url,omitempty"`
	Paper    string `yaml:"paper,omitempty" json:"paper,omitempty"`
	Title    string `yaml:"title,omitempty" json:"title,omitempty"`
	Locator  string `yaml:"locator,omitempty" json:"locator,omitempty"`
	Accessed string `yaml:"accessed,omitempty" json:"accessed,omitempty"`
	Note     string `yaml:"note,omitempty" json:"note,omitempty"`
}

// Lit is a list of strings that also accepts one string, such as `languages: tier-1+2`.
type Lit []string

// UnmarshalYAML accepts a scalar or a sequence.
func (l *Lit) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		*l = Lit{n.Value}
		return nil
	}
	var v []string
	if err := n.Decode(&v); err != nil {
		return err
	}
	*l = v
	return nil
}

// Front is the front matter of node.md (spec 04.5). The schema is schema/node.v1.json in tamnd/programming. A field that does not apply to a kind is absent.
type Front struct {
	Tag         string            `yaml:"tag" json:"tag"`
	Kind        string            `yaml:"kind" json:"kind"`
	Slug        string            `yaml:"slug" json:"slug"`
	Name        string            `yaml:"name" json:"name"`
	Aliases     []string          `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	Area        string            `yaml:"area" json:"area"`
	Areas       []string          `yaml:"areas,omitempty" json:"areas,omitempty"`
	Level       string            `yaml:"level" json:"level"`
	Summary     string            `yaml:"summary" json:"summary"`
	Solves      []Solve           `yaml:"solves,omitempty" json:"solves,omitempty"`
	Implements  []string          `yaml:"implements,omitempty" json:"implements,omitempty"`
	Techniques  []string          `yaml:"techniques,omitempty" json:"techniques,omitempty"`
	Uses        []string          `yaml:"uses,omitempty" json:"uses,omitempty"`
	VariantOf   string            `yaml:"variant_of,omitempty" json:"variant_of,omitempty"`
	Generalizes []string          `yaml:"generalizes,omitempty" json:"generalizes,omitempty"`
	Model       string            `yaml:"model,omitempty" json:"model,omitempty"`
	Complexity  *Complexity       `yaml:"complexity,omitempty" json:"complexity,omitempty"`
	Properties  map[string]any    `yaml:"properties,omitempty" json:"properties,omitempty"`
	Constraints map[string]any    `yaml:"constraints,omitempty" json:"constraints,omitempty"`
	History     *History          `yaml:"history,omitempty" json:"history,omitempty"`
	Prereqs     []string          `yaml:"prereqs,omitempty" json:"prereqs,omitempty"`
	SameAs      map[string]any    `yaml:"same_as,omitempty" json:"same_as,omitempty"`
	Lean        map[string]string `yaml:"lean,omitempty" json:"lean,omitempty"`
	Sources     []Source          `yaml:"sources,omitempty" json:"sources,omitempty"`
	Exercises   []string          `yaml:"exercises,omitempty" json:"exercises,omitempty"`
	People      []string          `yaml:"people,omitempty" json:"people,omitempty"`
	SeeAlso     []string          `yaml:"see_also,omitempty" json:"see_also,omitempty"`

	// Owner is the language, tool or library that owns a feature, release, diagnostic, proposal or api node (spec 04.3).
	Owner string `yaml:"owner,omitempty" json:"owner,omitempty"`

	// concept
	AppliesTo string `yaml:"applies_to,omitempty" json:"applies_to,omitempty"`
	Viz       string `yaml:"viz,omitempty" json:"viz,omitempty"`
	Formal    string `yaml:"formal,omitempty" json:"formal,omitempty"`

	// language
	Class          string         `yaml:"class,omitempty" json:"class,omitempty"`
	Tier           int            `yaml:"tier,omitempty" json:"tier,omitempty"`
	Paradigms      []string       `yaml:"paradigms,omitempty" json:"paradigms,omitempty"`
	Typing         map[string]any `yaml:"typing,omitempty" json:"typing,omitempty"`
	Memory         string         `yaml:"memory,omitempty" json:"memory,omitempty"`
	Execution      []string       `yaml:"execution,omitempty" json:"execution,omitempty"`
	Appeared       int            `yaml:"appeared,omitempty" json:"appeared,omitempty"`
	Designers      []string       `yaml:"designers,omitempty" json:"designers,omitempty"`
	InfluencedBy   []string       `yaml:"influenced_by,omitempty" json:"influenced_by,omitempty"`
	Standard       string         `yaml:"standard,omitempty" json:"standard,omitempty"`
	Reference      map[string]any `yaml:"reference,omitempty" json:"reference,omitempty"`
	ReferenceImpl  string         `yaml:"licence_of_reference_impl,omitempty" json:"licence_of_reference_impl,omitempty"`
	FileExtensions []string       `yaml:"file_extensions,omitempty" json:"file_extensions,omitempty"`
	HelloWorld     string         `yaml:"hello_world,omitempty" json:"hello_world,omitempty"`

	// feature
	Realizes  []string       `yaml:"realizes,omitempty" json:"realizes,omitempty"`
	Grammar   string         `yaml:"grammar,omitempty" json:"grammar,omitempty"`
	Spec      map[string]any `yaml:"spec,omitempty" json:"spec,omitempty"`
	Proposals []string       `yaml:"proposals,omitempty" json:"proposals,omitempty"`
	Compat    string         `yaml:"compat,omitempty" json:"compat,omitempty"`
	Since     string         `yaml:"since,omitempty" json:"since,omitempty"`

	// release
	Version    string           `yaml:"version,omitempty" json:"version,omitempty"`
	Released   string           `yaml:"released,omitempty" json:"released,omitempty"`
	EOL        string           `yaml:"eol,omitempty" json:"eol,omitempty"`
	Patches    []map[string]any `yaml:"patches,omitempty" json:"patches,omitempty"`
	Adds       []string         `yaml:"adds,omitempty" json:"adds,omitempty"`
	Removes    []string         `yaml:"removes,omitempty" json:"removes,omitempty"`
	Deprecates []string         `yaml:"deprecates,omitempty" json:"deprecates,omitempty"`
	Notes      map[string]any   `yaml:"notes,omitempty" json:"notes,omitempty"`

	// task
	Contract       string              `yaml:"contract,omitempty" json:"contract,omitempty"`
	Category       string              `yaml:"category,omitempty" json:"category,omitempty"`
	Languages      Lit                 `yaml:"languages,omitempty" json:"languages,omitempty"`
	UsesFeatures   map[string][]string `yaml:"uses_features,omitempty" json:"uses_features,omitempty"`
	UsesAlgorithms []string            `yaml:"uses_algorithms,omitempty" json:"uses_algorithms,omitempty"`

	// pattern, technique and practice
	Family       string            `yaml:"family,omitempty" json:"family,omitempty"`
	Polarity     string            `yaml:"polarity,omitempty" json:"polarity,omitempty"`
	Intent       string            `yaml:"intent,omitempty" json:"intent,omitempty"`
	AlsoKnownAs  []string          `yaml:"also_known_as,omitempty" json:"also_known_as,omitempty"`
	Related      []string          `yaml:"related,omitempty" json:"related,omitempty"`
	BuiltIn      map[string]string `yaml:"languages_where_built_in,omitempty" json:"languages_where_built_in,omitempty"`
	Refactorings []string          `yaml:"refactorings,omitempty" json:"refactorings,omitempty"`

	// tool and library
	Role               Lit            `yaml:"role,omitempty" json:"role,omitempty"`
	ImplementsLanguage string         `yaml:"implements_language,omitempty" json:"implements_language,omitempty"`
	Licence            string         `yaml:"licence,omitempty" json:"licence,omitempty"`
	Repo               string         `yaml:"repo,omitempty" json:"repo,omitempty"`
	Homepage           string         `yaml:"homepage,omitempty" json:"homepage,omitempty"`
	FirstRelease       string         `yaml:"first_release,omitempty" json:"first_release,omitempty"`
	AlternativeTo      []string       `yaml:"alternative_to,omitempty" json:"alternative_to,omitempty"`
	Latest             map[string]any `yaml:"latest,omitempty" json:"latest,omitempty"`

	// format and standard
	MediaType  string   `yaml:"media_type,omitempty" json:"media_type,omitempty"`
	Extensions []string `yaml:"extensions,omitempty" json:"extensions,omitempty"`
	Specs      []string `yaml:"specs,omitempty" json:"specs,omitempty"`
	Publisher  string   `yaml:"publisher,omitempty" json:"publisher,omitempty"`

	// weakness
	CWE               int      `yaml:"cwe,omitempty" json:"cwe,omitempty"`
	Abstraction       string   `yaml:"abstraction,omitempty" json:"abstraction,omitempty"`
	LanguagesAffected []string `yaml:"languages_affected,omitempty" json:"languages_affected,omitempty"`
	OWASPTop10        string   `yaml:"owasp_top10_2025,omitempty" json:"owasp_top10_2025,omitempty"`
	CWETop25          int      `yaml:"cwe_top25_2025,omitempty" json:"cwe_top25_2025,omitempty"`
	DetectedBy        []string `yaml:"detected_by,omitempty" json:"detected_by,omitempty"`
	Prevents          []string `yaml:"prevents,omitempty" json:"prevents,omitempty"`

	// diagnostic, proposal and api, when they are promoted to nodes
	Code     string   `yaml:"code,omitempty" json:"code,omitempty"`
	Severity string   `yaml:"severity,omitempty" json:"severity,omitempty"`
	Messages []string `yaml:"message_templates,omitempty" json:"message_templates,omitempty"`
	Concepts []string `yaml:"concepts,omitempty" json:"concepts,omitempty"`
	Features []string `yaml:"features,omitempty" json:"features,omitempty"`
	Path     string   `yaml:"path,omitempty" json:"api_path,omitempty"`
	Entity   string   `yaml:"entity,omitempty" json:"entity,omitempty"`
	Module   string   `yaml:"module,omitempty" json:"module,omitempty"`
	// Signature is the signature of an api entry. ProposalStatus, ProposalType, Created and Release are the facts of a proposal; the key is not status, because a status is computed (rule T05).
	Signature      string `yaml:"signature,omitempty" json:"signature,omitempty"`
	ProposalStatus string `yaml:"proposal_status,omitempty" json:"proposal_status,omitempty"`
	ProposalType   string `yaml:"proposal_type,omitempty" json:"proposal_type,omitempty"`
	Created        string `yaml:"created,omitempty" json:"created,omitempty"`
	Release        string `yaml:"release,omitempty" json:"release,omitempty"`

	// work, person and organization
	Authors []string `yaml:"authors,omitempty" json:"authors,omitempty"`
	Year    int      `yaml:"year,omitempty" json:"year,omitempty"`
	Born    string   `yaml:"born,omitempty" json:"born,omitempty"`
	Died    string   `yaml:"died,omitempty" json:"died,omitempty"`
	Founded string   `yaml:"founded,omitempty" json:"founded,omitempty"`
	URL     string   `yaml:"url,omitempty" json:"url,omitempty"`
}

// File is one file of an implementation.
type File struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"-"`
	Test    bool   `json:"test"`
}

// Impl is the implementation of a node in one language (spec 08.1).
type Impl struct {
	Lang  string `json:"lang"`
	Files []File `json:"files"`
}

// Main is the first file that is not a test.
func (im Impl) Main() *File {
	for i := range im.Files {
		if !im.Files[i].Test {
			return &im.Files[i]
		}
	}
	return nil
}

// HasTests tells if the implementation has native tests.
func (im Impl) HasTests() bool {
	for _, f := range im.Files {
		if f.Test {
			return true
		}
	}
	return false
}

// ProofFront is the front matter of a proof file (spec 06.5).
type ProofFront struct {
	Title  string   `yaml:"title" json:"title"`
	Proves []string `yaml:"proves" json:"proves"`
	Status string   `yaml:"status" json:"status"` // full or sketch
	Lean   string   `yaml:"lean,omitempty" json:"lean,omitempty"`
}

// Proof is one file of proofs/.
type Proof struct {
	ID    string     `json:"id"`
	Path  string     `json:"path"`
	Front ProofFront `json:"front"`
	Body  string     `json:"-"`
}

// Preset is a named input of a scene (spec 10.4).
type Preset struct {
	Name  string         `yaml:"name" json:"name"`
	Input map[string]any `yaml:"input" json:"input"`
	Note  string         `yaml:"note,omitempty" json:"note,omitempty"`
}

// Example is one file of examples/: the section of a concept or a feature for one language (spec 04.4).
type Example struct {
	Lang string `json:"lang"`
	Path string `json:"path"`
	Body string `json:"-"`
	Line int    `json:"-"` // the line of the file where the body starts
}

// Node is one node directory.
type Node struct {
	Front
	// Local is the slug as the front matter writes it. For an owned kind, Slug is the full slug, such as python/match-statement, and Local is match-statement.
	Local     string    `json:"local,omitempty"`
	Examples  []Example `json:"examples,omitempty"`
	Body      string    `json:"-"`
	BodyLine  int       `json:"-"` // the line of node.md where the body starts
	Dir       string    `json:"-"`
	Rel       string    `json:"path"` // the directory, relative to the root
	RawFront  []byte    `json:"-"`
	Forbidden []string  `json:"-"` // front matter keys that rule T05 rejects
	Impls     []Impl    `json:"impls,omitempty"`
	Proofs    []Proof   `json:"proofs,omitempty"`
	Scene     string    `json:"scene,omitempty"` // path of viz/scene.ts, relative to the root
	Presets   []Preset  `json:"presets,omitempty"`
}

// Lang returns the implementation in a language, or nil.
func (n *Node) Lang(id string) *Impl {
	for i := range n.Impls {
		if n.Impls[i].Lang == id {
			return &n.Impls[i]
		}
	}
	return nil
}

// forbiddenKeys are status fields. The build computes them, and nobody writes them (rule T05).
var forbiddenKeys = []string{"status", "ladder", "ladders", "badge", "badges", "evidence", "complete", "reviewed", "popularity"}

// SplitFront splits a Markdown file into its YAML front matter and its body.
// It returns the line number where the body starts.
func SplitFront(b []byte) (front, body []byte, bodyLine int, err error) {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(b, []byte("---\n")) {
		return nil, b, 1, errors.New("the file does not start with a front matter block")
	}
	rest := b[4:]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		if bytes.HasSuffix(rest, []byte("\n---")) {
			end = len(rest) - 4
		} else {
			return nil, b, 1, errors.New("the front matter block has no end line")
		}
	}
	front = rest[:end+1]
	after := min(end+5, len(rest))
	body = rest[after:]
	bodyLine = 2 + bytes.Count(front, []byte("\n")) + 1
	return front, body, bodyLine, nil
}

// Read reads one node directory. It does not check the node; package check does that.
func Read(root, dir string) (*Node, error) {
	path := filepath.Join(dir, "node.md")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	front, body, line, err := SplitFront(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	n := &Node{Body: string(body), BodyLine: line, Dir: dir, RawFront: front}
	n.Rel, _ = filepath.Rel(root, dir)
	n.Rel = filepath.ToSlash(n.Rel)

	// The forbidden keys come out of the YAML tree, not out of a map, so each scalar keeps its text: "3.10" stays "3.10".
	var doc yaml.Node
	if err := yaml.Unmarshal(front, &doc); err != nil {
		return nil, fmt.Errorf("%s: front matter: %w", path, err)
	}
	if len(doc.Content) == 1 && doc.Content[0].Kind == yaml.MappingNode {
		m := doc.Content[0]
		var keep []*yaml.Node
		for i := 0; i+1 < len(m.Content); i += 2 {
			if slices.Contains(forbiddenKeys, m.Content[i].Value) {
				n.Forbidden = append(n.Forbidden, m.Content[i].Value)
				continue
			}
			keep = append(keep, m.Content[i], m.Content[i+1])
		}
		m.Content = keep
		clean, err := yaml.Marshal(&doc)
		if err != nil {
			return nil, err
		}
		dec := yaml.NewDecoder(bytes.NewReader(clean))
		dec.KnownFields(true)
		if err := dec.Decode(&n.Front); err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("%s: front matter: %w", path, err)
		}
	}
	n.Local = n.Slug
	if n.Local == "" && n.Kind == "release" {
		n.Local = n.Version
	}
	if n.Local == "" && n.Kind == "diagnostic" {
		n.Local = n.Code
	}
	if n.Local == "" && n.Kind == "api" {
		n.Local = n.Path
	}
	n.Slug = FullSlug(n.Kind, n.Owner, n.Local)
	if err := readExamples(root, n); err != nil {
		return nil, err
	}

	if err := readImpls(root, n); err != nil {
		return nil, err
	}
	if err := readProofs(root, n); err != nil {
		return nil, err
	}
	if err := readScene(root, n); err != nil {
		return nil, err
	}
	return n, nil
}

// readExamples reads examples/<lang>.md. A file may start with a front matter block; the body is the part after it.
func readExamples(root string, n *Node) error {
	paths, err := filepath.Glob(filepath.Join(n.Dir, "examples", "*.md"))
	if err != nil {
		return err
	}
	sort.Strings(paths)
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
		ex := Example{Lang: strings.TrimSuffix(filepath.Base(p), ".md"), Body: string(b), Line: 1}
		if bytes.HasPrefix(b, []byte("---\n")) {
			if _, body, line, err := SplitFront(b); err == nil {
				ex.Body, ex.Line = string(body), line
			}
		}
		rel, _ := filepath.Rel(root, p)
		ex.Path = filepath.ToSlash(rel)
		n.Examples = append(n.Examples, ex)
	}
	return nil
}

// Example returns the example file of a language, or nil.
func (n *Node) Example(lang string) *Example {
	for i := range n.Examples {
		if n.Examples[i].Lang == lang {
			return &n.Examples[i]
		}
	}
	return nil
}

// IsTestFile tells if a file name is a native test, by the conventions of spec 08.3.
// Examples: test_binary_search.py, binary_search_test.go, binarySearch.test.js, BinarySearchTest.java, binary_search_tests.rs.
func IsTestFile(name string) bool {
	low := strings.ToLower(name)
	stem := strings.TrimSuffix(low, filepath.Ext(low))
	stem = strings.TrimSuffix(stem, ".test")
	switch {
	case strings.HasPrefix(low, "test_"), strings.Contains(low, "_test."), strings.Contains(low, ".test."), strings.Contains(low, "_spec."):
		return true
	case strings.HasSuffix(stem, "test"), strings.HasSuffix(stem, "tests"):
		return true
	}
	return false
}

func readImpls(root string, n *Node) error {
	base := filepath.Join(n.Dir, "impl")
	langs, err := os.ReadDir(base)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, l := range langs {
		if !l.IsDir() {
			continue
		}
		im := Impl{Lang: l.Name()}
		dir := filepath.Join(base, l.Name())
		err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
				return nil
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, p)
			name, _ := filepath.Rel(dir, p)
			im.Files = append(im.Files, File{
				Name:    filepath.ToSlash(name),
				Path:    filepath.ToSlash(rel),
				Content: string(b),
				Test:    IsTestFile(d.Name()),
			})
			return nil
		})
		if err != nil {
			return err
		}
		sort.SliceStable(im.Files, func(i, j int) bool {
			if im.Files[i].Test != im.Files[j].Test {
				return !im.Files[i].Test
			}
			return im.Files[i].Name < im.Files[j].Name
		})
		if len(im.Files) > 0 {
			n.Impls = append(n.Impls, im)
		}
	}
	return nil
}

func readProofs(root string, n *Node) error {
	paths, err := filepath.Glob(filepath.Join(n.Dir, "proofs", "*.md"))
	if err != nil {
		return err
	}
	sort.Strings(paths)
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		front, body, _, err := SplitFront(b)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		var pr Proof
		if err := yaml.Unmarshal(front, &pr.Front); err != nil {
			return fmt.Errorf("%s: front matter: %w", p, err)
		}
		pr.ID = strings.TrimSuffix(filepath.Base(p), ".md")
		rel, _ := filepath.Rel(root, p)
		pr.Path = filepath.ToSlash(rel)
		pr.Body = string(body)
		n.Proofs = append(n.Proofs, pr)
	}
	return nil
}

func readScene(root string, n *Node) error {
	p := filepath.Join(n.Dir, "viz", "scene.ts")
	if _, err := os.Stat(p); err == nil {
		rel, _ := filepath.Rel(root, p)
		n.Scene = filepath.ToSlash(rel)
	}
	b, err := os.ReadFile(filepath.Join(n.Dir, "viz", "presets.yaml"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &n.Presets); err != nil {
		return fmt.Errorf("%s/viz/presets.yaml: %w", n.Rel, err)
	}
	return nil
}

// BuildFiles are the names of the files that the generator writes (spec 08.5). Rule I15 refuses them in the content repo.
var BuildFiles = []string{
	"go.mod", "go.sum", "go.work", "Cargo.toml", "Cargo.lock", "package.json", "package-lock.json", "tsconfig.json",
	"pyproject.toml", "setup.py", "setup.cfg", "conftest.py", "pom.xml", "build.gradle", "build.gradle.kts",
	"settings.gradle", "settings.gradle.kts", "CMakeLists.txt", "Makefile", "Directory.Build.props",
}

// IsBuildFile tells if a file name is a build file of a language, such as go.mod or a .csproj file.
func IsBuildFile(name string) bool {
	b := filepath.Base(name)
	for _, f := range BuildFiles {
		if b == f {
			return true
		}
	}
	return strings.HasSuffix(b, ".csproj") || strings.HasSuffix(b, ".sln") || strings.HasSuffix(b, ".fsproj")
}
