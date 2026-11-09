# 09. Compat and versions

Programming knowledge goes out of date. A feature gets added, deprecated and removed, a default changes, an error message changes, a library breaks its API. This document makes versions into data that the build checks, so the site can say "this is true for Python 3.10 to 3.15" and prove it. The model is MDN's browser-compat-data (BCD, CC0) and the Baseline status of web-features, applied to languages, tools and libraries.

## 9.1 The parts

| Part | What | Where |
|---|---|---|
| Release register | every release of every language, tool and library that the site covers, with dates | `nodes/languages/<lang>/releases.yaml` and the release cards (04.2) |
| Compat entries | for each feature, api card and diagnostic: the version where it was added, changed, deprecated and removed, in each implementation | `compat.yaml` in the feature node; the `since`, `deprecated` and `removed` fields of cards |
| Probes | a small program per compat entry that fails before the version and passes from it | `probes/` in the feature node |
| Support window | the range of versions that a page's claims are about | computed per language (9.3) |
| Status | a computed label such as "available in all supported versions" | computed (9.5) |

## 9.2 Compat entries

The format follows BCD's `__compat` object, with implementations instead of browsers.

```yaml
# nodes/languages/python/control-flow/match-statement/compat.yaml
match-statement:
  description: "The match statement"
  spec: {url: "https://docs.python.org/3/reference/compound_stmts.html#the-match-statement"}
  source: {proposal: pep-634, release_notes: "https://docs.python.org/3.10/whatsnew/3.10.html#pep-634-structural-pattern-matching"}
  support:
    cpython: {version_added: "3.10", probe: probes/match_basic.py}
    pypy:    {version_added: "7.3.8", note: "PyPy 3.9+ language level 3.10", probe: probes/match_basic.py}
    micropython: {version_added: false, note: "Not supported in MicroPython 1.29", source: "https://docs.micropython.org/en/latest/genrst/core_language.html"}
  status: {standard: true, experimental: false, deprecated: false}

match-statement.class-patterns:
  description: "Class patterns with __match_args__"
  support:
    cpython: {version_added: "3.10", probe: probes/match_class.py}

match-statement.soft-keyword:
  description: "match and case are soft keywords; they stay valid as names"
  support:
    cpython: {version_added: "3.10", probe: probes/match_soft_keyword.py}
```

Fields of a support statement:

| Field | Meaning |
|---|---|
| `version_added` | the first release where it works without a flag; `false` for "not supported"; `preview` for preview or nightly; never `true` or `null` at K1 (a K1 entry has a real version) |
| `version_removed` | the first release where it no longer works |
| `version_deprecated` | the first release where it gives a deprecation warning or is deprecated in the docs |
| `flags` | the flag, option, edition or `from __future__` import needed before `version_added` (`--enable-preview` in Java, `-std=c++2b`, `#![feature(...)]` in nightly Rust, `from __future__ import annotations`) |
| `partial_implementation` | `true` with a `note` |
| `alternative_name` | the old name, for renamed features |
| `edition` | for Rust and other languages with editions: the edition range where the statement holds |
| `std` | for C and C++: the language standard where the feature is (`c++20`), separate from the compiler version where it is implemented |
| `note` | a short note, plain English |
| `probe` | the probe file (9.4) |
| `source` | a primary source for the statement (6.6) |

Implementations are columns. `registers/implementations.yaml` lists them per language:

| Language | Columns |
|---|---|
| Python | cpython, pypy, graalpy, micropython |
| JavaScript | the ECMAScript edition (the standard as a column), v8 (Node.js), spidermonkey, javascriptcore, bun, deno; browser support is imported from BCD, not probed by us |
| TypeScript | tsc |
| Java | the JLS version (the standard), openjdk |
| C | the ISO C standard (c89, c99, c11, c17, c23, c2y), gcc, clang, msvc |
| C++ | the ISO C++ standard (c++98 to c++23, and c++26 marked `draft` until ISO publishes it), gcc, clang, msvc, and the library implementations libstdc++, libc++, msvc-stl as separate columns for library features |
| C# | the C# language version, dotnet (Roslyn and the runtime) |
| Go | gc, gccgo, tinygo |
| Rust | rustc, with the edition |
| Kotlin | kotlin-jvm, kotlin-native, kotlin-js, kotlin-wasm |
| SQL | iso-sql, postgresql, mysql, mariadb, sqlite, sqlserver, oracle, duckdb |
| Shell | posix, bash, dash, zsh, ksh93, busybox |
| Regex | pcre2, ecmascript, python-re, java, dotnet, re2, rust-regex, posix-ere, oniguruma, icu |

For C and C++ the standard column and the compiler columns differ on purpose: "added in C++20" is a fact about the standard; "GCC 11 supports it with `-std=c++20`" is a fact about the compiler. The source for the compiler columns is each compiler's status page (the GCC "C++ Standards Support" page, the Clang "C++ Support" page, the MSVC conformance page), and the probe checks it.

## 9.3 The support window

Each language has a **support window**: the versions that the site's claims cover. A claim with no `since` or `until` is a claim about every version in the window.

