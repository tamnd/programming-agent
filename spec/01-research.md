# 01. Research: the state of the field, November 2026

This document is the evidence for the design. The research was done on 2026-11-09. Each section ends with what it means for programming. A number marked "Estimate" is to be measured in M0 to M1 (21). Counts are as published by each source on or near the date given; they move every week, and the harvest (13) re-measures them.

## 1.1 What exists, in one paragraph

Programming knowledge is the most written-about subject on the web, and it is also the most scattered. It lives in official references (one per language, each with its own structure), in standards behind paywalls (ISO C and C++), in a few great reference sites (MDN for the web, cppreference for C and C++, the Python and Rust docs), in comparison sites (Hyperpolyglot, Rosetta Code, Programming-Idioms, Learn X in Y minutes), in catalogues (PLDB, Wikidata, Linguist), in books under copyright, and, until 2024, in Stack Overflow. No source has all of: one page per idea, every mainstream language side by side, examples that are executed for the versions they claim, a version history per feature for languages other than the web, and citations to the normative source. MDN with BCD does this for the web platform only. The gap is the joining, the execution and the versions.

## 1.2 Counts: the size of the field

### Language catalogues

| Source | Count (date) | Licence | Notes |
|---|---|---|---|
| PLDB (pldb.io) | 5,167 entries, of which 3,471 are tagged as programming languages; 395 fields (Nov 2026) | public domain | `pldb.io/pldb.json` is 11.4 MB; a CSV is on Hugging Face; 174 features in its full feature list, 65 curated; **no version data per feature** |
| HOPL (hopl.info) | 8,945 languages | no licence, no dump | names and facts by hand only |
| Wikidata | 1,916 items under "programming language" (Q9143) and subclasses, 1,245 with an English Wikipedia article | CC0 | the backbone for `same_as` ids |
| GitHub Linguist 9.7.0 | 845 entries, 569 of type programming | MIT | names, aliases, extensions |
| Esolang wiki | 7,899 languages | CC0 | esoteric languages |
| Rosetta Code | 1,363 tasks, 423 draft tasks, about 1,023 languages (about Sept 2026) | GFDL 1.2 | on Miraheze; answered 403 to every automated request in Nov 2026 |
| 99 Bottles of Beer | site gone (NXDOMAIN) | | |

### Package registries (package counts, Nov 2026)

| Registry | Packages |
|---|---|
| npm | about 4.6 million |
| PyPI | 937,574 |
| Maven Central | 661,801 |
| NuGet | 590,605 |
| Packagist | 479,556 |
| crates.io | 373,351 |
| RubyGems | 198,469 |
| pub.dev | 95,025 |
| Hex | 26,966 |
| CRAN | about 25,700 |
| Hackage | 19,664 |
| opam | 4,726 |

Software Heritage archives 31.1 billion unique source files.

### Taxonomies, curricula and catalogues of concepts

| Source | Scale | Licence | Use |
|---|---|---|---|
| ACM/IEEE-CS/AAAI CS2023 | 17 knowledge areas, 162 knowledge units, about 1,400 topics | all rights reserved | names and locators for the universe |
| SWEBOK v4.0a | 18 knowledge areas, about 730 headings | no redistribution | names and locators |
| ACM CCS 2012 | 2,113 concepts (Software and its engineering 215; Semantics and reasoning 22) | ACM terms | names and the crosswalk (5.7) |
| Van Roy, "Programming paradigms for dummies" | 27 boxes, about 30 paradigms | all rights reserved | the paradigm structure (5.2) |
| TAPL, PFPL 2nd ed., Mitchell's *Concepts in Programming Languages* (14 chapters), PLAI 3.2.5 (CC BY-NC-SA) | | as stated | the concept universe; about 260 concepts as a first candidate list |
| Exercism | 85 tracks, 24 with a syllabus, 882 concepts | MIT | concept names per language |
| CWE 5.0 (2026-10-28) | 956 weaknesses, 396 groupings (the old "categories"), 58 views | MITRE terms | weakness cards |
| CWE Top 25 (2025) | 79, 89, 352, 862, 787, 22, 416, 125, 78, 94, 120, 434, 476, 121, 502, 122, 863, 20, 284, 200, 306, 918, 77, 639, 770 | MITRE terms | the security core |
| CAPEC 3.9 | 559 attack patterns | MITRE terms | `see_also` targets |
| OWASP Top 10:2025 | A01 Broken Access Control to A10 Mishandling of Exceptional Conditions | CC BY-SA 4.0 | the security core |
| OWASP ASVS 5.0; Cheat Sheet Series | about 345 requirements; 137 cheat sheets | CC BY-SA 4.0 | sources |
| Patterns | GoF 23; PoEAA 51; EIP 65; Fowler's refactorings 66 and smells 24; Azure cloud design patterns 44 (CC BY 4.0); microservices.io 53; the c2 wiki 36,852 page names | as stated | the pattern universe |
| Glossaries | FOLDOC 15,284 entries (GFDL); the Jargon File (public domain); ISO/IEC/IEEE 24765 SEVOCAB 4,633 terms | as stated | names |
| Unicode 18.0 (2026-09-16) | 172,808 characters | Unicode License v3 | the text chapter |
| man-pages 6.20; Arch Wiki | about 1,567 pages; 4,463 articles | various; GFDL | sources for Unix and tools |
| tldr-pages | 7,614 English pages | CC BY 4.0 | command names |
| Free-programming-books | about 10,738 entries | CC BY 4.0 | a checklist |
| Stack Overflow | 24.2 million questions, 65,782 tags | CC BY-SA | demand signal (18.6) |
| HOPL I to IV | 59 papers; no HOPL V | ACM | history |

