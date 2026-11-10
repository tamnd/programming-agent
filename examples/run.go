// SPDX-License-Identifier: Apache-2.0

package examples

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tamnd/programming-agent/lang"
)

// Result is the result of one case.
type Result struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	ID     string `json:"id"`
	Lang   string `json:"lang"`
	Mode   string `json:"mode"`
	Status string `json:"status"` // pass, fail or skip
	Why    string `json:"why,omitempty"`
	// Reason is a short code for the cause of a skip, such as deps when prog run cannot get the packages of the block here (spec 07.7).
	Reason string `json:"reason,omitempty"`
	Tool   string `json:"tool,omitempty"`
	Ms     int64  `json:"ms"`
	// Updated is set by Update when the output blocks of the case were written from this run.
	Updated bool `json:"updated,omitempty"`
	// Got is the output of the run, for -update.
	Got lang.RunResult `json:"-"`
	// Case is the case of the result.
	Case *Case `json:"-"`
}

// Runner runs the cases. It keeps the version of each toolchain.
type Runner struct {
	R  *lang.Runner
	mu sync.Mutex
	// versions caches the toolchain version of each language.
	versions map[string]string
	ready    map[string]error
}

// versioned are the languages whose toolchain version is the version of the language, so that since=, until= and version= can be checked against it.
var versioned = []string{"python", "go", "rust", "java", "kotlin", "bash", "powershell", "sql", "javascript"}

func (x *Runner) toolchain(p lang.Program) (string, error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if x.versions == nil {
		x.versions, x.ready = map[string]string{}, map[string]error{}
	}
	if err, ok := x.ready[p.ID]; ok {
		return x.versions[p.ID], err
	}
	err := x.R.ProgramReady(p)
	x.ready[p.ID] = err
	if err == nil {
		x.versions[p.ID] = x.R.ToolVersion(lang.Recipe{ID: p.ID, Tools: p.Tools, Version: p.Version})
	}
	return x.versions[p.ID], err
}

// Run runs one case and checks it.
func (x *Runner) Run(c *Case) Result {
	b := c.Block
	res := Result{File: c.File, Line: b.Line, ID: c.ID(), Lang: b.Lang, Mode: c.Mode(), Case: c}
	skip := func(why string) Result { res.Status, res.Why = "skip", why; return res }
	fail := func(format string, a ...any) Result {
		res.Status, res.Why = "fail", fmt.Sprintf(format, a...)
		return res
	}
	p, ok := lang.Programs[b.Lang]
	if !ok {
		return skip("prog run has no runner for " + b.Lang)
	}
	if pl := b.Attrs["platform"]; pl != "" && pl != "linux" {
		return skip("the block runs in the scheduled job for " + pl)
	}
	if impl := b.Attrs["impl"]; impl != "" {
		return skip("the block runs on the implementation " + impl + " in the version matrix")
	}
	skipDeps := func(err error) Result {
		var ds lang.DepsSkip
		if errors.As(err, &ds) {
			res.Reason = "deps"
			return skip("deps: " + ds.Why)
		}
		return fail("%v", err)
	}
	lock := c.LockPath()
	if lock != "" {
		// A language with no support for deps= is a skip, not a failure, so look before the toolchain.
		if _, err := lang.DepsFor(p.ID, x.R.Box != nil); err != nil {
			return skipDeps(err)
		}
	}
	tool, err := x.toolchain(p)
	if err != nil {
		return skip(err.Error())
	}
	res.Tool = tool
	if slices.Contains(versioned, p.ID) {
		have := versionOf(tool)
		if v := b.Attrs["since"]; v != "" && have != "" && compareVersions(have, v) < 0 {
			return skip(fmt.Sprintf("the block needs version %s and the toolchain is %s", v, have))
		}
		if v := b.Attrs["until"]; v != "" && have != "" && compareVersions(have, v) > 0 {
			return skip(fmt.Sprintf("the block works up to version %s and the toolchain is %s", v, have))
		}
		if v := b.Attrs["version"]; v != "" && v != "matrix" && have != "" && !sameVersion(have, v) {
			return skip(fmt.Sprintf("the block is for version %s and the toolchain is %s", v, have))
		}
	}
	p = adjust(p, b)
	if lock != "" {
		if p, err = x.R.WithDeps(p, lock); err != nil {
			return skipDeps(err)
		}
	}
	timeout := 10 * time.Second
	if t, err := strconv.Atoi(b.Attrs["timeout"]); err == nil && t > 0 && t <= 30 {
		timeout = time.Duration(t) * time.Second
	}
	got, err := x.R.RunProgram(p, c.Files(p), c.Stdin, strings.Fields(b.Attrs["args"]), timeout)
	res.Got, res.Ms = got, got.Ms
	if err != nil {
		return fail("%v", err)
	}
	switch c.Mode() {
	case "run", "wrong":
		if got.TimedOut {
			return fail("the program ran for more than %s", timeout)
		}
		if got.Exit != 0 {
			return fail("the program failed in the %s step with exit status %d:\n%s", got.Phase, got.Exit, tail(got.Stderr+got.Stdout))
		}
	case "fails":
		if why := checkFailure(b.Attrs["expect"], got); why != "" {
			return fail("%s", why)
		}
	}
	for _, o := range c.Outputs {
		if why := compare(o, got); why != "" {
			return fail("%s", why)
		}
	}
	res.Status = "pass"
	return res
}

