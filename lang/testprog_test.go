// SPDX-License-Identifier: Apache-2.0

package lang

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func steps(p Program) string {
	var out []string
	for _, s := range p.Steps {
		out = append(out, strings.Join(s, " "))
	}
	return strings.Join(out, " | ")
}

func TestTestProgramFor(t *testing.T) {
	for lang, want := range map[string]string{
		"python":     "python3 -m unittest -v {test}",
		"javascript": "node --test --test-reporter=spec {test}",
		"typescript": "node --test --test-reporter=spec {test}",
		"go":         "go mod init example | go test -v -count=1 .",
		"rust":       "rustc --edition 2024 --test -o main_test.exe {test} | ./main_test.exe --test-threads=1",
		"java":       "javac -d classes -cp /opt/prog/junit.jar {sources} | java -cp classes:/opt/prog/junit.jar org.junit.platform.console.ConsoleLauncher execute --disable-banner --disable-ansi-colors --details=tree --scan-class-path",
	} {
		p, err := TestProgramFor(lang, false)
		if err != nil {
			t.Fatalf("%s: %v", lang, err)
		}
		if got := steps(p); got != want || !p.Test {
			t.Errorf("%s: %s, Test %v; want %s", lang, got, p.Test, want)
		}
	}
	if TestPrograms["python"].Report != "stderr" || TestPrograms["go"].Setup != 1 {
		t.Error("python reports on stderr, and go has one setup step")
	}
	for lang, reason := range map[string]string{"c": "runner", "kotlin": "runner", "bash": "runner", "csharp": "deps", "cpp": "deps"} {
		_, err := TestProgramFor(lang, false)
		var ts TestSkip
		if !errors.As(err, &ts) || ts.Reason != reason || ts.Why == "" {
			t.Errorf("%s: %v; want a skip with the reason %s", lang, err, reason)
		}
	}
	for _, lang := range []string{"csharp", "cpp"} {
		if _, err := TestProgramFor(lang, true); err != nil {
			t.Errorf("%s with deps=: %v", lang, err)
		}
	}
	p, _ := TestProgramFor("rust", false)
	p.Steps[0][0] = "changed"
	if TestPrograms["rust"].Steps[0][0] != "rustc" {
		t.Error("TestProgramFor gives the steps of the table, not a copy")
	}
}

func TestWithFiles(t *testing.T) {
	p, _ := TestProgramFor("java", false)
	q := p.WithFiles(map[string]string{"CartTest.java": "", "Cart.java": "", "notes.txt": ""}, "CartTest.java")
	if got := strings.Join(q.Steps[0], " "); got != "javac -d classes -cp /opt/prog/junit.jar Cart.java CartTest.java" {
		t.Errorf("java: %s", got)
	}
	if p.Steps[0][len(p.Steps[0])-1] != ArgSources {
		t.Error("WithFiles changed the program it was called on")
	}
	p, _ = TestProgramFor("python", false)
	if got := steps(p.WithFiles(map[string]string{"test_cart.py": "", "cart.py": ""}, "test_cart.py")); got != "python3 -m unittest -v test_cart.py" {
		t.Errorf("python: %s", got)
	}
}

func TestResolveTest(t *testing.T) {
	p, _ := TestProgramFor("java", false)
	box := NewSandbox("docker", map[string]string{"java": "img"})
	q, err := (&Runner{Box: box}).ResolveTest(p)
	if err != nil || steps(q) != steps(p) {
		t.Errorf("in the sandbox, the steps use the jar of the image: %s, %v", steps(q), err)
	}
	cache := t.TempDir()
	jar := filepath.Join(cache, "junit-platform-console-standalone-"+JUnitVersion+".jar")
	if err := os.WriteFile(jar, []byte("jar"), 0o644); err != nil {
		t.Fatal(err)
	}
	q, err = (&Runner{Cache: cache}).ResolveTest(p)
	if err != nil {
		t.Fatal(err)
	}
	if q.Steps[0][4] != jar || q.Steps[1][2] != "classes"+string(os.PathListSeparator)+jar {
		t.Errorf("on this machine: %v", q.Steps)
	}
	py, _ := TestProgramFor("python", false)
	if q, _ := (&Runner{Cache: cache}).ResolveTest(py); steps(q) != steps(py) {
		t.Errorf("python changed: %s", steps(q))
	}
}

