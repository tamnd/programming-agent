// SPDX-License-Identifier: Apache-2.0

package site

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tamnd/programming-agent/node"
)

// cardsPerPage is the most cards that one collection page shows (spec 16.7).
const cardsPerPage = 200

// smallModule is the smallest API module that gets its own collection page. The cards of the smaller modules of an owner go on the collection page otherModules, so that the number of files of the site stays under the limit of the host. Each module still has a tile on the API index of its owner.
const (
	smallModule  = 20
	otherModules = "other-modules"
)

// CardView is one card on a collection page.
type CardView struct {
	Anchor  string
	Heading string
	Local   string
	Summary string
	Kind    string
	Code    string // the signature of an api entry, shown as code
	Facts   []Fact
	Links   []Link
	Sources []Link
	Record  string // the URL of the JSON record
	Help    string // the URL of the page that tells how to write the full page
}

// collectionNames are the headings of the card collections of an owner.
var collectionNames = map[string]string{
	"errors":    "Errors and diagnostics",
	"lints":     "Lints",
	"proposals": "Proposals",
	"releases":  "Releases",
	"api":       "Standard library reference",
}

// cardIndex groups the cards by the URL of their collection, and gives each card the URL of the page and the anchor that shows it. A collection with more than cardsPerPage cards has more pages: /<collection>/page/2/ and on.
func (b *builder) cardIndex() {
	b.cardsAt = map[string][]*node.Card{}
	b.cardHref = map[string]string{}
	b.cardDone = map[string]bool{}
	b.cardColl = map[string]string{}
	for _, k := range b.c.Cards {
		if k.Local() == "" || b.c.BySlug[k.Slug()] != nil {
			continue
		}
		u := k.Collection()
		b.cardsAt[u] = append(b.cardsAt[u], k)
	}
	for u, ks := range b.cardsAt {
		if k0 := ks[0]; k0.Kind == "api" && len(ks) < smallModule {
			o := "/" + k0.Owner + "/api/" + otherModules + "/"
			if o != u {
				b.cardsAt[o] = append(b.cardsAt[o], ks...)
				delete(b.cardsAt, u)
			}
		}
	}
	for u, ks := range b.cardsAt {
		if strings.HasSuffix(u, "/"+otherModules+"/") {
			// The cards of one module stay together, in the order of the module names.
			sort.SliceStable(ks, func(i, j int) bool {
				gi, gj := strings.ToLower(ks[i].Group()), strings.ToLower(ks[j].Group())
				if gi != gj {
					return naturalLess(gi, gj)
				}
				if ks[i].Group() != ks[j].Group() {
					return ks[i].Group() < ks[j].Group()
				}
				return cardLess(ks[i], ks[j])
			})
		} else {
			sort.SliceStable(ks, func(i, j int) bool { return cardLess(ks[i], ks[j]) })
		}
		for i, k := range ks {
			b.cardHref[k.Slug()] = pageURL(u, i/cardsPerPage) + "#" + k.Anchor()
			b.cardColl[k.Slug()] = u
		}
	}
}

// cardLess sorts cards by their local slug, with the numbers in order: E0382 before E0499, pep-8 before pep-634.
func cardLess(a, b *node.Card) bool {
	if a.Kind == "release" && b.Kind == "release" && a.Date != b.Date {
		// Releases are newest first.
		return a.Date > b.Date
	}
	x, y := a.Local(), b.Local()
	if strings.EqualFold(x, y) {
		return x < y
	}
	return naturalLess(strings.ToLower(x), strings.ToLower(y))
}

func naturalLess(x, y string) bool {
	for x != "" && y != "" {
		dx, dy := digits(x), digits(y)
		if dx > 0 && dy > 0 {
			nx, ny := strings.TrimLeft(x[:dx], "0"), strings.TrimLeft(y[:dy], "0")
			if len(nx) != len(ny) {
				return len(nx) < len(ny)
			}
			if nx != ny {
				return nx < ny
			}
			x, y = x[dx:], y[dy:]
			continue
		}
		if x[0] != y[0] {
			return x[0] < y[0]
		}
		x, y = x[1:], y[1:]
	}
	return len(x) < len(y)
}

