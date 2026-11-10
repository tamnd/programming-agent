// SPDX-License-Identifier: Apache-2.0

package examples

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/tamnd/programming-agent/lang"
)

// The report of a test runner holds timings, and some runners print paths and hashes of the build. The normalizers of a test block replace them, so that the output block is the same on each run (spec 07.11). The paths of the work directory are $WORK already, because RunProgram replaces them.

// sub is one replacement of a normalizer.
type sub struct {
	re   *regexp.Regexp
	with string
}

func subs(pairs ...string) []sub {
	out := make([]sub, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, sub{regexp.MustCompile(pairs[i]), pairs[i+1]})
	}
	return out
}

// testNormalizers are the normalizers of the reports of the test runners, by language.
var testNormalizers = map[string][]sub{
	// --- PASS: TestTotal (0.00s), and ok  	example	0.312s at the end.
	"go": subs(
		`(?m)^(\s*--- (?:PASS|FAIL|SKIP): .+ )\(\d+(?:\.\d+)?s\)$`, `${1}($$TIME)`,
		`(?m)^((?:ok|FAIL)\s+\S+\s+)\d+(?:\.\d+)?s$`, `${1}$$TIME`,
	),
	// Ran 3 tests in 0.001s
	"python": subs(`(?m)^(Ran \d+ tests? in )\d+(?:\.\d+)?s$`, `${1}$$TIME`),
	// ✔ total (0.41ms), ℹ duration_ms 74.05, and the frames of Node itself in a stack trace.
	"javascript": nodeTest,
	"typescript": nodeTest,
	// test result: ok. 3 passed; ... finished in 0.00s. Cargo also prints the build time and the hash of the test binary.
	"rust": subs(
		`finished in \d+(?:\.\d+)?s`, `finished in $$TIME`,
		`(?m)^(\s*Finished .* in )\d+(?:\.\d+)?m?s$`, `${1}$$TIME`,
		`(target/debug/deps/[\w-]+)-[0-9a-f]{16}\b`, `${1}-$$HASH`,
		// thread 'bad' (1335517) panicked at main.rs:6:5. The number is the ID of the thread.
		`(?m)^(thread '[^']*' )\(\d+\)( panicked)`, `${1}($$ID)$2`,
	),
	// Test run finished after 119 ms, and the frames of the JDK and of JUnit in a stack trace.
	"java": subs(
		`Test run finished after \d+ ms`, `Test run finished after $$TIME`,
		`(?m)^[ \t]+(?:java\.base/|jdk\.internal\.|java\.lang\.reflect\.|org\.junit\.|org\.opentest4j\.).*\n`, ``,
		`(?m)^([ \t]+)\[\.\.\.\] *\n(?:[ \t]+\[\.\.\.\] *\n)+`, "$1[...]\n",
	),
	// Duration: 49 ms, [3 ms], 00:00:00.12, and the frames of .NET in a stack trace.
	"csharp": subs(
		`Duration: (?:< )?\d+(?:\.\d+)? ?(?:ms|s|m)\b`, `Duration: $$TIME`,
		`\[(?:< )?\d+(?:\.\d+)? ?(?:ms|s|m)\]`, `[$$TIME]`,
		`\b\d{1,2}:\d{2}:\d{2}(?:\.\d+)?\b`, `$$TIME`,
		`(?m)^[ \t]+at (?:System\.|InvokeStub_|Xunit\.).*\n`, ``,
	),
	// [       OK ] CartTest.Total (0 ms), and the path of the gtest_main.cc of the build of GoogleTest. Catch2 prints a random seed.
	"cpp": subs(
		`\((\d+) ms( total)?\)`, `($$TIME$2)`,
		`(?m)^Running main\(\) from \S*gtest_main\.cc$`, `Running main() from gtest_main.cc`,
		`(?m)^(Randomness seeded to: )\d+$`, `${1}$$SEED`,
		`\b\d+(?:\.\d+)? s: `, `$$TIME: `,
	),
}

