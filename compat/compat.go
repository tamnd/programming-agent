// SPDX-License-Identifier: Apache-2.0

// Package compat reads the compat entries of a feature (spec 09.2), the probes that prove them (09.4), and the register of the implementations of each language.
package compat

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Entry is one compat entry. The key of the first entry is the slug of the feature. The keys of the other entries add a dot and a sub-feature, such as match-statement.class-patterns.
type Entry struct {
	Key         string            `json:"key"`
	Description string            `json:"description"`
	Spec        map[string]string `json:"spec,omitempty"`
	Source      map[string]string `json:"source,omitempty"`
	// Support maps an implementation column to its statements. Most columns have one statement. A column has more than one when the answer differs by edition or by flag.
	Support map[string][]Statement `json:"support"`
	// Columns are the keys of Support in the order of the file.
	Columns []string        `json:"-"`
	Status  map[string]bool `json:"status,omitempty"`
	Line    int             `json:"-"`
	// Unknown are the keys of the entry that the schema does not have.
	Unknown []string `json:"-"`
}

// Statement is one support statement of an implementation (spec 09.2).
type Statement struct {
	// Added is the version_added: a version, "false" for not supported, or "preview". It is empty when the field is missing, and "true" or "null" when the file has these values, so that rule K01 can find them.
	Added      string   `json:"version_added"`
	Removed    string   `json:"version_removed,omitempty"`
	Deprecated string   `json:"version_deprecated,omitempty"`
	Flags      []string `json:"flags,omitempty"`
	Partial    bool     `json:"partial_implementation,omitempty"`
	AltName    string   `json:"alternative_name,omitempty"`
	Edition    string   `json:"edition,omitempty"`
	Std        string   `json:"std,omitempty"`
	Note       string   `json:"notes,omitempty"`
	Probe      string   `json:"probe,omitempty"`
	Source     string   `json:"source,omitempty"`
	Line       int      `json:"-"`
	Unknown    []string `json:"-"`
}

// Supported tells if the statement names a version where the feature works.
func (s Statement) Supported() bool {
	return s.Added != "" && s.Added != "false" && s.Added != "true" && s.Added != "null"
}

// IsVersion tells if Added is a real version, not false or preview.
func (s Statement) IsVersion() bool { return s.Supported() && s.Added != "preview" }

// Read reads a compat file. The entries keep the order of the file.
func Read(path string) ([]Entry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(doc.Content) == 0 {
		return nil, nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s:%d: the file is not a map of entries", path, root.Line)
	}
	var out []Entry
	for i := 0; i+1 < len(root.Content); i += 2 {
		k, v := root.Content[i], root.Content[i+1]
		e, err := readEntry(k.Value, v)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %s: %w", path, k.Line, k.Value, err)
		}
		e.Line = k.Line
		out = append(out, e)
	}
	return out, nil
}

func readEntry(key string, v *yaml.Node) (Entry, error) {
	e := Entry{Key: key, Support: map[string][]Statement{}}
	if v.Kind != yaml.MappingNode {
		return e, errors.New("an entry is a map")
	}
	for i := 0; i+1 < len(v.Content); i += 2 {
		k, f := v.Content[i], v.Content[i+1]
		var err error
		switch k.Value {
		case "description":
			e.Description = f.Value
		case "spec":
			err = f.Decode(&e.Spec)
		case "source":
			err = f.Decode(&e.Source)
		case "status":
			err = f.Decode(&e.Status)
		case "support":
			if f.Kind != yaml.MappingNode {
				return e, errors.New("support is a map of implementations")
			}
			for j := 0; j+1 < len(f.Content); j += 2 {
				col, sv := f.Content[j], f.Content[j+1]
				items := []*yaml.Node{sv}
				if sv.Kind == yaml.SequenceNode {
					items = sv.Content
				}
				for _, it := range items {
					s, err := readStatement(it)
					if err != nil {
						return e, fmt.Errorf("support.%s: %w", col.Value, err)
					}
					e.Support[col.Value] = append(e.Support[col.Value], s)
				}
				e.Columns = append(e.Columns, col.Value)
			}
		default:
			e.Unknown = append(e.Unknown, k.Value)
		}
		if err != nil {
			return e, fmt.Errorf("%s: %w", k.Value, err)
		}
	}
	return e, nil
}