func digits(s string) int {
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	return i
}

func pageURL(u string, page int) string {
	if page == 0 {
		return u
	}
	return fmt.Sprintf("%spage/%d/", u, page+1)
}

// cardViews are the cards of one collection page.
func (b *builder) cardViews(ks []*node.Card) []CardView {
	out := make([]CardView, 0, len(ks))
	for _, k := range ks {
		v := CardView{Anchor: k.Anchor(), Heading: k.Heading(), Local: k.Local(), Summary: k.Summary, Kind: k.Kind, Code: k.Signature}
		v.Record = recordPage + "#" + b.cardColl[k.Slug()]
		v.Help = "/contribute/#promote-a-card"
		add := func(key, val string, code bool) {
			if val != "" {
				v.Facts = append(v.Facts, Fact{Key: key, Value: val, Code: code})
			}
		}
		switch k.Kind {
		case "diagnostic":
			add("Code", k.Code, true)
			add("Severity", k.Severity, false)
			for _, m := range k.MessageTemplates {
				add("Message", m, true)
			}
		case "proposal":
			add("Status", k.Status, false)
			add("Type", k.Type, false)
			add("Created", k.Created, false)
			add("Release", k.Release, false)
			add("Authors", strings.Join(k.Authors, ", "), false)
		case "api":
			add("Kind", k.Entity, false)
			add("Module", k.Module, true)
			add("Deprecated", k.Deprecated, false)
		case "weakness":
			add("Abstraction", k.Abstraction, false)
			add("Status", k.Status, false)
		case "licence":
			add("SPDX id", k.Name, true)
			add("OSI approved", yes(k.OSIApproved), false)
			add("FSF free", yes(k.FSFLibre), false)
			add("Deprecated", k.Deprecated, false)
		case "release":
			add("Released", k.Date, false)
			add("Status", k.Status, false)
			if k.LTS {
				add("Long-term support", "yes", false)
			}
			add("Latest patch", k.Latest, true)
			add("End of life", k.EOL, false)
		case "person":
			add("Description", k.Class, false)
			add("Born", k.Born, false)
			add("Died", k.Died, false)
			add("Also known as", strings.Join(k.Aka, ", "), false)
		case "organization":
			add("Description", k.Class, false)
			add("Founded", k.Date, false)
			add("Also known as", strings.Join(k.Aka, ", "), false)
		case "work":
			add("Type", k.Class, false)
			add("Authors", strings.Join(k.Authors, ", "), false)
			add("Published", k.Date, false)
		case "language", "tool", "library", "format":
			add("Class", k.Class, false)
			if k.Appeared > 0 {
				add("Appeared", strconv.Itoa(k.Appeared), false)
			}
			add("Creators", strings.Join(k.Creators, ", "), false)
			add("Also known as", strings.Join(k.Aka, ", "), false)
			add("File extensions", strings.Join(k.Extensions, " "), true)
			add("Latest version", k.Version, true)
		}
		add("Since", k.Since, false)
		add("Removed", k.Removed, false)
		for _, s := range append(append(append(append([]string(nil), k.Parents...), k.Features...), k.Concepts...), k.Tasks...) {
			v.Links = append(v.Links, b.link(s))
		}
		if node.OwnsReleases(k) {
			if rs := b.cardsAt["/"+k.ID+"/releases/"]; len(rs) > 0 {
				v.Links = append(v.Links, Link{Name: plural(len(rs), "release", "releases"), URL: "/" + k.ID + "/releases/"})
			}
		}
		wd := ""
		if k.Wikidata != "" {
			wd = "https://www.wikidata.org/wiki/" + k.Wikidata
		}
		for _, s := range []string{k.Source.URL, k.URL, k.Doc.URL, k.Website, k.Wikipedia, wd} {
			if s != "" && !hasURL(v.Sources, s) {
				v.Sources = append(v.Sources, Link{Name: hostOf(s), URL: s})
			}
		}
		out = append(out, v)
	}
	return out
}

func hasURL(ls []Link, u string) bool {
	for _, l := range ls {
		if l.URL == u {
			return true
		}
	}
	return false
}