- The **ceiling** is the newest stable release.
- The **floor** is the oldest release that its maintainers still support (from the language's own policy: Python's release-cycle.json, the Java LTS list, Go's "two most recent releases" policy extended by us to 8 minors for the matrix, Node's LTS schedule, .NET's LTS list), and at least the release 5 years before the ceiling for languages without a policy. `registers/support-policy.yaml` holds the rule per language with its source.
- endoflife.date (MIT, 489 products in API v1) is the cross-check for dates. It does not cover Swift, Dart, Zig, C, C++, TypeScript, OCaml, GCC or LLVM (01.4), so those languages use their own feeds.
- The window moves on its own when a release comes out or reaches its end of life. When the floor moves up, pages do not change, but some `since` attributes become redundant and the linter removes them (the "version floor" cleanup, 14.9).

## 9.4 Probes

A probe is the proof of a compat entry. It is a small program that uses the feature in the smallest way, with an expected result.

- The probe must **pass** (exit 0 and the expected output) on `version_added` and on every version after it, up to `version_removed`.
- The probe must **fail** on the release just before `version_added`, with the declared failure signature (a syntax error, a compile error with a code, a missing name). Rule K05: the failure signature is stored, so that an unrelated failure (a missing image, an import error of something else) does not count.
- The release before comes from the release register. Images of old versions are kept (old Docker tags stay pullable after end of life; we mirror the digests we use to GHCR so they stay available).
- One probe may prove several entries.
- A K2 entry needs a probe that passed both checks in the current images. Entries for versions older than the image matrix (for example "added in Python 2.2") are K1 with a source, not K2. The page shows "verified" or "from the documentation" for each.

The K ladder (04.9) counts this per feature.

## 9.5 Computed status

Like Baseline for the web, the build computes a status for each feature from its compat entries and the support window. The definition is versioned (`status-v1`) and published at `/about/status/`.

| Status | Rule (status-v1) | Badge text |
|---|---|---|
| `everywhere` | supported in every version of the support window, in every main implementation of the language | "Available in all supported versions" |
| `newly` | supported in the ceiling, but not in the floor | "Available since <v> (<date>)" |
| `preview` | only behind a flag or in a preview | "Preview in <v>" |
| `deprecated` | `version_deprecated` is in the window | "Deprecated in <v>" |
| `removed` | `version_removed` is at or below the ceiling | "Removed in <v>" |
| `partial` | some main implementations do not support it | "Not in <impl>" |
| `standard-only` | in the language standard but in no compiler at the ceiling (C and C++) | "In <std>, not yet in compilers" |

For HTML, CSS and JavaScript web APIs the site shows the web-features Baseline status itself (Apache-2.0 data) and the BCD tables, with credit, and does not compute its own.

## 9.6 Imported compat data

| Data | Licence | Use |
|---|---|---|
| MDN browser-compat-data (`@mdn/browser-compat-data`) | CC0 | the browser columns for JavaScript, HTML, CSS and Web API features; imported as-is with the BCD version recorded |
| web-features (W3C WebDX CG) | Apache-2.0 | the Baseline status and the feature grouping for web features |
| caniuse | CC BY 4.0 | a cross-check only, with credit |
| Python `versionadded`, `versionchanged`, `deprecated`, `deprecated-removed` directives in the CPython docs | PSF licence for the docs (19.3) | harvested as compat facts (facts are not copyrightable; we keep no text) |
| `RELEASES.md` of rust-lang/rust, releases.rs | MIT/Apache-2.0 | Rust stabilization versions |
| the JEP index (openjdk.org/jeps/0) | facts | the release field of each JEP |
| the Go release notes, the Kotlin "What's new" pages, the C# "What's new" pages, TypeScript release notes | facts | added versions |
| the GCC, Clang and MSVC conformance tables | facts | the compiler columns of C and C++ |
| cppreference's compiler support pages | CC BY-SA 3.0 and GFDL | a cross-check only; facts, never text |
| node.green | MIT | a cross-check for the Node.js column |

The harvester (13.2) reads each source on a schedule and writes K1 entries. A disagreement between two sources opens a gap item for the compat-prober agent (14.2), which writes the probe and settles it.

## 9.7 Releases as content

- Each release of a tier-1 or tier-2 language in the support window is a node (04.2). Its "What is new" section is generated from the compat entries with `version_added` equal to the release, and the writer adds one sentence per item. A feature in the official release notes with no compat entry is a gap (rule K06), so the release node is how the site finds the features it has not covered.
- **Freshness target**: a new release of a tier-1 language has its release node at W2 and its new features as cards with K1 within 7 days of the release, and at W3 with K2 within 30 days (Default). Pre-releases (betas, release candidates) get a release card marked `preview` and K1 entries with `version_added: preview`, so the site is ready on the release day.
- Each language page shows the release timeline (12.6) and the support window.

## 9.8 Libraries and tools

The same model covers tools and libraries, with lighter rules:

- A tool or library node has a release register (from its GitHub releases, its package registry, or endoflife.date).
- Its features and api cards carry `since` and `removed` from its changelog, at K1. Probes (K2) are required only for the standard libraries of tier-1 languages and for the core tools (Git, the main build tools).
- Breaking changes of major versions are migration nodes (`<library>/migrating-from-<n>-to-<m>`), required at W3 for the core libraries.

## 9.9 Version facts in the API

The site publishes the compat data as a dataset (16.10): `build/data/compat/<lang>.json` in BCD's shape, under CC0, with a version and a changelog, and on Hugging Face. To our knowledge in Nov 2026 no open dataset of language × feature × version exists (PLDB records whether a language has a feature, not when; 01.3), so this dataset is a contribution on its own.