var nodeTest = subs(
	`\(\d+(?:\.\d+)?ms\)`, `($$TIME)`,
	`(?m)^(\s*(?:ℹ|#) duration_ms:? )\d+(?:\.\d+)?$`, `${1}$$TIME`,
	`(?m)^(\s*duration_ms: )\d+(?:\.\d+)?$`, `${1}$$TIME`,
	`(?m)^[ \t]+at (?:.*\()?node:[\w/]+:\d+:\d+\)?\n`, ``,
	// The last frame can open the properties of the error with " {". The brace then goes to the end of the frame before it.
	`\n[ \t]+at (?:.*\()?node:[\w/]+:\d+:\d+\)? \{\n`, " {\n",
)

// NormalizeReport applies the normalizers of the test runner of a language to a report.
func NormalizeReport(lang, s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	for _, n := range testNormalizers[lang] {
		s = n.re.ReplaceAllString(s, n.with)
	}
	return s
}

// testFailed finds a failed test in the report of a test runner. A test block with fails needs it (spec 07.11).
var testFailed = map[string]*regexp.Regexp{
	"python":     regexp.MustCompile(`(?m)^(?:FAIL|ERROR): test\w*`),
	"go":         regexp.MustCompile(`(?m)^\s*--- FAIL: `),
	"javascript": regexp.MustCompile(`(?m)^[ℹ#] fail [1-9]`),
	"typescript": regexp.MustCompile(`(?m)^[ℹ#] fail [1-9]`),
	"rust":       regexp.MustCompile(`(?m)^test result: FAILED\.|^test .* \.\.\. FAILED$`),
	"java":       regexp.MustCompile(`\[\s*[1-9]\d* tests failed\s*\]`),
	"csharp":     regexp.MustCompile(`Failed:\s+[1-9]`),
	"cpp":        regexp.MustCompile(`(?m)^\[  FAILED  \] |^test cases: .*\b[1-9]\d* failed`),
}

// checkTests checks the run of a test block. A test block passes when every test passes. A test block with fails passes when the tests ran and at least one failed. It returns why the run does not match, or "".
func checkTests(c *Case, got lang.RunResult) string {
	report := got.Stdout + got.Stderr
	if got.TimedOut {
		return "the tests ran for more than the time limit"
	}
	if got.Phase == "build" {
		return "the test file did not build:\n" + tail(got.Stderr+got.Stdout)
	}
	if !c.Block.TestFails() {
		if got.Exit != 0 {
			return "a test failed, or the test runner failed, with exit status " + strconv.Itoa(got.Exit) + ":\n" + tail(report)
		}
		return ""
	}
	if got.Exit == 0 {
		return "every test passed, but the block expects a test that fails (expect=" + TestFailure + ")"
	}
	if re := testFailed[c.Block.Lang]; re != nil && !re.MatchString(report) {
		return "the test runner failed, but the report shows no test that fails:\n" + tail(report)
	}
	return ""
}

// reportStream is the stream of an output block of a case: the stream= of the block, or the stream of the report of the test runner for a test block.
func reportStream(c *Case, o *Block) string {
	if s := o.Attrs["stream"]; s != "" {
		return s
	}
	if c != nil && c.Mode() == "test" {
		return lang.TestPrograms[c.Block.Lang].Report
	}
	return ""
}

// NorunShare counts the blocks of a Markdown file for rule X08 (spec 07.8). total is the number of code blocks with a run mode, and norun is the number of them that are norun. A norun block of a group that runs is a part of a program that runs, so it does not count as norun.
func NorunShare(src string) (norun, total int) {
	blocks := Parse(src, 1)
	runs := map[string]bool{}
	for i := range blocks {
		b := &blocks[i]
		if g := b.Attrs["group"]; g != "" && b.Runs() {
			runs[b.Lang+"\x00"+g] = true
		}
	}
	for i := range blocks {
		b := &blocks[i]
		switch m := b.Mode(); {
		case m == "norun":
			total++
			if g := b.Attrs["group"]; g == "" || !runs[b.Lang+"\x00"+g] {
				norun++
			}
		case b.Runs():
			total++
		}
	}
	return norun, total
}
