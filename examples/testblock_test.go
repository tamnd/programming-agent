// SPDX-License-Identifier: Apache-2.0

package examples

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/tamnd/programming-agent/lang"
)

func block(info string) *Block {
	b := &Block{Info: info}
	fillInfo(b)
	return b
}

func TestTestMode(t *testing.T) {
	for _, tc := range []struct {
		info, mode  string
		fails, runs bool
	}{
		{"python test", "test", false, true},
		{"python test fails expect=test-failure", "test", true, true},
		{"python fails test expect=test-failure", "test", true, true},
		{"python fails expect=exception", "fails", false, true},
		{"python run test", "many", false, false},
		{"python test norun reason=fragment", "many", false, false},
		{"python norun reason=fragment", "norun", false, false},
		{"python", "", false, false},
	} {
		b := block(tc.info)
		if b.Mode() != tc.mode || b.TestFails() != tc.fails || b.Runs() != tc.runs {
			t.Errorf("%s: mode %q, fails %v, runs %v", tc.info, b.Mode(), b.TestFails(), b.Runs())
		}
	}
}

func TestTestRules(t *testing.T) {
	for _, tc := range []struct {
		src, rules string
	}{
		{"```python test\nimport unittest\n```\n\n```output stream=stderr\nOK\n```\n", ""},
		// A test block needs an output block, as a run block does. A test block with fails does not.
		{"```python test\nimport unittest\n```\n", "X02"},
		{"```python test fails expect=test-failure\nimport unittest\n```\n", ""},
		{"```python test fails\nimport unittest\n```\n", "X04"},
		{"```python test fails expect=exception\nimport unittest\n```\n", "X04"},
		{"```python test expect=test-failure match=any\nimport unittest\n```\n", "X04"},
		{"```python fails expect=test-failure\nimport unittest\n```\n", "X04"},
		{"```python test norun reason=fragment\nimport unittest\n```\n", "X01"},
		// A test file can come from a group: the fragment is norun, and the test block owns the run.
		{"```go norun reason=fragment group=g file=cart.go\npackage cart\n```\n\n```go test group=g match=any\npackage cart\n```\n", ""},
	} {
		_, issues := Collect("p.md", tc.src, 1)
		var got []string
		for _, i := range issues {
			got = append(got, i.Rule)
		}
		if strings.Join(got, " ") != tc.rules {
			t.Errorf("%q: rules %v, want %q", tc.src, issues, tc.rules)
		}
	}
	src := "```go norun reason=fragment group=g file=cart.go\npackage cart\n```\n\n```go test fails expect=test-failure group=g\npackage cart\n```\n"
	cases, _ := Collect("p.md", src, 1)
	if len(cases) != 1 || cases[0].Label() != "test fails" || cases[0].Block.Line != 5 {
		t.Fatalf("cases %+v", cases)
	}
	f := cases[0].Files(lang.TestPrograms["go"])
	if f["cart.go"] == "" || f["main_test.go"] == "" || cases[0].TestFile(lang.TestPrograms["go"]) != "main_test.go" {
		t.Errorf("files %q", f)
	}
	cases, _ = Collect("p.md", "```java test file=CartTest.java match=any\nclass CartTest {}\n```\n", 1)
	if got := cases[0].TestFile(lang.TestPrograms["java"]); got != "CartTest.java" {
		t.Errorf("test file %q", got)
	}
}

