# 10. Tasks and testing

A **task** is a common programming job with a contract, solved in many languages: "read a file line by line", "parse JSON into a typed record", "send an HTTP GET with a timeout", "format a date in ISO 8601", "reverse a string by grapheme clusters". Tasks are the phrasebook of the site, the answer to "how do I X in L". This document defines the task universe, the contract, the implementations and their tests. It reuses algorithm's problem machinery (2150/09: contracts, cases, checkers, the differential and mutation tests) and says what is new.

## 10.1 The task universe

| Source | What we take | Licence and rule |
|---|---|---|
| Rosetta Code | task **names** and the one-line idea, to find gaps (1,363 tasks and 423 drafts, about Sept 2026, 01.3) | GFDL: names and facts only, no task text, no code (19.3) |
| Programming-Idioms | idiom names and ids, as a crosswalk | CC BY-SA: names only by default; text adapted with credit only when needed |
| Go by Example (85 pages), Rust by Example, Kotlin by Example, Swift by Sundell topics, Python Module of the Week | titles, as a check of what a language's users need | names only |
| Exercism problem-specifications | the canonical-data.json cases of exercises (MIT) | the cases may be reused with credit, as test data for tasks that match |
| the Stack Overflow "how do I" questions of the top tags | the questions with the most votes per tag, as a demand signal (CC BY-SA: titles as facts only) | names only |
| the official tutorials and cookbooks of tier-1 languages | the jobs that they teach | names only |
| our own list | jobs that the matrix and the concept pages need | ours |

The curator merges these into `registers/tasks.yaml` with groups (10.6). Estimate: 1,500 tasks in the universe, 300 in the core (5.4).

**Task sets by tier**:

| Set | Size | Languages that must implement it |
|---|---|---|
| core tasks | about 300 | tier 1 and tier 2 |
| tier-3 task set | 20: hello world, FizzBuzz, read a line from stdin, sum a list of numbers, reverse a string, read a file, write a file, command-line arguments, a record type, a map, sort a list, string formatting, error handling, a loop with an index, a recursive function, a higher-order function, a closure, string to number, current time, exit code | tier 3 |
| full set | about 1,500 | tier 1 at M4; tier 2 where the task applies |

A task that does not apply to a language (for example "free a block of memory" in a garbage-collected language) has a matrix-style decision with a reason, as for concepts (4.8).

## 10.2 The contract

Three contract types:

| Type | The implementation is | The harness checks | Example |
|---|---|---|---|
| `function` | a function with a given signature, in the language's types | calls it with the cases; compares the result with the checker | "reverse a string by grapheme clusters": `reverse(s: string) -> string` |
| `program` | a complete program | runs it with `args`, `stdin` and fixture files; compares stdout, stderr, the exit status and the output files | "count the words in a file given as an argument" |
| `io` | a program that talks to a local service the harness starts | runs it against the fixture service; checks the requests the service saw and the program's output | "send an HTTP GET with a 2 second timeout and print the status" |

`contract.yaml`:

```yaml
type: function
signature:
  name: reverse_graphemes
  params: [{name: s, type: string}]
  returns: string
types:                 # per-language type mapping; the harness writes the adapters
  python: {string: str}
  rust: {string: "&str -> String"}
  java: {string: String}
checker: exact          # exact | unordered | float:1e-9 | custom:checker.py
cases: [cases/basic.jsonl, cases/unicode.jsonl]
limits: {time_s: 2, memory_mb: 256}
notes: >-
  A grapheme cluster is as defined by UAX #29 extended grapheme clusters. Languages whose standard library has no grapheme segmentation may use the library named in matrix.yaml.
```

`cases/*.jsonl` has one case per line: `{"id": "zwj-family", "in": {"s": "a👨‍👩‍👧b"}, "out": "b👨‍👩‍👧a", "why": "a ZWJ sequence is one grapheme"}`. Every case has a `why`. The edge cases are the content: the task page renders them as the "Edge cases" section (6.3).

## 10.3 Implementations and their tests

