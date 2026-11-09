# 13. Sources and harvest

Most programming facts already exist in structured form: catalogues of languages, release feeds, proposal indexes, compiler error indexes, stdlib inventories, CWE, SPDX. The harvester turns them into cards and compat entries without a model, or with a small model for extraction only. This document defines the sources, the adapters, the raw snapshots, the canon, the licence classes, and the citation check that ties every citation to a stored copy of its source.

## 13.1 Principles

1. **Harvest before you write.** A card made from a structured source is cheaper and more correct than a page that a model writes from memory. The catalog step (14.4) starts from the harvest.
2. **Every fact has a snapshot.** The harvester stores the raw bytes of every source response in R2 `programming-raw`, keyed by BLAKE3. A card's `source` names the snapshot hash. The citation check reads the snapshot, not the live web (13.6).
3. **Licence decides what we keep.** Each source has a licence class (13.5). A class D source gives names, ids and facts only, never text.
4. **Be a good citizen.** Adapters use official APIs and dumps first, respect robots.txt and rate limits, identify themselves with a generic project User-Agent and a project contact address (never a personal address), and cache. A source that blocks automated access (Rosetta Code returned 403 to every automated request in Nov 2026, 01.2) is class E and is not scraped around; we use its names from other sources or by hand.

## 13.2 Sources and adapters

Each adapter lives in `programming-agent/harvest/<source>/` (20) and writes `cards/<kind>/*.jsonl`, `registers/universe/<source>.jsonl`, and compat entries. Many sources already have a local `*-cli` tool in the author's workspace (devdocs, pypi, cratesio, npm, wikidata, stackoverflow and more); an adapter may wrap such a tool, but the adapter is the contract.

| Source | Gives | Access | Class (13.5) | Schedule |
|---|---|---|---|---|
| PLDB | language cards: 5,167 entries, 3,471 tagged as programming languages, 395 fields | `pldb.io/pldb.json` (11.4 MB) | A | weekly |
| Wikidata | language, person, organization, software items; `same_as` ids | SPARQL, dumps | A | weekly |
| GitHub Linguist | language names, aliases, extensions, colours, type (845 entries) | `languages.yml` | C | each release |
| Esolang wiki | esoteric language cards (7,899) | MediaWiki API | A | monthly |
| endoflife.date | release and support dates (489 products) | API v1 | C | daily |
| language release feeds | releases and pre-releases: python.org and `peps.python.org/api/release-cycle.json`, go.dev `dl/?mode=json`, nodejs.org `dist/index.json`, the .NET `releases-index.json`, openjdk.org, GitHub releases of rust-lang, Kotlin, Swift, TypeScript, ziglang.org `download/index.json`, and the rest per language | official APIs | facts | daily (release-watcher, 14.2) |
| proposal indexes | PEPs (`peps.json`), JEPs (openjdk.org/jeps/0), Rust RFCs, TC39 proposals, WG21 and WG14 paper lists, Swift Evolution, KEEPs, C# proposals, PHP RFCs | APIs, git repos, HTML | facts (texts link out) | weekly |
| stdlib indexes | api cards: Sphinx `objects.inv` (Python and every Sphinx-documented project), rustdoc JSON, `go doc` and pkg.go.dev, Javadoc element lists, .NET API docs (the dotnet-api-docs repo), MDN Web API tree, the C and C++ library clause lists, POSIX interface index | tooling in the images and official repos | facts | per release |
| compiler and linter diagnostic indexes | rustc error codes, TypeScript `diagnosticMessages.json`, Roslyn error codes, javac keys, Python exception hierarchy, Go vet analyzers, GCC and Clang warning options, Ruff, ESLint, Clippy, golangci-lint, ShellCheck, RuboCop, Checkstyle, PMD, SpotBugs rule lists | the tools' repos and docs | facts; message templates from the tools' source under their licences | per release |
| MDN BCD | browser compat for web features (20,933 features, 17 browsers, v8.1.8) | npm package | A | each release |
| web-features | Baseline status | npm package | C | each release |
| CWE | weakness cards: CWE 5.0 (2026-10-28) has 956 active weaknesses, 396 groupings, 58 views | XML download | C for ids and names, D for text (MITRE terms, notice) | each release |
| CAPEC | attack patterns (559 in 3.9) | XML | D until checked | each release |
| OWASP | Top 10:2025, ASVS 5.0, Cheat Sheets (137) | git | B (CC BY-SA 4.0) | each release |
| SPDX | licence and exception cards | JSON | A | each release |
| IANA | media types, URI schemes, port numbers used by programmers | registries | facts | monthly |
| package registries | library cards with the scope rule (18.5): npm, PyPI, crates.io, Maven Central, NuGet, RubyGems, Packagist, Go proxy, Hex, CRAN, pub.dev, Hackage, opam; ecosyste.ms as a cross-check | APIs | facts | weekly |
| Stack Overflow | tag names and counts (65,782 tags), question titles of the top questions per tag as a demand signal | Stack Exchange API (not the dumps, whose terms changed in 2024) | B for text; names and counts as facts | monthly |
| Exercism | concepts per track (882 concepts in 24 tracks), problem-specifications cases | git | C (MIT) | monthly |
| tldr-pages | command names for the command-line tool cards (7,614 English pages) | git | C (CC BY 4.0) | monthly |
| DevDocs | the list of documentation sets (247) as a checklist of tools and libraries | `docs.json` | facts | monthly |
| curricula and BoKs | CS2023 (17 knowledge areas, 162 units), SWEBOK v4 (18 knowledge areas), ACM CCS 2012 (SKOS file, 2,113 concepts) | PDFs and SKOS | D | once, then on new versions |
| glossaries | FOLDOC (15,284 entries, names only), Jargon File (public domain), SEVOCAB (names; definitions quotable with IEEE notice) | downloads | D / A / C (single definitions with the IEEE notice) | once |
| Wikipedia | article titles and facts with their citations, for language, concept and person cards; the comparison-of-languages series as a checklist | API | B (CC BY-SA 4.0); we take facts and follow their citations | monthly |
| the canon (13.3) | tables of contents and indexes, by hand or from publisher pages | by hand | D | once |
| research literature | papers for history and theory: DBLP, Semantic Scholar, ACM DL, arXiv, USENIX | APIs | facts and links | on demand |