func TestTestDepsPrograms(t *testing.T) {
	toml := map[string][]byte{"Cargo.toml": []byte("[package]\nname = \"example\"\n")}
	rust, _ := TestProgramFor("rust", true)
	p, err := rustTestDeps(rust, toml, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Main != "src/main.rs" || p.Seed != "" || p.Steps[0][0] != "sh" || !strings.HasSuffix(p.Steps[0][2], "exec cargo test --offline --locked --no-run --quiet") || strings.Join(p.Steps[1], " ") != "cargo test --offline --locked -- --test-threads=1" {
		t.Errorf("rust in the image: %+v", p)
	}
	if p, _ = rustTestDeps(rust, toml, "/cache/rust-1"); p.Seed != "/cache/rust-1" || steps(p) != "cargo test --offline --locked --no-run --quiet | cargo test --offline --locked -- --test-threads=1" {
		t.Errorf("rust on this machine: %+v", p)
	}

	cpp, _ := TestProgramFor("cpp", true)
	for manifest, want := range map[string]string{
		`{"name": "x", "dependencies": ["fmt", "gtest"]}`:          "c++ -std=c++23 -Wall -I/opt/prog/include -L/opt/prog/lib -o main_test.exe {sources} -lgtest_main -lgtest -pthread | ./main_test.exe",
		`{"dependencies": [{"name": "catch2", "version>=": "3"}]}`: "c++ -std=c++23 -Wall -I/opt/prog/include -L/opt/prog/lib -o main_test.exe {sources} -lCatch2Main -lCatch2 | ./main_test.exe --order decl",
	} {
		p, err := cppTestDeps(cpp, map[string][]byte{"vcpkg.json": []byte(manifest)}, "")
		if err != nil || steps(p) != want {
			t.Errorf("cpp %s: %s, %v", manifest, steps(p), err)
		}
	}
	if _, err := cppTestDeps(cpp, map[string][]byte{"vcpkg.json": []byte(`{"dependencies": ["fmt"]}`)}, ""); err == nil || !strings.Contains(err.Error(), "no test framework") {
		t.Errorf("a vcpkg.json with no test framework: %v", err)
	}
	p, _ = cppTestDeps(cpp, map[string][]byte{"vcpkg.json": []byte(`["gtest"]`)}, "")
	if len(p.Steps) != 0 {
		t.Errorf("a bad vcpkg.json gives steps: %v", p.Steps)
	}

	cs, _ := TestProgramFor("csharp", true)
	proj := map[string][]byte{"Cart.Tests.csproj": []byte("<Project/>")}
	p, err = csharpTestDeps(cs, proj, "")
	if err != nil {
		t.Fatal(err)
	}
	if want := "dotnet build -nologo -v q -clp:NoSummary;ErrorsOnly -p:NuGetAudit=false -p:RestoreSources=/opt/nuget | dotnet test --no-build --nologo -v q"; steps(p) != want || string(p.Files["Cart.Tests.csproj"]) != "<Project/>" {
		t.Errorf("csharp in the image: %s, %v", steps(p), p.Files)
	}
	t.Setenv("NUGET_PACKAGES", "/nuget")
	if p, _ = csharpTestDeps(cs, proj, "/cache/cs-1"); !strings.Contains(steps(p), "-p:RestoreSources=/nuget |") {
		t.Errorf("csharp on this machine: %s", steps(p))
	}

	java, _ := TestProgramFor("java", true)
	p, _ = javaDeps(java, nil, "")
	if want := "javac -d classes -cp /opt/prog/junit.jar:/opt/prog/lib/* {sources} | java -cp classes:/opt/prog/junit.jar:/opt/prog/lib/* " + junitLauncher; !strings.HasPrefix(steps(p), want) {
		t.Errorf("java test in the image: %s", steps(p))
	}
}

func TestWithDepsTest(t *testing.T) {
	dir := t.TempDir()
	csproj := filepath.Join(dir, "Cart.Tests.csproj")
	if err := os.WriteFile(csproj, []byte("<Project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The dotnet image has the packages, so the block needs no image of its own, and no engine runs here.
	r := &Runner{Box: NewSandbox("docker", map[string]string{"csharp": "img"})}
	cs, _ := TestProgramFor("csharp", true)
	p, err := r.WithDeps(cs, csproj)
	if err != nil {
		t.Fatal(err)
	}
	if p.Image != "" || len(p.Steps) != 2 || p.Files["Cart.Tests.csproj"] == nil {
		t.Errorf("csharp in the sandbox: %+v", p)
	}
	if _, err := r.WithDeps(cs, filepath.Join(dir, "packages.lock.json")); err == nil || !strings.Contains(err.Error(), "*.csproj") {
		t.Errorf("a lock file that is not a .csproj must fail: %v", err)
	}
	if _, err := r.WithDeps(Programs["csharp"], csproj); !errors.As(err, new(DepsSkip)) {
		t.Errorf("a C# run block with deps= must skip: %v", err)
	}
}

func TestDepsMissing(t *testing.T) {
	missing := RunResult{Phase: "build", Exit: 1, Stdout: "/w/Cart.Tests.csproj : error NU1101: Unable to find package xunit.v3. No packages exist with this id in source(s): nuget\n"}
	if !DepsMissing("csharp", missing) {
		t.Error("NU1101 is a missing package")
	}
	if DepsMissing("csharp", RunResult{Phase: "build", Exit: 1, Stdout: "CartTests.cs(3,1): error CS1002: ; expected\n"}) {
		t.Error("a compile error is not a missing package")
	}
	if missing.Phase = "run"; DepsMissing("csharp", missing) {
		t.Error("a run step does not restore")
	}
	if DepsMissing("rust", RunResult{Phase: "build", Stderr: "error NU1101"}) {
		t.Error("rust has no Missing")
	}
}

// TestTestProgramsLocal runs a passing and a failing test file with the test runner of each language on this machine. A language whose toolchain is not here is skipped.
func TestTestProgramsLocal(t *testing.T) {
	for _, tc := range []struct {
		lang, file, pass, failing, want string
	}{
		{"python", "test_main.py", "import unittest\n\nclass T(unittest.TestCase):\n    def test_one(self):\n        self.assertEqual(1 + 1, 2)\n", "import unittest\n\nclass T(unittest.TestCase):\n    def test_one(self):\n        self.assertEqual(1 + 1, 3)\n", "test_one (test_main.T.test_one) ... ok"},
		{"javascript", "main.test.js", "import test from 'node:test';\nimport assert from 'node:assert/strict';\ntest('one', () => assert.equal(1 + 1, 2));\n", "import test from 'node:test';\nimport assert from 'node:assert/strict';\ntest('one', () => assert.equal(1 + 1, 3));\n", "✔ one"},
		{"typescript", "main.test.ts", "import test from 'node:test';\nimport assert from 'node:assert/strict';\nconst n: number = 2;\ntest('one', () => assert.equal(1 + 1, n));\n", "import test from 'node:test';\nimport assert from 'node:assert/strict';\nconst n: number = 3;\ntest('one', () => assert.equal(1 + 1, n));\n", "✔ one"},
		{"go", "main_test.go", "package cart\n\nimport \"testing\"\n\nfunc TestOne(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal(\"bad\")\n\t}\n}\n", "package cart\n\nimport \"testing\"\n\nfunc TestOne(t *testing.T) {\n\tt.Fatal(\"bad\")\n}\n", "--- PASS: TestOne"},
		{"rust", "main.rs", "#[test]\nfn one() {\n    assert_eq!(1 + 1, 2);\n}\n", "#[test]\nfn one() {\n    assert_eq!(1 + 1, 3);\n}\n", "test one ... ok"},
	} {
		t.Run(tc.lang, func(t *testing.T) {
			r := &Runner{TempDir: t.TempDir()}
			p, err := TestProgramFor(tc.lang, false)
			if err != nil {
				t.Fatal(err)
			}
			if err := r.ProgramReady(p); err != nil {
				t.Skip(err)
			}
			for _, failing := range []bool{false, true} {
				code := tc.pass
				if failing {
					code = tc.failing
				}
				files := map[string]string{tc.file: code}
				got, err := r.RunProgram(p.WithFiles(files, tc.file), files, "", nil, time.Minute)
				if err != nil {
					t.Fatal(err)
				}
				report := got.Stdout + got.Stderr
				if failing {
					if got.Exit == 0 || got.Phase != "run" {
						t.Errorf("a failed test gives %+v", got)
					}
					continue
				}
				if got.Exit != 0 || !strings.Contains(report, tc.want) {
					t.Errorf("the tests give %+v; want %q", got, tc.want)
				}
			}
		})
	}
}

// TestJavaTestLocal runs a JUnit test on this machine. It needs javac and java on the PATH and the JUnit jar in the cache of prog, so it does not download it.
func TestJavaTestLocal(t *testing.T) {
	r := &Runner{TempDir: t.TempDir()}
	p, _ := TestProgramFor("java", false)
	if err := r.ProgramReady(p); err != nil {
		t.Skip(err)
	}
	cache, err := r.cacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(cache, "junit-platform-console-standalone-"+JUnitVersion+".jar")); err != nil {
		t.Skip("the JUnit jar is not in the cache")
	}
	if p, err = r.ResolveTest(p); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"CartTest.java": "import org.junit.jupiter.api.Test;\nimport static org.junit.jupiter.api.Assertions.assertEquals;\n\nclass CartTest {\n    @Test\n    void total() {\n        assertEquals(2, 1 + 1);\n    }\n}\n"}
	got, err := r.RunProgram(p.WithFiles(files, "CartTest.java"), files, "", nil, time.Minute)
	if err != nil || got.Exit != 0 || !strings.Contains(got.Stdout, "total() ✔") || !strings.Contains(got.Stdout, "1 tests successful") {
		t.Errorf("the tests give %+v, %v", got, err)
	}
}

// TestCsharpTestLocal runs an xUnit test with dotnet test on this machine. It needs dotnet on the PATH and the packages of the .csproj in the NuGet folder. When they are not there, the block must be a skip.
func TestCsharpTestLocal(t *testing.T) {
	if _, err := exec.LookPath("dotnet"); err != nil {
		t.Skip("no dotnet")
	}
	lock, err := filepath.Abs(filepath.Join("testdata", "deps", "csharp", "Cart.Tests.csproj"))
	if err != nil {
		t.Fatal(err)
	}
	r := &Runner{Cache: t.TempDir(), TempDir: t.TempDir()}
	cs, _ := TestProgramFor("csharp", true)
	p, err := r.WithDeps(cs, lock)
	if errors.As(err, new(DepsSkip)) {
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"CartTests.cs": "using Xunit;\n\npublic class CartTests\n{\n    [Fact]\n    public void Total() => Assert.Equal(2, 1 + 1);\n}\n"}
	got, err := r.RunProgram(p.WithFiles(files, "CartTests.cs"), files, "", nil, 2*time.Minute)
	if err != nil || got.Exit != 0 || !strings.Contains(got.Stdout, "Passed!") {
		t.Errorf("the tests give %+v, %v", got, err)
	}
}

// TestRustTestDepsLocal runs a test block with the crate of testdata/deps/rust through cargo test on this machine. It needs cargo and the crate in the cargo cache.
func TestRustTestDepsLocal(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("no cargo")
	}
	lock, err := filepath.Abs(filepath.Join("testdata", "deps", "rust", "Cargo.lock"))
	if err != nil {
		t.Fatal(err)
	}
	r := &Runner{Cache: t.TempDir(), TempDir: t.TempDir()}
	rust, _ := TestProgramFor("rust", true)
	p, err := r.WithDeps(rust, lock)
	if errors.As(err, new(DepsSkip)) {
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	code := "pub fn show(n: i64) -> String {\n    itoa::Buffer::new().format(n).to_string()\n}\n\n#[test]\nfn negative() {\n    assert_eq!(show(-42), \"-42\");\n}\n"
	files := map[string]string{p.Main: code}
	got, err := r.RunProgram(p.WithFiles(files, p.Main), files, "", nil, 2*time.Minute)
	if err != nil || got.Exit != 0 || !strings.Contains(got.Stdout, "test negative ... ok") {
		t.Errorf("the tests give %+v, %v", got, err)
	}
}