### Popularity (Nov 2026)

| Source | Result |
|---|---|
| TIOBE, Nov 2026 | Python 20.86%, C 11.48%, Java 8.89%, C++ 8.35%, C# 3.39%, JavaScript 2.42%, then R, SQL, Visual Basic, Go (10th), Rust (11th); Dart enters the top 20 |
| PYPL | Python about 52% share of tutorial searches |
| RedMonk, Jan 2026 | JavaScript, Python, Java, PHP, C#, TypeScript, ... |
| GitHub Octoverse 2025 | TypeScript became the most used language on GitHub |
| GitHub Octoverse 2026 | 225 million developers; 40% of pull requests start as agent drafts |
| Stack Overflow survey 2026 | JavaScript 62%, SQL 58.4%, HTML/CSS 58.2%, Python 58% |
| JetBrains State of Developer Ecosystem 2025 | 24,534 respondents |
| IEEE Spectrum | no 2026 edition found |

**What it means.** The language universe is about 12,000 items after merging (18.4), and every one can have a card at no model cost, because PLDB and Wikidata are public domain and CC0. The concept universe is small (a few thousand) and is defined by sources we cannot copy, so their names and locators define the universe and the text is ours (13.5). The long tail (libraries, files) needs the scope rule of 18.5. Tier 1 (8.2) follows the surveys: Python, C, C++, Java, C#, JavaScript, TypeScript, Go, Rust, Kotlin, plus SQL, shell, regex, HTML/CSS and PowerShell.

## 1.3 Reference sites and how they are built

| Site | What we learn | Licence |
|---|---|---|
| MDN | page types per kind (reference, guide, glossary), a strict front matter, the compat table from BCD; 14,739 en-US pages | prose CC BY-SA 2.5; code CC0; BCD CC0 |
| MDN BCD v8.1.8 | 20,933 features across 17 browsers; `version_added`, flags, partial support, notes | CC0 |
| web-features | Baseline: "Newly available" when in all core browsers, "Widely available" 30 months later | Apache-2.0 |
| cppreference | the deepest C and C++ reference, with a "since C++17" mark on every item; read-only from March 2025, moved to MediaWiki on 2026-04-17 with help from isocpp | CC BY-SA 3.0 and GFDL |
| Python docs | `versionadded`, `versionchanged`, `deprecated` directives on every item; Sphinx `objects.inv` | PSF v2; code 0BSD |
| The Rust Reference | rule ids on every normative paragraph, linked to tests; the Ferrocene Language Specification adopted by the Rust project on 2025-03-26 | MIT/Apache-2.0 |
| Go by Example | 85 pages, one idea per page, executed output | CC BY 3.0 |
| Hyperpolyglot; Learn X in Y minutes | side-by-side tables; one page per language | CC BY-SA 3.0 |
| Programming-Idioms | one idiom × many languages; the site says it does not guarantee that code is correct | CC BY-SA |
| Rosetta Code | 1,363 tasks; code not tested by the site | GFDL 1.2 |
| DevDocs | 247 documentation sets in one UI with offline search | MPL-2.0 (the app); docs under their own licences |
| Exercism | concept documents with "about" and "introduction" files; a syllabus graph | MIT |
| Diátaxis | four forms: tutorial, how-to, reference, explanation | CC BY-SA 4.0 |
| SEP; nLab | one entry per concept, long-lived, versioned, cited | as stated |

