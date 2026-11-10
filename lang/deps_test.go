// SPDX-License-Identifier: Apache-2.0

package lang

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestDepsTag(t *testing.T) {
	lock := map[string][]byte{"Cargo.lock": []byte("version = 4\n"), "Cargo.toml": []byte("[package]\nname = \"x\"\n")}
	tag := DepsTag("rust", pinned, "FROM x\n", lock)
	if !regexp.MustCompile(`^prog-deps-rust:[0-9a-f]{24}$`).MatchString(tag) {
		t.Fatalf("the tag %q has the wrong form", tag)
	}
	if again := DepsTag("rust", pinned, "FROM x\n", map[string][]byte{"Cargo.toml": lock["Cargo.toml"], "Cargo.lock": lock["Cargo.lock"]}); again != tag {
		t.Errorf("the same inputs give %s and %s", tag, again)
	}
	other := map[string][]byte{"Cargo.lock": []byte("version = 3\n"), "Cargo.toml": lock["Cargo.toml"]}
	for name, got := range map[string]string{
		"another lock file":   DepsTag("rust", pinned, "FROM x\n", other),
		"another image":       DepsTag("rust", strings.Replace(pinned, "930e", "930f", 1), "FROM x\n", lock),
		"another Dockerfile":  DepsTag("rust", pinned, "FROM y\n", lock),
		"another language":    DepsTag("java", pinned, "FROM x\n", lock),
		"a file name changes": DepsTag("rust", pinned, "FROM x\n", map[string][]byte{"Cargo.lock2": lock["Cargo.lock"], "Cargo.toml": lock["Cargo.toml"]}),
	} {
		_, h, _ := strings.Cut(got, ":")
		if h == tag[len("prog-deps-rust:"):] {
			t.Errorf("%s gives the same hash %s", name, got)
		}
	}
}

func TestDepsDockerfile(t *testing.T) {
	const base = "prog-base-x:1"
	for lang, want := range map[string][]string{
		"rust": {"FROM " + base + "\n", "COPY Cargo.toml Cargo.lock /opt/prog/deps/", "CARGO_NET_OFFLINE=false cargo vendor --locked --versioned-dirs /opt/prog/vendor > .cargo/config.toml", "cargo build --offline --locked", "chmod -R a+rX /opt/prog", "USER 1000:1000"},
		"java": {"FROM " + base + " AS fetch\n", "ADD --checksum=sha256:" + mavenSHA256 + " https://repo1.maven.org/", "apache-maven-" + mavenVersion + "-bin.tar.gz", "COPY pom.xml /opt/prog/deps/pom.xml", mavenCopyDependencies + " -DincludeScope=test " + mavenExclude + " -DoutputDirectory=/opt/prog/lib", "\nFROM " + base + "\nCOPY --from=fetch /opt/prog/lib /opt/prog/lib"},
		"cpp":  {"FROM " + base + "\n", "COPY vcpkg*.json /opt/prog/deps/", "command -v vcpkg", `echo "prog: the C++ image has no vcpkg`, "vcpkg install --x-manifest-root=/opt/prog/deps --x-install-root=/opt/prog/vcpkg", "ln -s \"$triplet/include\" /opt/prog/include", "USER 1000:1000"},
	} {
		df := DepsDockerfile(DepsSupport[lang], base)
		for _, w := range want {
			if !strings.Contains(df, w) {
				t.Errorf("the %s Dockerfile has no %q:\n%s", lang, w, df)
			}
		}
		if strings.Contains(df, "%!") {
			t.Errorf("the %s Dockerfile has a bad format verb:\n%s", lang, df)
		}
	}
	// The rust image sets CARGO_NET_OFFLINE=true. Only the step that fetches may turn it off, and the image must not keep it off.
	if df := DepsDockerfile(DepsSupport["rust"], base); strings.Count(df, "CARGO_NET_OFFLINE") != 1 || strings.Contains(df, "ENV CARGO") {
		t.Errorf("the rust Dockerfile changes CARGO_NET_OFFLINE for the image:\n%s", df)
	}
}