// The fixtures are real reports of the test runners, with the paths of the work directory already $WORK.
func TestNormalizeReport(t *testing.T) {
	for _, tc := range []struct {
		lang, in, want string
	}{
		{"go",
			"=== RUN   TestBad\n    main_test.go:8: got 2, want 3\n--- FAIL: TestBad (0.00s)\nFAIL\nFAIL\texample\t0.160s\nFAIL\n",
			"=== RUN   TestBad\n    main_test.go:8: got 2, want 3\n--- FAIL: TestBad ($TIME)\nFAIL\nFAIL\texample\t$TIME\nFAIL\n"},
		{"go",
			"=== RUN   TestTotal\n=== RUN   TestTotal/empty\n--- PASS: TestTotal (0.01s)\n    --- PASS: TestTotal/empty (0.00s)\nPASS\nok  \texample\t0.312s\n",
			"=== RUN   TestTotal\n=== RUN   TestTotal/empty\n--- PASS: TestTotal ($TIME)\n    --- PASS: TestTotal/empty ($TIME)\nPASS\nok  \texample\t$TIME\n"},
		{"python",
			"test_bad (test_main.T.test_bad) ... FAIL\n\n----------------------------------------------------------------------\nRan 2 tests in 0.001s\n\nFAILED (failures=1)\n",
			"test_bad (test_main.T.test_bad) ... FAIL\n\n----------------------------------------------------------------------\nRan 2 tests in $TIME\n\nFAILED (failures=1)\n"},
		{"python", "Ran 1 test in 12.5s\r\n\r\nOK\r\n", "Ran 1 test in $TIME\n\nOK\n"},
		{"javascript",
			"✔ ok (0.466333ms)\n✖ bad (0.5075ms)\nℹ tests 2\nℹ fail 1\nℹ duration_ms 57.464333\n\n✖ failing tests:\n\ntest at main.test.js:4:1\n✖ bad (0.5075ms)\n" +
				"  AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:\n  \n  2 !== 3\n  \n" +
				"      at TestContext.<anonymous> (file://$WORK/main.test.js:4:26)\n" +
				"      at Test.runInAsyncScope (node:async_hooks:226:14)\n" +
				"      at Test.run (node:internal/test_runner/test:1402:25)\n" +
				"      at Test.processPendingSubtests (node:internal/test_runner/test:974:18)\n" +
				"      at async startSubtestAfterBootstrap (node:internal/test_runner/harness:387:3) {\n" +
				"    generatedMessage: true,\n    actual: 2,\n  }\n",
			"✔ ok ($TIME)\n✖ bad ($TIME)\nℹ tests 2\nℹ fail 1\nℹ duration_ms $TIME\n\n✖ failing tests:\n\ntest at main.test.js:4:1\n✖ bad ($TIME)\n" +
				"  AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:\n  \n  2 !== 3\n  \n" +
				"      at TestContext.<anonymous> (file://$WORK/main.test.js:4:26) {\n" +
				"    generatedMessage: true,\n    actual: 2,\n  }\n"},
		{"typescript", "✔ total (12ms)\nℹ duration_ms 74.05\n", "✔ total ($TIME)\nℹ duration_ms $TIME\n"},
		{"rust",
			"running 2 tests\ntest bad ... FAILED\ntest ok ... ok\n\nfailures:\n\n---- bad stdout ----\n\nthread 'bad' (1335517) panicked at main.rs:6:5:\nassertion `left == right` failed\n\ntest result: FAILED. 1 passed; 1 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.00s\n",
			"running 2 tests\ntest bad ... FAILED\ntest ok ... ok\n\nfailures:\n\n---- bad stdout ----\n\nthread 'bad' ($ID) panicked at main.rs:6:5:\nassertion `left == right` failed\n\ntest result: FAILED. 1 passed; 1 failed; 0 ignored; 0 measured; 0 filtered out; finished in $TIME\n"},
		{"rust",
			"    Finished `test` profile [unoptimized + debuginfo] target(s) in 0.52s\n     Running unittests src/main.rs (target/debug/deps/cart-1a2b3c4d5e6f7a8b)\n",
			"    Finished `test` profile [unoptimized + debuginfo] target(s) in $TIME\n     Running unittests src/main.rs (target/debug/deps/cart-$HASH)\n"},
		{"java",
			"│     └─ bad() ✘ expected: <3> but was: <2>\n\nFailures (1):\n  JUnit Jupiter:CartTest:bad()\n" +
				"    => org.opentest4j.AssertionFailedError: expected: <3> but was: <2>\n" +
				"       org.junit.jupiter.api.AssertionFailureBuilder.build(AssertionFailureBuilder.java:151)\n" +
				"       org.junit.jupiter.api.Assertions.assertEquals(Assertions.java:531)\n" +
				"       CartTest.bad(CartTest.java:10)\n" +
				"       java.base/java.lang.reflect.Method.invoke(Method.java:583)\n" +
				"       java.base/java.util.ArrayList.forEach(ArrayList.java:1612)\n" +
				"\nTest run finished after 127 ms\n[         1 tests failed          ]\n",
			"│     └─ bad() ✘ expected: <3> but was: <2>\n\nFailures (1):\n  JUnit Jupiter:CartTest:bad()\n" +
				"    => org.opentest4j.AssertionFailedError: expected: <3> but was: <2>\n" +
				"       CartTest.bad(CartTest.java:10)\n" +
				"\nTest run finished after $TIME\n[         1 tests failed          ]\n"},
		{"csharp",
			"Test run for $WORK/bin/Debug/net10.0/Cart.Tests.dll (.NETCoreApp,Version=v10.0)\n[xUnit.net 00:00:00.18]     CartTests.Bad [FAIL]\n\nFailed!  - Failed:     1, Passed:     1, Skipped:     0, Total:     2, Duration: 18 ms - Cart.Tests.dll (net10.0)\n",
			"Test run for $WORK/bin/Debug/net10.0/Cart.Tests.dll (.NETCoreApp,Version=v10.0)\n[xUnit.net $TIME]     CartTests.Bad [FAIL]\n\nFailed!  - Failed:     1, Passed:     1, Skipped:     0, Total:     2, Duration: $TIME - Cart.Tests.dll (net10.0)\n"},
		{"csharp",
			"  Passed CartTests.Ok [3 ms]\n  Stack Trace:\n     at CartTests.Bad() in $WORK/CartTests.cs:line 8\n     at System.RuntimeMethodHandle.InvokeMethod(Object target)\n     at InvokeStub_CartTests.Bad(Object, Object, IntPtr*)\n",
			"  Passed CartTests.Ok [$TIME]\n  Stack Trace:\n     at CartTests.Bad() in $WORK/CartTests.cs:line 8\n"},
		{"cpp",
			"Running main() from /opt/homebrew/Cellar/googletest/1.17.0/share/googletest/src/gtest_main.cc\n[       OK ] CartTest.Total (0 ms)\n[==========] 2 tests from 1 test suite ran. (3 ms total)\n",
			"Running main() from gtest_main.cc\n[       OK ] CartTest.Total ($TIME)\n[==========] 2 tests from 1 test suite ran. ($TIME total)\n"},
		{"cpp", "Randomness seeded to: 3141592\n", "Randomness seeded to: $SEED\n"},
		// A language with no normalizers keeps its report.
		{"kotlin", "took 3 ms\n", "took 3 ms\n"},
	} {
		if got := NormalizeReport(tc.lang, tc.in); got != tc.want {
			t.Errorf("%s:\n got %q\nwant %q", tc.lang, got, tc.want)
		}
	}
}