// adjust applies the std=, flags= and sanitizer settings of a block to the build step of the language.
func adjust(p lang.Program, b *Block) lang.Program {
	steps := make([][]string, len(p.Steps))
	for i, s := range p.Steps {
		steps[i] = slices.Clone(s)
	}
	p.Steps = steps
	first := p.Steps[0]
	if std := b.Attrs["std"]; std != "" && (p.ID == "c" || p.ID == "cpp") {
		for i, a := range first {
			if strings.HasPrefix(a, "-std=") {
				first[i] = "-std=" + std
			}
		}
	}
	var extra []string
	if kind, ok := strings.CutPrefix(b.Attrs["expect"], "sanitizer:"); ok && (p.ID == "c" || p.ID == "cpp") {
		extra = append(extra, "-g", "-fsanitize="+kind, "-fno-sanitize-recover=all")
	}
	flags := strings.Fields(b.Attrs["flags"])
	extra = append(extra, flags...)
	// A flag of the block replaces the default flag of the same option, such as --edition 2021 or -std=c17.
	for _, fl := range flags {
		switch {
		case fl == "--edition" || strings.HasPrefix(fl, "--edition="):
			first = dropOption(first, "--edition")
		case strings.HasPrefix(fl, "-std="):
			first = dropOption(first, "-std=")
		}
	}
	if len(extra) > 0 {
		first = append(first[:1], append(extra, first[1:]...)...)
	}
	p.Steps[0] = first
	return p
}

// dropOption removes an option from a command: "--edition 2024", "--edition=2024" or "-std=c23".
func dropOption(c []string, opt string) []string {
	var out []string
	for i := 0; i < len(c); i++ {
		switch {
		case c[i] == opt:
			i++
		case strings.HasPrefix(c[i], opt) && (strings.HasSuffix(opt, "=") || strings.HasPrefix(c[i], opt+"=")):
		default:
			out = append(out, c[i])
		}
	}
	return out
}

