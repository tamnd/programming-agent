# 02. The goal

## 2.1 One sentence

**Every idea, language, feature, task, pattern, practice, tool and fact of programming, in one place, with every example executed for the versions it claims, every fact cited, every language compared side by side, and every gap measured, at https://programming.tamnd.com.**

## 2.2 What counts as knowledge about programming

The user's goal is "ALL human knowledge related to programming, NO MISSING". The scope is everything a programmer knows or looks up, from the first program to language design, in seven groups:

| Group | What | Examples | Spec |
|---|---|---|---|
| Ideas | concepts, definitions, paradigms, the theory of programming languages | closures, ownership, monads, variance, tail calls, garbage collection, type inference | 05.2 group A, and B for implementation |
| Languages | every language as a catalogue entry, and the main ones as full guides with every feature and release | Python and its `match` statement; Rust 1.99; COBOL as a card; Brainfuck as a card | 08 |
| Doing | tasks with implementations in many languages, idioms, techniques, refactorings | read a file line by line; parse JSON; extract function | 10 |
| Building software | design, patterns, architecture, testing, debugging, performance, security, process and practice; data, systems, web and cloud programming; the domains | observer; hexagonal architecture; property-based testing; CWE-89 | 05.2 groups C, D and E |
| Tools and ecosystems | compilers, build tools, package managers, version control, editors, debuggers, CI, the core libraries and frameworks | Git; Cargo; GDB; React as a core framework | 05.2 `tools-and-workflow`, and the tool sub-areas of each language tree (5.3) |
| Formats and standards | data formats, encodings, protocols programmers use, the language standards themselves | JSON; UTF-8; ISO C23; ECMA-262 | 05.2 `foundations`, `data-and-persistence`, `web-development`; standards on language pages |
| Context | history, people, organizations, works, licences, the profession | HOPL; Grace Hopper; the WG21 committee; MIT licence | 05.2 group F |

Algorithms and data structures belong to `algorithm.tamnd.com`, and mathematics to math2; the border is 5.8. Programming links to them with `al:` and `m2:` links and does not repeat them.

"No missing" has a precise meaning (18.1): every item that any trusted source names has an outcome (a node, a card, an alias, a section, a merge, or out of scope with a published rule), and the site shows the gaps. It does not mean a long page for each of four million npm packages; it means that the rule that decides which ones get a card is published, and every package that the rule selects has one (18.5).

## 2.3 The quality promise

A page on this site says what it claims and proves it with evidence a reader can check:

1. **Every example runs.** Every code block is executed by the build in a pinned image for each version it claims, and its output is the recorded output (07). A block that is not run says why (7.8).
2. **Every version claim is computed.** "Since Python 3.10" comes from a compat entry with a probe that passes on 3.10 and fails on 3.9 (09.4), never from prose.
3. **Every normative claim cites its source** (the language reference, the standard clause, the PEP) at the version it is about, and the citation is checked against a stored snapshot (13.6).
4. **Every cell is decided.** For every concept and every tier-1 and tier-2 language, the page says whether the language has it natively, in the stdlib, in a library, as an idiom, or not at all (04.8).
5. **Formal claims are labelled.** "Proved" means a Lean theorem; models are not confused with languages (11.1).
6. **Agents write, checks gate, people review.** Every page shows what an agent wrote, what the checks verified, and what a person reviewed (15.5, Q-17).
7. **Nothing stale.** Pages do not say "latest" or "currently"; the facts box and the status labels are computed from the release register (6.7, 9.5).

## 2.4 The scale

| Item | At M5 (Estimate, 18.8) |
|---|---|
| nodes (full pages) | about 60,000 |
| cards (catalogue entries) | about 550,000 |
| languages with a card | about 12,000 |
| languages with a full guide (tier 1, D4) | 15 |
| languages with a guide (tier 2, D3) | about 34 |
| executed code blocks | about 400,000 |
| compat entries | about 100,000 |
| runner images | about 400 |

This is about ten times the size of algorithm in nodes and a hundred times in catalogue entries, which drives the card model (04.1), the hosting (16.20) and the search design (17).

## 2.5 Who it is for

| Reader | Needs | Served by |
|---|---|---|
| learner | a path from zero, in one language, with runnable examples | paths (16.12), level `intro` (5.6), the playground (7.10) |
| working programmer | "how do I X in Y", "what does this error mean", "since which version" | tasks (10), diagnostic cards, compat tables, search with error paste (17.6) |
| polyglot | "what is the Rust for this Python idea" | concept pages with the matrix, comparisons (8.8) |
| language designer and researcher | the design space, history, semantics | concept pages, proposals, formal developments (11) |
| teacher | course material with exercises | exercises (10.7), paths |
| agent | correct, versioned, machine-readable facts | `.md` and `.json` twins, datasets, MCP (16.10) |

## 2.6 Non-goals

- **Not a Q&A site.** No questions, answers or votes. Reader reports fix pages (16.19).
- **Not a library documentation host.** Library cards link to their official docs; only core frameworks get nodes (Q-18).
- **Not a code host or package index.** Code lives in examples and tasks only.
- **Not a news site.** Releases are content with a fixed form (9.7), not news.
- **No rankings of quality.** Popularity is shown as computed data with sources and dates (Q-12), never as "best".
- **No server execution in v1.** Readers run code in the browser only (Q-07).
- **No copying.** Text from sources we may not adapt is never copied (19).
- **English only until M5** (Q-14).

## 2.7 Relation to other tamnd projects

| Project | Relation |
|---|---|
| algorithm (2150) | the model for everything: pipeline, packages, tags, ladders, hosting; owns algorithms and data structures (5.8) |
| math2 (2148) | owns mathematics; the source of `markup`, `copycheck`, `prov` and the MCP pattern |
| the `*-cli` tools in the author's workspace | wrapped by harvest adapters (13.2) |