func TestTestFailed(t *testing.T) {
	for _, tc := range []struct {
		lang, failed, passed string
	}{
		{"python", "FAIL: test_bad (test_main.T.test_bad)\n", "test_ok (test_main.T.test_ok) ... ok\n\nOK\n"},
		{"python", "ERROR: test_bad (test_main.T.test_bad)\n", "Ran 0 tests in $TIME\n"},
		{"go", "--- FAIL: TestBad ($TIME)\n", "--- PASS: TestOK ($TIME)\nFAIL\texample [build failed]\n"},
		{"javascript", "ℹ pass 1\nℹ fail 1\n", "ℹ pass 2\nℹ fail 0\n"},
		{"typescript", "# fail 2\n", "# fail 0\n"},
		{"rust", "test bad ... FAILED\n", "test ok ... ok\n\ntest result: ok. 1 passed\n"},
		{"rust", "test result: FAILED. 1 passed; 1 failed\n", "error[E0425]: cannot find value `x`\n"},
		{"java", "[         1 tests failed          ]\n", "[         0 tests failed          ]\n"},
		{"csharp", "Failed!  - Failed:     1, Passed:     1\n", "Passed!  - Failed:     0, Passed:     2\n"},
		{"cpp", "[  FAILED  ] CartTest.Bad (0 ms)\n", "[  PASSED  ] 2 tests.\n"},
		{"cpp", "test cases: 2 | 1 passed | 1 failed\n", "All tests passed (3 assertions in 2 test cases)\n"},
	} {
		re := testFailed[tc.lang]
		if !re.MatchString(tc.failed) {
			t.Errorf("%s: no failed test in %q", tc.lang, tc.failed)
		}
		if re.MatchString(tc.passed) {
			t.Errorf("%s: a failed test in %q", tc.lang, tc.passed)
		}
	}
}