- `impl/<lang>/` holds the implementation and its test file, as in algorithm 2150/09. The test file is generated from the cases by `prog gen tests` with the language's test framework (from `registers/languages.yaml`), and may have extra hand-written tests.
- The implementation must be idiomatic for its language at the language's ceiling version, follow `style.yaml` (8.5), pass the formatter and linter, and use only the standard library unless the task's matrix cell is `library` (7.7).
- A task page shows the idiomatic implementation first. Other ways (older versions, a library, a lower-level way) are in `impl/<lang>/alt-<name>.<ext>` and appear under "Alternatives", each with its own test.
- An implementation that needs a newer version than the window floor declares `since` in its header comment, and the matrix run (X3) checks that.
- `prog test <task>` runs every implementation in its image (07.4). `-kind diff` runs every language on the same generated inputs and compares the outputs (the differential test of algorithm). `-kind mutate` checks that the cases catch mutants of the reference implementation.

## 10.4 Concept and feature tests

Concept and feature pages do not need contracts, but their examples run (07). Two extra kinds of test come from tasks:

- **Matrix tests**: for each matrix cell that is `native`, `stdlib`, `library` or `idiom`, the example in the concept's language section is a run block, so the cell's evidence is executed.
- **Behaviour tests**: a feature node may have `tests/` with cases for each rule of its Semantics section (6.3, `r1`, `r2`, ...). A rule with a test is shown with a "tested" mark that links to the test. This is the model of the Rust Reference, which links its rule ids to the rustc tests that exercise them (01.3).

## 10.5 Testing as a subject

The `testing` area (05.2) is a subject of its own. It has concept nodes (unit tests, test doubles, property-based testing, fuzzing, mutation testing), technique nodes, and tool and library nodes for every test framework of tiers 1 and 2. Its examples run like every other example: a page about pytest runs pytest in the image, and the output block is pytest's output, normalized.

## 10.6 Task groups and URLs

- Tasks are in the topic tree under the area of their subject (`io/files/read-a-file-line-by-line` is under `systems-programming/files`, or `data-and-persistence/files`; the curator decides at M1 with one rule: a task's area is the area of the concept it mostly uses).
- `/tasks/` lists all tasks by group: text and strings, numbers, collections, dates and times, files and paths, processes and the OS, networking, concurrency, data formats, errors, testing, command-line programs, classes and records, functional idioms, generics, memory, interop, databases, the web, security.
- The URL is the task's global slug: `/read-a-file-line-by-line/`, and `/read-a-file-line-by-line/rust/` for one language's implementation (an anchor-backed view of the same page).

## 10.7 Exercises

Exercises (04.2) are cards that point to practice problems on other sites: Exercism (MIT problem specifications, so the card can carry the problem's name and its case count), LeetCode, HackerRank, Advent of Code, Project Euler, Codewars. A card holds the name, the site, the id, the URL, the difficulty in the site's own scale, the concepts it practices (linked by the curator), and the languages the site supports. It never holds the problem text of a site whose terms forbid it. Algorithm-heavy exercises link to algorithm.tamnd.com problem nodes first.

## 10.8 Security examples

Weakness nodes (04.2) show vulnerable code. That code runs only in the **security harness**: a variant of the sandbox with an exploit input and an oracle.

- `vuln/<lang>/vulnerable.<ext>` and `vuln/<lang>/fixed.<ext>`, with `vuln/exploit.jsonl`: inputs that show the weakness (an SQL injection string against a SQLite fixture, a path traversal against a fixture directory, a buffer overflow under ASan, a ReDoS input with a timeout).
- The harness checks that the exploit **succeeds** against the vulnerable code (the oracle sees the effect: data that should not be read, a sanitizer report, a timeout) and **fails** against the fixed code (the oracle sees no effect, and the normal cases pass).
- The exploits stay inside the sandbox: local fixtures only, no network, no real targets. They are teaching examples of defects, not attack tools: no exploit targets a real product version, and CVE items link to the vendor advisory, not to exploit code (19.7).
