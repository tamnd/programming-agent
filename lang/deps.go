// SPDX-License-Identifier: Apache-2.0

package lang

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

// DepsTimeout is the longest time that the packages of one lock file can take to build, in an image or on this machine.
var DepsTimeout = 20 * time.Minute

// Deps is how a language gets the third-party packages of a lock file (spec 07.7). The image of the language does not have them. In a sandbox, the runner builds an image from the image of the language with the packages in it. The image build is the only step with network, and the block runs in the new image with no network. On this machine, the runner gets the packages with the local toolchain and no network, or the block is skipped.
type Deps struct {
	// Lock is the name of the lock file, such as Cargo.lock.
	Lock string
	// With are the other files next to the lock file that the build needs, such as Cargo.toml.
	With []string
	// Optional are files next to the lock file that the build uses when they are there.
	Optional []string
	// Dockerfile is the text of the image with the packages. %[1]s is the image of the language. The build context holds the lock file and the files of With and Optional.
	Dockerfile string
	// Program changes the program of a block so that it uses the packages. files are the lock file and the files of With and Optional, by name. dir is the directory of the packages that Prepare made on this machine. It is "" in the image.
	Program func(p Program, files map[string][]byte, dir string) (Program, error)
	// Prepare gets the packages into dir on this machine, with no network. dir holds the lock file and the files of With and Optional. Prepare is nil when the language has no local support.
	Prepare func(ctx context.Context, dir string) error
}

// DepsSkip is the error of a block with deps= that cannot get its packages here. The block is skipped with the reason deps. It is not a failure.
type DepsSkip struct{ Lang, Why string }

func (e DepsSkip) Error() string { return e.Lang + ": " + e.Why }

// The Maven that the Java image build downloads, and the plugin that copies the jars. The checksum is of the tar file on Maven Central.
const (
	mavenVersion          = "3.9.16"
	mavenSHA256           = "80ffca22aed9e8b9713a232f3394fd81d7f20322df75efdb2b047dbd3e3a23bb"
	mavenCopyDependencies = "org.apache.maven.plugins:maven-dependency-plugin:3.11.0:copy-dependencies"
)

// The paths of the packages in the images that the Dockerfiles below build.
const (
	imageDepsDir    = "/opt/prog/deps"
	imageJavaLib    = "/opt/prog/lib"
	imageCppInclude = "/opt/prog/include"
	imageCppLib     = "/opt/prog/lib"
)

const (
	// depsReady marks a directory of packages on this machine that is complete.
	depsReady = ".prog-ready"
	// depsNoLocalSupportHint ends the reason of a skip on this machine.
	depsNoLocalSupportHint = "use -sandbox docker"
)

// rustDockerfile builds the crates of Cargo.lock on top of the Rust image. The base image sets CARGO_NET_OFFLINE, so the step that fetches sets it to false for itself only.
const rustDockerfile = `FROM %[1]s
USER root
COPY Cargo.toml Cargo.lock /opt/prog/deps/
WORKDIR /opt/prog/deps
# cargo vendor fetches the crates of the lock file into /opt/prog/vendor. This is the only step with network. The config that it prints tells cargo to read the crates from there.
# The build of an empty main compiles the crates once. Then a block compiles only its own code.
RUN mkdir src .cargo && echo 'fn main() {}' > src/main.rs \
 && CARGO_HOME=/tmp/cargo CARGO_NET_OFFLINE=false cargo vendor --locked --versioned-dirs /opt/prog/vendor > .cargo/config.toml \
 && CARGO_HOME=/tmp/cargo cargo build --offline --locked \
 && rm -rf src /tmp/cargo \
 && chmod -R a+rX /opt/prog
WORKDIR /
USER 1000:1000
`

// javaDockerfile copies the jars of pom.xml into /opt/prog/lib. The JVM image has no Maven, so the first stage downloads it, pinned by checksum. The image of the block has the jars and not Maven.
const javaDockerfile = `FROM %[1]s AS fetch
USER root
ADD --checksum=sha256:` + mavenSHA256 + ` https://repo1.maven.org/maven2/org/apache/maven/apache-maven/` + mavenVersion + `/apache-maven-` + mavenVersion + `-bin.tar.gz /tmp/maven.tar.gz
RUN mkdir -p /opt/maven /opt/prog/lib && tar -xzf /tmp/maven.tar.gz -C /opt/maven --strip-components=1
COPY pom.xml /opt/prog/deps/pom.xml
# Maven copies the jars of the pom and the jars that they need into /opt/prog/lib. This is the only step with network.
RUN /opt/maven/bin/mvn -B -q -f /opt/prog/deps/pom.xml -Dmaven.repo.local=/tmp/m2 ` + mavenCopyDependencies + ` -DincludeScope=runtime -DoutputDirectory=/opt/prog/lib

FROM %[1]s
COPY --from=fetch /opt/prog/lib /opt/prog/lib
`