func TestCheckTests(t *testing.T) {
	test := &Case{Block: block("go test")}
	fails := &Case{Block: block("go test fails expect=test-failure")}
	for _, tc := range []struct {
		c    *Case
		got  lang.RunResult
		want string
	}{
		{test, lang.RunResult{Phase: "run", Stdout: "--- PASS: TestOK ($TIME)\nok\n"}, ""},
		{test, lang.RunResult{Phase: "run", Exit: 1, Stdout: "--- FAIL: TestBad ($TIME)\n"}, "a test failed, or the test runner failed, with exit status 1"},
		{test, lang.RunResult{Phase: "build", Exit: 1, Stderr: "syntax error"}, "the test file did not build"},
		{test, lang.RunResult{Phase: "run", Exit: 1, TimedOut: true}, "more than the time limit"},
		{fails, lang.RunResult{Phase: "run", Exit: 1, Stdout: "--- FAIL: TestBad ($TIME)\n"}, ""},
		{fails, lang.RunResult{Phase: "run", Stdout: "--- PASS: TestOK ($TIME)\n"}, "every test passed"},
		{fails, lang.RunResult{Phase: "run", Exit: 1, Stdout: "FAIL\texample [setup failed]\n"}, "the report shows no test that fails"},
		{fails, lang.RunResult{Phase: "build", Exit: 1, Stderr: "syntax error"}, "the test file did not build"},
	} {
		why := checkTests(tc.c, tc.got)
		if (tc.want == "") != (why == "") || !strings.Contains(why, tc.want) {
			t.Errorf("%s on %+v: %q, want %q", tc.c.Label(), tc.got, why, tc.want)
		}
	}
}

func TestReportStream(t *testing.T) {
	src := "```python test id=a\nimport unittest\n```\n\n```output\nOK\n```\n\n```go test id=b\npackage cart\n```\n\n```output\nok\n```\n\n```python run id=c\nprint(1)\n```\n\n```output\n1\n```\n"
	cases, issues := Collect("p.md", src, 1)
	if len(issues) != 0 || len(cases) != 3 {
		t.Fatalf("%d cases, issues %v", len(cases), issues)
	}
	for i, want := range []string{"stderr", "", ""} {
		if got := reportStream(&cases[i], cases[i].Outputs[0]); got != want {
			t.Errorf("case %s: stream %q, want %q", cases[i].ID(), got, want)
		}
	}
	if got := reportStream(&cases[0], block("output stream=stdout")); got != "stdout" {
		t.Errorf("stream= of the output block: %q", got)
	}
	// The recorded report goes through the normalizers too, so an old report with a time still matches.
	got := lang.RunResult{Stderr: "Ran 1 test in $TIME\n\nOK\n"}
	if why := compare(block("output"), got, &cases[0]); why == "" {
		t.Error("an empty output block must not match")
	}
	o := block("output")
	o.Code = "Ran 1 test in 0.004s\n\nOK\n"
	if why := compare(o, got, &cases[0]); why != "" {
		t.Error(why)
	}
}

