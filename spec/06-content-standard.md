# 06. The content standard

This document says what a page must contain and how it must be written. It extends 2150/06 (algorithm). Where this document is silent, 2150/06 applies.

## 6.1 Principles

1. **Each page is one of four forms, and says which.** We use the Diátaxis split (diataxis.fr, CC BY-SA 4.0): reference, explanation, how-to, tutorial. The kind fixes the form (6.3). A concept page explains. A feature page is reference first, with an explanation section. A task page is a how-to. A learning path is a tutorial. A page does not mix a tutorial into a reference.
2. **Every behaviour claim is checked or cited.** A sentence that says what a program does is shown by an executed example (07), a compat entry (09), or a citation to a primary source with a locator (6.6). A sentence that is none of these is opinion and must read as opinion ("Many style guides prefer ...", with a source).
3. **Version is part of every fact.** A behaviour claim is about a language version range. The page states the range in its compat entries, not in loose prose. When no range is stated, the claim is about every version in the support window of the language (09.3).
4. **One fact, one place.** A rule of Python lives on the Python feature page. The concept page shows it through the `realizes` edge, with a short summary, and does not repeat the rules.
5. **The reader can run it.** Every example is complete, or has a declared prelude (07.3). The reader can copy it, or open it in the playground (16.8), and get the shown output.
6. **Plain English.** ASD-STE100 style, as in algorithm 06.2: short sentences, one idea per sentence, active voice, the same word for the same thing, no idiom, no filler, no marketing words. No em dash and no en dash. No line break inside a sentence. Code names in code spans.
7. **Neutral on languages.** No language is "better". A comparison states differences, with evidence. Benchmarks are allowed only with the benchmark code in the node, run in the CI images, with the hardware stated (07.9). Popularity numbers come from the build (04.5), with the date.
8. **No invented facts.** Never invent a version number, a date, an API name, a flag, an error code, a page number or a section number. If you are not sure, you leave it out and open a question in the node's `notes.md`. The checks of 15 find invented API names (X09), invented versions (K04) and invented locators (G05).

## 6.2 Markup

The markup is the algorithm markup (2150/06.3, the `markup` package): CommonMark with GFM tables, KaTeX math, `[[slug]]` and `[[slug|text]]` links, `[[slug#anchor]]`, callouts (`> [!note]`, `> [!warning]`, `> [!pitfall]`, `> [!history]`, `> [!since v="3.10"]`), and figure blocks. New for programming:

| Markup | Meaning |
|---|---|
| `[[python/match-statement]]` | a link to an owned item: `<owner>/<slug>` |
| `[[python/api/str.split]]` | a link to an api card; renders as code with the link |
| `[[rust/errors/E0382]]` | a link to a diagnostic card |
| `[[python/proposals/pep-634]]` | a link to a proposal card |
| `[[al:merge-sort]]`, `[[m2:XXXXXX]]` | a link to algorithm or math2 (5.8) |
| `` `py:str.split` `` | an inline code span with a language prefix: the build resolves it to an api card and links it (rule X09 checks it) |
| `` `user:parse_config` `` | an inline code span of user code: the build does not resolve it |
| ```` ```python run ```` | an executed example (07.2) |
| `{rule:destructors.scope.nesting}` | a reference to a rule id of the language's reference, when the language has rule ids (the Rust Reference, the FLS) |
| `{std:c23/6.5.7p4}` | a reference to a clause of a standard: the standard slug, the clause, the paragraph |
| `{since:python/match-statement}` | renders the version where the feature was added, from compat (never typed by hand) |

A page body may not contain raw HTML, except the `<details>` element for a hidden prelude. Scenes and playgrounds are inserted by the template, not by HTML in the body.

## 6.3 Sections per kind

A section in **bold** is required at W2. A section in normal type is required at W3. A section in *italic* is optional. Headings are `##`. The renderer adds the title, the summary, the badges and the facts box from the front matter; the body does not repeat them.

**concept** (form: explanation)