## 13.3 The canon

The canon is the list of works whose tables of contents and indexes define what "the field knows". It is used as a checklist (the canon match step, 14.1) and as secondary sources. It is never copied. `registers/canon.yaml` lists each work with its ISBN or DOI, edition, and the parts we index. The first list (Draft):

| Area | Works |
|---|---|
| Programming and languages | *Structure and Interpretation of Computer Programs*; *How to Design Programs*; *Concepts, Techniques, and Models of Computer Programming* (Van Roy and Haridi); *Types and Programming Languages* (Pierce); *Practical Foundations for Programming Languages* (Harper); *Concepts in Programming Languages* (Mitchell); *Programming Language Pragmatics* (Scott); *Essentials of Programming Languages*; *Crafting Interpreters* |
| Language books | *The C Programming Language*; *The C++ Programming Language*; *Effective* series (C++, Java, Modern C++); *The Rust Programming Language*; *The Go Programming Language*; *Fluent Python*; *JavaScript: The Definitive Guide*; *Programming in Haskell*; *Real World OCaml*; *Programming Erlang* |
| Design and construction | *Design Patterns* (GoF); *Pattern-Oriented Software Architecture* 1 and 2; *Patterns of Enterprise Application Architecture*; *Enterprise Integration Patterns*; *Refactoring* 2nd ed.; *Code Complete* 2nd ed.; *The Pragmatic Programmer* 20th anniversary ed.; *A Philosophy of Software Design*; *Working Effectively with Legacy Code*; *Domain-Driven Design* |
| Process and practice | *The Mythical Man-Month*; *Peopleware*; *Accelerate*; *Software Engineering at Google*; *Extreme Programming Explained* |
| Testing and debugging | *Test-Driven Development: By Example*; *xUnit Test Patterns*; *Why Programs Fail* |
| Systems | *Computer Systems: A Programmer's Perspective*; *Operating Systems: Three Easy Pieces*; *The Linux Programming Interface*; *Advanced Programming in the UNIX Environment* |
| Concurrency and performance | *The Art of Multiprocessor Programming*; *Java Concurrency in Practice*; *Systems Performance* |
| Security | *The Art of Software Security Assessment*; *Secure Programming Cookbook*; OWASP materials |
| Compilers | *Compilers: Principles, Techniques, and Tools*; *Engineering a Compiler*; *Modern Compiler Implementation* |
| History | the HOPL I to IV proceedings (59 papers; HOPL IV in PACMPL vol. 4, 2020, open access) |