func hostOf(u string) string {
	s := strings.TrimPrefix(strings.TrimPrefix(u, "https://"), "http://")
	if i := strings.IndexByte(s, '/'); i > 0 {
		s = s[:i]
	}
	return s
}

func yes(v bool) string {
	if v {
		return "yes"
	}
	return ""
}

// cardLead is the lead of a collection page.
func (b *builder) cardLead(owner, part string, n int) (string, string) {
	if owner == "" {
		// The cards of a kind with no owner are in a sub-area, such as the weakness catalogue.
		name, lead := humanize(part), ""
		if a := b.area(part); a != nil {
			name, lead = a.Name, a.Summary
		} else if t, ok := cardCollections[part]; ok {
			name, lead = t[0], t[1]
		}
		return name, strings.TrimSpace(fmt.Sprintf("%s This page has %s.", lead, plural(n, "card", "cards")))
	}
	name := owner
	if l := b.c.LangByID[owner]; l != nil {
		name = l.Name
	} else if m := b.c.BySlug[owner]; m != nil {
		name = m.Name
	}
	switch part {
	case "errors":
		return name + " errors and diagnostics", fmt.Sprintf("The error codes and warnings of %s, with the message that each one prints and what causes it. This page has %s.", name, plural(n, "card", "cards"))
	case "lints":
		return name + " lints", fmt.Sprintf("The lint rules of %s, with what each one finds. This page has %s.", name, plural(n, "card", "cards"))
	case "releases":
		return name + " releases", fmt.Sprintf("Each release of %s since the first, newest first, with its date, its status, its latest patch and its end of life. This page has %s.", name, plural(n, "card", "cards"))
	case "proposals":
		return name + " proposals", fmt.Sprintf("The proposals that changed %s, with their status and the release that took them. This page has %s.", name, plural(n, "card", "cards"))
	}
	return name, fmt.Sprintf("This page has %s.", plural(n, "card", "cards"))
}

// cardCollections are the names and leads of the card collections that are not areas.
var cardCollections = map[string][2]string{
	"all":           {"All languages", "Each programming language, markup language, query language and notation in the Programming Language Database (PLDB), with the year it appeared, its creators and its file extensions. The languages that have a full guide on this site are on the languages page, not on this list."},
	"tools":         {"Tools", "The editors, compilers, package managers, databases and other programs in the Programming Language Database (PLDB), with the year each one appeared and its creators."},
	"libraries":     {"Libraries and frameworks", "The libraries and frameworks in the Programming Language Database (PLDB), with the year each one appeared and its creators."},
	"people":        {"People", "The people who made or named items of programming, with the years of their life and what they are known for."},
	"organizations": {"Organizations", "The companies, foundations, committees and standards bodies that made or named items of programming."},
	"works":         {"Works", "The books, papers, essays and talks of programming, with their authors and the year of publication."},
}

// cardParents are the pages that show a tile for a card collection that is not an area.
var cardParents = map[string]string{"/languages/all/": "/languages/", "/tools/": "/tools-and-workflow/", "/libraries/": "/tools-and-workflow/", "/people/": "/history-and-culture/", "/organizations/": "/history-and-culture/", "/works/": "/history-and-culture/"}

// collectionTiles are the tiles of the card collections whose parent is the page u.
func (b *builder) collectionTiles(u string, hue int) []Tile {
	var urls []string
	for c, p := range cardParents {
		if p == u && len(b.cardsAt[c]) > 0 {
			urls = append(urls, c)
		}
	}
	sort.Strings(urls)
	var out []Tile
	for _, c := range urls {
		ks := b.cardsAt[c]
		parts := strings.Split(strings.Trim(c, "/"), "/")
		name, lead := b.cardLead("", parts[len(parts)-1], len(ks))
		out = append(out, Tile{Name: name, URL: c, Summary: lead, Hue: hue, Count: len(ks), Unit: "card"})
	}
	return out
}