**What it means.** The best references attach a version to every item (MDN, cppreference, Python) and a rule id to every normative sentence (Rust). Nobody does the version table for every language: there is no open dataset of language × feature × version (9.9). The comparison sites do not run their code. Programming takes MDN's page types and compat model, the Rust Reference's rule ids (6.3), Go by Example's executed output, and Diátaxis for the form of each page.

## 1.4 Execution, runners and versions

### Online execution services

| Service | State (Nov 2026) |
|---|---|
| Judge0 CE | GPL-3.0, 90+ languages, self-hostable |
| Piston | MIT; the public API closed on 2026-02-15; self-hostable |
| Compiler Explorer | 103 languages; compile and run |
| Wandbox | active |
| TIO (Try It Online) | dormant since 2021 |
| Cloudflare Containers | GA on 2026-04-13 |
| Cloudflare Sandbox SDK | 1.0 on 2026-09-30 |
| e2b | Apache-2.0 |
| Daytona | repository archived on 2026-10-03 |

### WebAssembly runtimes for the browser (sizes measured, gzip)

| Runtime | Version | Size |
|---|---|---|
| Pyodide | 314.0.7 | about 6 MB |
| MicroPython (wasm) | 1.29 | 197 KB |
| ruby.wasm | | 5.1 MB |
| wasmoon (Lua) | | 111 KB |
| sqlite-wasm | 3.54.1 | 407 KB |
| PGlite | 0.5.15 | about 5 MB |
| wasm_of_ocaml | 6.5 | about 9 MB |
| CheerpJ 5.0 (Java) | | needs a commercial licence |

Kotlin and Go have no practical in-browser compiler; their playgrounds compile on a server.

### Testing documentation code

Every serious ecosystem has a doc-test tool, and each runs only its own language: rustdoc tests, Python `doctest`, Go `Example` functions, Elixir ExUnit doctests, Haskell `doctest`, `prysk` (cram tests for shells), OCaml `mdx` 2.7, Go `txtar`, `nbval` for notebooks. None runs a page with many languages and many versions. Programming needs its own harness (07).

### Pinning toolchains

| Tool | State |
|---|---|
| Docker official images | 141 |
| Nix flakes | still experimental in upstream Nix (Oct 2026) |
| mise | v2026.11.5 |
| uv | manages Python versions |
| asdf | 0.20.2 |
| GitHub Actions runners | `ubuntu-latest` is 24.04; `macos-latest` is macOS 26 on arm64; arm64 Linux runners free for public repositories since 2025-08-07 |
| endoflife.date | 489 products; does not cover Swift, Dart, Zig, C, C++, TypeScript, OCaml, GCC or LLVM |

### Syntax and grammars

tree-sitter 0.27.2 with about 530 parsers; Shiki 4.6 with 242 languages; Chroma 2.27; highlight.js 11.12 with 193 languages; Linguist 845 entries.

### Language versions on 2026-11-09

| Language | Current | Notes |
|---|---|---|
| Python | 3.15.0 (2026-10-09) | 3.10 is end of life; 3.16 alphas |
| Go | 1.27 | |
| Rust | 1.99 | edition 2024 is the newest; no 2027 edition announced |
| Java | JDK 27 (2026-09-15) | JDK 25 is the LTS |
| Kotlin | 2.4.21 | |
| C# / .NET | .NET 10 LTS; .NET 11 GA planned 2026-11-10 | C# 15 brings unions |
| TypeScript | 7.0 | the native compiler port in Go; GA 2026-07-08 |
| JavaScript | ES2026 (7 features; Temporal not included) | Node 26 is the active LTS |
| C | C23 (ISO/IEC 9899:2024) | C2y working draft N3998 |
| C++ | C++23 | C++26: WG21 finished the text in March 2026, the DIS ballot failed in Oct 2026 (9 of 22 national bodies voted no, over contracts) |
| Swift | 6.4 | |
| Dart | 3.13 | |
| Ruby | 4.0 | |
| PHP | 8.5 | |
| Scala | 3.10 | |
| Haskell (GHC) | 9.14 | |
| OCaml | 5.5 | |
| Elixir | 1.20 | |
| Zig | 0.17 | |
| Julia | 1.13 | |
| Lua | 5.5 | |
| R | 4.6.2 | |

