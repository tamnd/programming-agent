# 04. The object model

This document defines what the content is made of: the kinds of item, the difference between a node and a card, the identifiers, the directory of a node, the front matter of each kind, the edges, and the ladders that the build computes.
It follows 2150/04 (algorithm) where it can, and says where it differs.

## 4.1 Nodes and cards

Programming has two orders of magnitude more named things than algorithms: about 5,000 languages, tens of thousands of error codes and lint rules, thousands of change proposals, and hundreds of thousands of standard-library entities.
A page written by hand for each would be thin, would never finish, and would drown the pages that matter.
So an item exists in one of two forms.

| Form | What it is | Where it lives | Rendered as | Written by |
|---|---|---|---|---|
| **card** | a row of checked facts with sources: names, dates, versions, links, a one-sentence summary, the edges | a JSONL register under `cards/<kind>/` | an entry on a collection page, with a stable anchor and its own JSON record | the harvester from structured sources (13), checked by rules; a person or an agent only to fix it |
| **node** | a directory with a written page, examples, tests, and optional proofs, scenes and Lean | `nodes/<area>/<sub-area>/<slug>/` | a page of its own | writers, people and agents (14) |

Rules:

- Every kind can exist as a card. Most kinds start as cards and are **promoted** to nodes when someone writes them. `prog promote <id>` makes the node directory from the card. It moves the facts into the front matter and removes the row from the register, in one commit. The tag does not change.
- A card is never a stub page. A page exists only for a node at W1 or more (06.8). Below that, the reader sees the card on its collection page, and the card links to "Help write this page" (16.19).
- Four kinds are **card-only by default**: `diagnostic`, `proposal`, `api` and `exercise` (4.2). They become nodes only when the editors decide that the item needs prose. Examples are a famous error such as E0382 "borrow of moved value", a landmark proposal such as PEP 484, or a central API such as `fetch`.
- The coverage levels (18.2) count cards and nodes on one scale. A card is C1. A node starts at C2.

## 4.2 The kinds

