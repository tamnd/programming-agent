// SPDX-License-Identifier: Apache-2.0

package check

import (
	"strings"
	"testing"
)

func TestX05(t *testing.T) {
	c := fixture(t, nil)
	n := c.BySlug["python/closures-and-nonlocal"]
	if n == nil {
		t.Fatal("no feature python/closures-and-nonlocal in the test corpus")
	}
	n.Body = "## Semantics\n\n```rust fails expect=error:E0382 id=known\nfn main() {}\n```\n\n```rust fails expect=error:E9999 id=unknown\nfn main() {}\n```\n\n```rust fails expect=compile id=plain\nfn main() {}\n```\n"
	got := fired(Run(c, Options{Main: true, Only: map[string]bool{n.Slug: true}}), "X05")
	if len(got) != 1 || !strings.Contains(got[0], "expect=error:E9999 names no diagnostic: add a card for rust/errors/E9999") {
		t.Fatalf("X05 = %q, want one finding for E9999", got)
	}
}
