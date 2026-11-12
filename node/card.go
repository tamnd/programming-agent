// SPDX-License-Identifier: Apache-2.0

package node

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Card is one row of a card register, cards/<kind>/<owner>/<pp>.jsonl (spec 04.1 and 20.3). A card is a set of checked facts with sources. It has no page of its own: the site shows it on the collection page of its kind and owner.
type Card struct {
	Tag      string   `json:"tag"`
	Kind     string   `json:"kind"`
	Owner    string   `json:"owner"`
	Title    string   `json:"title,omitempty"`
	Name     string   `json:"name,omitempty"`
	Summary  string   `json:"summary,omitempty"`
	URL      string   `json:"url,omitempty"`
	Licence  string   `json:"licence,omitempty"`
	Since    string   `json:"since,omitempty"`
	Removed  string   `json:"removed,omitempty"`
	Features []string `json:"features,omitempty"`
	Concepts []string `json:"concepts,omitempty"`
	Tasks    []string `json:"tasks,omitempty"`

	// diagnostic
	Code             string   `json:"code,omitempty"`
	Severity         string   `json:"severity,omitempty"`
	MessageTemplates []string `json:"message_templates,omitempty"`
	// Type is the type of a proposal, such as standards-track, or lint for a diagnostic that is a lint rule.
	Type string `json:"type,omitempty"`

	// proposal
	ID      string   `json:"id,omitempty"`
	Status  string   `json:"status,omitempty"`
	Created string   `json:"created,omitempty"`
	Release string   `json:"release,omitempty"`
	Authors []string `json:"authors,omitempty"`

	// api
	Path       string `json:"path,omitempty"`
	Entity     string `json:"entity,omitempty"`
	Signature  string `json:"signature,omitempty"`
	Module     string `json:"module,omitempty"`
	Deprecated string `json:"deprecated,omitempty"`
	Doc        struct {
		URL string `json:"url,omitempty"`
	} `json:"doc,omitzero"`

	// weakness
	Abstraction string   `json:"abstraction,omitempty"`
	Parents     []string `json:"parents,omitempty"`

	// licence
	OSIApproved bool `json:"osi_approved,omitempty"`
	FSFLibre    bool `json:"fsf_libre,omitempty"`

	// language, tool, library and format: the facts of a PLDB record (spec 08.3).
	Class      string   `json:"class,omitempty"`
	Appeared   int      `json:"appeared,omitempty"`
	Creators   []string `json:"creators,omitempty"`
	Aka        []string `json:"aka,omitempty"`
	Extensions []string `json:"extensions,omitempty"`
	Version    string   `json:"version,omitempty"`
	Website    string   `json:"website,omitempty"`
	Wikipedia  string   `json:"wikipedia,omitempty"`

	// release: the version is the local slug, such as 3.12 or es2015.
	Date   string `json:"date,omitempty"`
	EOL    string `json:"eol,omitempty"`
	Latest string `json:"latest,omitempty"`
	LTS    bool   `json:"lts,omitempty"`

	Source struct {
		URL     string `json:"url,omitempty"`
		Licence string `json:"licence,omitempty"`
	} `json:"source,omitzero"`

	// Raw is the row as it is in the file. The site writes it as the JSON record of the card.
	Raw map[string]any `json:"-"`
	// File is the path of the shard, relative to the root, and Line is the line of the row, from 1.
	File string `json:"-"`
	Line int    `json:"-"`
}

// CardKinds are the kinds that are cards by default (spec 04.1), with the field that holds the local slug. A weakness card is a CWE entry, such as cwe-89, and a licence card is an SPDX licence, such as apache-2.0. A language, tool, library or format card is a PLDB record, such as cobol. These kinds have no owner, so the local slug is the full slug.
var CardKinds = map[string]string{"diagnostic": "code", "proposal": "id", "api": "path", "exercise": "id", "weakness": "id", "licence": "id", "release": "version", "language": "id", "tool": "id", "library": "id", "format": "id"}

// CardAreas are the pages that show the cards of the kinds that have no owner (spec 05.2 and 16.7). Most are sub-areas; languages/all, tools and libraries are collections with no area.
var CardAreas = map[string]string{"weakness": "weakness-catalogue", "licence": "licences", "language": "languages/all", "tool": "tools", "library": "libraries", "format": "data-formats"}

