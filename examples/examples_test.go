// SPDX-License-Identifier: Apache-2.0

package examples

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tamnd/programming-agent/lang"
)

const page = "# Title\n\n" +
	"```python run id=hello\nprint('hello')\n```\n\n" +
	"```output for=hello\nhello\n```\n\n" +
	"````markdown\n```python\nnot a block\n```\n````\n\n" +
	"```py fails expect=exception:ZeroDivisionError\nprint(1 / 0)\n```\n\n" +
	"```python norun reason=\"fragment\" group=g\ndef f():\n    return 41\n```\n\n" +
	"```python run group=g\nprint(f() + 1)\n```\n\n" +
	"```output\n42\n```\n\n" +
	"```text\nnot code\n```\n"

func TestParse(t *testing.T) {
	bs := Parse(page, 1)
	if len(bs) != 8 {
		t.Fatalf("got %d blocks, want 8", len(bs))
	}
	b := bs[0]
	if b.Lang != "python" || !b.Has("run") || b.Attrs["id"] != "hello" || b.Code != "print('hello')\n" || b.Line != 3 {
		t.Errorf("first block: %+v", b)
	}
	if bs[2].Lang != "markdown" || !strings.Contains(bs[2].Code, "not a block") {
		t.Errorf("a longer fence must hold a shorter one: %+v", bs[2])
	}
	if bs[3].Lang != "python" || bs[3].Attrs["expect"] != "exception:ZeroDivisionError" {
		t.Errorf("the alias py: %+v", bs[3])
	}
	if bs[4].Attrs["reason"] != "fragment" {
		t.Errorf("a quoted value: %+v", bs[4].Attrs)
	}
	if got := page[bs[1].Start:bs[1].End]; got != "hello\n" {
		t.Errorf("offsets give %q", got)
	}
}

func TestCollect(t *testing.T) {
	cases, issues := Collect("p.md", page, 1)
	if len(issues) != 0 {
		t.Fatalf("issues: %v", issues)
	}
	if len(cases) != 3 {
		t.Fatalf("got %d cases, want 3", len(cases))
	}
	if cases[0].ID() != "hello" || len(cases[0].Outputs) != 1 {
		t.Errorf("case 1: %+v", cases[0])
	}
	g := cases[2]
	if len(g.Parts) != 2 || len(g.Outputs) != 1 || g.Outputs[0].Code != "42\n" {
		t.Errorf("group: parts %d outputs %d", len(g.Parts), len(g.Outputs))
	}
	if f := g.Files(lang.Programs["python"]); !strings.Contains(f["main.py"], "def f") || !strings.Contains(f["main.py"], "print(f() + 1)") {
		t.Errorf("group files: %q", f)
	}
}

func TestRules(t *testing.T) {
	src := "```python\nx = 1\n```\n\n```go run\npackage main\n```\n\n```rust fails expect=nope\nfn main() {}\n```\n\n```bash norun reason=lazy\nls\n```\n\n```output for=missing\n1\n```\n\n```\nno language\n```\n"
	_, issues := Collect("p.md", src, 1)
	var got []string
	for _, i := range issues {
		got = append(got, i.Rule)
	}
	want := "X01 X02 X04 X01 X02 X01"
	if strings.Join(got, " ") != want {
		t.Errorf("rules %v, want %s\n%v", got, want, issues)
	}
}

func TestNormalize(t *testing.T) {
	if got := Normalize("a  \r\n\x1b[31mb\x1b[0m\n\n\n", nil); got != "a\nb\n" {
		t.Errorf("default: %q", got)
	}
	if got := Normalize("<obj at 0x7f3a2b1c>", []string{"addr"}); got != "<obj at 0x…>\n" {
		t.Errorf("addr: %q", got)
	}
	if got := Normalize("0.30000000000000004", []string{"float:3"}); got != "0.3\n" {
		t.Errorf("float: %q", got)
	}
}

func TestCompare(t *testing.T) {
	out := func(info, code string) *Block {
		b := &Block{Info: info, Code: code}
		fillInfo(b)
		return b
	}
	got := lang.RunResult{Stdout: "b\na\n", Stderr: "warning 12\n"}
	if why := compare(out("output", "b\na\n"), got, nil); why != "" {
		t.Error(why)
	}
	if why := compare(out("output", "a\nb\n"), got, nil); why == "" {
		t.Error("exact match must see the order")
	}
	if why := compare(out("output match=unordered", "a\nb\n"), got, nil); why != "" {
		t.Error(why)
	}
	if why := compare(out("output stream=stderr match=regex", `warning \d+`+"\n"), got, nil); why != "" {
		t.Error(why)
	}
}

