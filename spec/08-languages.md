# 08. Languages

This document defines how the site covers languages: the language universe, the tiers, the depth bands, the language register, the style rules per language, the domain languages, and the comparison system. The counts in this document are from 01.2. Where 01.2 has no number, the number is an Estimate.

## 8.1 What counts as a language

A **language** item (04.2) is any notation that a person writes to tell a computer what to do, and that has a name, a definition or an implementation, and at least one source. This includes:

| Class | Examples | In the universe |
|---|---|---|
| `programming` | general-purpose languages | yes |
| `query` | SQL and its dialects, SPARQL, Cypher, GraphQL, XQuery, Datalog, jq, KQL | yes |
| `shell` | sh, Bash, zsh, fish, PowerShell, cmd, Nushell | yes |
| `markup` | HTML, XML, Markdown, LaTeX (as a language with macros), Typst | only when it has semantics a program runs (HTML with the DOM, TeX with macros); plain markup is a `format` |
| `config` | Nix, Dhall, CUE, Jsonnet, HCL, Starlark, Pkl | yes when it has computation; YAML, TOML, INI are `format` |
| `hdl` | Verilog, SystemVerilog, VHDL, Chisel, SpinalHDL, Bluespec | yes |
| `dsl` | Make, CMake, regex (as a notation family), AWK, sed, Gherkin, OpenSCAD, PLC languages | yes |
| `shading` | GLSL, HLSL, MSL, WGSL, Slang | yes |
| `assembly` | x86-64, AArch64, RISC-V, WebAssembly text, LLVM IR, JVM bytecode | yes, one item per instruction set's assembly; IRs are `language` with class `assembly` |
| `visual` | Scratch, Blockly, LabVIEW G, Simulink, Ladder logic | yes |
| `esoteric` | Brainfuck, Befunge, Malbolge, INTERCAL, Piet | yes, as cards; a node only for the canon of esoteric languages |
| `proof` | Lean, Rocq, Agda, Isabelle/Isar, Idris, F*, Dafny | yes as languages; their mathematical use is math2's subject |
| `historical` | Plankalkül, FLOW-MATIC, Short Code, IAL | yes; the class is a flag `historical: true` on top of the class above |

A dialect with its own name is a `language` with `variant_of` (04.2). A version is a `release`, not a language. A language that exists only in one paper with no implementation is a card with `implemented: false`.

## 8.2 The tiers

| Tier | Meaning | Members | Depth band (8.3) | Images (07.4) | Examples in concepts | Tasks |
|---|---|---|---|---|---|---|
| **1** | the languages most programmers use; every core concept and task has a section in each | Python, JavaScript, TypeScript, Java, C, C++, C#, Go, Rust, Kotlin (the 10 of algorithm), plus the domain languages SQL, Bash, regex, HTML/CSS, PowerShell | D4 | full version matrix | all cells of the matrix decided; examples where `native`, `stdlib`, `idiom` | all core tasks |
| **2** | languages with a large community or a distinct model that teaches something | Swift, Ruby, PHP, Dart, Scala, Haskell, OCaml, F#, Elixir, Julia, R, Zig, Lua (the 13 of algorithm), plus Objective-C, Perl, Erlang, Clojure, Common Lisp, Scheme (Racket), Fortran, COBOL, MATLAB/Octave, Nim, Crystal, Gleam, Mojo, Solidity, Ada, Prolog, Groovy, Visual Basic .NET, Assembly (x86-64 and AArch64) | D3 | newest and oldest supported, full matrix at M4 | all cells decided at W4 | core tasks at M3 |
| **3** | any language with a runnable open toolchain that we can put in an image | e.g. D, Pascal (Free Pascal), Tcl, Smalltalk (Pharo), Forth, APL (GNU APL, Dyalog has no open image), J, K (ngn/k), Raku, Hack, V, Odin, Carbon (when it runs), Idris, Agda, Lean, Elm, PureScript, ReScript, Haxe, Vala, Ballerina, Pony, Io, Factor, Red, Janet, Fennel, Hy, Awk, sed, Make, CMake, jq, Nix, Starlark, Verilog (Icarus, Verilator), VHDL (GHDL), Brainfuck, Befunge, ... | D2 | newest only | none required | the "hello world", "FizzBuzz" and 10 to 20 more tasks of the tier-3 task set (10.1) |
| **4** | every other language in the universe | all the rest | D1 | none | none | none |

Rules:

- A language moves up a tier by a register change with a reason, reviewed by an editor. The reasons that count are: the surveys of 01.2 (a rank in the top 30 of two of TIOBE, PYPL, the Stack Overflow survey, the GitHub Octoverse, the IEEE Spectrum ranking, the RedMonk ranking), or a model that the concept pages need (for example Prolog for logic programming, APL for array programming, Erlang for actors, Smalltalk for message passing).
- A tier-1 or tier-2 language never moves down because its rank falls. It may become `historical`.
- Tier 2 grows from algorithm's 13 to about 34. Algorithm keeps its own 23 languages (10 + 13). Programming's tier 2 additions do not change algorithm (5.8).

