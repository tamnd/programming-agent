// SPDX-License-Identifier: Apache-2.0

package lang

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Program is how one language runs one example block (spec 07.4). The block is the file Main. The steps run in order in a new directory; the last step runs the program, and the steps before it build it.
type Program struct {
	ID    string
	Main  string
	Tools []string
	Steps [][]string
	// Version prints the version of the toolchain, for the record of the run.
	Version []string
	// Image is the image that runs the steps in a sandbox. It is empty for the image of the language. A block with deps= runs in an image that has the packages of its lock file (spec 07.7).
	Image string
	// Seed is a directory on this machine. RunProgram copies its files into the work directory before the files of the block. A block with deps= gets the packages of its lock file this way when it runs on this machine.
	Seed string
	// Files are more files that RunProgram writes into the work directory after the seed and before the files of the block, such as the .csproj of a C# test block.
	Files map[string][]byte
	// Test is true for the test runner of a test block (spec 07.11). A block with deps= then uses the Test function of DepsSupport.
	Test bool
	// Report is the stream of the report of a test runner: "" for stdout, or stderr. It is the default stream of the output block of a test block.
	Report string
	// Setup is the number of first steps that only prepare the directory, such as go mod init. The flags of a block go to the step after them.
	Setup int
}

// Programs are the languages that prog run can run. Regex and HTML have no runner, so their blocks are norun. The test blocks use TestPrograms.
var Programs = map[string]Program{
	"python":     {ID: "python", Main: "main.py", Tools: []string{"python3"}, Steps: [][]string{{"python3", "main.py"}}, Version: []string{"python3", "--version"}},
	"javascript": {ID: "javascript", Main: "main.js", Tools: []string{"node"}, Steps: [][]string{{"node", "main.js"}}, Version: []string{"node", "--version"}},
	"typescript": {ID: "typescript", Main: "main.ts", Tools: []string{"node"}, Steps: [][]string{{"node", "main.ts"}}, Version: []string{"node", "--version"}},
	"go":         {ID: "go", Main: "main.go", Tools: []string{"go"}, Steps: [][]string{{"go", "build", "-o", "main", "main.go"}, {"./main"}}, Version: []string{"go", "version"}},
	"rust":       {ID: "rust", Main: "main.rs", Tools: []string{"rustc"}, Steps: [][]string{{"rustc", "--edition", "2024", "-o", "main", "main.rs"}, {"./main"}}, Version: []string{"rustc", "--version"}},
	"c":          {ID: "c", Main: "main.c", Tools: []string{"cc"}, Steps: [][]string{{"cc", "-std=c23", "-Wall", "-o", "main", "main.c", "-lm"}, {"./main"}}, Version: []string{"cc", "--version"}},
	"cpp":        {ID: "cpp", Main: "main.cpp", Tools: []string{"c++"}, Steps: [][]string{{"c++", "-std=c++23", "-Wall", "-o", "main", "main.cpp"}, {"./main"}}, Version: []string{"c++", "--version"}},
	"java":       {ID: "java", Main: "Main.java", Tools: []string{"javac", "java"}, Steps: [][]string{{"javac", "-d", "classes", "Main.java"}, {"java", "-cp", "classes", "Main"}}, Version: []string{"javac", "-version"}},
	"kotlin":     {ID: "kotlin", Main: "main.kt", Tools: []string{"kotlinc", "java"}, Steps: [][]string{{"kotlinc", "-nowarn", "main.kt", "-include-runtime", "-d", "main.jar"}, {"java", "-jar", "main.jar"}}, Version: []string{"kotlinc", "-version"}},
	"csharp":     {ID: "csharp", Main: "main.cs", Tools: []string{"dotnet"}, Steps: [][]string{{"dotnet", "build", "main.cs", "-p:PublishAot=false", "-o", "out", "-v", "q", "-nologo", "-clp:NoSummary;ErrorsOnly"}, {"dotnet", "out/main.dll"}}, Version: []string{"dotnet", "--version"}},
	"bash":       {ID: "bash", Main: "main.sh", Tools: []string{"bash"}, Steps: [][]string{{"bash", "main.sh"}}, Version: []string{"bash", "--version"}},
	"powershell": {ID: "powershell", Main: "main.ps1", Tools: []string{"pwsh"}, Steps: [][]string{{"pwsh", "-NoProfile", "-NonInteractive", "-File", "main.ps1"}}, Version: []string{"pwsh", "--version"}},
	"sql":        {ID: "sql", Main: "main.sql", Tools: []string{"sqlite3"}, Steps: [][]string{{"sqlite3", "-bail", ":memory:", ".read main.sql"}}, Version: []string{"sqlite3", "--version"}},
}

// ProgramReady tells if the runner can run the blocks of the language. In a sandbox, the image of the language replaces the tools of this machine.
func (r *Runner) ProgramReady(p Program) error {
	return r.Available(Recipe{ID: p.ID, Tools: p.Tools, Version: p.Version})
}