| Kind | It is | Example | Default form | Has examples | Has compat | Has Lean | Has a scene |
|---|---|---|---|---|---|---|---|
| `area` | a node of the topic tree or a language tree (05) | Concurrency; Python › Functions | node | no | no | no | optional |
| `concept` | a language-independent idea of programming | closure, garbage collection, ownership, type inference, async/await | node | yes, in each language where it applies | through the features that realize it | for PL-theory concepts (11) | optional |
| `definition` | a precise term, as in a glossary | lvalue, side effect, referential transparency, ABI | card, then node | optional | no | optional | no |
| `language` | a programming, query, shell, markup-with-semantics, configuration, hardware-description or esoteric language, or a named dialect | Python, Rust, COBOL, SQL, Bash, Verilog, Brainfuck | card for most, node for tiers 1 to 3 | hello world and characteristic examples | the release table | no | no |
| `feature` | a construct, rule or facility of one language or its standard library | Python `match`; Rust `?`; C++ `constexpr`; Go `defer`; Java records; JS optional chaining | node | yes | yes | no | optional |
| `release` | a version of a language, a standard or a tool | Python 3.15; Java 27; C++23; ECMAScript 2026; Go 1.27 | card, node for the releases of tiers 1 and 2 | the new features have them | it is a column of the compat tables | no | no |
| `task` | a common programming job with a contract, solved in many languages | read a file line by line; parse JSON; send an HTTP GET; format a date; reverse a string by grapheme | node | an implementation in each language | through the features it uses | no | optional |
| `pattern` | a named, reusable design: a design pattern, an architectural or integration pattern, a concurrency pattern, or a negative one (an anti-pattern or a code smell) | Observer; Hexagonal architecture; Circuit breaker; God object; Long method | node | yes | no | no | optional |
| `technique` | a named method with steps that a programmer applies to code or to a program's history | Extract Function; property-based testing; git bisect; memoization in practice; rubber duck debugging | node | before and after code | no | no | optional |
| `practice` | a principle, a law, a guideline, a methodology or a workflow | SOLID; DRY; Conway's law; Hyrum's law; TDD; trunk-based development; semantic versioning | node | optional | no | no | no |
| `tool` | software that programmers use to make software: compilers, interpreters, runtimes, VMs, build systems, package managers, version control, editors, debuggers, profilers, linters, formatters, CI services, AI coding agents | GCC; CPython; V8; Cargo; Git; Vim; GDB; Ruff; GitHub Actions | card, node for notable tools | usage examples | the release table | no | optional |
| `library` | a library or a framework, including a standard library as a whole | React; Django; NumPy; Tokio; Spring; the C++ standard library | card, node for notable ones | yes | the release table | no | no |
| `format` | a data format, an encoding, a wire protocol as programmers use it, or a notation | JSON; YAML; Protocol Buffers; UTF-8; HTTP/1.1 messages; semver strings; regular-expression syntax | card, node for notable ones | parse and emit examples | no | optional (grammars) | optional |
| `standard` | a specification document | ISO/IEC 9899:2024; ECMA-262; POSIX.1-2024; RFC 8259; Unicode 18 | card, node for the core | no | it defines a release | no | no |
| `weakness` | a class of defect: a security weakness (CWE) or a general bug class | SQL injection (CWE-89); use after free (CWE-416); off-by-one; data race | card from CWE, node for the core | vulnerable and fixed code | no | no | optional |
| `diagnostic` | an error, warning, exception type or lint rule of a language or a tool | Rust E0382; TypeScript TS2339; Python `KeyError`; Clippy `needless_range_loop`; Ruff `E501` | card | the code that triggers it, and the fix | the versions where it exists | no | no |
| `proposal` | a change proposal of a language or a tool, with its status | PEP 634; JEP 444; Rust RFC 2094; TC39 "Temporal"; WG21 P2300 | card | no | it links the release that shipped it | no | no |
| `api` | an entity of a standard library or a platform interface: a module, a type, a function, a method, a constant, a system call, a web API | Python `str.split`; Java `java.util.HashMap`; POSIX `fork`; Web `fetch`; C `memcpy` | card | linked from features and tasks | `since`, `deprecated`, `removed` | no | no |
| `work` | a book, a paper, an essay or a talk that is part of the canon of programming | *The C Programming Language*; *SICP*; "Go To Statement Considered Harmful"; *Design Patterns* | card, node for the canon | no | no | no | no |
| `person` | a person who made or named items | Grace Hopper; Dennis Ritchie; Barbara Liskov | card, node for the notable | no | no | no | no |
| `organization` | a company, a foundation, a committee or a standards body | PSF; WG21; TC39; the Rust Foundation; Bell Labs | card | no | no | no | no |
| `exercise` | a practice problem on a public site, as a link with metadata | Exercism `python/two-fer`; a LeetCode problem; an Advent of Code day | card | no | no | no | no |

Notes:

- **Algorithms and data structures are not kinds here.** They live on algorithm.tamnd.com. A page here links to them with `[[al:merge-sort]]` (06.2). A *task* can use an algorithm. For example, "sort a list of records by two keys" links to merge sort for the stability argument. The task page is about the language facilities, not the algorithm.
- **Mathematics is not a kind here.** It lives on math2.tamnd.com and is linked with `[[m2:<tag>]]`.
- **A concept and a feature are different things.** The concept `closure` is one page with a section for each language where closures exist. The feature `python/closures-and-nonlocal` is the Python-specific page: the rules of `nonlocal`, the late-binding pitfall in loops, and `__closure__`. A feature `realizes` one or more concepts (4.6). Most concepts get their language sections from the features that realize them, so the facts are written once.
- **Variants.** A dialect is a `language` with `variant_of` (MicroPython → Python). A fork of a tool is a `tool` with `variant_of` (OpenTofu → Terraform). An implementation of a language is a `tool` with `role: implementation` and `implements_language` (CPython, PyPy, GraalPy → Python).
- **A release is not a page per patch version.** Releases are minor or major versions (Python 3.15, Go 1.27, Rust 1.99). Patch versions are rows in the release's front matter, unless a patch changed behaviour that a compat entry needs (then the compat entry names the patch).
- **Proofs** exist only for the concepts of 11 (formal semantics) and as arguments inside pages. A proof is not a node, as in algorithm.

