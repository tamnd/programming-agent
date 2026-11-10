// SPDX-License-Identifier: Apache-2.0

package lang

import (
	"os"
	"path"
	"slices"
)

// The placeholders in the steps of a test runner. WithFiles replaces them with the names of the files of the block.
const (
	// ArgTest is the file of the test block.
	ArgTest = "{test}"
	// ArgSources is the files of the block that have the extension of Main, in sorted order.
	ArgSources = "{sources}"
)

// junitLauncher is the main class of the JUnit Platform console launcher.
const junitLauncher = "org.junit.platform.console.ConsoleLauncher"

// javaTestSteps compiles the test classes against the JUnit jar and runs the console launcher. The launcher finds the tests in the directories of the class path, so a jar on the class path is not scanned. sep is the separator of the class path.
func javaTestSteps(jar, sep string) [][]string {
	return [][]string{
		{"javac", "-d", "classes", "-cp", jar, ArgSources},
		{"java", "-cp", "classes" + sep + jar, junitLauncher, "execute", "--disable-banner", "--disable-ansi-colors", "--details=tree", "--scan-class-path"},
	}
}

// TestPrograms are the test runners of the test blocks (spec 07.11). A test block is a test file, not a program with a main. The last step runs the tests, and its exit status is 0 when every test passes. The steps can hold ArgTest and ArgSources. C# and C++ have no steps here, because their test frameworks are third-party packages: a block of these languages runs only with deps=, and the steps come from DepsSupport.
var TestPrograms = map[string]Program{
	// unittest writes its report to stderr.
	"python":     {ID: "python", Main: "test_main.py", Tools: []string{"python3"}, Steps: [][]string{{"python3", "-m", "unittest", "-v", ArgTest}}, Version: []string{"python3", "--version"}, Test: true, Report: "stderr"},
	"javascript": {ID: "javascript", Main: "main.test.js", Tools: []string{"node"}, Steps: [][]string{{"node", "--test", "--test-reporter=spec", ArgTest}}, Version: []string{"node", "--version"}, Test: true},
	// Node runs a TypeScript test file by type stripping, as it runs the run blocks.
	"typescript": {ID: "typescript", Main: "main.test.ts", Tools: []string{"node"}, Steps: [][]string{{"node", "--test", "--test-reporter=spec", ArgTest}}, Version: []string{"node", "--version"}, Test: true},
	// go mod init writes the go.mod of the module with the version of the toolchain, so that the language version is the version of the toolchain. -count=1 turns off the cache of the test results.
	"go":     {ID: "go", Main: "main_test.go", Tools: []string{"go"}, Steps: [][]string{{"go", "mod", "init", "example"}, {"go", "test", "-v", "-count=1", "."}}, Version: []string{"go", "version"}, Test: true, Setup: 1},
	"rust":   {ID: "rust", Main: "main.rs", Tools: []string{"rustc"}, Steps: [][]string{{"rustc", "--edition", "2024", "--test", "-o", "main_test", ArgTest}, {"./main_test", "--test-threads=1"}}, Version: []string{"rustc", "--version"}, Test: true},
	"java":   {ID: "java", Main: "MainTest.java", Tools: []string{"javac", "java"}, Steps: javaTestSteps(imageJUnit, ":"), Version: []string{"javac", "-version"}, Test: true},
	"csharp": {ID: "csharp", Main: "MainTests.cs", Tools: []string{"dotnet"}, Version: []string{"dotnet", "--version"}, Test: true},
	"cpp":    {ID: "cpp", Main: "main_test.cpp", Tools: []string{"c++"}, Version: []string{"c++", "--version"}, Test: true},
}

// testNeedsDeps gives the reason of a skip for a test block with no deps= in a language whose test framework is a package.
var testNeedsDeps = map[string]string{
	"csharp": "a C# test block needs deps= with a .csproj that names the test framework, such as xunit, because .NET has no test framework in its base class library",
	"cpp":    "a C++ test block needs GoogleTest or Catch2 through deps= with a vcpkg.json, because the C++ standard library has no test framework",
}

// noTestRunner gives the reason of a skip for a test block of a language with no test runner here.
var noTestRunner = map[string]string{
	"c": "C has no test runner; a C test is a program that uses assert, so write it as a run block (spec 07.11)",
}

// TestSkip is the error of a test block that prog run cannot run. Reason is runner when the language has no test runner here, or deps when the test framework is a package and the block has no deps=. It is a skip, not a failure.
type TestSkip struct{ Lang, Reason, Why string }

func (e TestSkip) Error() string { return e.Lang + ": " + e.Why }

// TestProgramFor gives the test runner of a language. deps tells if the block has deps=. It returns a TestSkip when prog run has no test runner for the language, or when the runner needs deps= and the block has none.
func TestProgramFor(lang string, deps bool) (Program, error) {
	p, ok := TestPrograms[lang]
	if !ok {
		why := noTestRunner[lang]
		if why == "" {
			why = "prog run has no test runner for " + lang
		}
		return p, TestSkip{lang, "runner", why}
	}
	if len(p.Steps) == 0 && !deps {
		return p, TestSkip{lang, "deps", testNeedsDeps[lang]}
	}
	return cloneSteps(p), nil
}

// ResolveTest puts the paths of this machine into a test runner, such as the path of the JUnit jar. In a sandbox, the paths of the image stay. It downloads the JUnit jar once when it is not in the cache.
func (r *Runner) ResolveTest(p Program) (Program, error) {
	if r.Box != nil || p.ID != "java" {
		return p, nil
	}
	jar, err := r.junit()
	if err != nil {
		return p, err
	}
	p.Steps = javaTestSteps(jar, string(os.PathListSeparator))
	return p, nil
}

// WithFiles replaces ArgTest with test, the file of the test block, and ArgSources with the files that have the extension of Main.
func (p Program) WithFiles(files map[string]string, test string) Program {
	var sources []string
	ext := path.Ext(p.Main)
	for name := range files {
		if path.Ext(name) == ext {
			sources = append(sources, name)
		}
	}
	slices.Sort(sources)
	p = cloneSteps(p)
	for i, s := range p.Steps {
		var out []string
		for _, a := range s {
			switch a {
			case ArgTest:
				out = append(out, test)
			case ArgSources:
				out = append(out, sources...)
			default:
				out = append(out, a)
			}
		}
		p.Steps[i] = out
	}
	return p
}