func TestDepsFor(t *testing.T) {
	for _, tc := range []struct {
		lang    string
		sandbox bool
		test    bool
		skip    bool
	}{
		{"python", true, false, true},
		{"python", false, false, true},
		{"python", true, true, true},
		{"kotlin", true, false, true},
		{"cpp", false, false, true},
		{"cpp", true, false, false},
		{"cpp", true, true, false},
		{"cpp", false, true, true},
		{"rust", false, false, false},
		{"rust", true, false, false},
		{"rust", false, true, false},
		{"java", false, false, false},
		{"java", false, true, false},
		{"csharp", false, true, false},
		{"csharp", true, true, false},
		// C# runs only the test blocks with deps=.
		{"csharp", true, false, true},
	} {
		_, err := DepsFor(tc.lang, tc.sandbox, tc.test)
		var ds DepsSkip
		if got := errors.As(err, &ds); got != tc.skip {
			t.Errorf("DepsFor(%s, %v, %v) = %v", tc.lang, tc.sandbox, tc.test, err)
		}
	}
}

func TestWithDepsSkips(t *testing.T) {
	r := &Runner{Cache: t.TempDir()}
	if _, err := r.WithDeps(Programs["python"], "deps/python/requirements.lock"); !errors.As(err, new(DepsSkip)) {
		t.Errorf("python with deps= must skip: %v", err)
	}
	if _, err := r.WithDeps(Programs["cpp"], "deps/cpp/vcpkg.json"); !errors.As(err, new(DepsSkip)) {
		t.Errorf("cpp with deps= on this machine must skip: %v", err)
	}
	if _, err := r.WithDeps(Programs["rust"], "deps/rust/requirements.lock"); err == nil || errors.As(err, new(DepsSkip)) || !strings.Contains(err.Error(), "Cargo.lock") {
		t.Errorf("a lock file of another tool must fail: %v", err)
	}
	if _, err := r.WithDeps(Programs["rust"], filepath.Join(t.TempDir(), "Cargo.lock")); err == nil || errors.As(err, new(DepsSkip)) {
		t.Errorf("a missing lock file must fail: %v", err)
	}
}

