package node

import (
	"regexp"
	"slices"
)

// Kind describes one node kind of spec 04.2: its names, the slug form and the sections of spec 06.3.
type Kind struct {
	ID     string
	One    string // the name of one item, such as "design pattern"
	Many   string // the plural, such as "design patterns"
	Owned  bool   // the URL has the owner first (spec 04.3)
	Infix  string // the path between the owner and the local slug, such as "errors" for a diagnostic
	Blurb  string // one sentence for the index page of the kind
	Weight int    // the order of the kind in lists, high first
	// Sections are the headings of the body in their order. Need is 2 for a section that W2 needs, 3 for W3 and 0 for an optional section.
	Sections []Section
}

// Section is one heading of a kind.
type Section struct {
	Name string
	Need int
}

func s2(names ...string) []Section { return secs(2, names...) }
func s3(names ...string) []Section { return secs(3, names...) }
func s0(names ...string) []Section { return secs(0, names...) }

func secs(need int, names ...string) []Section {
	out := make([]Section, len(names))
	for i, n := range names {
		out[i] = Section{Name: n, Need: need}
	}
	return out
}

func cat(parts ...[]Section) []Section {
	var out []Section
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// KindList is the list of the 22 kinds, in the order of spec 04.2.
var KindList = []Kind{
	{ID: "area", One: "area", Many: "areas", Weight: 5},
	{ID: "concept", One: "concept", Many: "concepts", Weight: 50,
		Blurb:    "A concept is an idea of programming that does not depend on one language. Each page shows the idea in each language where it exists.",
		Sections: cat(s2("Idea", "How it works", "In languages"), s3("Variations and related ideas", "Pitfalls", "Cost", "History"), s0("Theory", "Exercises"), s2("References"))},
	{ID: "definition", One: "definition", Many: "definitions", Weight: 30,
		Blurb:    "A definition gives the exact meaning of a term.",
		Sections: cat(s2("Definition"), s3("Examples and non-examples"), s0("Related definitions"), s2("References"))},
	{ID: "language", One: "language", Many: "languages", Weight: 48,
		Blurb:    "A language page is the front page of the guide to that language.",
		Sections: cat(s2("Overview", "A first program", "The language at a glance", "The guide"), s3("Versions and support", "Implementations", "Ecosystem", "History and influence"), s0("Coming from another language"), s2("References"))},
	{ID: "feature", One: "feature", Many: "features", Owned: true, Weight: 46,
		Blurb:    "A feature is a construct, a rule or a facility of one language. Each page gives the syntax, the rules and the versions.",
		Sections: cat(s2("Synopsis", "Semantics", "Examples"), s3("Versions", "Pitfalls", "Comparison"), s0("Implementation notes"), s2("References"))},
	{ID: "release", One: "release", Many: "releases", Owned: true, Weight: 40,
		Blurb:    "A release is a version of a language, a standard or a tool.",
		Sections: cat(s2("What is new", "Removed and deprecated"), s3("Migration", "Support dates"), s2("References"))},
	{ID: "task", One: "task", Many: "tasks", Weight: 44,
		Blurb:    "A task is a common job of a programmer. Each page has a contract and a solution in many languages, and each solution passes the same tests.",
		Sections: cat(s2("Task", "Solutions"), s3("Edge cases", "Pitfalls"), s0("Related tasks"), s2("References"))},
	{ID: "pattern", One: "pattern", Many: "patterns", Weight: 42,
		Blurb:    "A pattern is a named design that people use again and again. Anti-patterns and code smells are patterns too.",
		Sections: cat(s2("Intent", "Problem", "Solution", "Example"), s3("Consequences", "When not to use it", "How to fix it", "Variants and related patterns"), s0("Known uses"), s2("References"))},
	{ID: "technique", One: "technique", Many: "techniques", Weight: 41,
		Blurb:    "A technique is a named method with steps, such as a refactoring or a way to debug.",
		Sections: cat(s2("Purpose", "Steps", "Example"), s3("When to use", "Pitfalls", "Tools that do it"), s2("References"))},
	{ID: "practice", One: "practice", Many: "practices", Weight: 40,
		Blurb:    "A practice is a principle, a law, a guideline or a method of work.",
		Sections: cat(s2("Statement", "Meaning"), s3("Examples", "Criticism and limits", "History"), s2("References"))},
	{ID: "tool", One: "tool", Many: "tools", Weight: 39,
		Blurb:    "A tool is software that programmers use to make software, such as a compiler, a build tool or a debugger.",
		Sections: cat(s2("What it is", "Install and first use"), s3("Concepts of the tool", "Common tasks", "Configuration", "Alternatives", "History and versions"), s2("References"))},
	{ID: "library", One: "library", Many: "libraries", Weight: 38,
		Blurb:    "A library is code that a program uses, from a standard library to a framework.",
		Sections: cat(s2("What it is", "Install and first use"), s3("Concepts of the tool", "Common tasks", "Configuration", "Alternatives", "History and versions"), s2("References"))},
	{ID: "format", One: "format", Many: "formats", Weight: 37,
		Blurb:    "A format is a way to write data, such as JSON or UTF-8.",
		Sections: cat(s2("What it is", "Syntax or structure", "Examples"), s3("Versions", "Pitfalls", "Tools and libraries"), s2("References"))},
	{ID: "standard", One: "standard", Many: "standards", Weight: 36,
		Blurb:    "A standard is a specification document, such as ISO C or ECMA-262.",
		Sections: cat(s2("What it is", "Syntax or structure", "Examples"), s3("Versions", "Pitfalls", "Tools and libraries"), s2("References"))},
	{ID: "weakness", One: "weakness", Many: "weaknesses", Weight: 35,
		Blurb:    "A weakness is a class of defect, such as a security weakness of CWE or a common bug.",
		Sections: cat(s2("Description", "Vulnerable code", "Fixed code"), s3("How to detect it", "How to prevent it", "Real cases"), s2("References"))},
	{ID: "diagnostic", One: "diagnostic", Many: "diagnostics", Owned: true, Infix: "errors", Weight: 34,
		Blurb:    "A diagnostic is an error, a warning or a lint rule of a language or a tool.",
		Sections: cat(s2("Message", "Cause", "Example", "Fix"), s3("Related diagnostics"), s2("References"))},
	{ID: "proposal", One: "proposal", Many: "proposals", Owned: true, Infix: "proposals", Weight: 33,
		Blurb:    "A proposal is a change request of a language or a tool, such as a PEP or a JEP.",
		Sections: cat(s2("Summary"), s3("Motivation", "What changed", "History"), s2("References"))},
	{ID: "api", One: "API entry", Many: "API entries", Owned: true, Infix: "api", Weight: 32,
		Blurb:    "An API entry is a module, a type, a function or a constant of a standard library or a platform.",
		Sections: cat(s2("Signature", "Semantics", "Examples"), s3("Versions", "Pitfalls"), s2("References"))},
	{ID: "work", One: "work", Many: "works", Weight: 20,
		Blurb:    "A work is a book, a paper, an essay or a talk that is part of the canon of programming.",
		Sections: cat(s2("Summary"), s3("Contents", "Influence"), s2("References"))},
	{ID: "person", One: "person", Many: "people", Weight: 15,
		Blurb:    "A person made or named items of programming.",
		Sections: cat(s2("Biography"), s3("Work"), s2("References"))},
	{ID: "organization", One: "organization", Many: "organizations", Weight: 14,
		Blurb:    "An organization is a company, a foundation, a committee or a standards body.",
		Sections: cat(s2("About"), s3("Work"), s2("References"))},
	{ID: "exercise", One: "exercise", Many: "exercises", Owned: true, Weight: 10,
		Blurb: "An exercise is a practice problem on a public site."},
}

// Kinds are the ids of the 22 kinds of spec 04.2.
var Kinds = func() []string {
	out := make([]string, len(KindList))
	for i, k := range KindList {
		out[i] = k.ID
	}
	return out
}()

// KindOf returns the description of a kind, or nil.
func KindOf(id string) *Kind {
	for i := range KindList {
		if KindList[i].ID == id {
			return &KindList[i]
		}
	}
	return nil
}

// IsKind tells if k is one of the 22 kinds.
func IsKind(k string) bool {
	return slices.Contains(Kinds, k)
}

// IsOwned tells if the URL of a kind has the owner first.
func IsOwned(k string) bool {
	if d := KindOf(k); d != nil {
		return d.Owned
	}
	return false
}

// The forms of a local slug. A feature slug is kebab-case. A release slug is a version. A diagnostic slug is a code. An api slug is a path in the library.
var (
	KebabRE   = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	VersionRE = regexp.MustCompile(`^[0-9][0-9a-z]*(?:[.-][0-9a-z]+)*$`)
	CodeRE    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)
	APIRE     = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(?:[.:/][A-Za-z_][A-Za-z0-9_]*)*$`)
)

// LocalOK tells if a local slug has the form of its kind.
func LocalOK(kind, local string) bool {
	if len(local) == 0 || len(local) > 80 {
		return false
	}
	switch kind {
	case "release":
		return VersionRE.MatchString(local)
	case "diagnostic":
		return CodeRE.MatchString(local)
	case "api":
		return APIRE.MatchString(local)
	}
	return KebabRE.MatchString(local) && len(local) <= 60
}

// FullSlug joins the owner, the infix of the kind and the local slug. The result is the path of the URL without the slashes at the ends.
func FullSlug(kind, owner, local string) string {
	if owner == "" || !IsOwned(kind) {
		return local
	}
	if kind == "exercise" {
		return "exercises/" + owner + "/" + local
	}
	if in := KindOf(kind).Infix; in != "" {
		return owner + "/" + in + "/" + local
	}
	return owner + "/" + local
}