// cardPages writes each collection page that no area page wrote, the index /<owner>/api/ of each owner, and the JSON records of the cards of each collection in the record packs (packRecords).
func (b *builder) cardPages() error {
	urls := make([]string, 0, len(b.cardsAt))
	for u := range b.cardsAt {
		urls = append(urls, u)
	}
	sort.Strings(urls)
	apiGroups := map[string][]string{}
	for _, u := range urls {
		ks := b.cardsAt[u]
		k0 := ks[0]
		if k0.Kind == "api" {
			apiGroups[k0.Owner] = append(apiGroups[k0.Owner], u)
		}
		if b.c.BySlug[strings.Trim(u, "/")] != nil {
			continue
		}
		part := strings.Split(strings.Trim(u, "/"), "/")
		var heading, lead string
		crumbs := b.ownerCrumbs(k0.Owner)
		if p, ok := cardParents[u]; ok && k0.Owner == "" {
			if a := b.area(strings.Trim(p, "/")); a != nil {
				crumbs = []Link{{Name: a.Name, URL: p}}
			}
		}
		if k0.Kind == "api" && strings.HasSuffix(u, "/"+otherModules+"/") {
			heading = "Other modules"
			lead = fmt.Sprintf("The entries of the %s of the %s reference that have fewer than %d entries each, in the order of the module names. This page has %s.", plural(len(groupsOf(ks)), "module", "modules"), b.ownerName(k0.Owner), smallModule, plural(len(ks), "card", "cards"))
			crumbs = append(crumbs, Link{Name: collectionNames["api"], URL: "/" + k0.Owner + "/api/"})
		} else if k0.Kind == "api" {
			heading = k0.Group()
			lead = fmt.Sprintf("The entries of %s in the %s reference. This page has %s.", k0.Group(), b.ownerName(k0.Owner), plural(len(ks), "card", "cards"))
			crumbs = append(crumbs, Link{Name: collectionNames["api"], URL: "/" + k0.Owner + "/api/"})
		} else {
			heading, lead = b.cardLead(k0.Owner, part[len(part)-1], len(ks))
		}
		first := 0
		if b.cardDone[u] {
			// An area page is the first page.
			first = 1
		}
		if err := b.cardList(u, heading, lead, crumbs, ks, first); err != nil {
			return err
		}
	}
	owners := make([]string, 0, len(apiGroups))
	for o := range apiGroups {
		owners = append(owners, o)
	}
	sort.Strings(owners)
	for _, o := range owners {
		u := "/" + o + "/api/"
		if b.c.BySlug[o+"/api"] != nil {
			continue
		}
		name := b.ownerName(o)
		p := ListPage{Base: b.base(name+" "+strings.ToLower(collectionNames["api"]), "The modules, types and functions of the "+name+" standard library.", u)}
		p.Section, p.Heading = "languages", name+" "+strings.ToLower(collectionNames["api"])
		p.Lead = fmt.Sprintf("The modules, types and functions of the %s standard library. Each entry gives the signature, the release that added it and a link to the official documentation.", name)
		p.Crumbs = b.ownerCrumbs(o)
		p.Hue = b.hueOf("languages")
		for _, gu := range apiGroups[o] {
			ks := b.cardsAt[gu]
			if !strings.HasSuffix(gu, "/"+otherModules+"/") {
				p.Tiles = append(p.Tiles, Tile{Name: ks[0].Group(), URL: gu, Hue: p.Hue, Count: len(ks), Unit: "card"})
				continue
			}
			// A small module has a tile that points at its first card on the page of the other modules.
			for _, g := range groupsOf(ks) {
				p.Tiles = append(p.Tiles, Tile{Name: g[0].Group(), URL: b.cardHref[g[0].Slug()], Hue: p.Hue, Count: len(g), Unit: "card"})
			}
		}
		sort.SliceStable(p.Tiles, func(i, j int) bool {
			return naturalLess(strings.ToLower(p.Tiles[i].Name), strings.ToLower(p.Tiles[j].Name))
		})
		if err := b.page("list", u, p); err != nil {
			return err
		}
	}
	// The records of the cards of a collection are one JSON array. The arrays go into a few packs, so that the number of files stays far under the limit of the host.
	recs := make([]record, 0, len(urls))
	for _, u := range urls {
		rs := make([]map[string]any, 0, len(b.cardsAt[u]))
		for _, k := range b.cardsAt[u] {
			rec := map[string]any{}
			for key, val := range k.Raw {
				rec[key] = val
			}
			rec["url"] = b.opt.BaseURL + b.cardHref[k.Slug()]
			if k.URL != "" {
				rec["source_url"] = k.URL
			}
			rs = append(rs, rec)
		}
		data, err := json.Marshal(rs)
		if err != nil {
			return err
		}
		recs = append(recs, record{URL: u, Data: data})
	}
	return b.packRecords(recs)
}

