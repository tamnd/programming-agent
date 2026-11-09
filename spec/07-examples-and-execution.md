# 07. Examples and execution

Every piece of code on the site is executed by the build, or is marked with the reason why it is not. This is the main difference between this site and the reference sites that disclaim correctness (Rosetta Code shows output pasted by hand; Programming-Idioms says it gives no guarantee, 01.3). This document defines the example markup, the runner, the images, the version matrix, the output rules, and the reader's playground.

## 7.1 Decisions

| Decision | Choice | Why |
|---|---|---|
| One harness or each language's doctest? | **One harness** (`prog run`) for all languages | 25 doctest dialects (rustdoc, Python doctest, Go `Example`, ExUnit doctest, sol/doctest, mdx, ...) have 25 sets of rules; one harness gives one set of attributes, one output normalizer and one cache |
| Where do examples live? | In Markdown, in fenced blocks with attributes, parsed with a CommonMark AST (goldmark, as in the `markup` package) | the page source is the test source; no copy to drift |
| Where do they run in CI? | In OCI images pinned by digest, one per language version, under gVisor (`runsc`) with no network and a read-only root | the algorithm sandbox, extended to many versions |
| Is output typed by hand? | Never. `prog run -update` writes the `output` blocks, as `dune promote` and `prysk -i` do | output typed by hand drifts |
| Can readers run code? | Yes, in the browser, where a WebAssembly runtime exists (7.10). No server execution for readers in v1 (Default; 22, Q-07) | static site, no cost per reader, no abuse surface |
| Pinning | image digests plus `mise` or `uv` inside the image; Nix optional | Nix flakes are still experimental in upstream Nix in Oct 2026 (01.4) |

## 7.2 Block attributes

A fenced code block's info string is the language name and then attributes. The language name is the `id` of `registers/languages.yaml` (`python`, `rust`, `cpp`, `sql`, `bash`, ...) or an alias of it.

| Attribute | Meaning | Default |
|---|---|---|
| `run` | run the block; it must exit 0 and its output must match the next `output` block | none |
| `fails` | run the block; it must fail as declared by `expect=` | none |
| `expect=<sig>` | the failure signature: `compile`, `runtime`, `exit:<n>`, `exception:<Type>`, `error:<code>` (a diagnostic code such as `E0382` or `TS2339`), `panic`, `timeout`, `sanitizer:<kind>` | none |
| `wrong` | run the block; it exits 0, but the output is the wrong answer that the page explains | none |
| `norun reason="..."` | do not run; the reason is required and must be one of the codes of 7.8 | none |
| `id=<name>` | a stable id, unique on the page; it is the anchor and the key of the recorded output | the build assigns `ex1`, `ex2`, ... and the linter asks for an id at W3 |
| `since=<v>` | the example needs version v or later; the K check (09.4) runs it on the version before and expects failure | the language's support window floor |
| `until=<v>` | the example works up to version v; it fails after (for removed features) | none |
| `version=<v>` | run only on version v (for side-by-side behaviour changes) | none |
| `version=matrix` | run on every version of the support window (the X3 check) | X3 runs the matrix for every `run` block anyway; this attribute marks a block that is about the matrix |
| `impl=<id>` | run on this implementation only (`pypy`, `graalpy`, `clang`, `msvc`, `bun`, ...) | the language's default implementation |
| `std=<x>` | the language standard flag (`c17`, `c23`, `c++20`, `c++23`, `es2026`) | the newest standard in the image |
| `flags="..."` | extra compiler or interpreter flags from the allow list of the language | none |
| `file=<name>` | the block is a file of a multi-file example, with its path | `main.<ext>` |
| `group=<name>` | the blocks with the same group are one program (multi-file, or a program shown in parts) | none |
| `prelude=<id>` | prepend a hidden block (shown in a `<details>`) | none |
| `stdin=<id>` | feed the named `input` block to stdin | empty |
| `args="..."` | command-line arguments | none |
| `deps=<lock>` | the example uses packages pinned in a lock file of the node (7.7) | none |
| `timeout=<s>` | wall-time limit in seconds, up to 30 | 10 |
| `browser` | the reader may run it in the playground (7.10); the build checks that the browser runtime gives the same output | on when the language has a browser runtime and the block has no `deps`, `file` or OS use |
| `lint="rule,rule"` | lint rules turned off for this block, each with a reason in the node's notes | none |
| `platform=<os>` | `linux` (default), `windows`, `macos`; non-Linux blocks run in the scheduled jobs of 7.5 | `linux` |

Output blocks:

````markdown
```python run id=split-basic
print("a,b,,c".split(","))
```

```output for=split-basic
['a', 'b', '', 'c']
```
````