func readStatement(v *yaml.Node) (Statement, error) {
	s := Statement{Line: v.Line}
	if v.Kind != yaml.MappingNode {
		return s, errors.New("a support statement is a map")
	}
	for i := 0; i+1 < len(v.Content); i += 2 {
		k, f := v.Content[i], v.Content[i+1]
		var err error
		switch k.Value {
		case "version_added":
			s.Added = f.Value
			if f.Tag == "!!null" {
				s.Added = "null"
			}
		case "version_removed":
			s.Removed = f.Value
		case "version_deprecated":
			s.Deprecated = f.Value
		case "flags":
			if f.Kind == yaml.ScalarNode {
				s.Flags = []string{f.Value}
			} else {
				err = f.Decode(&s.Flags)
			}
		case "partial_implementation":
			err = f.Decode(&s.Partial)
		case "alternative_name":
			s.AltName = f.Value
		case "edition":
			s.Edition = f.Value
		case "std":
			s.Std = f.Value
		case "note":
			s.Note = f.Value
		case "probe":
			s.Probe = f.Value
		case "source":
			s.Source = f.Value
		default:
			s.Unknown = append(s.Unknown, k.Value)
		}
		if err != nil {
			return s, fmt.Errorf("%s: %w", k.Value, err)
		}
	}
	return s, nil
}

// Expect is what a probe must do (spec 09.4). It comes from the file probes/probes.yaml of the feature, keyed by the name of the probe file.
type Expect struct {
	// Output is the standard output that the probe must give where it passes. When it is empty, exit status 0 is enough.
	Output string `yaml:"output,omitempty" json:"output,omitempty"`
	// FailsWith is a regular expression that the output of the probe must match on the release before version_added (rule K05).
	FailsWith string `yaml:"fails_with,omitempty" json:"fails_with,omitempty"`
}

// ReadExpect reads probes/probes.yaml in dir. A missing file gives an empty map.
func ReadExpect(dir string) (map[string]Expect, error) {
	out := map[string]Expect{}
	b, err := os.ReadFile(filepath.Join(dir, "probes", "probes.yaml"))
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("%s: %w", filepath.Join(dir, "probes", "probes.yaml"), err)
	}
	return out, nil
}

// Impl is one implementation column of a language, from registers/implementations.yaml.
type Impl struct {
	ID   string `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
	// Standard is true for a column that is the language standard, such as c++20, not a program.
	Standard bool `yaml:"standard,omitempty" json:"standard,omitempty"`
	// Main is true for the implementations that the computed status counts (spec 09.5).
	Main bool `yaml:"main,omitempty" json:"main,omitempty"`
	// Image is the container image of a release, with {v} for the version, such as python:{v}-slim. Probe runs a probe in it.
	Image string `yaml:"image,omitempty" json:"image,omitempty"`
	// Run is the command that runs a probe, with {file} for the path of the probe in the container.
	Run []string `yaml:"run,omitempty" json:"run,omitempty"`
	// Releases are the releases of the implementation, oldest first. Rule K02 checks each version_added against this list, and probe takes the release before from it.
	Releases []string `yaml:"releases,omitempty" json:"releases,omitempty"`
}

// Register maps a language id to its implementation columns, in the order of the file.
type Register map[string][]Impl

// ReadRegister reads registers/implementations.yaml under root. A missing file gives an empty register.
func ReadRegister(root string) (Register, error) {
	out := Register{}
	p := filepath.Join(root, "registers", "implementations.yaml")
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("%s: %w", p, err)
	}
	return out, nil
}

// Impl finds the column of a language.
func (r Register) Impl(lang, id string) (Impl, bool) {
	for _, im := range r[lang] {
		if im.ID == id {
			return im, true
		}
	}
	return Impl{}, false
}

// Before returns the release just before v, and false when v is the first release or is not in the list.
func (im Impl) Before(v string) (string, bool) {
	i := slices.Index(im.Releases, v)
	if i <= 0 {
		return "", false
	}
	return im.Releases[i-1], true
}

// Known tells if v is a release in the list.
func (im Impl) Known(v string) bool { return slices.Contains(im.Releases, v) }

// Ref returns the image of release v.
func (im Impl) Ref(v string) string { return strings.ReplaceAll(im.Image, "{v}", v) }