// checkFailure checks the failure signature of a fails block (rule X04). It returns why the run does not match, or "".
func checkFailure(sig string, got lang.RunResult) string {
	kind, arg, _ := strings.Cut(sig, ":")
	if got.TimedOut {
		if kind == "timeout" {
			return ""
		}
		return "the program ran out of time, but the block expects " + sig
	}
	if kind == "timeout" {
		return "the block expects a timeout, but the program stopped"
	}
	if got.Exit == 0 {
		return fmt.Sprintf("the block expects the failure %s, but the program exited with status 0", sig)
	}
	errText := got.Stderr + "\n" + got.Stdout
	switch kind {
	case "compile":
		if got.Phase == "build" || compileErrorRE.MatchString(errText) {
			return ""
		}
		return "the block expects a compile error, but the program failed when it ran:\n" + tail(errText)
	case "runtime":
		if got.Phase == "run" {
			return ""
		}
		return "the block expects a runtime error, but the build failed:\n" + tail(errText)
	case "exit":
		if strconv.Itoa(got.Exit) == arg {
			return ""
		}
		return fmt.Sprintf("the block expects the exit status %s, but the status is %d", arg, got.Exit)
	case "exception":
		re := regexp.MustCompile(`(?m)(^|[\s.:'"(])` + regexp.QuoteMeta(arg) + `($|[\s:'"),])`)
		if got.Phase == "run" && re.MatchString(errText) {
			return ""
		}
		return fmt.Sprintf("the block expects the exception %s:\n%s", arg, tail(errText))
	case "error":
		re := regexp.MustCompile(`(^|\W)` + regexp.QuoteMeta(arg) + `(\W|$)`)
		if re.MatchString(errText) {
			return ""
		}
		if codes := diagCodeRE.FindAllString(errText, 5); len(codes) > 0 {
			return fmt.Sprintf("the block expects the error %s, but the error is %s", arg, strings.Join(codes, ", "))
		}
		return fmt.Sprintf("the block expects the error %s:\n%s", arg, tail(errText))
	case "panic":
		if strings.Contains(errText, "panic:") || strings.Contains(errText, "panicked at") {
			return ""
		}
		return "the block expects a panic:\n" + tail(errText)
	case "sanitizer":
		mark := map[string]string{"address": "AddressSanitizer", "undefined": "runtime error:", "leak": "LeakSanitizer", "thread": "ThreadSanitizer", "memory": "MemorySanitizer"}[arg]
		if mark == "" {
			mark = arg
		}
		if strings.Contains(errText, mark) {
			return ""
		}
		return fmt.Sprintf("the block expects a report of the %s sanitizer:\n%s", arg, tail(errText))
	}
	return "the failure signature " + sig + " is not known"
}

var (
	compileErrorRE = regexp.MustCompile(`(?m)^\s*(SyntaxError|IndentationError|TabError)\b|error\[E\d{4}\]|error TS\d+|error CS\d+`)
	diagCodeRE     = regexp.MustCompile(`\b(E\d{4}|TS\d{4,5}|CS\d{4}|SC\d{4})\b`)
)

// compare checks the output of a run against an output block. It returns why they differ, or "".
func compare(o *Block, got lang.RunResult) string {
	match := o.Attrs["match"]
	if match == "any" {
		return ""
	}
	actual := Stream(got, o.Attrs["stream"])
	norm := strings.Split(o.Attrs["normalize"], ",")
	want, have := Normalize(o.Code, norm), Normalize(actual, norm)
	switch match {
	case "", "exact":
		if want == have {
			return ""
		}
	case "unordered":
		if sortedLines(want) == sortedLines(have) {
			return ""
		}
	case "regex":
		wl, hl := lines(want), lines(have)
		if len(wl) == len(hl) {
			ok := true
			for i := range wl {
				re, err := regexp.Compile("^(?:" + wl[i] + ")$")
				if err != nil {
					return fmt.Sprintf("line %d of the output block is not a regular expression: %v", i+1, err)
				}
				if !re.MatchString(hl[i]) {
					ok = false
					break
				}
			}
			if ok {
				return ""
			}
		}
	default:
		return "match is exact, regex, unordered or any, not " + match
	}
	return fmt.Sprintf("the output at line %d differs from the run:\n%s", o.Line, diff(want, have))
}

// Stream gives the text of a stream of the run: stdout, stderr or both.
func Stream(got lang.RunResult, stream string) string {
	switch stream {
	case "stderr":
		return got.Stderr
	case "both":
		if got.Stdout != "" && got.Stderr != "" && !strings.HasSuffix(got.Stdout, "\n") {
			return got.Stdout + "\n" + got.Stderr
		}
		return got.Stdout + got.Stderr
	}
	return got.Stdout
}

