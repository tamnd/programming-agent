// SPDX-License-Identifier: Apache-2.0

package check

import (
	"strings"
	"testing"
)

func TestCardRules(t *testing.T) {
	c := fixture(t, map[string]string{
		// R01: a row that does not parse, a kind that is not a card kind, a missing title, and a row in the wrong shard.
		"cards/diagnostic/rust/e0.jsonl": `{"tag":"9E0382","kind":"diagnostic","owner":"rust","code":"E0382","title":"A value was used after it was moved","source":{"url":"https://doc.rust-lang.org/error_codes/E0382.html"}}
{not json
{"tag":"NEW-x","kind":"feature","owner":"rust","code":"E0499"}
{"tag":"NEW-E0499","kind":"diagnostic","owner":"rust","code":"E0499","source":{"url":"https://doc.rust-lang.org/error_codes/E0499.html"}}
`,
		"cards/diagnostic/rust/xx.jsonl": `{"tag":"NEW-E0599","kind":"diagnostic","owner":"rust","code":"E0599","title":"No method of that name","source":{"url":"https://doc.rust-lang.org/error_codes/E0599.html"}}
`,
		// R02 and R05: no source, an owner that does not exist, and a feature that names nothing.
		"cards/api/nolang/st.jsonl": `{"tag":"NEW-str","kind":"api","owner":"nolang","path":"str.split","entity":"method","features":["python/no-such-feature"]}
`,
		// R05: a tool card can own a release, but a licence card cannot.
		"cards/release/vim-editor/9_.jsonl": `{"tag":"NEW-vim-9.1","kind":"release","owner":"vim-editor","version":"9.1","name":"Vim 9.1","source":{"url":"https://www.vim.org/"}}
`,
		"cards/release/mit/1_.jsonl": `{"tag":"NEW-mit-1.0","kind":"release","owner":"mit","version":"1.0","name":"MIT 1.0","source":{"url":"https://opensource.org/license/mit"}}
`,
		// R04: the same slug twice.
		"cards/api/python/st.jsonl": `{"tag":"9STRSP","kind":"api","owner":"python","path":"str.split","entity":"method","doc":{"url":"https://docs.python.org/3/library/stdtypes.html#str.split"}}
{"tag":"NEW-str2","kind":"api","owner":"python","path":"str.split","entity":"method","doc":{"url":"https://docs.python.org/3/library/stdtypes.html#str.split"}}
`,
	})
	fs := Run(c, Options{})
	want := map[string][]string{
		"R01": {"not valid JSON", `kind "feature"`, "needs the field title", "must be in cards/diagnostic/rust/e0.jsonl"},
		"R02": {"nolang/api/str.split has no source"},
		"R04": {"is also the card at cards/api/python/st.jsonl:1"},
		"R05": {`owner "nolang"`, `no node or card has the slug "python/no-such-feature"`},
	}
	for rule, subs := range want {
		got := strings.Join(fired(fs, rule), "\n")
		for _, s := range subs {
			if !strings.Contains(got, s) {
				t.Errorf("%s: no finding with %q in:\n%s", rule, s, got)
			}
		}
	}
	for _, f := range fired(fs, "R05") {
		if strings.Contains(f, "vim-editor") {
			t.Errorf("R05 fired on the release of a tool card: %s", f)
		}
	}
	if !strings.Contains(strings.Join(fired(fs, "R05"), "\n"), `owner "mit"`) {
		t.Errorf("R05 did not fire on the release of a licence card")
	}
	if got := fired(fs, "T02"); len(got) != 0 {
		t.Errorf("T02 fired on the cards: %v", got)
	}
}

func TestCardSlugsAndCollections(t *testing.T) {
	c := fixture(t, nil)
	cases := map[string]string{
		"rust/errors/E0382":         "/rust/errors/#E0382",
		"python/proposals/pep-634":  "/python/proposals/#pep-634",
		"python/api/str.split":      "/python/api/builtins/#str.split",
		"javascript/api/JSON.parse": "/javascript/api/JSON/#JSON.parse",
		"cwe-89":                    "/weakness-catalogue/#cwe-89",
		"mit":                       "/licences/#mit",
		"cobol":                     "/languages/all/#cobol",
		"python/3.12":               "/python/releases/#3.12",
		"vim-editor":                "/tools/#vim-editor",
	}
	for slug, want := range cases {
		k := c.CardBySlug[slug]
		if k == nil {
			t.Errorf("no card %s", slug)
			continue
		}
		if got := k.Collection() + "#" + k.Anchor(); got != want {
			t.Errorf("%s: got %s, want %s", slug, got, want)
		}
	}
	if len(c.CardErrors) != 0 {
		t.Errorf("errors in the test cards: %v", c.CardErrors)
	}
}