1. **Idea**: what it is and why it exists, in 1 to 3 paragraphs, with the smallest example in one language (the "anchor language" of the concept, the first tier-1 language where the concept is native, by the order of `registers/languages.yaml`).
2. **How it works**: the model behind it (the environment model for closures, the reachability model for GC), with a figure or a scene when `viz: expected`.
3. **In languages**: one subsection per language with a matrix cell (4.8), rendered from `examples/<lang>.md`. Each subsection: the verdict (`native`, `stdlib`, ...), the example with its output, the differences in 1 to 4 bullet points, the link to the feature page.
4. Variations and related ideas: the named variants and the differences between them.
5. Pitfalls: the common mistakes, each with a `fails` example or a wrong-output example and the fix.
6. Cost: what it costs in time and memory in the main implementations, with sources (6.6).
7. History: who and when, with sources.
8. *Theory*: the formal view, linked to the Lean part (11) when there is one.
9. *Exercises*: links to exercise cards and to tasks.
10. **References**: from the front matter `sources`, rendered by the build; the writer adds no other list.

**feature** (form: reference, then explanation)

1. **Synopsis**: the syntax or the signature, in the language's grammar notation or as code, with every form.
2. **Semantics**: what it does, rule by rule. Each rule is a numbered paragraph with an id (`r1`, `r2`, ...) and a citation (6.6). A rule that has a probe or an example links it.
3. **Examples**: executed examples, from simple to complete.
4. Versions: rendered from compat (09). The writer writes the migration notes only.
5. Pitfalls: as for concepts.
6. Comparison: how the same concept looks in 3 to 5 other languages, from the concept page (the build takes it from there; the writer writes no new examples here).
7. *Implementation notes*: how the main implementations do it, with sources.
8. **References**.