Rule: a canon entry needs an editor's reason. The canon is not a ranking; it is a checklist.

## 13.4 Raw snapshots

- Every adapter run writes its raw responses to R2 `programming-raw/<source>/<date>/<blake3>`, with the URL, the time, the HTTP headers and the licence class.
- The repository keeps only the manifest (`harvest/manifests/<source>.jsonl`: URL, time, hash), never the raw bytes (Default; 20).
- For citations of web pages (official docs, release notes, proposals), the citation check (13.6) fetches and snapshots the cited URL at the version the citation names, once, and keeps it.
- Snapshots of class D sources are kept for checking only. They are never published.

## 13.5 Licence classes

The classes are defined in 19.2 and each source's class is in 19.3; this is the summary the harvester uses. Every adapter writes the class into each snapshot and each card's `source`.

| Class | Licence of the source | What the harvester keeps | Examples (19.3) |
|---|---|---|---|
| A | public domain, CC0 | anything; data may enter the CC0 zone | PLDB, Wikidata, Esolang wiki, MDN BCD, SPDX list data, the Jargon File |
| B | share-alike (CC BY-SA) | facts, ids and names into cards; prose only into the prose zone with attribution, never into code (C06) or data (C07) | Wikipedia, MDN prose, OWASP, Stack Overflow, Learn X in Y minutes, Hyperpolyglot, Programming-Idioms text |
| C | permissive (CC BY, MIT, Apache-2.0, BSD, 0BSD, PSF, Ecma alternative copyright, W3C, Unicode License v3) | facts, structure, short quotes and adapted prose with attribution | Linguist, Exercism, web-features, endoflife.date, tldr-pages, Go docs, Python docs, the Rust Reference and Book, the Kotlin spec, ECMA-262, W3C, Unicode data; CWE ids and names |
| D | all rights reserved, personal-use only, GFDL, non-commercial, or unclear terms | names, ids, facts, section and clause numbers, quotes under G04 only | ISO standards, POSIX, SWEBOK, CS2023, ACM CCS, the canon (13.3), cppreference (by default), FOLDOC, Pro Git, PLAI, CWE text, CAPEC until checked |
| E | forbidden to harvest (terms forbid automated access, or the owner asked us to stop) | nothing by machine; a person may cite it | Rosetta Code while it returns 403 to bots (its text is class D in any case) |

Rule C01 (19.11) enforces it: the provenance of every change names its sources, and a span from a class D source above the quote limit, or an adaptation over the copy-check threshold, fails.

## 13.6 The citation check

The check behind G05 (6.6) and the verifier (15.4):

1. For each citation, resolve the source to a snapshot: the URL at the cited version (for versioned docs, the URL of that version, such as `docs.python.org/3.14/...`), or the standard's clause in our clause index.
2. Check the locator exists in the snapshot: the anchor id in the HTML, the section number in the table of contents, the rule id in the Rust Reference, the line range at the pinned commit.
3. Check that the claim's key terms (the API names, the keywords, the numbers) appear in the located section.
4. For class D sources that we cannot snapshot (ISO standards behind a paywall), check the clause number against the clause index from the free final drafts (for C23, the nearest free draft is N3220) and mark the citation `draft-checked`.

A citation that fails any step fails G05 and goes to the fixer.