// cppDockerfile installs the ports of vcpkg.json with vcpkg. The C++ runner image has no vcpkg. With such an image, the build stops with a message that starts with "prog: ".
const cppDockerfile = `FROM %[1]s
USER root
COPY vcpkg*.json /opt/prog/deps/
WORKDIR /opt/prog/deps
# vcpkg installs the ports into /opt/prog/vcpkg. This is the only step with network. A manifest with no baseline gets the baseline of the vcpkg in the image.
RUN if ! command -v vcpkg >/dev/null 2>&1; then echo "prog: the C++ image has no vcpkg, so it cannot install the packages of vcpkg.json; add vcpkg to the image, or mark the block norun" >&2; exit 1; fi \
 && if ! grep -q builtin-baseline vcpkg.json && [ ! -f vcpkg-configuration.json ]; then vcpkg x-update-baseline --add-initial-baseline; fi \
 && vcpkg install --x-manifest-root=/opt/prog/deps --x-install-root=/opt/prog/vcpkg \
 && triplet=$(find /opt/prog/vcpkg -mindepth 1 -maxdepth 1 -type d ! -name vcpkg | head -n 1) \
 && mkdir -p "$triplet/include" "$triplet/lib" \
 && ln -s "$triplet/include" /opt/prog/include && ln -s "$triplet/lib" /opt/prog/lib \
 && chmod -R a+rX /opt/prog
WORKDIR /
USER 1000:1000
`

// DepsSupport are the languages that can run a block with deps=. A block of another language with deps= is skipped with the reason deps.
var DepsSupport = map[string]Deps{
	"rust": {Lock: "Cargo.lock", With: []string{"Cargo.toml"}, Dockerfile: rustDockerfile, Program: rustDeps, Prepare: rustPrepare},
	"java": {Lock: "pom.xml", Dockerfile: javaDockerfile, Program: javaDeps, Prepare: javaPrepare},
	"cpp":  {Lock: "vcpkg.json", Optional: []string{"vcpkg-configuration.json"}, Dockerfile: cppDockerfile, Program: cppDeps},
}

// DepsFor returns how the language gets its packages. It returns a DepsSkip when the language has no support for deps=, or no support on this machine when sandbox is false.
func DepsFor(lang string, sandbox bool) (Deps, error) {
	d, ok := DepsSupport[lang]
	if !ok {
		return d, DepsSkip{lang, "prog run cannot get third-party packages for " + lang}
	}
	if !sandbox && d.Prepare == nil {
		return d, DepsSkip{lang, "prog run cannot get the packages of " + d.Lock + " on this machine; " + depsNoLocalSupportHint}
	}
	return d, nil
}

// DepsDockerfile gives the Dockerfile of the image with the packages, on top of the image base.
func DepsDockerfile(d Deps, base string) string {
	return fmt.Sprintf(d.Dockerfile, base)
}

// DepsTag gives the tag of the image with the packages. It is a hash of the image of the language, the Dockerfile and the files of the build context. The same inputs give the same tag, so a second run finds the image and does not build it again.
func DepsTag(lang, base, dockerfile string, files map[string][]byte) string {
	return "prog-deps-" + lang + ":" + depsHash(lang, base, dockerfile, files)
}

func depsHash(lang, base, dockerfile string, files map[string][]byte) string {
	h := sha256.New()
	for _, s := range []string{lang, base, dockerfile} {
		fmt.Fprintf(h, "%d\x00%s\x00", len(s), s)
	}
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	slices.Sort(names)
	for _, n := range names {
		fmt.Fprintf(h, "%d\x00%s\x00%d\x00", len(n), n, len(files[n]))
		h.Write(files[n])
	}
	return hex.EncodeToString(h.Sum(nil))[:24]
}