**What it means.** There is no hosted service to rely on: Piston closed, TIO is dormant, Daytona is archived, Judge0 is GPL. We run our own pinned OCI images under gVisor (07.4), with Docker images and `mise`/`uv` inside, not Nix. Readers run code in the browser only, with WebAssembly runtimes (7.10); a server run is Q-07. The release calendar moves every month (two major releases in Oct and Nov 2026 alone), so release handling must be automatic (9.7), and drafts such as C++26 need a `draft` marker (Q-21).

## 1.5 Formal semantics

| Development | State (Nov 2026) |
|---|---|
| K framework | 7.0; the K semantics of C was last updated in Feb 2022 |
| PLT Redex | λπ for Python; a Lua 5.2 model |
| WebAssembly | 3.0 (2025-09-17), specified with SpecTec, which generates the formal rules, the prose and an interpreter |
| CompCert | 3.18 (Aug 2026); non-commercial licence |
| CakeML | BSD; the Brack work at CPP 2026 |
| RustBelt and Iris | active |
| MiniRust; a-mir-formality | incomplete |
| Aeneas | SymCrypt: 237 KLOC of Lean for 16.7 KLOC of Rust; no soundness link to the real Rust semantics |
| Lean | v4.35.0 (2026-10-12) |
| CSLib | Apache-2.0, started June 2025; has STLC safety and strong normalization, System F-sub, CCS; no records module |

**What it means.** Only small, core languages have full mechanized semantics. For real languages the formal work covers fragments, models or translations. The site must not overclaim: three evidence labels (proved, executable, normative) and rule F02 (a theorem about a model is not a theorem about a language) (11.1). The Lean package `Prog` shares algorithm's pin and builds on CSLib (11.2).

## 1.6 What algorithm taught

From `tamnd/algorithm` and `tamnd/programming-agent` (2150), first week, 2026-11-03 to 2026-11-09:

| Fact | Value |
|---|---|
| nodes | 8,396 in 31 areas |
| written pages | 3,617 |
| Lean proofs | 1,428 |
| commits | about 137 |
| algorithm-agent | about 18,250 lines of Go |
| metered model spend | about $73 over 957 metered events; most writing ran on unmetered interactive sessions |
| hosting failure | GitHub Pages failed at about 3,500 nodes (about 19,500 files) |

**What it means.** Writing speed is not the limit; gates, evidence and review are (21). The pipeline, the tag and provenance packages, the audit JSON and the PR-per-sub-area flow work and are reused (14.1, 20.2). Hosting must be a Worker with R2 from day one (16.20). Stubs must not be published as if they were pages, which is why cards exist (04.1).

## 1.7 AI, agents and documentation in 2026

| Fact | Value |
|---|---|
| Stack Overflow new questions | 966 in Oct 2026, against about 207,000 per month at the 2014 peak |
| Tailwind CSS | docs traffic down about 40%; about 75% of the engineering team laid off; declined to ship llms.txt |
| llms.txt | on about 8.7% of the top 1,000 sites |
| Context7; DeepWiki | MCP servers that serve library docs to coding agents |
| MCP | spec 2026-07-28, with `server/discover` |
| Mintlify | serves `/mcp` for every docs site it hosts |
| Cloudflare "Markdown for Agents" | beta since Feb 2026: Markdown served on `Accept: text/markdown` |
| GitHub Octoverse 2026 | 40% of pull requests start as agent drafts |

**What it means.** Readers now ask agents, and agents read docs. The question-and-answer layer of the web (Stack Overflow) has almost stopped growing, so its knowledge must be captured by topic, not by question. Agents are a main audience: every page has `.md` and `.json` twins, `llms.txt` files, datasets, and an MCP server (16.10). Correctness matters more, not less, because agents repeat what they read; the execution and compat evidence is what makes the site worth reading for them.

## 1.8 The design in one list

1. Cards for the long tail, nodes for ideas (04.1), because the field has millions of named things and thousands of ideas.
2. Concept × language matrix with every cell decided (04.8), because comparison is the most asked question and no site does it with tested code.
3. Every example executed for its stated versions (07), because no reference site does it across languages.
4. Compat per feature per version with probes (09), because the dataset does not exist outside the web.
5. Harvest first (13), because most facts are in public-domain or CC0 catalogues.
6. Algorithm's pipeline and hosting lessons (14, 16.20).
7. Agents as readers (16.10).
