# 18. Coverage: what "no missing" means

The goal says "ALL human knowledge related to programming, NO MISSING". This document turns that into something the build can measure. The method is algorithm's (2150/16): a universe of every item that any trusted source names, an outcome for each, coverage levels, and a report that shows the gaps. What is new is the scale (millions of rows, most of them cards), the language dimension (the matrix), and the version dimension (the support window).

## 18.1 The claim

The site may say "complete" for a scope when, and only when:

1. every row of every source in the scope's source list (18.4) has an outcome in the universe register (no `pending`);
2. every row with the outcome `node` has a node at the level that the scope requires (18.2);
3. every row with the outcome `card` has a card that passes the R rules (15.3);
4. every matrix cell of the scope is decided (4.8);
5. every feature in the release notes of every release in the support window has a compat entry (K06).

The claim is per scope ("the Python language guide is complete at C4", "the core is complete at C5"), never for the whole site in general, and it is computed, with the date and the source versions. The home page shows the computed claims.

## 18.2 Coverage levels

| Level | Name | Meaning |
|---|---|---|
| C0 | registered | the item is in the universe with outcome `node` or `card` |
| C1 | card | a card exists and passes the R rules |
| C2 | page | a node at W2 and G1 |
| C3 | grounded and executed | W3, G2, X2 in every tier-1 language where its cells need an example, E11 for tier 1 |
| C4 | versioned | C3, and K2 for its compat entries, and X3 |
| C5 | complete | C4, and V2 or F2 where the kind expects them, and W4; tier-2 cells decided and X2 |
| C6 | reviewed | C5, and the reviews of 15.2 for W5, X4 and K3 |

Algorithm's C0 to C6 measured writing and proof. Programming's levels measure grounding, execution and versions, because those are the evidence of programming knowledge.

## 18.3 The universe register

`registers/universe/<source>.jsonl`, one file per source, one row per source item. Each row has the source id, the source's name for the item, the snapshot hash, and an **outcome**:

| Outcome | Meaning |
|---|---|
| `node` | the item gets a node (target) or has one; with the tag |
| `card` | the item is a card; with the tag |
| `alias` | another name of an item; with the tag |
| `section` | a section of a node; with the tag and the anchor |
| `merged` | merged into another item; with the tag and a reason |
| `disambiguation` | the name has several meanings; with the tags |
| `elsewhere:algorithm`, `elsewhere:math2` | belongs on the other site; with the target tag (5.8) |
| `out-of-scope` | not programming knowledge; with a reason code (`not-programming`, `product-marketing`, `personal`, `too-small` with a rule) |
| `not-an-item` | the row is not a thing (a list page, a navigation entry, a typo) |
| `pending` | not decided yet |

Rules:

- An outcome other than `pending` needs a reason or a target.
- `out-of-scope` with `too-small` needs a rule from `registers/scope-rules.yaml`, for example: "a library with fewer than N dependents in its registry and no mention in a canon source is not carded" (18.5).
- The curator decides outcomes in bulk with the matcher's scores; a person approves the rules, not each row.

## 18.4 Sources of the universe and the scopes

Each scope lists the sources whose rows it must cover. The counts are from 01.2, Nov 2026.