func TestUpdateTest(t *testing.T) {
	src := "```python test id=a\nimport unittest\n```\n\nText.\n\n```python test fails expect=test-failure id=b\nimport unittest\n```\n\n```python test id=c\nimport unittest\n```\n\n```go test id=d\npackage cart\n```\n\n```output\nold\n```\n"
	cases, _ := Collect("p.md", src, 1)
	if len(cases) != 4 {
		t.Fatalf("%d cases", len(cases))
	}
	pass := "test_ok (test_main.T.test_ok) ... ok\n\nRan 1 test in $TIME\n\nOK\n"
	failed := "FAIL: test_bad (test_main.T.test_bad)\n\nRan 1 test in $TIME\n\nFAILED (failures=1)\n"
	rs := []Result{
		{Status: "fail", Case: &cases[0], Got: lang.RunResult{Phase: "run", Stderr: pass}},
		{Status: "pass", Case: &cases[1], Got: lang.RunResult{Phase: "run", Exit: 1, Stderr: failed}},
		// A test block whose tests fail keeps its output block as it is.
		{Status: "fail", Case: &cases[2], Got: lang.RunResult{Phase: "run", Exit: 1, Stderr: failed}},
		{Status: "fail", Case: &cases[3], Got: lang.RunResult{Phase: "run", Stdout: "PASS\nok  \texample\t$TIME\n"}},
	}
	got, n := Update(src, rs)
	want := "```python test id=a\nimport unittest\n```\n\n```output stream=stderr\n" + pass + "```\n\nText.\n\n" +
		"```python test fails expect=test-failure id=b\nimport unittest\n```\n\n```output stream=stderr\n" + failed + "```\n\n" +
		"```python test id=c\nimport unittest\n```\n\n" +
		"```go test id=d\npackage cart\n```\n\n```output\nPASS\nok  \texample\t$TIME\n```\n"
	if got != want || n != 3 || !rs[0].Updated || !rs[1].Updated || rs[2].Updated || !rs[3].Updated {
		t.Errorf("got %d edits:\n%s", n, got)
	}
}

func TestNorunShare(t *testing.T) {
	src := "```python run\nprint(1)\n```\n\n```output\n1\n```\n\n" +
		// A fragment of a group that runs is a part of a program that runs.
		"```python norun reason=fragment group=g\ndef f(): pass\n```\n\n```python run group=g match=any\nf()\n```\n\n" +
		// A fragment of a group that does not run is norun.
		"```go norun reason=fragment group=g\npackage main\n```\n\n" +
		"```bash norun reason=platform\nls\n```\n\n" +
		"```go test match=any\npackage cart\n```\n\n" +
		"```text\nnot code\n```\n\n```json\n{}\n```\n"
	norun, total := NorunShare(src)
	if norun != 2 || total != 6 {
		t.Errorf("norun %d of %d, want 2 of 6", norun, total)
	}
	if norun, total := NorunShare("Text only.\n"); norun != 0 || total != 0 {
		t.Errorf("no blocks: %d of %d", norun, total)
	}
}

func TestAdjustTest(t *testing.T) {
	b := block(`go test flags="-race -short"`)
	p := adjust(lang.TestPrograms["go"], b)
	if got := strings.Join(p.Steps[0], " "); got != "go mod init example" {
		t.Errorf("setup step: %s", got)
	}
	if got := strings.Join(p.Steps[1], " "); got != "go test -race -short -v -count=1 ." {
		t.Errorf("test step: %s", got)
	}
	b = block(`rust test flags="--edition 2021"`)
	if got := strings.Join(adjust(lang.TestPrograms["rust"], b).Steps[0], " "); got != "rustc --edition 2021 --test -o main_test.exe {test}" {
		t.Errorf("rust: %s", got)
	}
	if got := strings.Join(lang.TestPrograms["go"].Steps[1], " "); got != "go test -v -count=1 ." {
		t.Errorf("adjust changed the test runner of the language: %s", got)
	}
}

// TestTestSkip checks that a test block that prog run cannot run is a skip with a reason, and not a failure. The toolchains need not be there.
func TestTestSkip(t *testing.T) {
	src := "```c test match=any\nint main(void) {}\n```\n\n" +
		"```kotlin test match=any\nfun main() {}\n```\n\n" +
		"```csharp test match=any\nusing Xunit;\n```\n\n" +
		"```cpp test match=any\n#include <gtest/gtest.h>\n```\n\n" +
		"```python test match=any deps=deps/python/requirements.lock\nimport pytest\n```\n\n" +
		"```cpp test match=any deps=deps/cpp/vcpkg.json\n#include <gtest/gtest.h>\n```\n"
	cases, issues := Collect("p.md", src, 1)
	if len(issues) != 0 || len(cases) != 6 {
		t.Fatalf("%d cases, issues %v", len(cases), issues)
	}
	x := &Runner{R: &lang.Runner{TempDir: t.TempDir(), Cache: t.TempDir()}}
	for i, want := range []string{"runner", "runner", "deps", "deps", "deps", "deps"} {
		r := x.Run(&cases[i])
		if r.Status != "skip" || r.Reason != want || r.Why == "" || r.Mode != "test" {
			t.Errorf("%s test block: %s %q %q", r.Lang, r.Status, r.Reason, r.Why)
		}
	}
}