var (
	ansiRE  = regexp.MustCompile(`\x1b\[[0-9;?]*[A-Za-z]`)
	addrRE  = regexp.MustCompile(`0x[0-9a-fA-F]{4,}`)
	timeRE  = regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:?\d{2})?|\b\d{2}:\d{2}:\d{2}(\.\d+)?\b`)
	tmpRE   = regexp.MustCompile(`(/private)?(/tmp|/var/folders)/[^\s'"):]+`)
	hashRE  = regexp.MustCompile(`\b[0-9a-f]{7,64}\b`)
	floatRE = regexp.MustCompile(`-?\d+\.\d+(e[+-]?\d+)?`)
)

// Normalize applies the default normalizers of spec 07.6 and the extra ones that are named: addr, time, tmp, hash and float:<digits>.
func Normalize(s string, extra []string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = ansiRE.ReplaceAllString(s, "")
	for _, n := range extra {
		name, arg, _ := strings.Cut(strings.TrimSpace(n), ":")
		switch name {
		case "addr":
			s = addrRE.ReplaceAllString(s, "0x…")
		case "time":
			s = timeRE.ReplaceAllString(s, "<time>")
		case "tmp":
			s = tmpRE.ReplaceAllString(s, "<tmp>")
		case "hash":
			s = hashRE.ReplaceAllString(s, "<hash>")
		case "float":
			d, err := strconv.Atoi(arg)
			if err != nil {
				d = 6
			}
			s = floatRE.ReplaceAllStringFunc(s, func(f string) string {
				v, err := strconv.ParseFloat(f, 64)
				if err != nil || math.IsInf(v, 0) {
					return f
				}
				return strconv.FormatFloat(v, 'g', d, 64)
			})
		}
	}
	ls := lines(s)
	for i := range ls {
		ls[i] = strings.TrimRight(ls[i], " \t")
	}
	for len(ls) > 0 && ls[len(ls)-1] == "" {
		ls = ls[:len(ls)-1]
	}
	if len(ls) == 0 {
		return ""
	}
	return strings.Join(ls, "\n") + "\n"
}

// Clean gives the text that -update writes in an output block: the run output with the line ends, the colour codes and the trailing spaces made the same as Normalize makes them.
func Clean(s string) string { return Normalize(s, nil) }

func lines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n")
}

func sortedLines(s string) string {
	l := lines(s)
	slices.Sort(l)
	return strings.Join(l, "\n")
}

// diff shows the first line where two outputs differ, with the lines around it.
func diff(want, have string) string {
	w, h := lines(want), lines(have)
	i := 0
	for i < len(w) && i < len(h) && w[i] == h[i] {
		i++
	}
	var b strings.Builder
	show := func(sign string, l []string) {
		for j := i; j < len(l) && j < i+3; j++ {
			fmt.Fprintf(&b, "  %s %d | %s\n", sign, j+1, l[j])
		}
		if i >= len(l) {
			fmt.Fprintf(&b, "  %s (no line %d)\n", sign, i+1)
		}
	}
	show("want", w)
	show("got ", h)
	return strings.TrimSuffix(b.String(), "\n")
}

func tail(s string) string {
	l := lines(strings.TrimSpace(s))
	if len(l) > 12 {
		l = append([]string{"…"}, l[len(l)-12:]...)
	}
	return "  " + strings.Join(l, "\n  ")
}

var verRE = regexp.MustCompile(`\d+(\.\d+)+|\d+`)

// versionOf finds the version number in the output of a version command, such as 3.15.0 in "Python 3.15.0rc1".
func versionOf(tool string) string {
	tool = strings.TrimPrefix(tool, "go version go")
	return verRE.FindString(tool)
}

func compareVersions(a, b string) int {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(as) && i < len(bs); i++ {
		x, _ := strconv.Atoi(as[i])
		y, _ := strconv.Atoi(bs[i])
		if x != y {
			return x - y
		}
	}
	return 0
}

// sameVersion tells if the toolchain version have is version v or a patch of it, so that 3.15.0 is 3.15.
func sameVersion(have, v string) bool {
	return have == v || strings.HasPrefix(have, v+".")
}