**language** (form: reference and explanation, the guide's front page)

1. **Overview**: what the language is for, the design goals, in the words of its designers with sources.
2. **A first program**: hello world and one characteristic program, both executed.
3. **The language at a glance**: the type system, the memory model, the execution model, the paradigms, in a table built from the front matter plus one paragraph each.
4. **The guide**: the table of the 25 sub-areas (5.3) with the feature counts and the coverage bar (computed).
5. Versions and support: the releases table and the support window, computed from `releases.yaml` (09.2).
6. Implementations: from `implementations.yaml`.
7. Ecosystem: the package manager, the build tools, the main libraries by need, as cards.
8. History and influence: the genealogy graph (12.6), the history with sources.
9. *Coming from another language*: generated links to the comparison pages (16.5).
10. **References**: the specification, the reference, the main books.

**release**: **What is new** (the list of feature changes, generated from compat, each with one sentence written by the writer), **Removed and deprecated** (generated), Migration (written), Support dates (computed), **References** (the release notes).

**task** (form: how-to)

1. **Task**: the problem in 1 to 3 sentences, and the contract (10.2) rendered.
2. **Solutions**: one tab per language, each an implementation that passes the cases, with notes on the choices it made (the idiomatic way first; other ways as `### Alternatives`).
3. Edge cases: the cases of the contract, explained, and how each language handles them.
4. Pitfalls.
5. *Related tasks*.
6. **References**.

**pattern** (form: explanation; follows the classic pattern form, without copying any book's text)

1. **Intent**: one or two sentences.
2. **Problem**: the forces, in plain words.
3. **Solution**: the structure, as a diagram (a generated class or sequence diagram, 12.5), and the participants.
4. **Example**: executed in at least 3 tier-1 languages from different paradigms at W3, and all tier-1 languages where the pattern applies at W4; languages where the pattern is built in show the built-in form instead.
5. Consequences: the trade-offs.
6. When not to use it: (for positive patterns) the costs and the anti-pattern it becomes; for negative patterns this section is **How to fix it**, with the refactorings that apply.
7. Variants and related patterns.
8. *Known uses*: open-source projects with links to a commit-pinned file.
9. **References**.

**technique** (form: how-to): **Purpose**, **Steps** (numbered), **Example** (before and after code, executed, with tests that pass before and after for a refactoring), When to use, Pitfalls, Tools that do it (as cards), **References**.

**practice** (form: explanation): **Statement** (the principle, the law or the method in one paragraph, quoted only when the source licence allows and with the source), **Meaning** (what it asks the programmer to do), Examples (code when it applies), Criticism and limits (with sources), History, **References**.

**tool** and **library** (form: reference): **What it is**, **Install and first use** (executed in the image when the tool is in the image list, 07.4), Concepts of the tool, Common tasks (each with a command or code), Configuration, Alternatives (computed from `alternative_to`), History and versions, **References**. Git is the model: it has a node tree of its own under `tools-and-workflow/version-control/git/`, with one node per command group and per concept (the object model, refs, the index, merging, rebasing).

**format** and **standard** (form: reference): **What it is**, **Syntax or structure** (with the grammar when it has one), **Examples** (parse and emit in at least 3 tier-1 languages, executed), Versions, Pitfalls (for formats: the security and interop traps), Tools and libraries, **References** (the specification with clauses).

**weakness** (form: explanation with a how-to): **Description** (in our words; the CWE text is linked, not copied), **Vulnerable code** (one `norun` or sandboxed example per relevant tier-1 language, 07.8), **Fixed code** (executed), How to detect it (tools and rules, as cards), How to prevent it (practices and features, as `prevents` edges), Real cases (CVE ids as cards, with dates and sources), **References**.

**definition**, **person**, **organization**, **work**: as in algorithm, with **Definition** / **Biography** / **About** / **Summary** and **References**. A work page summarizes; it never reproduces.

## 6.4 Examples in prose

Rules for the code that a page shows, in addition to 07:

- Every example is the smallest complete program that shows the point. It has a comment only where the reader would ask "why?".
- The example follows the language's style rules (08.4) and passes the formatter and the linter of `registers/languages.yaml`, with the lint exceptions of the node listed in the block attributes.
- Names are meaningful and neutral: `order`, `total`, `user`, `path`. No `foo`, `bar` or jokes, except on the pages about those names.
- The example shows its output in an `output` block, captured from the run, never typed by hand (07.3).
- An example that shows wrong code is marked `fails` (it must fail as declared) or `wrong` (it runs and prints the wrong answer, which the next block explains). The renderer marks both with a red edge and an icon, so that no reader copies them by mistake.

## 6.5 Words and names

- `registers/terms.yaml` lists the preferred term for each idea and the terms to avoid. For example: "function" (not "method") unless the language calls it a method; "error" for the general idea and the language's own word for the language's construct ("exception" in Java, "panic" in Go and Rust, "error" for `Result::Err`). The linter of 15 checks it (rule S05).
- On a language page, use the language's own words. The language template's `style.yaml` (08.4) maps the general term to the language's term, for example `concept: closures → python: "closure", "nested function"; java: "lambda", "anonymous class"`.
- Names of languages, tools and companies are written as their owners write them: "JavaScript", "TypeScript", "C#", "Node.js", "GitHub", "macOS". The register `registers/names.yaml` lists them, and rule S06 checks them.
- Trademarks: nominative use only (19.6). No logos.

## 6.6 Grounding and citations

The G ladder (04.9) needs citations at the claim level. A citation is `[^n]` in the text, with the footnote rendered from the `sources` list of the front matter, or an inline `{std:...}` or `{rule:...}` reference.

A citation is **primary** when it is one of:

| Primary source | Locator form | Example |
|---|---|---|
| an ISO, IEC, ECMA, IEEE, IETF or W3C standard, or a final committee draft | clause and paragraph | `{std:c23/6.5.7p4}`, `{std:ecma-262-2026/13.15.2}`, `{std:rfc9110/9.3.1}` |
| the official reference of the language | a stable anchor URL | `https://docs.python.org/3.14/reference/compound_stmts.html#the-match-statement` |
| a rule id of a reference that has them | the rule id | `{rule:rust/destructors.scope.nesting.function-body}` |
| the language's change proposal | id and section | `pep-634 §Semantics`, `jep-444`, `p2300r10 §4.2` |
| the release notes of the version | an anchor URL | "What's New in Python 3.14" with the section anchor |
| the source code of the reference implementation | a commit-pinned file and line range | `github.com/python/cpython/blob/<sha>/Objects/listobject.c#L120-L160` |
| the compiler's or tool's own documentation | a stable anchor URL | the GCC manual's option index |

A citation is **secondary** when it is a book, a paper, a talk, an article or a Q&A answer. Secondary sources are allowed for history, for explanation, for criticism and for practice, but a normative claim (behaviour, default, limit, version) needs a primary source for G2.

Rules:

- **G01**: every source has a locator: a section, a clause, a page, an anchor, or a pinned line range. "See the docs" is not a locator.
- **G02**: every footnote is used, and every use resolves.
- **G03**: no citation of a page that is not a source (an SEO article, a content farm, a page of AI-written text). `registers/sources-deny.yaml` lists the domains that may not be cited.
- **G04**: a quote is at most 50 words, marked as a quote, with the source and a licence that permits it (13.5).
- **G05**: no invented locator. The citation check (13.6) confirms that the anchor or the clause exists in the raw snapshot of the source for the cited version.
- **G06**: Stack Overflow answers, blog posts and issue threads may be cited only for "people report ..." or for history, never as the proof of a behaviour.

## 6.7 Version language

These rules turn version facts into data, so that the build can check them (K04).

- Never write "since Python 3.10" in prose. Write `{since:python/match-statement}`, which the build fills from compat. Or use the callout `> [!since v="3.10"]`, which the build checks against compat.
- "Modern", "new", "recent", "latest", "currently" and "nowadays" are banned words in feature and concept pages (rule S07), because they rot. The facts box shows "verified on <date> against <versions>" instead (16.4).
- A deprecated or removed feature is shown with its status from compat. The page does not say "do not use this" without the deprecation source.
- When the behaviour differs between versions, the page shows the versions side by side with one example per behaviour, each with `version=` (07.6).

## 6.8 The W ladder

As in algorithm 06.9, adapted:

| Level | Name | Meaning |
|---|---|---|
| W0 | stub | the directory exists with front matter; the page is not published (the card shows instead) |
| W1 | sketch | the summary and the first required section exist |
| W2 | draft | all bold sections of the kind exist (6.3); the format and the links pass; G1 |
| W3 | complete | all bold and normal sections exist; the matrix cells for tier 1 are decided (E11); X2 in the languages it shows; G2 |
| W4 | checked | W3, and the critic's audit passed with no blocker and no major issue (15.4), on this hash |
| W5 | reviewed | W4, and a person with the right role reviewed this exact hash (15.2) |

The build computes the level. A person or an agent never writes it.

## 6.9 Length

| Kind | Target words at W3 (excluding code) | Hard limit |
|---|---|---|
| concept | 800 to 2,500, plus 150 to 400 per language section | 8,000 |
| feature | 400 to 1,500 | 5,000 |
| language | 1,000 to 2,500 | 4,000 (the guide pages hold the rest) |
| task | 200 to 600, plus the implementations | 2,000 |
| pattern | 600 to 1,500 | 4,000 |
| technique, practice | 300 to 1,200 | 3,000 |
| tool, library, format, standard | 400 to 2,000 | 5,000 |
| weakness | 500 to 1,500 | 4,000 |
| definition | 50 to 400 | 1,000 |

A page over its hard limit fails rule S08. The fix is to split it into nodes, not to cut facts.

## 6.10 What a page must not do

- Must not hold a version claim in prose (K04), an unresolved API name (X09), or an example with no run mode (X01).
- Must not copy prose from any source except a short quote under G04. Must not paraphrase one source section by section ("close paraphrase" fails the critic, 15.4).
- Must not explain an algorithm that algorithm.tamnd.com covers. It links it.
- Must not give a benchmark result without the benchmark (07.9).
- Must not say "best", "fastest", "modern", "powerful", "simple", "easy", "just", "obviously", "of course", "note that", "it is worth noting", "in conclusion", "delve", "leverage", "robust", "seamless", "cutting-edge" (rule S02, `registers/banned-words.yaml`).
- Must not link to a page that does not exist (except `[[NEW-<slug>]]` in a draft, which blocks W3).