// readDeps reads the lock file and the files next to it that the build needs.
func readDeps(lock string, d Deps) (map[string][]byte, error) {
	files := map[string][]byte{}
	dir := filepath.Dir(lock)
	for i, name := range append(append([]string{filepath.Base(lock)}, d.With...), d.Optional...) {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			if i > len(d.With) && errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("deps: %w", err)
		}
		files[name] = b
	}
	return files, nil
}

// WithDeps changes the program of a block with deps= so that it uses the packages of the lock file (spec 07.7). In a sandbox, it builds the image with the packages, or reuses it. On this machine, it gets the packages with no network into the cache directory. It returns a DepsSkip when the packages cannot be had here.
func (r *Runner) WithDeps(p Program, lock string) (Program, error) {
	d, err := DepsFor(p.ID, r.Box != nil)
	if err != nil {
		return p, err
	}
	if filepath.Base(lock) != d.Lock {
		return p, fmt.Errorf("deps= names %s, but the %s runner reads %s", filepath.Base(lock), p.ID, d.Lock)
	}
	files, err := readDeps(lock, d)
	if err != nil {
		return p, err
	}
	if r.Box != nil {
		image, err := r.Box.Derive(p.ID, d, files)
		if err != nil {
			return p, err
		}
		q, err := d.Program(p, files, "")
		q.Image = image
		return q, err
	}
	dir, err := r.prepareDeps(p, d, files)
	if err != nil {
		return p, err
	}
	return d.Program(p, files, dir)
}

// depsJob makes the packages of one lock file once, for all the blocks that use it.
type depsJob struct {
	once sync.Once
	err  error
}

// depsJobs holds a job for each directory or image of packages. A Runner is copied by value in places, so the jobs are not in it.
var depsJobs sync.Map

func jobFor(key string) *depsJob {
	j, _ := depsJobs.LoadOrStore(key, &depsJob{})
	return j.(*depsJob)
}

// prepareDeps makes the packages of the lock file in the cache directory, once. The key has the version of the toolchain, so a new compiler makes the packages again.
func (r *Runner) prepareDeps(p Program, d Deps, files map[string][]byte) (string, error) {
	cache, err := r.cacheDir()
	if err != nil {
		return "", err
	}
	tool := Recipe{ID: p.ID, Tools: p.Tools, Version: p.Version}.ToolVersion()
	dir := filepath.Join(cache, "deps", p.ID+"-"+depsHash(p.ID, "local "+tool, "", files))
	j := jobFor(dir)
	j.once.Do(func() { j.err = makeDeps(d, files, dir) })
	return dir, j.err
}

func makeDeps(d Deps, files map[string][]byte, dir string) error {
	if _, err := os.Stat(filepath.Join(dir, depsReady)); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(filepath.Dir(dir), filepath.Base(dir)+".part-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	for name, b := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), b, 0o644); err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), DepsTimeout)
	defer cancel()
	if err := d.Prepare(ctx, tmp); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmp, depsReady), nil, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, dir); err != nil {
		// Another prog run made the same directory first.
		if _, serr := os.Stat(filepath.Join(dir, depsReady)); serr == nil {
			return nil
		}
		return err
	}
	return nil
}

// cacheDir is the cache directory of prog: Cache, or prog in the user cache directory.
func (r *Runner) cacheDir() (string, error) {
	if r.Cache != "" {
		return r.Cache, nil
	}
	c, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(c, "prog"), nil
}

// runIn runs a command in dir and returns its output.
func runIn(ctx context.Context, dir string, c ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, c[0], c[1:]...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), testEnv...)
	return cmd.CombinedOutput()
}

func cloneSteps(p Program) Program {
	steps := make([][]string, len(p.Steps))
	for i, s := range p.Steps {
		steps[i] = slices.Clone(s)
	}
	p.Steps = steps
	return p
}

// insertAfterTool puts args after the command name of a step.
func insertAfterTool(step []string, args ...string) []string {
	return append(append(slices.Clone(step[:1]), args...), step[1:]...)
}

// cargoBuild builds the block and its crates. The crates are compiled already, so only the block compiles. The profile is dev, so a block with deps= has the overflow checks of a block that rustc builds.
var cargoBuild = []string{"cargo", "build", "--offline", "--locked", "--quiet"}

