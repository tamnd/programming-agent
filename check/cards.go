// SPDX-License-Identifier: Apache-2.0

package check

import (
	"errors"
	"regexp"

	"github.com/tamnd/programming-agent/node"
)

// cardProvRE is a provisional id of a card. A card slug can hold capitals and dots, as in NEW-E0382 or NEW-str.split.
var cardProvRE = regexp.MustCompile(`^NEW-[A-Za-z0-9._:/-]+$`)

func init() {
	Rules["R01"] = "each card row is valid JSON for its kind, has the fields of its kind, and is in the shard of its kind, owner and slug"
	Rules["R02"] = "each card names its source with a URL"
	Rules["R04"] = "no two cards have the same slug, and no card has the slug of a node; prog promote removes the row"
	Rules["R05"] = "the owner of each card is a language or a node, and its features and concepts name nodes"
}

// cardFields are the fields that each kind of card must have, after the slug field.
var cardFields = map[string][]string{
	"diagnostic": {"title"},
	"proposal":   {"title", "status"},
	"api":        {"entity"},
	"exercise":   {"title"},
}

// cards runs the rules R01 to R05 on the card registers.
func (r *run) cards() {
	for _, err := range r.c.CardErrors {
		var e *node.CardError
		if errors.As(err, &e) {
			r.add("R01", Block, e.File, e.Line, "the row is not valid JSON for a card: %v", e.Err)
		} else {
			r.add("R01", Block, "cards", 0, "%v", err)
		}
	}
	reg := map[string]string{}
	for _, t := range r.c.Tags {
		reg[t.Tag] = t.Slug
	}
	seen := map[string]*node.Card{}
	for _, k := range r.c.Cards {
		slugField, ok := node.CardKinds[k.Kind]
		if !ok {
			r.add("R01", Block, k.File, k.Line, "the kind %q is not a card kind", k.Kind)
			continue
		}
		if k.Local() == "" {
			r.add("R01", Block, k.File, k.Line, "a %s card needs the field %s", k.Kind, slugField)
			continue
		}
		for _, f := range cardFields[k.Kind] {
			if v, _ := k.Raw[f].(string); v == "" {
				r.add("R01", Block, k.File, k.Line, "the %s card %s needs the field %s", k.Kind, k.Local(), f)
			}
		}
		if k.Owner == "" && node.IsOwned(k.Kind) {
			r.add("R01", Block, k.File, k.Line, "the %s card %s needs an owner", k.Kind, k.Local())
		}
		if want := node.ShardOf(k); want != k.File {
			r.add("R01", Block, k.File, k.Line, "the card %s must be in %s", k.Slug(), want)
		}
		switch {
		case k.Tag == "":
			r.add("T02", Block, k.File, k.Line, "the card %s has no tag", k.Slug())
		case !cardProvRE.MatchString(k.Tag):
			if s, ok := reg[k.Tag]; !ok {
				r.add("T02", Block, k.File, k.Line, "the tag %s is not in registers/tags.reg", k.Tag)
			} else if s != k.Slug() {
				r.add("T02", Block, k.File, k.Line, "registers/tags.reg gives the tag %s to %q, not %q", k.Tag, s, k.Slug())
			}
		}
		if k.Source.URL == "" && k.URL == "" && k.Doc.URL == "" {
			r.add("R02", Block, k.File, k.Line, "the card %s has no source: give source.url, url or doc.url", k.Slug())
		}
		if prev, dup := seen[k.Slug()]; dup {
			r.add("R04", Block, k.File, k.Line, "the slug %s is also the card at %s:%d", k.Slug(), prev.File, prev.Line)
		} else {
			seen[k.Slug()] = k
		}
		if n := r.c.BySlug[k.Slug()]; n != nil {
			r.add("R04", Block, k.File, k.Line, "the node %s has the slug %s; remove the card row, as prog promote does", n.Rel, k.Slug())
		}
		if k.Owner != "" && k.Kind != "exercise" && r.c.LangByID[k.Owner] == nil && r.c.BySlug[k.Owner] == nil {
			r.add("R05", Block, k.File, k.Line, "the owner %q is not a language of the register or a node", k.Owner)
		}
		r.cardLinks(k, "features", k.Features)
		r.cardLinks(k, "concepts", k.Concepts)
	}
}

func (r *run) cardLinks(k *node.Card, field string, slugs []string) {
	for _, s := range slugs {
		if r.c.BySlug[s] == nil && r.c.CardBySlug[s] == nil {
			r.add("R05", Warn, k.File, k.Line, "%s of %s: no node or card has the slug %q", field, k.Slug(), s)
		}
	}
}

// cardURL reports whether x, a path with no slashes at the ends, is the URL of a card or of a card collection.
func (r *run) cardURL(x string) bool {
	if r.c.CardBySlug[x] != nil {
		return true
	}
	for _, k := range r.c.Cards {
		if k.Collection() == "/"+x+"/" || (k.Kind == "api" && x == k.Owner+"/api") {
			return true
		}
	}
	return false
}