## 8.3 Depth bands

| Band | Name | What the language has | Pages |
|---|---|---|---|
| D1 | card | the language card: names, class, years, designers, paradigms, typing, influences, file extensions, hello world (from a source, not executed), links, `same_as` ids | 0 (an entry on the catalogue page and the family page) |
| D2 | profile | the language node at W3, hello world and 20 tasks executed, the 25 template sub-areas as README pages with the main features as cards, releases as cards | 1 node, 25 short area pages, feature cards |
| D3 | guide | D2, and the features of the language's own tutorial as nodes at W3, all matrix cells decided, the core tasks, the diagnostics and proposals as cards, the compat at K1 | about 150 to 400 nodes |
| D4 | complete guide | D3, and every feature of the language reference as a node (one per section of the reference at the least), the full stdlib as api cards, the compat at K2 over the support window, every diagnostic and lint of the main tools as a card, the pitfalls, the idioms, the migration guides, every release since the first as a card and the releases in the support window as nodes | about 500 to 2,000 nodes each, plus 2,000 to 40,000 api cards each |

"Every section of the reference" is the measurable meaning of "no missing" for a language at D4: the harvester parses the table of contents of the official reference (the Python Language Reference, the JLS, the ECMAScript specification, the C and C++ standards' clause list, the Go spec, the Rust Reference with its rule ids, the C# specification, the Kotlin spec, the TypeScript handbook since TypeScript has no current spec), and each section must map to a feature node, a section of a feature node, or an outcome in the universe (18.3).

## 8.4 The language register

`registers/languages.yaml` is a superset of algorithm's register (2150/07). New fields are marked ✱.

```yaml
- id: python
  name: Python
  tier: 1
  band: D4                                    # ✱
  class: programming                          # ✱
  ext: py
  file_case: snake
  toolchain: "CPython 3.15"                   # the default in examples (3.15.0 released 2026-10-09)
  window: {floor: "3.11", ceil: "3.15"}       # ✱ the support window (09.3), computed, stored for the build
  images: {default: "python:3.15", matrix: ["3.9","3.10","3.11","3.12","3.13","3.14","3.15"]}   # ✱ tags; digests in images.yaml
  implementations: [cpython, pypy, graalpy, micropython]   # ✱
  test: pytest
  format: "ruff format"
  lint: "ruff check"
  typecheck: "mypy --strict"                  # ✱ for languages with optional static checking
  run: "python3 {file}"
  repl: {prompt: ">>> ", cont: "... ", cmd: "python3 -i"}   # ✱
  comment: "#"
  highlight: python                           # Shiki id
  treesitter: python                          # ✱ for symbol extraction (15.3)
  browser: pyodide                            # ✱ 07.10
  playground: null                            # ✱ external playground URL template
  reference: {url: "https://docs.python.org/3/reference/", toc: "harvest/python-reference-toc"}   # ✱ the D4 section list
  stdlib_index: {kind: sphinx-inventory, url: "https://docs.python.org/3/objects.inv"}   # ✱ for api cards (13.2)
  diagnostics: [{tool: cpython, kind: exception-hierarchy}, {tool: ruff, kind: lint-rules}, {tool: mypy, kind: error-codes}, {tool: pylint, kind: messages}]   # ✱
  proposals: {kind: pep, index: "https://peps.python.org/api/peps.json"}   # ✱
  releases: {kind: release-cycle, url: "https://peps.python.org/api/release-cycle.json"}   # ✱
  style: {guide: pep-8, file: languages/python/style.yaml}   # ✱
  int: int
  list: list
  since: 2026-11
```

The register lives in `programming` and is the source of truth. Algorithm copies its 23 rows from it (5.8).

## 8.5 Style per language

`languages/<lang>/style.yaml` holds the conventions of examples in the language, so that examples are idiomatic and consistent across pages:

- the formatter and its config, the linter and its rule set, the type checker for languages with optional types;
- naming conventions (from the language's own guide: PEP 8, Effective Go, the Rust API guidelines, the Google Java Style Guide, the C# coding conventions, the Kotlin coding conventions, the C++ Core Guidelines for C++, MISRA only on the embedded pages);
- the default error style for examples (exceptions or result values), the default output call (`print`, `console.log`, `fmt.Println`, `println!`, `System.out.println`, `std::println` for C++23 and later with a `std::cout` variant for older standards);
- the general-term to language-term map (6.5);
- the version floor: the oldest version an example may assume without a `since` attribute (the support window floor, 09.3).

## 8.6 Domain languages

Tier 1 includes five domain languages because almost every programmer uses them, and their knowledge is as large as a general language's. Each has its own template (5.3) and its own implementations as compat columns.

| Language | Implementations as compat columns | Template sub-areas (beyond overview, getting started, pitfalls, versions) | Source of truth |
|---|---|---|---|
| SQL | ISO SQL:2023 (the standard as a column), PostgreSQL, MySQL, MariaDB, SQLite, SQL Server, Oracle, DuckDB, BigQuery, Snowflake | queries, joins, subqueries and CTEs, aggregation, window functions, set operations, NULL semantics, data types, DDL, constraints, indexes, transactions and isolation, procedural extensions, JSON in SQL, dialect differences, query planning | each engine's manual; ISO 9075 clauses by number only |
| Bash and the POSIX shell | POSIX sh (POSIX.1-2024), Bash, dash, zsh, ksh93, BusyBox ash, fish (as a separate language with its own tree) | commands and pipelines, quoting and expansion, parameters and variables, control flow, functions, redirection and file descriptors, here-documents, job control, traps and signals, the core utilities (one feature per utility: grep, sed, awk, find, xargs, sort, cut, tr, ...), portability | POSIX.1-2024 (cite by section), the Bash manual (GFDL, so facts and locators only, 19.3) |
| Regular expressions | PCRE2, ECMAScript (V8), Python `re`, Java, .NET, Go RE2, Rust `regex`, POSIX BRE and ERE, Oniguruma (Ruby), ICU | syntax, character classes, Unicode properties, quantifiers (greedy, lazy, possessive), groups and backreferences, lookaround, anchors and boundaries, flags, substitution, flavours, performance and ReDoS | each engine's documentation; Unicode TS #18 |
| HTML and CSS | browsers as columns from MDN BCD (CC0), plus the Baseline status from web-features | HTML elements and attributes (as api cards), the DOM, forms, accessibility, CSS selectors, the cascade, box model, layout (flow, flexbox, grid), units, typography, animations, responsive design, custom properties, at-rules | WHATWG HTML, the CSS specifications, MDN BCD |
| PowerShell | Windows PowerShell 5.1, PowerShell 7.x | the pipeline of objects, cmdlets, providers, scripting, modules, remoting, error handling, the cmdlet reference as api cards | Microsoft Learn (MIT for the docs repo, CC BY 4.0 for prose; checked in 19.3) |

## 8.7 The catalogue of all languages

The catalogue is the D1 layer for every language in the universe.

- **Sources** (13.2): PLDB (public domain; its CSV on Hugging Face), Wikidata (CC0; items that are an instance of "programming language" and its subclasses), GitHub Linguist (MIT; names, extensions, colours, aliases), the Esolang wiki (CC0 for content: names and facts only), HOPL (the History of Programming Languages online index, facts only), Rosetta Code's language list (names only), Wikipedia's lists (facts, with citations to the article). The counts of each source are in 01.2.
- **Merge**: the harvester matches the sources by `same_as` ids first, then by name and year and designer. A match below the threshold goes to the curator (14.3). The output is one card per language with its sources.
- **Family pages**: each family of the genealogy (12.6) has an area page that lists its languages as cards: the Fortran family, the Algol family, the Lisp family, the ML family, the C family, the Pascal and Modula family, Smalltalk and its children, the shells, the array languages, the stack languages, the logic languages, the BASIC family, the scripting family, the .NET languages, the JVM languages, the BEAM languages, the query languages, the HDLs, the shading languages, the esoteric languages.
- **Rule**: every language in PLDB with a rank, and every Wikidata item with a Wikipedia article in English that is an instance of a programming language class, is in the universe with an outcome (18.3). This is the "no missing" bar for the catalogue.

## 8.8 Comparisons

Comparison pages are generated, not written, from the facts the nodes already hold, so they cannot disagree with the language pages.

- `/compare/<a>-vs-<b>/` for every pair of tier-1 and tier-2 languages that a reader asks for (the pairs come from search logs and the list of "coming from X" pairs in `registers/compare-pairs.yaml`; at M3 about 200 pairs, Estimate). It shows: the facts box side by side, the concept matrix rows side by side with the examples, the tasks side by side, the differences that the matrix shows (the cells that differ), and the "coming from A" notes that a writer writes for each direction (a short node, `compare/<a>-to-<b>`, 300 to 800 words, required for the top 50 pairs at W3).
- `/compare/` lets the reader pick any languages (up to 6) and any concepts or tasks. This is the Hyperpolyglot model, made from executed examples.
- A comparison never ranks the languages. It never shows a benchmark that is not in a node (7.9).

## 8.9 Guides as reading paths

Each D3 and D4 language has two reading paths (16.12), computed from the tree and the prereqs:

- **Learn L**: the tutorial order for a new programmer, from the intro-level features.
- **L for programmers of M**: for each tier-1 language M, the order that skips what M already teaches and starts with the differences (from the matrix cells that differ between L and M).