// rustDeps makes the block the main of a Cargo project with the Cargo.toml and the Cargo.lock of the node. In the image, the build step copies the project with the compiled crates from /opt/prog/deps first. On this machine, the project comes from the cache directory as the seed of the run.
func rustDeps(p Program, files map[string][]byte, dir string) (Program, error) {
	bin, err := cargoBin(files["Cargo.toml"])
	if err != nil {
		return p, err
	}
	p = cloneSteps(p)
	p.Main = "src/main.rs"
	p.Tools = []string{"cargo"}
	build := cargoBuild
	if dir == "" {
		build = []string{"sh", "-c", "cp -pR " + imageDepsDir + "/. . && exec " + strings.Join(cargoBuild, " ")}
	} else {
		p.Seed = dir
	}
	p.Steps = [][]string{build, {"./target/debug/" + bin}}
	return p, nil
}

var (
	tomlSection = regexp.MustCompile(`^\s*\[\[?\s*([A-Za-z0-9_.-]+)\s*\]\]?`)
	tomlName    = regexp.MustCompile(`^\s*name\s*=\s*["']([^"']+)["']`)
)

// cargoBin gives the name of the binary of a Cargo.toml: the name of the first [[bin]], or the name of the package.
func cargoBin(toml []byte) (string, error) {
	section, pkg, bin := "", "", ""
	for _, line := range strings.Split(string(toml), "\n") {
		if m := tomlSection.FindStringSubmatch(line); m != nil {
			section = m[1]
			continue
		}
		if m := tomlName.FindStringSubmatch(line); m != nil {
			switch {
			case section == "package" && pkg == "":
				pkg = m[1]
			case section == "bin" && bin == "":
				bin = m[1]
			}
		}
	}
	if bin != "" {
		return bin, nil
	}
	if pkg == "" {
		return "", errors.New("deps: the Cargo.toml next to the Cargo.lock has no package name")
	}
	return pkg, nil
}

// rustPrepare compiles the crates of the lock file on this machine with no network. The crates must be in the cargo cache already.
func rustPrepare(ctx context.Context, dir string) error {
	if _, err := exec.LookPath("cargo"); err != nil {
		return DepsSkip{"rust", "cargo is not installed"}
	}
	if err := os.MkdirAll(filepath.Join(dir, "src"), 0o755); err != nil {
		return err
	}
	if err := write(filepath.Join(dir, "src"), "main.rs", "fn main() {}\n"); err != nil {
		return err
	}
	if _, err := runIn(ctx, dir, "cargo", "fetch", "--offline", "--locked"); err != nil {
		return DepsSkip{"rust", "the crates of Cargo.lock are not all in the cargo cache of this machine, and the run has no network; " + depsNoLocalSupportHint}
	}
	if out, err := runIn(ctx, dir, cargoBuild...); err != nil {
		return fmt.Errorf("deps: cargo build of the crates of Cargo.lock failed:\n%s", tail(string(out), 2000))
	}
	return os.RemoveAll(filepath.Join(dir, "src"))
}

// javaDeps puts the jars of the pom on the class path of javac and java.
func javaDeps(p Program, _ map[string][]byte, dir string) (Program, error) {
	lib, sep := imageJavaLib+"/*", ":"
	if dir != "" {
		lib, sep = filepath.Join(dir, "lib", "*"), string(os.PathListSeparator)
	}
	// The first step is javac and the last step is java. The command names can be paths, because of -tool.
	p = cloneSteps(p)
	p.Steps[0] = insertAfterTool(p.Steps[0], "-cp", lib)
	for _, s := range p.Steps[1:] {
		if k := slices.Index(s, "-cp"); k >= 0 && k+1 < len(s) {
			s[k+1] += sep + lib
		}
	}
	return p, nil
}

// javaPrepare copies the jars of the pom into dir/lib with the Maven of this machine, offline. The jars and the dependency plugin must be in the local Maven repository already.
func javaPrepare(ctx context.Context, dir string) error {
	if _, err := exec.LookPath("mvn"); err != nil {
		return DepsSkip{"java", "mvn is not installed"}
	}
	lib := filepath.Join(dir, "lib")
	if err := os.MkdirAll(lib, 0o755); err != nil {
		return err
	}
	if _, err := runIn(ctx, dir, "mvn", "-o", "-B", "-q", "-f", filepath.Join(dir, "pom.xml"), mavenCopyDependencies, "-DincludeScope=runtime", "-DoutputDirectory="+lib); err != nil {
		return DepsSkip{"java", "the jars of pom.xml are not all in the local Maven repository of this machine, and the run has no network; " + depsNoLocalSupportHint}
	}
	return nil
}