func TestDepsPrograms(t *testing.T) {
	toml := map[string][]byte{"Cargo.toml": []byte("[package]\nname = \"json-demo\"\nversion = \"0.1.0\"\n\n[dependencies]\nserde = { version = \"1\" }\n")}
	p, err := rustDeps(Programs["rust"], toml, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Main != "src/main.rs" || p.Seed != "" || len(p.Steps) != 2 || p.Steps[0][0] != "sh" || !strings.Contains(p.Steps[0][2], "cp -pR /opt/prog/deps/. . && exec cargo build --offline --locked --quiet") || p.Steps[1][0] != "./target/debug/json-demo" {
		t.Errorf("rust in the image: %+v", p)
	}
	p, err = rustDeps(Programs["rust"], toml, "/cache/rust-1")
	if err != nil {
		t.Fatal(err)
	}
	if p.Seed != "/cache/rust-1" || strings.Join(p.Steps[0], " ") != "cargo build --offline --locked --quiet" {
		t.Errorf("rust on this machine: %+v", p)
	}
	if Programs["rust"].Main != "main.rs" || Programs["rust"].Steps[0][0] != "rustc" {
		t.Error("rustDeps changed the program of the language")
	}

	p, _ = javaDeps(Programs["java"], nil, "")
	if got := strings.Join(p.Steps[0], " ") + " | " + strings.Join(p.Steps[1], " "); got != "javac -cp /opt/prog/lib/* -d classes Main.java | java -cp classes:/opt/prog/lib/* Main" {
		t.Errorf("java in the image: %s", got)
	}
	p, _ = javaDeps(Programs["java"], nil, "/cache/java-1")
	if lib := filepath.Join("/cache/java-1", "lib", "*"); p.Steps[0][2] != lib || p.Steps[1][2] != "classes"+string(os.PathListSeparator)+lib {
		t.Errorf("java on this machine: %v", p.Steps)
	}
	if Programs["java"].Steps[1][2] != "classes" {
		t.Error("javaDeps changed the program of the language")
	}

	p, _ = cppDeps(Programs["cpp"], nil, "")
	if got := strings.Join(p.Steps[0], " "); got != "c++ -I/opt/prog/include -L/opt/prog/lib -std=c++23 -Wall -o main main.cpp" {
		t.Errorf("cpp: %s", got)
	}
}

func TestCargoBin(t *testing.T) {
	for toml, want := range map[string]string{
		"[package]\nname = \"example\"\nversion = \"0.1.0\"\n\n[dependencies]\nname = \"not this\"\n": "example",
		"[package]\nname = \"pkg\"\n\n[[bin]]\nname = \"tool\"\npath = \"src/main.rs\"\n":             "tool",
		"[dependencies.serde]\nname = \"x\"\n[package]\n  name='late'\n":                              "late",
	} {
		if got, err := cargoBin([]byte(toml)); err != nil || got != want {
			t.Errorf("cargoBin(%q) = %q, %v; want %q", toml, got, err, want)
		}
	}
	if _, err := cargoBin([]byte("[dependencies]\nname = \"x\"\n")); err == nil {
		t.Error("a Cargo.toml with no package must fail")
	}
}

func TestReadDeps(t *testing.T) {
	dir := t.TempDir()
	lock := filepath.Join(dir, "Cargo.lock")
	if err := os.WriteFile(lock, []byte("version = 4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readDeps(lock, DepsSupport["rust"]); err == nil {
		t.Error("a Cargo.lock with no Cargo.toml must fail")
	}
	vc := filepath.Join(dir, "vcpkg.json")
	if err := os.WriteFile(vc, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	files, err := readDeps(vc, DepsSupport["cpp"])
	if err != nil || len(files) != 1 || string(files["vcpkg.json"]) != "{}\n" {
		t.Errorf("readDeps = %v, %v", files, err)
	}
}

func TestBuildWhy(t *testing.T) {
	out := "#5 [2/3] RUN if ! command -v vcpkg; then echo \"prog: the C++ image has no vcpkg, so ...\" >&2; exit 1; fi\n" +
		"#5 0.215 prog: the C++ image has no vcpkg, so it cannot install the packages of vcpkg.json\n#5 ERROR: process did not complete successfully: exit code: 1\n"
	if got := buildWhy(out, errors.New("exit status 1")); got != "the C++ image has no vcpkg, so it cannot install the packages of vcpkg.json" {
		t.Errorf("docker: %q", got)
	}
	if got := buildWhy("STEP 3/4: RUN ...\nprog: no vcpkg\nError: building at STEP\n", errors.New("exit status 1")); got != "no vcpkg" {
		t.Errorf("podman: %q", got)
	}
	if got := buildWhy("error: failed to download serde\n", errors.New("exit status 101")); !strings.Contains(got, "exit status 101") || !strings.Contains(got, "failed to download serde") {
		t.Errorf("no message: %q", got)
	}
}

func TestCopyTree(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "target", "debug"), 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(src, "target", "debug", "lib.rlib")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := os.Chtimes(f, old, old); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, depsReady), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(dst, "target", "debug", "lib.rlib"))
	if err != nil || !info.ModTime().Equal(old) {
		t.Errorf("the copy does not keep the time: %v, %v", info, err)
	}
	if _, err := os.Stat(filepath.Join(dst, depsReady)); err == nil {
		t.Error("the copy has the ready mark")
	}
}

