package lang

import (
	"testing"

	"github.com/tamnd/programming-agent/node"
)

// The C test build links libm after the sources, so a test of code that uses math.h links on glibc.
func TestClangLinksLibm(t *testing.T) {
	im := node.Impl{Lang: "c", Files: []node.File{{Name: "m.c"}, {Name: "test_m.c", Test: true}}}
	steps, err := clang(nil, t.TempDir(), nil, im)
	if err != nil {
		t.Fatal(err)
	}
	build := steps[0]
	if build[len(build)-5] != "test_m.c" || build[len(build)-1] != "-lm" {
		t.Fatalf("want the sources, then -o test_main, then -lm; got %v", build)
	}
}