// WithTool gives a copy of the program that uses the command to in place of the tool from, such as python3.15 in place of python3.
func (p Program) WithTool(from, to string) Program {
	swap := func(c []string) []string {
		c = append([]string(nil), c...)
		if len(c) > 0 && c[0] == from {
			c[0] = to
		}
		return c
	}
	p.Tools = swap(p.Tools)
	for i := range p.Tools {
		if p.Tools[i] == from {
			p.Tools[i] = to
		}
	}
	steps := make([][]string, len(p.Steps))
	for i, s := range p.Steps {
		steps[i] = swap(s)
	}
	p.Steps = steps
	p.Version = swap(p.Version)
	return p
}

// RunResult is what one run of a block gives.
type RunResult struct {
	Stdout, Stderr string
	Exit           int
	// Phase is "build" when a build step failed, and "run" otherwise.
	Phase    string
	TimedOut bool
	Ms       int64
}

// BuildTimeout is the longest time of one build step of a program, such as javac or cargo build. The timeout of a block applies to the last step only, which runs the program.
var BuildTimeout = 2 * time.Minute

// RunProgram writes the files to a new directory and runs the steps of the language. files maps a file name to its content; the block itself is under the Main name of the language. The output has the paths of the directory removed.
func (r *Runner) RunProgram(p Program, files map[string]string, stdin string, args []string, timeout time.Duration) (RunResult, error) {
	res := RunResult{Phase: "run"}
	dir, err := os.MkdirTemp(r.TempDir, "prog-run-"+p.ID+"-")
	if err != nil {
		return res, err
	}
	if !r.Keep {
		defer os.RemoveAll(dir)
	}
	defer r.reap(dir)
	if p.Seed != "" {
		if err := copyTree(p.Seed, dir); err != nil {
			return res, err
		}
	}
	for name, body := range p.Files {
		if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(name)), body, 0o644); err != nil {
			return res, err
		}
	}
	for name, body := range files {
		full := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return res, err
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			return res, err
		}
	}
	start := time.Now()
	for i, step := range p.Steps {
		last := i == len(p.Steps)-1
		// The block timeout limits the run of the program. A build step has its own longer limit, because a compiler on a cold cache can take more than the time of the whole run (spec 07.2).
		limit := BuildTimeout
		if last {
			limit = timeout
		}
		ctx, cancel := context.WithTimeout(context.Background(), limit)
		c := step
		if last {
			c = append(append([]string(nil), step...), args...)
		}
		argv := c
		if r.Box != nil {
			image := p.Image
			if image == "" {
				image = r.Box.Images[p.ID]
			}
			argv = r.Box.ArgvImage(image, last && stdin != "", []Mount{{Dir: dir}}, c)
		}
		cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), testEnv...)
		var out, errb bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errb
		if last && stdin != "" {
			cmd.Stdin = strings.NewReader(stdin)
		}
		err := cmd.Run()
		timedOut := errors.Is(ctx.Err(), context.DeadlineExceeded)
		cancel()
		res.Stdout, res.Stderr = clean(out.String(), dir), clean(errb.String(), dir)
		if err != nil {
			if timedOut {
				res.TimedOut = true
			}
			var ee *exec.ExitError
			if !errors.As(err, &ee) && !res.TimedOut {
				return res, fmt.Errorf("%s: %w", strings.Join(c, " "), err)
			}
			res.Exit = 1
			if ee != nil && ee.ExitCode() > 0 {
				res.Exit = ee.ExitCode()
			}
			if !last {
				res.Phase = "build"
				// A build tool such as dotnet writes the errors to stdout and its notices to stderr. Keep both streams for the signature check.
				res.Stderr = strings.TrimSpace(res.Stderr + "\n" + res.Stdout)
				res.Stdout = ""
			}
			break
		}
	}
	res.Ms = time.Since(start).Milliseconds()
	return res, nil
}

// rustcSource matches the path of the standard library in a message of rustc. The path holds the commit of the compiler, which changes with each release.
var rustcSource = regexp.MustCompile(`/rustc/[0-9a-f]{40}/`)

// clean replaces the path of the run directory with $WORK, so that the output is the same on each run (spec 07.6). It also replaces the commit in the paths of the Rust standard library with $RUSTC, so that the output is the same with each release of rustc.
func clean(s, dir string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = rustcSource.ReplaceAllString(s, "/rustc/$$RUSTC/")
	// macOS gives the temporary directory as /var/... or /tmp/..., and the tools may print the real path /private/var/... or /private/tmp/.... The real path goes first, because it holds the other one.
	if strings.HasPrefix(dir, "/var/") || strings.HasPrefix(dir, "/tmp/") {
		s = strings.ReplaceAll(s, "/private"+dir, "$WORK")
	}
	for _, d := range []string{dir, filepath.ToSlash(dir)} {
		s = strings.ReplaceAll(s, d, "$WORK")
	}
	return s
}