// TestRustDepsLocal runs a block with the crate of testdata/deps/rust on this machine. It needs cargo and the crate in the cargo cache. With an empty cargo cache, the block must be a skip.
func TestRustDepsLocal(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("no cargo")
	}
	if (&Runner{}).ProgramReady(Programs["rust"]) != nil {
		t.Skip("no rustc")
	}
	lock, err := filepath.Abs(filepath.Join("testdata", "deps", "rust", "Cargo.lock"))
	if err != nil {
		t.Fatal(err)
	}
	main := "fn main() {\n    let mut b = itoa::Buffer::new();\n    println!(\"{}\", b.format(-42i64));\n}\n"

	t.Run("empty cache", func(t *testing.T) {
		t.Setenv("CARGO_HOME", t.TempDir())
		r := &Runner{Cache: t.TempDir(), TempDir: t.TempDir()}
		_, err := r.WithDeps(Programs["rust"], lock)
		if !errors.As(err, new(DepsSkip)) {
			t.Fatalf("an empty cargo cache must give a skip: %v", err)
		}
	})

	r := &Runner{Cache: t.TempDir(), TempDir: t.TempDir()}
	p, err := r.WithDeps(Programs["rust"], lock)
	if errors.As(err, new(DepsSkip)) {
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		got, err := r.RunProgram(p, map[string]string{p.Main: main}, "", nil, time.Minute)
		if err != nil || got.Exit != 0 || got.Stdout != "-42\n" {
			t.Fatalf("the block gives %+v, %v", got, err)
		}
	}
	// A second WithDeps finds the directory of the packages.
	q, err := r.WithDeps(Programs["rust"], lock)
	if err != nil || q.Seed != p.Seed {
		t.Errorf("the second WithDeps gives %q, %v; want %q", q.Seed, err, p.Seed)
	}
}

// TestDepsImage builds the images with the packages of testdata/deps and runs a block in each one with no network. Set PROG_TEST_RUST_IMAGE or PROG_TEST_JVM_IMAGE to the runner image to run it. The image build needs the network.
func TestDepsImage(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("no docker")
	}
	for _, tc := range []struct {
		env, lang, lock, main, want string
	}{
		{"PROG_TEST_RUST_IMAGE", "rust", "rust/Cargo.lock", "fn main() {\n    let mut b = itoa::Buffer::new();\n    println!(\"{}\", b.format(-42i64));\n}\n", "-42\n"},
		{"PROG_TEST_JVM_IMAGE", "java", "java/pom.xml", "public class Main {\n    public static void main(String[] args) {\n        System.out.println(new com.google.gson.Gson().toJson(java.util.List.of(1, 2)));\n    }\n}\n", "[1,2]\n"},
	} {
		t.Run(tc.lang, func(t *testing.T) {
			ref := os.Getenv(tc.env)
			if ref == "" {
				t.Skip(tc.env + " is not set")
			}
			box := NewSandbox("docker", map[string]string{tc.lang: ref})
			r := &Runner{Box: box, TempDir: t.TempDir()}
			if err := r.ProgramReady(Programs[tc.lang]); err != nil {
				t.Fatal(err)
			}
			lock, err := filepath.Abs(filepath.Join("testdata", "deps", filepath.FromSlash(tc.lock)))
			if err != nil {
				t.Fatal(err)
			}
			p, err := r.WithDeps(Programs[tc.lang], lock)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(p.Image, "prog-deps-"+tc.lang+":") {
				t.Fatalf("the program runs in %q", p.Image)
			}
			got, err := r.RunProgram(p, map[string]string{p.Main: tc.main}, "", nil, 2*time.Minute)
			if err != nil || got.Exit != 0 || got.Stdout != tc.want {
				t.Fatalf("the block gives %+v, %v", got, err)
			}
			// A new sandbox finds the image by its tag and does not build it again.
			start := time.Now()
			r2 := &Runner{Box: NewSandbox("docker", map[string]string{tc.lang: ref}), TempDir: t.TempDir()}
			q, err := r2.WithDeps(Programs[tc.lang], lock)
			if err != nil || q.Image != p.Image {
				t.Fatalf("the second run gives %q, %v; want %q", q.Image, err, p.Image)
			}
			if d := time.Since(start); d > 30*time.Second {
				t.Errorf("the second run took %s; it must reuse the image", d)
			}
		})
	}
}