func TestCheckFailure(t *testing.T) {
	for _, tc := range []struct {
		sig  string
		got  lang.RunResult
		pass bool
	}{
		{"error:E0382", lang.RunResult{Exit: 1, Phase: "build", Stderr: "error[E0382]: borrow of moved value: `s`"}, true},
		{"error:E0499", lang.RunResult{Exit: 1, Phase: "build", Stderr: "error[E0382]: borrow of moved value: `s`"}, false},
		{"compile", lang.RunResult{Exit: 1, Phase: "build"}, true},
		{"compile", lang.RunResult{Exit: 1, Phase: "run", Stderr: "  File \"main.py\", line 1\nSyntaxError: invalid syntax"}, true},
		{"runtime", lang.RunResult{Exit: 1, Phase: "run"}, true},
		{"exception:ZeroDivisionError", lang.RunResult{Exit: 1, Phase: "run", Stderr: "ZeroDivisionError: division by zero"}, true},
		{"exception:KeyError", lang.RunResult{Exit: 1, Phase: "run", Stderr: "ZeroDivisionError: division by zero"}, false},
		{"exception:NullPointerException", lang.RunResult{Exit: 1, Phase: "run", Stderr: "Exception in thread \"main\" java.lang.NullPointerException: Cannot invoke"}, true},
		{"exit:3", lang.RunResult{Exit: 3, Phase: "run"}, true},
		{"exit:3", lang.RunResult{Exit: 0, Phase: "run"}, false},
		{"panic", lang.RunResult{Exit: 2, Phase: "run", Stderr: "panic: runtime error: index out of range [5] with length 3"}, true},
		{"timeout", lang.RunResult{Exit: 1, TimedOut: true}, true},
	} {
		why := checkFailure(tc.sig, tc.got)
		if (why == "") != tc.pass {
			t.Errorf("%s on %+v: %q", tc.sig, tc.got, why)
		}
	}
}

func TestUpdate(t *testing.T) {
	src := "```python run\nprint(1)\n```\n\nText.\n\n```python run id=b\nprint(2)\n```\n\n```output for=b\nold\n```\n"
	cases, _ := Collect("p.md", src, 1)
	rs := []Result{
		{Status: "fail", Case: &cases[0], Got: lang.RunResult{Stdout: "1\n"}},
		{Status: "fail", Case: &cases[1], Got: lang.RunResult{Stdout: "2\r\n"}},
	}
	got, n := Update(src, rs)
	want := "```python run\nprint(1)\n```\n\n```output\n1\n```\n\nText.\n\n```python run id=b\nprint(2)\n```\n\n```output for=b\n2\n```\n"
	if got != want || n != 2 || !rs[0].Updated || !rs[1].Updated {
		t.Errorf("got %d edits:\n%s", n, got)
	}
	if _, issues := Collect("p.md", got, 1); len(issues) != 0 {
		t.Errorf("issues after the update: %v", issues)
	}
}

// TestRunPython runs real blocks when python3 is on this machine.
func TestRunPython(t *testing.T) {
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("no python3")
	}
	cases, issues := Collect("p.md", page, 1)
	if len(issues) != 0 {
		t.Fatal(issues)
	}
	x := &Runner{R: &lang.Runner{TempDir: t.TempDir()}}
	for i := range cases {
		if r := x.Run(&cases[i]); r.Status != "pass" {
			t.Errorf("%s: %s %s", cases[i].ID(), r.Status, r.Why)
		}
	}
	bad := "```python run\nprint('x')\n```\n\n```output\ny\n```\n"
	cases, _ = Collect("p.md", bad, 1)
	if r := x.Run(&cases[0]); r.Status != "fail" || !strings.Contains(r.Why, "differs") {
		t.Errorf("a wrong output must fail: %+v", r)
	}
}