## 4.3 Tags, slugs, namespaces and URIs

- **Tag.** 6 characters of Crockford base32, the alphabet of math2 and algorithm, for example `7QKX2M`. A tag is permanent and never reused. `registers/tags.reg` lists each tag. A new item has the provisional id `NEW-<slug>` until `prog mint` gives it a tag before the merge. Cards get their tags at harvest time, in bulk, from the same minting code.
- **URI.** `pg://7QKX2M`, and `pg://7QKX2M@2027.03` for the item in a release. Fragments name parts: `pg://7QKX2M#pitfall-late-binding`, `pg://7QKX2M/rust`. The scheme differs from `al://` (algorithm) and `m2://` (math2), so a link always says which project it means (22, Q-02).
- **Global slugs.** These kinds have one slug that is unique over the global namespace and is the URL: `area` (top level), `concept`, `definition`, `language`, `task`, `pattern`, `technique`, `practice`, `tool`, `library`, `format`, `standard`, `weakness`, `work`, `person`, `organization`. Examples: `closures`, `python`, `read-a-file-line-by-line`, `observer-pattern`, `git`, `react`, `json`, `sql-injection`.
- **Owned slugs.** These kinds are named inside the language, tool or library that owns them, and their URL has the owner first:

| Kind | Owner | URL form | Example |
|---|---|---|---|
| `feature` | a language or a library | `/<owner>/<slug>/` | `/python/match-statement/`, `/rust/question-mark-operator/` |
| `release` | a language, standard, tool or library | `/<owner>/<version>/` | `/python/3.14/`, `/cpp/23/`, `/node/24/` |
| `diagnostic` | a language or a tool | `/<owner>/errors/<code>/` or `/<owner>/lints/<code>/` | `/rust/errors/E0382/`, `/ruff/lints/E501/` |
| `api` | a library (a standard library included) | `/<owner>/api/<path>/` | `/python/api/str.split/`, `/posix/api/fork/` |
| `proposal` | a language or a tool | `/<owner>/proposals/<id>/` | `/python/proposals/pep-634/`, `/java/proposals/jep-444/` |
| `exercise` | a site | `/exercises/<site>/<id>/` | `/exercises/exercism/python-two-fer/` |

  An owned slug is unique inside its owner. The owner slug must be a node or card of a kind that can own (language, tool, library, standard, organization for `exercise` sites).
- **Reserved words.** No global slug may be a reserved word of the site (16.2): `api`, `errors`, `lints`, `proposals`, `compare`, `compat`, `tasks`, `concepts`, `languages`, `search`, `about`, `spec`, `t`, `x`, and the rest of `registers/reserved.yaml`. No owned slug may be a version string unless it is a release.
- **Aliases.** `aliases` lists other names, each with a language when it is not English. A rename moves the old slug to `registers/aliases.reg`, and the site redirects it, as in algorithm. CWE ids, PEP numbers and error codes are aliases that search resolves directly: `cwe-89`, `pep-634` and `E0382` go straight to their item.
- **Disambiguation.** Shared names get qualified slugs and an entry in `registers/disambiguation.yaml`. Examples: `go` the language and `go-tool` the command; `rust` and `rust-game`; `swift` the language and `swift-parallel-scripting`; `d` the language and `d-programming-language` as an alias; `c` and the many C-named things; "R"; "Julia"; "heap" (memory) on this site against "heap" (data structure) on algorithm. Single-letter and common-word language names always carry `-language` in the aliases and are matched only with a context word in search (17.6).

## 4.4 The node directory

A node lives at `nodes/<top>/<sub>/<slug>/` for the topic tree, and at `nodes/languages/<lang>/<sub>/<slug>/` for the language trees (05.3). A concept directory:

```
nodes/functions/closures-and-scope/closures/
  node.md                    front matter and the text (06)
  examples/
    python.md                the Python section: prose, run blocks, output blocks (07)
    javascript.md
    rust.md
    go.md
    …                        one file for each language where the concept applies
  matrix.yaml                the decisions for languages with no example: absent, library, idiom (4.8)
  viz/                       optional scene (12)
  figures/                   optional SVG
```

A feature directory:

```
nodes/languages/python/control-flow/match-statement/
  node.md
  examples/python.md         the examples; a feature has its owner's language only
  compat.yaml                the compat entries of this feature (09), with probes
  probes/
    match_basic.py           a probe: fails before 3.10, passes from 3.10 (09.4)
```

A task directory has a contract and implementations, with the same machinery as an algorithm problem (10):

```
nodes/io/files/read-a-file-line-by-line/
  node.md
  contract.yaml              the signature or program contract, fixtures, the checker (10.2)
  cases/basic.jsonl
  cases/edge.jsonl           empty file, no final newline, CRLF, a 1 MB line, invalid UTF-8
  fixtures/                  input files that the cases name
  impl/
    python/read_lines.py
    python/test_read_lines.py
    rust/read_lines.rs
    …
```

A language directory is the root of its language tree. It holds the overview page and the registers of the language:

```
nodes/languages/python/
  node.md                    the language page (06.3)
  releases.yaml              every release with dates, end of life and the compat column id (09.2)
  implementations.yaml       CPython, PyPy, GraalPy, MicroPython, … with the compat columns they get
  concepts.yaml              the matrix row: for each core concept, native, library, idiom, absent, n/a (4.8)
  style.yaml                 the conventions of examples in this language (08.4)
  <sub-area>/<feature>/      the feature nodes, in the fixed sub-areas of the language template (05.3)
```

The Lean part, when there is one, lives in the Lean package, `lean/Prog/<Area>/<Name>.lean`, as in algorithm (11).

## 4.5 The front matter

The schema is `schema/node.v1.json`, one definition per kind. Fields that do not apply to a kind are absent. Common fields for every kind:

```yaml
---
tag: 7QKX2M
kind: concept
slug: closures
name: Closure
aliases: [lexical closure, function closure, "閉包" (zh), "замыкание" (ru)]
area: functions/closures-and-scope          # primary area (05.2)
areas: [paradigms/functional-programming]   # secondary areas
level: intro                                # intro | core | advanced | research (05.6)
summary: >-
  A function together with the variables of the scope where it was made, which it can still read and change after that scope has ended.
prereqs: [first-class-functions, lexical-scope, call-stack]
see_also: [lambda-expressions, partial-application, objects-as-closures]
history:
  year: 1964
  people: [peter-landin]
  first: {source: landin-1964, note: "the term in 'The mechanical evaluation of expressions'"}
same_as:
  wikidata: Q330263
  wikipedia: Closure_(computer_programming)
  pldb: null
  algorithm: null                             # an al:// tag when the item is also on algorithm
  math2: null
sources:
  - {work: "isbn:9780262510875", locator: "SICP 2e, Section 3.2, The environment model of evaluation"}
  - {url: "https://docs.python.org/3/reference/executionmodel.html#resolution-of-names", locator: "Python Language Reference, 4.2.2 Resolution of names", accessed: 2026-11-20}
---
```

Kind-specific fields:

**concept**

```yaml
applies_to: core-languages        # which matrix rows must decide this concept (4.8)
viz: expected                     # expected | optional | none: whether a scene is part of "complete"
formal: none                      # expected | optional | none: whether a Lean part is part of "complete" (11)
```

**language**

```yaml
kind: language
slug: python
name: Python
class: programming                 # programming | query | shell | markup | config | hdl | esoteric | dsl | assembly | visual
tier: 1                            # 1 | 2 | 3 | 4 (08.2); absent for cards of tier 4
paradigms: [imperative, object-oriented, functional, procedural]
typing: {discipline: dynamic, strength: strong, checking: [runtime, optional-static], inference: n/a}
memory: garbage-collected          # garbage-collected | reference-counted | manual | ownership | region | mixed
execution: [interpreted-bytecode, jit]   # compiled-native | interpreted | interpreted-bytecode | jit | transpiled | vm
appeared: 1991
designers: [guido-van-rossum]
influenced_by: [abc-language, modula-3, c, lisp, icon, haskell]
standard: null                     # the standard slug, when the language has one
reference: {url: "https://docs.python.org/3/reference/", licence: psf-2.0}
licence_of_reference_impl: psf-2.0
file_extensions: [py, pyi, pyw]
hello_world: examples/python.md#hello
popularity: {tiobe: 1, pypl: 1, octoverse: 2, so_survey: 3, as_of: 2026-11}   # build-filled, never by hand (rule T05)
same_as: {wikidata: Q28865, wikipedia: Python_(programming_language), pldb: python, linguist: Python, rosetta: Python, hopl: null}
```