// testCode are a passing and a failing test file for each language that has a test runner in its standard library.
var testCode = map[string][2]string{
	"python": {
		"import unittest\n\n\nclass T(unittest.TestCase):\n    def test_ok(self):\n        self.assertEqual(1 + 1, 2)\n",
		"import unittest\n\n\nclass T(unittest.TestCase):\n    def test_bad(self):\n        self.assertEqual(1 + 1, 3)\n",
	},
	"javascript": {
		"import test from 'node:test';\nimport assert from 'node:assert/strict';\n\ntest('ok', () => assert.equal(1 + 1, 2));\n",
		"import test from 'node:test';\nimport assert from 'node:assert/strict';\n\ntest('bad', () => assert.equal(1 + 1, 3));\n",
	},
	"typescript": {
		"import test from 'node:test';\nimport assert from 'node:assert/strict';\n\ntest('ok', () => {\n  const n: number = 1 + 1;\n  assert.equal(n, 2);\n});\n",
		"import test from 'node:test';\nimport assert from 'node:assert/strict';\n\ntest('bad', () => {\n  const n: number = 1 + 1;\n  assert.equal(n, 3);\n});\n",
	},
	"go": {
		"package cart\n\nimport \"testing\"\n\nfunc TestOK(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Error(\"sum\")\n\t}\n}\n",
		"package cart\n\nimport \"testing\"\n\nfunc TestBad(t *testing.T) {\n\tt.Errorf(\"got %d, want %d\", 2, 3)\n}\n",
	},
	"rust": {
		"#[test]\nfn ok() {\n    assert_eq!(1 + 1, 2);\n}\n",
		"#[test]\nfn bad() {\n    assert_eq!(1 + 1, 3);\n}\n",
	},
}

// TestRunTestLocal runs real test blocks with the test runners of this machine. A language whose toolchain is not here is skipped. A passing test block gets its output block from Update, and then it must pass again with it, so the normalizers make the report the same on each run.
func TestRunTestLocal(t *testing.T) {
	for _, l := range []string{"python", "javascript", "typescript", "go", "rust"} {
		t.Run(l, func(t *testing.T) {
			if _, err := exec.LookPath(lang.TestPrograms[l].Tools[0]); err != nil {
				t.Skip("no " + lang.TestPrograms[l].Tools[0])
			}
			ok, bad := testCode[l][0], testCode[l][1]
			fence := func(info, code string) string { return "```" + l + " " + info + " timeout=30\n" + code + "```\n\n" }
			src := fence("test id=pass", ok) + fence("test id=fail match=any", bad) + fence("test fails expect=test-failure id=fails", bad) + fence("test fails expect=test-failure id=nofail", ok)
			cases, _ := Collect("p.md", src, 1)
			x := &Runner{R: &lang.Runner{TempDir: t.TempDir()}}
			var rs []Result
			for i, want := range []string{"pass", "fail", "pass", "fail"} {
				r := x.Run(&cases[i])
				if r.Status != want {
					t.Errorf("%s: %s, want %s: %s\n%s%s", r.ID, r.Status, want, r.Why, r.Got.Stdout, r.Got.Stderr)
				}
				rs = append(rs, r)
			}
			if r := rs[3]; !strings.Contains(r.Why, "every test passed") {
				t.Errorf("nofail: %q", r.Why)
			}
			if r := rs[1]; !strings.Contains(r.Why, "a test failed") {
				t.Errorf("fail: %q", r.Why)
			}
			got, n := Update(src, rs[:1])
			if n != 1 {
				t.Fatalf("Update made %d edits", n)
			}
			cases, issues := Collect("p.md", got, 1)
			if len(issues) != 0 || len(cases[0].Outputs) != 1 {
				t.Fatalf("after the update: %v\n%s", issues, got)
			}
			out := cases[0].Outputs[0].Code
			if !strings.Contains(out, "$TIME") {
				t.Errorf("the report is not normalized:\n%s", out)
			}
			if r := x.Run(&cases[0]); r.Status != "pass" {
				t.Errorf("the second run: %s %s", r.Status, r.Why)
			}
		})
	}
}