func TestAdjust(t *testing.T) {
	b := &Block{Info: `rust fails flags="--edition 2021"`}
	fillInfo(b)
	got := strings.Join(adjust(lang.Programs["rust"], b).Steps[0], " ")
	if got != "rustc --edition 2021 -o main main.rs" {
		t.Errorf("rust: %s", got)
	}
	b = &Block{Info: `c fails expect=sanitizer:address flags=-std=c17`}
	fillInfo(b)
	got = strings.Join(adjust(lang.Programs["c"], b).Steps[0], " ")
	if got != "cc -g -fsanitize=address -fno-sanitize-recover=all -std=c17 -Wall -o main main.c -lm" {
		t.Errorf("c: %s", got)
	}
	if strings.Join(lang.Programs["rust"].Steps[0], " ") != "rustc --edition 2024 -o main main.rs" {
		t.Error("adjust changed the program of the language")
	}
}

func TestDeps(t *testing.T) {
	src := "```rust run id=a deps=deps/rust/Cargo.lock\nfn main() {}\n```\n\n```output\n```\n\n" +
		"```rust norun reason=fragment group=g file=util.rs deps=deps/rust/Cargo.lock\npub fn f() {}\n```\n\n" +
		"```rust run group=g match=any\nmod util;\nfn main() {}\n```\n\n" +
		"```python run match=any\nprint(1)\n```\n"
	cases, issues := Collect("p.md", src, 1)
	if len(issues) != 0 {
		t.Fatal(issues)
	}
	if len(cases) != 3 {
		t.Fatalf("got %d cases, want 3", len(cases))
	}
	cases[0].Dir = filepath.Join("nodes", "json")
	if cases[0].Block.Deps() != "deps/rust/Cargo.lock" || cases[0].LockPath() != filepath.Join("nodes", "json", "deps", "rust", "Cargo.lock") {
		t.Errorf("case 1: deps %q, lock %q", cases[0].Deps(), cases[0].LockPath())
	}
	if cases[1].Deps() != "deps/rust/Cargo.lock" {
		t.Errorf("a group takes deps= from any of its blocks: %q", cases[1].Deps())
	}
	if cases[2].Deps() != "" || cases[2].LockPath() != "" {
		t.Errorf("a block with no deps=: %q", cases[2].Deps())
	}
	// A Cargo project has its files in src.
	p := lang.Programs["rust"]
	p.Main = "src/main.rs"
	if f := cases[1].Files(p); f["src/util.rs"] == "" || !strings.Contains(f["src/main.rs"], "mod util;") {
		t.Errorf("group files: %q", f)
	}
	if f := cases[1].Files(lang.Programs["rust"]); f["util.rs"] == "" || f["main.rs"] == "" {
		t.Errorf("group files with rustc: %q", f)
	}

	for _, bad := range []string{"../other/Cargo.lock", "/etc/Cargo.lock", "deps/../../x/Cargo.lock", "deps//Cargo.lock", `deps\rust\Cargo.lock`, ""} {
		_, issues := Collect("p.md", "```rust run match=any deps=\""+bad+"\"\nfn main() {}\n```\n", 1)
		if len(issues) != 1 || issues[0].Rule != "X01" || !strings.Contains(issues[0].Msg, "deps=") {
			t.Errorf("deps=%q gives %v", bad, issues)
		}
	}
}

// TestDepsSkip checks that a block with deps= is a skip with the reason deps, and not a failure, when prog run cannot get its packages. The toolchain need not be there.
func TestDepsSkip(t *testing.T) {
	src := "```python run match=any deps=deps/python/requirements.lock\nimport requests\n```\n\n" +
		"```cpp run match=any deps=deps/cpp/vcpkg.json\n#include <nlohmann/json.hpp>\nint main() {}\n```\n\n" +
		"```kotlin run match=any deps=deps/kotlin/pom.xml\nfun main() {}\n```\n"
	cases, issues := Collect("p.md", src, 1)
	if len(issues) != 0 || len(cases) != 3 {
		t.Fatalf("%d cases, issues %v", len(cases), issues)
	}
	local := &Runner{R: &lang.Runner{TempDir: t.TempDir(), Cache: t.TempDir()}}
	box := &Runner{R: &lang.Runner{TempDir: t.TempDir(), Box: lang.NewSandbox("docker", map[string]string{})}}
	for _, tc := range []struct {
		x *Runner
		i int
	}{
		{local, 0}, {local, 1}, {local, 2}, {box, 0}, {box, 2},
	} {
		r := tc.x.Run(&cases[tc.i])
		if r.Status != "skip" || r.Reason != "deps" || !strings.HasPrefix(r.Why, "deps: ") {
			t.Errorf("%s block (sandbox %v): %s %q %q", r.Lang, tc.x.R.Box != nil, r.Status, r.Reason, r.Why)
		}
	}
}