| Scope | Sources | Rows (approx.) | Typical outcome |
|---|---|---|---|
| languages | PLDB (5,167 entries, 3,471 programming languages; public domain), Wikidata (1,916 items under Q9143), Linguist (845 entries, 569 programming), Esolang wiki (7,899 languages), Rosetta Code language list (about 1,023), HOPL (8,945, names and facts by hand only, no dump) | about 12,000 after merging (Estimate) | `card` for most; `node` for tiers 1 to 3 |
| releases | release feeds of every tier 1 to 3 language, tool and core library; endoflife.date (489 products) | about 30,000 | `card`; `node` for tier-1 and tier-2 releases in the window |
| concepts | CS2023 knowledge units (SDF, FPL, SE, PDC, SEC, SF), SWEBOK v4 topics, the canon indexes (13.3), Wikipedia's comparison-of-languages series, Van Roy's paradigm taxonomy, TAPL, PFPL, Mitchell, Exercism concepts (882 concepts in 24 tracks with a syllabus) | about 5,000 rows, about 2,500 items after merging (Estimate) | `node` |
| features | the table of contents of each tier-1 and tier-2 language's reference (8.3) | about 40,000 rows (Estimate) | `node` or `section` |
| api | the stdlib indexes (Sphinx inventories, Javadoc, rustdoc JSON, Go's package index, .NET API docs, MDN Web API pages, POSIX, the C and C++ library clause lists) | about 400,000 (Estimate) | `card` |
| diagnostics | rustc error codes, TypeScript diagnostics, the C# compiler errors, the Java compiler keys, Python's exception hierarchy, Go vet checks, the GCC and Clang warning flags, the main linters' rule lists (Ruff, ESLint, Clippy, golangci-lint linters, Checkstyle, PMD, SpotBugs, RuboCop, ShellCheck, and others) | about 30,000 (Estimate) | `card` |
| proposals | PEPs, JEPs, Rust RFCs, TC39 proposals, WG21 papers that shipped, WG14 papers that shipped, Swift Evolution, Kotlin KEEPs, C# proposals, Go proposals that shipped, PHP RFCs, Ruby feature tickets that shipped | about 15,000 (Estimate) | `card` |
| tasks | 10.1 | about 3,000 rows, 1,500 tasks | `node` |
| patterns and practices | GoF, POSA, PoEAA, EIP, the cloud design patterns catalogues, Fowler's refactorings and smells, the laws of software, the agile and process vocabularies, Wikipedia's lists of patterns and anti-patterns | about 2,000 rows, 1,200 items (Estimate) | `node` |
| weaknesses | CWE 5.0 (956 active weaknesses, 396 groupings, 58 views; "categories" are called "groupings" since 5.0), CAPEC 3.9 (559 attack patterns, as `see_also` targets), OWASP Top 10:2025, the CWE Top 25 (2025), OWASP ASVS 5.0 | about 1,500 | `card`, `node` for the core |
| tools | the tool rows of the developer surveys, the package managers and build tools of every tier 1 to 3 language, the linters, the debuggers and profilers, Wikidata software items with a "programming tool" class | about 5,000 (Estimate) | `card`, `node` for the core |
| libraries | the registries (npm, PyPI, crates.io, Maven Central, NuGet, RubyGems, Packagist, Go modules, Hex, CRAN, pub.dev, Hackage, opam), filtered by the scope rule of 18.5 | about 20,000 after the rule (Estimate) | `card` |
| formats and standards | IANA media types, the RFCs of the IETF's application and transport areas that programmers use, ISO/IEC JTC 1/SC 22 standards, ECMA standards, W3C and WHATWG specifications, the Unicode standard and its annexes | about 3,000 (Estimate) | `card`, `node` for the core |
| licences | SPDX licence list and exceptions | about 700 | `card` |
| people, organizations, works | Wikidata items linked from our nodes, the ACM Turing Award list, the canon of 13.3, the standards bodies and foundations | about 10,000 (Estimate) | `card`, `node` for the canon |
| terms | FOLDOC (names only), the Jargon File (names only), ISO/IEC/IEEE 24765 SEVOCAB terms (names only), the glossaries of the tier-1 references | about 20,000 rows | `node` for definitions in the core; `alias` or `out-of-scope` for the rest |

## 18.5 The scope rule for the long tail

There are millions of libraries and billions of source files (Software Heritage holds 31.1 billion unique files, 01.2). "No missing" cannot mean a page for each. The rule:

- **Libraries** get a card when they are a dependency of at least N other packages in their registry (Default: N = 100 for npm and PyPI, 50 for other registries, adjusted per registry so that each registry yields at most about 5,000 cards), or when a canon source, a tier-1 language's official docs, or one of our nodes names them.
- **Tools** get a card when a developer survey lists them, when a tier-1 to tier-3 language's official docs name them, or when a node names them.
- **Languages** get a card when any source of 18.4 lists them. There is no size threshold for languages: every language that a trusted catalogue names is in.
- Everything else from a registry has the outcome `out-of-scope` with `too-small`. The rule is published, so that the claim "no missing" is exact: "every library with at least 100 dependents on npm is carded".

## 18.6 Demand signals

Coverage is about sources. Demand is about readers. The two meet in the queue (14.8):

- Stack Overflow tag counts (the crosswalk, 5.7): a tag with many questions and no matching item is a gap of high priority.
- The developer surveys (01.2): languages, tools and frameworks in the top lists are the core.
- Our search queries with no result (17.9).
- Reader reports.

## 18.7 The coverage report

`prog report coverage` writes `/about/coverage/`:

- per scope: rows, outcomes, the level histogram (C0 to C6), the share at the milestone target;
- per language: the template sub-areas, the features against the reference's table of contents, the matrix row, the compat K levels, the task count;
- per concept: the matrix column;
- the gap list with links to the queue;
- the history (one line per day), so that growth and decay are visible.

As in algorithm, the report is computed from the repository, never typed.

## 18.8 Scale estimate

| Item class | Count at M5 (Estimate) |
|---|---|
| nodes | about 60,000 (concepts 2,500; features 25,000 for tiers 1 and 2; tasks 1,500; patterns and practices 1,200; releases 1,500; tools and libraries 3,000; weaknesses 300; definitions 3,000; the rest: language-tree pages, formats, standards, people, works, comparisons) |
| cards | about 550,000 (api 400,000; diagnostics 30,000; libraries 20,000; proposals 15,000; releases 30,000; languages 12,000; the rest) |
| executed blocks | about 400,000 |
| compat entries | about 100,000 |
| runner images | about 400 (language versions and implementations) |

For comparison: algorithm has 8,396 nodes after its first week (01.6), DevDocs aggregates 247 documentation sets, and MDN has tens of thousands of pages. These counts make programming the largest of the three sites by an order of magnitude, which drives the hosting choice (16.20) and the card model (04.1).