// cppDeps puts the headers and the libraries of the ports on the command of the compiler. A block that links a library gives -l with flags=.
func cppDeps(p Program, _ map[string][]byte, _ string) (Program, error) {
	p = cloneSteps(p)
	p.Steps[0] = insertAfterTool(p.Steps[0], "-I"+imageCppInclude, "-L"+imageCppLib)
	return p, nil
}

// derivedImage is the build of one image with packages, once for all the blocks that use it.
type derivedImage struct {
	once sync.Once
	err  error
}

// Derive builds the image of the language with the packages of a lock file, or reuses it (spec 07.7). The tag comes from DepsTag. The build is the only step with network.
func (s *Sandbox) Derive(lang string, d Deps, files map[string][]byte) (string, error) {
	base, ok := s.Images[lang]
	if !ok {
		return "", Skip{lang, "image"}
	}
	// The Dockerfile names the image of the language by a local tag. Some builders cannot read an image ID in FROM, and the image workflow tests an image by its ID.
	alias := "prog-base-" + lang + ":" + depsHash(lang, base, "", nil)
	df := DepsDockerfile(d, alias)
	tag := DepsTag(lang, base, df, files)
	s.mu.Lock()
	if s.derived == nil {
		s.derived = map[string]*derivedImage{}
	}
	j := s.derived[tag]
	if j == nil {
		j = &derivedImage{}
		s.derived[tag] = j
	}
	s.mu.Unlock()
	j.once.Do(func() { j.err = s.build(base, alias, tag, df, d.Lock, files) })
	return tag, j.err
}

func (s *Sandbox) build(base, alias, tag, dockerfile, lock string, files map[string][]byte) error {
	if exec.Command(s.Engine, "image", "inspect", tag).Run() == nil {
		return nil
	}
	if out, err := exec.Command(s.Engine, "tag", base, alias).CombinedOutput(); err != nil {
		return fmt.Errorf("tag %s: %w: %s", base, err, strings.TrimSpace(string(out)))
	}
	dir, err := os.MkdirTemp("", "prog-deps-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	for name, b := range files {
		if err := os.WriteFile(filepath.Join(dir, name), b, 0o644); err != nil {
			return err
		}
	}
	df := filepath.Join(dir, "Dockerfile")
	if err := os.WriteFile(df, []byte(dockerfile), 0o644); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), DepsTimeout)
	defer cancel()
	args := []string{"build", "--label", "prog.deps=1", "-t", tag, "-f", df}
	if filepath.Base(s.Engine) == "docker" {
		args = append(args, "--progress=plain")
	}
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, s.Engine, append(args, dir)...)
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("the image with the packages of %s did not build: %s", lock, buildWhy(out.String(), err))
	}
	return nil
}

// progLine finds a message of a Dockerfile of this package in the build output. Docker shows it as "#7 0.25 prog: ...", and Podman shows the line as it is.
var progLine = regexp.MustCompile(`(?m)^(?:#\d+ [\d.]+ )?prog: (.+)$`)

// buildWhy gives the reason of a failed image build: the message of the Dockerfile, or the end of the output.
func buildWhy(out string, err error) string {
	if m := progLine.FindStringSubmatch(out); m != nil {
		return strings.TrimSpace(m[1])
	}
	return fmt.Sprintf("%v\n%s", err, tail(out, 2000))
}

// copyTree copies the files of src into dst. It keeps the modes and the modification times, so that a build tool such as cargo sees the compiled packages as fresh.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		if rel == depsReady {
			return nil
		}
		to := filepath.Join(dst, rel)
		info, err := e.Info()
		if err != nil {
			return err
		}
		switch {
		case e.IsDir():
			if err := os.MkdirAll(to, info.Mode().Perm()|0o700); err != nil {
				return err
			}
		case info.Mode()&fs.ModeSymlink != 0:
			link, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(link, to)
		default:
			if err := copyFile(p, to, info.Mode().Perm()|0o600); err != nil {
				return err
			}
		}
		return os.Chtimes(to, info.ModTime(), info.ModTime())
	})
}

func copyFile(from, to string, mode fs.FileMode) error {
	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