**feature**

```yaml
kind: feature
slug: match-statement
owner: python
name: The match statement
area: languages/python/control-flow
realizes: [pattern-matching, algebraic-data-types-in-practice]
grammar: "match_stmt in the Python grammar"       # where the syntax is defined
spec: {standard: null, url: "https://docs.python.org/3/reference/compound_stmts.html#the-match-statement", locator: "8.6"}
proposals: [pep-634, pep-635, pep-636]
compat: compat.yaml
since: "3.10"                                     # copied from compat by the build; a hand-written value must agree (rule K03)
```

**release**

```yaml
kind: release
owner: python
version: "3.14"
released: 2025-10-07
eol: 2030-10                                      # from the owner's policy; endoflife.date as a cross-check (13.2)
status: supported                                 # build-computed from the dates, never by hand
patches: [{v: 3.14.0, date: 2025-10-07}, {v: 3.14.1, date: ...}]
adds: [template-strings, free-threaded-build-supported, deferred-annotations]   # feature slugs; the build checks each against compat
removes: [...]
deprecates: [...]
notes: {url: "https://docs.python.org/3.14/whatsnew/3.14.html"}
```

**task**

```yaml
kind: task
slug: read-a-file-line-by-line
contract: contract.yaml
category: io/files
languages: tier-1+2                # the languages that must have an implementation
uses_features: {python: [python/with-statement, python/file-objects], rust: [rust/bufread-lines], go: [go/bufio-scanner]}
uses_algorithms: []                # al:// links
```

**pattern**, **technique**, **practice**

```yaml
kind: pattern
family: behavioral-design          # creational-design | structural-design | behavioral-design | architectural | integration | concurrency | cloud | resilience | testing | code-smell | anti-pattern | ...
polarity: positive                 # positive | negative (anti-patterns and smells)
intent: "Let objects be told when another object's state changes, without coupling them to it."
also_known_as: [publish-subscribe (loosely), dependents]
related: [mediator-pattern, event-emitter]
languages_where_built_in: {csharp: csharp/events, javascript: javascript/eventtarget}
```

**tool**, **library**

```yaml
kind: tool
role: [version-control]           # implementation | compiler | interpreter | runtime | vm | build | package-manager | version-control | editor | ide | debugger | profiler | linter | formatter | static-analyzer | test-runner | ci | container | ai-assistant | shell | terminal | documentation | ...
implements_language: null         # for role implementation
languages: []                     # the languages it serves
licence: gpl-2.0-only             # SPDX id
repo: https://github.com/git/git
first_release: 2005-04-07
latest: {version: "2.52.0", date: ..., checked: 2026-11-20}   # build-filled from the release feed
```

**weakness**

```yaml
kind: weakness
cwe: 89
abstraction: base                 # pillar | class | base | variant | compound (as in CWE)
languages_affected: [any-with-sql]
owasp_top10_2025: A05
cwe_top25_2025: 2
detected_by: [semgrep-rule-x, codeql-query-y]   # tool rule slugs, as cards
```

**diagnostic** (card)

```json
{"tag":"3MFR8Q","kind":"diagnostic","owner":"rust","code":"E0382","title":"borrow of moved value",
 "severity":"error","since":"1.0.0","removed":null,
 "source":{"url":"https://doc.rust-lang.org/error_codes/E0382.html","licence":"MIT OR Apache-2.0"},
 "concepts":["move-semantics","ownership"],"features":["rust/ownership","rust/borrowing"],
 "message_templates":["borrow of moved value: `{}`","use of moved value: `{}`"],
 "trigger":"probes/E0382.rs","fix":"probes/E0382_fixed.rs"}
```

**proposal** (card)