| Output attribute | Meaning |
|---|---|
| `for=<id>` | the block whose output this is; required |
| `stream=stdout\|stderr\|both` | which stream; `both` interleaves as the terminal showed it | `stdout` |
| `match=exact\|regex\|unordered\|any` | exact after normalization (7.6); each line a regex; the same lines in any order (sets, maps, threads); any output (only the exit status counts) | `exact` |
| `normalize="addr,time,tmp,hash,float:6"` | extra normalizers: memory addresses, timestamps, temp paths, hash ids, float precision | the defaults of 7.6 |
| `version=<v>` | the output on version v, when it differs between versions | none |

Rules:

- **X01**: every fenced code block has a language and exactly one of `run`, `fails`, `wrong`, `norun`. A block of a language that is not in the register (for example a pseudo-language or a grammar) must be `text`, `ebnf`, `abnf`, `peg` or another notation id of `registers/notations.yaml`.
- **X02**: every `run` and `wrong` block has an `output` block, unless `match=any`.
- **X03**: every `output` block was written by `prog run -update`. The build stores the hash of each recorded output in `build/runs/` and fails if the text in the page differs from the last recorded run.
- **X04**: a `fails` block fails with exactly the declared signature. A compile error with another code, or a runtime error of another type, fails the check.
- **X05**: an `expect=error:<code>` points to a diagnostic card that exists (04.2), so that every diagnostic shown on the site links to its card.

## 7.3 Programs, fragments and preludes

Most examples are complete programs. Some languages need a lot of ceremony for a complete program (Java's class and `main`, C's includes). The harness supports **wrappers** so the page can show the important lines only:

- Each language has wrappers in `languages/<lang>/wrappers/`: `program` (none), `main` (put the block inside `main`), `expr` (print the value of the expression, as a REPL would), `repl` (run as a REPL session: lines that start with the language's prompt, `>>>`, `>`, `iex>`, `ghci>`, `scala>`, `jshell>`, are input and the rest is output). A block selects a wrapper with `wrap=<name>`.
- `wrap=repl` blocks have their output inside the block, as a REPL session shows it. The harness checks the output lines against the run. This is the doctest model, made language-neutral.
- A hidden prelude (`prelude=<id>`) is shown under a `<details>` element named "Setup code". The copy button copies the prelude and the block together, so that the reader always gets a complete program.
- **X06**: the copied text of every block (prelude, wrapper, block) is a program that runs. The build tests it in this form.

## 7.4 The runner and the images

`prog run` extracts the blocks with the CommonMark AST, groups them, wraps them, and runs each one in its image.

**Images.** One OCI image per language version (and per implementation where the page uses it), named `ghcr.io/tamnd/prog-runner/<lang>:<version>`, used only by digest. `registers/images.yaml` lists the digests. The image contains the toolchain, the language's formatter and linter at pinned versions, the test framework of 10.3, and nothing else.

| Source of the base | When |
|---|---|
| the Docker official image (python, pypy, golang, rust, node, ruby, php, gcc, haskell, elixir, erlang, swift, julia, perl, r-base, dart, clojure, groovy, nim, eclipse-temurin, ...) pinned by digest | when one exists; 141 official images in Nov 2026 |
| the vendor's image (`mcr.microsoft.com/dotnet/sdk`, `ocaml/opam`) | when the vendor publishes one |
| our own Dockerfile with `mise` (Zig, Kotlin, Scala, Lua, Deno, Bun, ...) or `uv python install` | otherwise |

**The version matrix** of tier 1 at M0 (Draft; 09.3 sets the support window rule):

| Language | Versions in the matrix | Implementations |
|---|---|---|
| Python | 3.9 (for the `since` checks of 3.10 features) to 3.15, and 3.16 alphas in a scheduled job | CPython; PyPy 3.11 for the pages that name it |
| JavaScript | Node.js 18 to 26; the ECMAScript edition gates use Node's `--harmony` flags only when the page says so | V8 (Node), Bun, Deno for pages that name them; browsers are covered by BCD data, not by our runs (09.6) |
| TypeScript | 4.9, each minor of 5.x, 6.0 (the last compiler written in TypeScript), 7.0 and later (the native Go port, GA 2026-07-08), on the newest Node | tsc (both lines); the type checking result is the output of `fails expect=error:TS<n>` blocks |
| Java | 8, 11, 17, 21, 25 (LTS) and the newest feature release (JDK 27 in Nov 2026) | Temurin |
| C | GCC 9 to 15 and Clang 12 to 21, each with every `-std` it supports | GCC, Clang; MSVC in the scheduled Windows job |
| C++ | as C, with `-std=c++11` to `c++23`, and `c++26` (`c++2c`) as a draft standard: WG21 finished the text in March 2026, but the ISO DIS ballot failed in Oct 2026 (01.4) | GCC, Clang, MSVC (scheduled) |
| C# | .NET 6, 8, 10 (LTS) and the newest STS (.NET 11, GA planned 2026-11-10) | .NET SDK |
| Go | the newest 8 minor versions (1.20 to 1.27) | gc; TinyGo for pages that name it |
| Rust | 1.56 (edition 2021), each sixth stable release after it, and the newest stable (1.99 in Nov 2026); editions 2015, 2018, 2021, 2024 | rustc |
| Kotlin | 1.9, 2.0 to the newest, on the JVM | kotlinc |
| SQL | SQLite 3.x newest, PostgreSQL 13 to 18, MySQL 8.0, 8.4 and 9.x, MariaDB LTS, DuckDB newest | each dialect is an implementation column of the SQL compat (09.2) |
| Bash | 3.2 (the macOS system bash), 4.4, 5.0 to 5.3; dash and zsh for the portability pages; BusyBox ash | none |
| Regex | engines as implementations: PCRE2, Python `re`, JS (V8), Java, .NET, Go RE2, Rust `regex`, POSIX ERE (GNU grep) | none |
| HTML/CSS | rendered in headless Chromium and Firefox for pages that show rendering; the compat facts come from BCD | none |
| PowerShell | 5.1 (scheduled Windows job), 7.4 LTS and newest | none |

Tier 2 has the newest version and the oldest supported one at M3, and the full matrix at M4.

**Sandbox.** As in algorithm (`algo test`): no network, a read-only root with a tmpfs work directory, a user with no privileges, cgroup limits (CPU 1 core, memory 512 MiB by default, PIDs 64), and the wall timeout. CI runs the images under gVisor (`runsc`). Local runs use Podman rootless or Docker, and `-unsafe-local` runs on the host for development only, as in algorithm.

**Cache.** A result is keyed by BLAKE3(image digest, wrapper, flags, stdin, the block group's text). A run whose key is in the cache is skipped. A change of an image digest re-runs every block that uses it, which is how a new compiler release finds regressions in our pages.

**Sharding.** CI shards by image across the free `ubuntu-24.04` and `ubuntu-24.04-arm` runners, as algorithm does by language. Examples with `platform=windows` or `platform=macos` run in a scheduled weekly job on `windows-latest` (Server 2025) and `macos-latest` (macOS 26 arm64), not on each PR (Default).

## 7.5 What examples may not do

- **No network.** An example that shows an HTTP client runs against a local server that the harness starts in the same sandbox (`fixture=http:<name>`, from `fixtures/http/`), or it is `norun reason=network` with recorded output from a scheduled job.
- **No clock or random dependence in exact output.** Use a fixed seed, or `match=regex`, or the `time` normalizer. Concurrency examples whose order varies use `match=unordered`.
- **No files outside the work directory.** Fixtures are copied into it.
- **No GUI.** GUI examples are `norun reason=gui` with a screenshot taken in a scheduled job with a virtual display (12.7).
- **No undefined behaviour shown as a working program.** A C or C++ example that has UB is `fails expect=sanitizer:<kind>` under UBSan or ASan, or `norun reason=ub` with the explanation. A page never shows "the output" of a UB program as if it were defined.

## 7.6 Output normalization and the X ladder

Default normalizers, applied to both the recorded and the actual output before the compare: CRLF to LF, trailing whitespace on each line removed, a final newline added, ANSI colour codes removed, the work-directory path replaced by `$WORK`, the image's home directory by `$HOME`. Per-language normalizers in `languages/<lang>/normalize.yaml`, for example Python's `<object at 0x...>` addresses, Java's `@1b6d3586` identity hashes, Go's goroutine ids in panics, Rust's `thread 'main' panicked at src/main.rs:LL:CC`.

**X3** needs each `run` block to pass on every version of the language's support window (09.3) that is not before its `since` or after its `until`. When the output differs between versions, the page must show both outputs with `version=` (6.7). A difference that the page does not show fails X3: this is how the build finds behaviour changes that a page did not mention (for example dict ordering, integer division, a changed error message).

Error messages change between versions. A `fails` block checks the **signature** (the code or the exception type), not the message text. The output block of a `fails` block shows the message of the newest version, with `match=regex` on the stable part.

## 7.7 Third-party packages

- An example on a page that is not about a library uses the standard library only (X07).
- A page about a library (a `library` node, or a task that names a library as the idiomatic way) may use it with `deps=<lock>`. The node holds the lock file (`requirements.lock`, `Cargo.lock`, `package-lock.json`, `go.sum`, `pom.xml` with versions, ...). The image fetches packages once at build time from a mirror into a cache layer; the run itself still has no network.
- The library's version is part of the example's facts and shows on the page.

## 7.8 Reasons for `norun`

| Code | Meaning | What the build does instead |
|---|---|---|
| `network` | needs a real network | scheduled job with network, recorded output, marked "recorded on <date>" |
| `gui` | needs a display | scheduled job with Xvfb, screenshot |
| `platform` | needs Windows or macOS | scheduled job on that runner |
| `hardware` | needs a GPU, a microcontroller, special hardware | none; the page says so; the example must still compile (`compile-only`) where a toolchain exists |
| `ub` | undefined behaviour shown on purpose | compile only, with the sanitizer check where possible |
| `vulnerable` | a weakness example (06.3, weakness) | runs only inside the security harness with the exploit input (10.8) |
| `fragment` | a part of a larger program, for reading only | the group that holds the full program must run |
| `pseudo` | not code of a real language (should be rare) | none; the linter asks to change it to `text` |
| `long` | needs more than 30 seconds | scheduled job |
| `interactive` | needs a person at a terminal | a scripted session with `expect`, when possible |

`compile-only` is a mode of its own: the block is compiled (or type-checked, or parsed) but not run. It is the default for `norun` with `hardware`, `ub`, `fragment` when the language has a compiler.

Rule X08: `norun` blocks are at most 10% of the blocks of a node, except for nodes in `embedded-and-hardware`, `graphics-and-games-programming`, `mobile-and-desktop` and `distributed-and-cloud`, where the limit is 40%.

## 7.9 Benchmarks

A benchmark claim on a page ("X is about 3 times faster than Y for this input") needs a benchmark block: `bench` with the harness of the language (pytest-benchmark, criterion, JMH, Go `testing.B`, BenchmarkDotNet, google-benchmark), run in a scheduled job on a fixed runner class, with the result recorded with the hardware, the versions and the date. The page shows the result as a range, with "recorded on", and never as a fact about the languages in general. Rule X10.

## 7.10 The reader's playground

The reader can edit and run an example in the browser where a WebAssembly runtime runs it. The runtime loads in a Web Worker on the first click of "Run", never on page load. The site serves the runtimes itself, from R2, under cross-origin isolation headers (COOP and COEP) on the playground routes only (16.8).

| Language | Browser runtime (Nov 2026) | Size, gzip (measured, 01.4) | Default |
|---|---|---|---|
| Python | Pyodide 314.0.x (CPython 3.14) | about 6 MB with the stdlib zip | on |
| Python (quick) | MicroPython 1.29 | about 200 KB | on for `intro` pages, with a "full Python" switch |
| JavaScript | the browser itself, in a sandboxed iframe or a Worker; QuickJS-wasm for isolation | none | on |
| TypeScript | the `typescript` compiler in a Worker, then JS | about 3 MB | on |
| SQL | sqlite-wasm 3.54 (about 400 KB); PGlite 0.5 for PostgreSQL pages (about 5 MB) | | on |
| Lua | wasmoon (Lua 5.4) | about 110 KB | on |
| Ruby | ruby.wasm, minimal build | about 5 MB | on |
| PHP | php-wasm | not measured | on |
| OCaml | wasm_of_ocaml toplevel 6.5 | about 9 MB | on |
| Regex | each engine as wasm or JS: PCRE2 wasm, the JS engine, RE2 wasm, Rust `regex` wasm | small | on |
| Bash | a POSIX shell in wasm with BusyBox (Default: off until tested) | | off |
| C, C++ | Clang in wasm (Runno, xeus-cpp) | large | off (experimental flag) |
| Java | TeaVM with javac in the browser (subset of the JDK) | not measured | off (experimental flag) |
| Go, Rust, Kotlin, Swift, C#, Haskell, Elixir, Zig, others | no maintained in-browser compiler in Nov 2026 | none | "Open in the official playground" with the code passed in the URL where the playground supports it (go.dev/play, play.rust-lang.org, play.kotlinlang.org, Compiler Explorer, SwiftFiddle) |

Rules:

- **X11**: a block marked `browser` gives the same normalized output in its browser runtime as in CI. The build runs the browser runtimes headless (Node or Chromium) on every block marked `browser`. A mismatch removes the `browser` mark and opens a gap item.
- The playground shows the recorded CI output next to the reader's output. When the reader has not changed the code and the outputs differ, it shows a notice ("this runtime differs from the CI run") with the version of each.
- The playground never sends code to our servers. "Open in" links send code to the third party only when the reader clicks, with a notice.
- A server-side "Run" for compiled languages (Judge0 CE or Piston self-hosted, or Cloudflare Sandbox SDK with our images) is **not in v1**. 22, Q-07 holds the question. Piston's public API closed in Feb 2026 and TIO is dormant, so no public free service is assumed.

## 7.11 The examples dataset

The build emits `build/data/examples.jsonl`: one row per block with the page, the id, the language, the version range, the implementation, the code (as the reader copies it), the recorded output, the run status per version, and the hashes. It is published under MIT (the code) and CC0 (the metadata) on the site and on Hugging Face (16.10). It is the "language × feature × version" dataset that no one publishes today (01.3).