// groupsOf splits the sorted cards of the page of the other modules into the cards of each module.
func groupsOf(ks []*node.Card) [][]*node.Card {
	var out [][]*node.Card
	for i, k := range ks {
		if i == 0 || k.Group() != ks[i-1].Group() {
			out = append(out, nil)
		}
		out[len(out)-1] = append(out[len(out)-1], k)
	}
	return out
}

// cardList writes the pages of one collection, from the page first (0 is the first page).
func (b *builder) cardList(u, heading, lead string, crumbs []Link, ks []*node.Card, first int) error {
	pages := (len(ks) + cardsPerPage - 1) / cardsPerPage
	for i := first; i < pages; i++ {
		end := min((i+1)*cardsPerPage, len(ks))
		pu := pageURL(u, i)
		h := heading
		if pages > 1 {
			h = fmt.Sprintf("%s, page %d of %d", heading, i+1, pages)
		}
		p := ListPage{Base: b.base(h, lead, pu)}
		p.Section, p.Heading, p.Lead = "languages", h, lead
		p.Crumbs = crumbs
		p.Hue = b.hueOf("languages")
		p.Cards = b.cardViews(ks[i*cardsPerPage : end])
		p.Pager = pager(u, i, pages)
		if err := b.page("list", pu, p); err != nil {
			return err
		}
	}
	return nil
}

func pager(u string, i, pages int) []Link {
	if pages < 2 {
		return nil
	}
	var out []Link
	for j := 0; j < pages; j++ {
		l := Link{Name: fmt.Sprint(j + 1), URL: pageURL(u, j)}
		if j == i {
			l.Note = "current"
		}
		out = append(out, l)
	}
	return out
}

func (b *builder) ownerName(o string) string {
	if l := b.c.LangByID[o]; l != nil {
		return l.Name
	}
	if n := b.c.BySlug[o]; n != nil {
		return n.Name
	}
	if k := b.c.CardBySlug[o]; node.OwnsReleases(k) {
		return k.Heading()
	}
	return humanize(o)
}

func (b *builder) ownerCrumbs(o string) []Link {
	if b.c.LangByID[o] != nil {
		return []Link{{Name: "Languages", URL: "/languages/"}, {Name: b.ownerName(o), URL: "/" + o + "/"}}
	}
	if n := b.c.BySlug[o]; n != nil {
		return []Link{{Name: n.Name, URL: "/" + n.Slug + "/"}}
	}
	if k := b.c.CardBySlug[o]; node.OwnsReleases(k) && b.cardHref[o] != "" {
		return []Link{{Name: k.Heading(), URL: b.cardHref[o]}}
	}
	return nil
}

// cardTiles are the tiles of the card collections of an owner that no sub-area of the language guide shows, for the language page.
func (b *builder) cardTiles(owner string, shown map[string]bool) []Tile {
	hue := b.hueOf("languages")
	count := map[string]int{}
	for u, ks := range b.cardsAt {
		k0 := ks[0]
		if k0.Owner != owner {
			continue
		}
		if k0.Kind == "api" {
			u = "/" + owner + "/api/"
		}
		count[u] += len(ks)
	}
	urls := make([]string, 0, len(count))
	for u := range count {
		if !shown[u] {
			urls = append(urls, u)
		}
	}
	sort.Strings(urls)
	var out []Tile
	for _, u := range urls {
		part := strings.Split(strings.Trim(u, "/"), "/")
		name := collectionNames[part[len(part)-1]]
		if name == "" {
			name = humanize(part[len(part)-1])
		}
		out = append(out, Tile{Name: name, URL: u, Hue: hue, Count: count[u], Unit: "card"})
	}
	return out
}