```json
{"tag":"1VQ7KD","kind":"proposal","owner":"python","id":"pep-634","title":"Structural Pattern Matching: Specification",
 "status":"final","type":"standards-track","created":"2020-09-12","release":"python/3.10",
 "authors":["brandt-bucher","guido-van-rossum"],"features":["python/match-statement"],
 "url":"https://peps.python.org/pep-0634/","licence":"public-domain-or-cc0"}
```

**api** (card)

```json
{"tag":"9ZQ1HT","kind":"api","owner":"python","path":"str.split","entity":"method",
 "signature":"str.split(sep=None, maxsplit=-1)","module":"builtins",
 "since":"1.6","deprecated":null,"removed":null,
 "doc":{"url":"https://docs.python.org/3/library/stdtypes.html#str.split"},
 "source":{"inventory":"objects.inv", "version":"3.14"},
 "features":["python/str-methods"],"tasks":["split-a-string"]}
```

Rules, as in algorithm and beyond it:

- No status, level, badge or popularity in front matter written by hand (T05).
- Every date, version and number in front matter has a source in `sources` or comes from a register with a source (G01).
- `since`, `deprecated` and `removed` come from compat entries (09). A hand-written value that disagrees with compat fails rule K03.
- Every claim of time or space cost of an API or a construct (for example "`list.append` is amortized O(1)") names its model and case, as in algorithm 06.6. The build parses it with the complexity grammar of 2150/04.4.

## 4.6 Edges

Edges are typed fields of the front matter. The build derives the reverse edges.

| Edge | From → to | Meaning |
|---|---|---|
| `realizes` | feature → concept | the feature is the language's form of the concept |
| `owner` | feature, release, diagnostic, api, proposal → language, tool, library, standard | the namespace and the parent |
| `implements_language` | tool → language | the tool is an implementation of the language |
| `variant_of` | language, tool, library, pattern, technique → same kind | a dialect, a fork, a named variant |
| `influenced_by` | language → language | a design influence, with a source (the genealogy, 12.6) |
| `standardized_by` | language, format → standard | the standard that defines it |
| `defines` | standard → release, feature, format | the standard defines this release or construct |
| `proposes` | proposal → feature | the proposal added or changed the feature |
| `ships_in` | proposal, feature → release | the release where it first shipped (from compat) |
| `uses_features` | task, pattern, example → feature | the solution uses the feature |
| `uses_algorithms` | any → algorithm (`al://`) | a link to algorithm.tamnd.com |
| `uses_math` | any → math2 (`m2://`) | a link to math2.tamnd.com |
| `prevents` | practice, technique, feature, tool → weakness | the item helps to prevent the weakness |
| `detects` | tool, diagnostic → weakness | the tool or rule finds the weakness |
| `triggers` | weakness, feature → diagnostic | the code shape that causes the diagnostic |
| `alternative_to` | library, tool, pattern → same kind | a choice for the same need, shown in comparisons |
| `prereqs` | any → any | read the target first; a DAG (rule E01) |
| `people` | any → person | made, named or changed it |
| `organizations` | any → organization | steward, maker or standards body |
| `exercises` | concept, feature, task, pattern → exercise | practice for the target |
| `see_also` | any → any | at most 8 |

The graph rules of algorithm (no cycle in `prereqs`, `variant_of`; each target of the right kind) apply. Two rules are new. Every `feature` realizes at least one concept, or carries `realizes: []` with the reason `language-specific` (E10). Every concept of the core has a decision for each language of its `applies_to` (E11, 4.8).

## 4.7 Provenance

As in algorithm 04.7: `prov/<tag>.jsonl`, one line per change of a file: the date, the file, the operation, the author (a GitHub login or an agent run id), the reviewer, the sources used. Cards have provenance per harvest run: `prov/cards/<kind>.jsonl` records the run id, the source snapshot (BLAKE3 of the raw capture) and the rows that the run wrote.

## 4.8 The concept × language matrix

The matrix answers "how does language L do concept C?" for every pair that matters. It is also the main tool for "no missing" in the language trees.

- **Rows** are languages of tiers 1 and 2. **Columns** are the concepts with `applies_to: core-languages` (about 400 at the start, 05.4).
- **A cell** is one of:

| Value | Meaning | Evidence that the build needs |
|---|---|---|
| `native` | the language has the concept as a construct or a rule | a feature that `realizes` it, with an executed example |
| `stdlib` | the standard library provides it | an `api` card or a feature, with an executed example |
| `library` | a well-known library provides it | a `library` card and an executed example with the library pinned (07.7) |
| `idiom` | there is no construct, but a standard way to get the effect | an executed example in the concept's language section, marked as an idiom |
| `absent` | the language cannot express it, or only by building an interpreter | one line of reason, and a source when the claim is not obvious |
| `n/a` | the question does not apply (garbage collection in a language with no heap) | one line of reason |

- `native` and `stdlib` cells come from the features and api cards. `library`, `idiom`, `absent` and `n/a` cells are written in the language's `concepts.yaml` or the concept's `matrix.yaml`.
- Rule E11: a concept of the core with an empty cell in a tier-1 row fails at W3, and in a tier-2 row at W4.
- The site shows the matrix on each concept page (a row per language), on each language page (a column per concept area), and as the comparison pages (16.5).

## 4.9 The computed ladders

The build computes these ladders. No person and no agent writes them (T05). They replace the H, L, I ladders of algorithm with the evidence that programming knowledge needs.

**W. Writing** (06.8): W0 stub to W5 reviewed, as in algorithm, with the section lists of 06.3.

**G. Grounding** (for every node)

| Level | Meaning |
|---|---|
| G0 | no sources |
| G1 | at least one source with a locator for each section that states facts |
| G2 | every normative claim (behaviour, default, limit, version, complexity, history) carries a citation to a primary source: a standard clause, an official reference anchor, a proposal or the release notes (06.6) |
| G3 | G2, and the citation check passed: each cited anchor exists in the raw snapshot of its source, and the claim's key terms appear there (13.6); and a person reviewed the claims of this exact hash |

**X. Execution** (for each language section and each implementation)

| Level | Meaning |
|---|---|
| X0 | no code |
| X1 | each `run` block runs in its pinned image and ends as declared (exit 0, or the declared error) |
| X2 | X1, and each output block matches, each task implementation passes the shared cases, and format and lint are clean |
| X3 | X2, and the examples run on every version of the language's support window that the page claims (07.6), and each version statement on the page is a K2 compat entry |
| X4 | X3, and a person reviewed this exact file hash |

**K. Compat** (for features, api cards and diagnostics)

| Level | Meaning |
|---|---|
| K0 | no compat entry |
| K1 | compat entries with sources for each implementation column of the language (09.2) |
| K2 | K1, and each entry with `version_added` has a probe that fails on the release before and passes from the release named, in the matrix of images (09.4) |
| K3 | K2, and a person reviewed the entries |

**F. Formal** (only for nodes with `formal: expected` or `optional`, 11): F0 none; F1 the Lean definitions exist and a person reviewed that they say what the page says; F2 the theorems of the page are proved with no `sorry` and only the standard axioms; F3 F2, and the Lean interpreter of the calculus agrees with a real language on the generated programs of 11.6.

**V. Visualization** (12): V0 none; V1 static figures; V2 an interactive scene or a recorded trace with step back and code sync that passes the checks; V3 reviewed.

A node is **complete** when it has W4, G2, X3 in each tier-1 language where its matrix cell needs an example, X2 in each tier-2 language where the cell needs one, K2 (features), V2 (`viz: expected`), F2 (`formal: expected`) and the reviews that its kind needs (15.2). A card is complete when its facts pass the card rules (15.3 rule family R) and its sources resolve.

A review counts only for the hash that it reviewed, as in algorithm and math2.

## 4.10 Prohibitions

- No status, level, badge or popularity number written by hand (T05).
- No code block without a run mode. Each fenced code block is `run`, `fails`, `norun` with a reason, or `output` (07.2, rule X01).
- No version claim in prose ("since 3.10", "added in Java 16", "deprecated in C++17") that is not a compat entry (K04).
- No inline API name in a language section that does not resolve to an `api` card of that language and version, unless it is marked as user code (X09).
- No copied text or code from a source whose licence forbids it (13.5, rule C01).
- No example that needs the network at run time (07.5).
- No example that depends on a package outside the standard library, unless the page is about that package and pins it (07.7).