// Local is the slug of the card inside its owner: the code of a diagnostic, the id of a proposal, an exercise, a weakness or a licence, the path of an api entry, the version of a release.
func (c *Card) Local() string {
	switch CardKinds[c.Kind] {
	case "code":
		return c.Code
	case "id":
		return c.ID
	case "path":
		return c.Path
	case "version":
		return c.Version
	}
	return ""
}

// Slug is the full slug of the card, as a node of the same kind would have it: rust/errors/E0382.
func (c *Card) Slug() string { return FullSlug(c.Kind, c.Owner, c.Local()) }

// Heading is the text that names the card: the title, or the name, or the local slug.
func (c *Card) Heading() string {
	switch {
	case c.Title != "":
		return c.Title
	case c.Name != "":
		return c.Name
	}
	return c.Local()
}

// Group is the part of the collection that holds the card. An api card is in the group of its module, class or package: str for str.split, os.path for os.path.join, std::vec::Vec for std::vec::Vec::push. A card of another kind is in one group, "".
func (c *Card) Group() string {
	if c.Kind != "api" {
		return ""
	}
	p := c.Path
	for _, sep := range []string{"::", ".", "/", "#"} {
		if i := strings.LastIndex(p, sep); i > 0 {
			return p[:i]
		}
	}
	if c.Module != "" {
		return c.Module
	}
	return "builtins"
}

// Collection is the URL of the collection page that shows the card.
func (c *Card) Collection() string {
	switch c.Kind {
	case "diagnostic":
		if c.Type == "lint" {
			return "/" + c.Owner + "/lints/"
		}
		return "/" + c.Owner + "/errors/"
	case "proposal":
		return "/" + c.Owner + "/proposals/"
	case "release":
		return "/" + c.Owner + "/releases/"
	case "api":
		return "/" + c.Owner + "/api/" + c.Group() + "/"
	case "exercise":
		return "/exercises/" + c.Owner + "/"
	}
	if a, ok := CardAreas[c.Kind]; ok {
		return "/" + a + "/"
	}
	return "/" + c.Owner + "/"
}

// Anchor is the id of the card on its collection page.
func (c *Card) Anchor() string { return c.Local() }

// CardError is a row of a shard that could not be read.
type CardError struct {
	File string
	Line int
	Err  error
}

func (e *CardError) Error() string { return fmt.Sprintf("%s:%d: %v", e.File, e.Line, e.Err) }

// ReadCards reads each shard under cards/. A missing directory gives no cards. A row that does not parse is an error of the list, and the reading goes on.
func ReadCards(root string) ([]*Card, []error) {
	var cards []*Card
	var errs []error
	base := filepath.Join(root, "cards")
	err := filepath.WalkDir(base, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".jsonl") {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		sc := bufio.NewScanner(bytes.NewReader(b))
		sc.Buffer(make([]byte, 0, 64*1024), 4<<20)
		line := 0
		for sc.Scan() {
			line++
			text := bytes.TrimSpace(sc.Bytes())
			if len(text) == 0 {
				continue
			}
			c := &Card{File: rel, Line: line}
			if err := json.Unmarshal(text, c); err != nil {
				errs = append(errs, &CardError{rel, line, err})
				continue
			}
			if err := json.Unmarshal(text, &c.Raw); err != nil {
				errs = append(errs, &CardError{rel, line, err})
				continue
			}
			cards = append(cards, c)
		}
		return sc.Err()
	})
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		errs = append(errs, err)
	}
	sort.SliceStable(cards, func(i, j int) bool { return cards[i].Slug() < cards[j].Slug() })
	return cards, errs
}

// ShardOf is the path of the shard that must hold the card: cards/<kind>/<owner>/<pp>.jsonl, where <owner> is _ for a card with no owner and <pp> is the first two characters of the local slug in lower case (spec 20.3).
func ShardOf(c *Card) string {
	owner := c.Owner
	if owner == "" {
		owner = "_"
	}
	pp := strings.ToLower(c.Local())
	if len(pp) > 2 {
		pp = pp[:2]
	}
	pp = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, pp)
	return "cards/" + c.Kind + "/" + owner + "/" + pp + ".jsonl"
}
