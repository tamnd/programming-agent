# Spec 2150p: programming

**All human knowledge about programming in one place.** Each concept, language, feature, release, task, pattern, practice, tool, format, standard and weakness has a page or a catalogue card. Every example is executed for the versions it claims. Every version fact comes from a probe. Every normative claim cites its source. Every concept is compared across every main language. Everything is published at https://programming.tamnd.com.

- Repositories: `tamnd/programming` (content) and `tamnd/programming-agent` (the Go toolchain and the agents, binary `prog`). This follows algorithm and algorithm-agent (spec 2150), which followed math2 (2148).
- Date: 2026-11-09. Status: **Draft**, research complete. Nothing is created or deployed yet.

## Documents

| # | Document | What it settles |
|---|---|---|
| 00 | README | this index, the settled decisions, the conventions |
| 01 | [research](01-research.md) | the field in November 2026: catalogues and counts, reference sites, execution services, runtimes, versions, formal semantics, algorithm's lessons, AI and docs |
| 02 | [the goal](02-the-goal.md) | the one sentence, what counts, the quality promise, the scale (about 60,000 nodes, 550,000 cards), readers, non-goals |
| 03 | [architecture](03-architecture.md) | the two repositories, the stores, the packages of `prog`, the seven flows, deployment |
| 04 | [object model](04-object-model.md) | nodes and cards, the 22 kinds, tags and slugs, front matter, edges, the concept × language matrix, the ladders W G X K F V |
| 05 | [taxonomy](05-taxonomy.md) | the topic tree (36 areas in 6 groups), the language template (25 sub-areas), the core list, levels, crosswalks, the border with algorithm and math2 |
| 06 | [content standard](06-content-standard.md) | Diátaxis forms, markup, sections per kind, citations, version language, the W ladder |
| 07 | [examples and execution](07-examples-and-execution.md) | block attributes, the runner and images, the version matrix, normalization, the X ladder, the in-browser playground |
| 08 | [languages](08-languages.md) | language classes, tiers 1 to 4, depth bands D1 to D4, the register, style per language, domain languages, comparisons |
| 09 | [compat and versions](09-compat-and-versions.md) | compat entries with implementation columns, probes, the support window, computed status, imports, releases as content |
| 10 | [tasks and testing](10-tasks-and-testing.md) | the task universe, contracts, implementations, behaviour tests, exercises, the security harness |
| 11 | [formal semantics](11-formal-semantics.md) | evidence labels, the `Prog` Lean package, notation, grammars, differential tests |
| 12 | [visualization](12-visualization.md) | scenes, recorded traces, generated diagrams, genealogy, the V ladder |
| 13 | [sources and harvest](13-sources-and-harvest.md) | the adapters, the canon, raw snapshots, licence classes, the citation check |
| 14 | [agents and pipeline](14-agents-and-pipeline.md) | roles, the pipeline, harvest-first catalog, models, the writer ladder, gates, queue, budget |
| 15 | [quality and audit](15-quality-and-audit.md) | the threat model, reviews, rule families T E S G X K F V C R L, the audit and verifier, metrics |
| 16 | [site](16-site.md) | URLs, templates, comparisons, the playground, the agent surface (.md, llms.txt, datasets, MCP), paths, hosting and deploy |
| 17 | [search](17-search.md) | query types, the sharded index, the palette, error paste, ranking, evaluation, gaps |
| 18 | [coverage](18-coverage.md) | what "no missing" means: the universe register, outcomes, scopes, the scope rule, C0 to C6, the report |
| 19 | [legal and licensing](19-legal-and-licensing.md) | licence zones, the source table, AI labels, trademarks, takedown, security policy |
| 20 | [repository layout](20-repository-layout.md) | the trees of both repositories, copied packages, the `prog` command tree, workflows |
| 21 | [milestones](21-milestones.md) | M0 (2026-11) to M5 (2028-12) with exit criteria, the gold set, risks |
| 22 | [open questions](22-open-questions.md) | Q-01 to Q-22, each with a Default |

Reading order: 02, 03 and 04 first. Then 05 and 06 (where a page sits and what it is), 07, 08, 09 and 10 (the evidence: execution, languages, versions, tasks), 11 and 12 (formal and visual), 13, 14 and 18 (how the corpus gets complete), 15 (how it stays correct), 16 and 17 (how people and agents use it), and 19 to 21 (how to build it).

## Settled decisions

1. **Nodes for ideas, cards for the long tail.** A node is a directory with a full page. A card is a row in a JSONL register: a language from PLDB, a release, a proposal, a diagnostic, an api symbol, a library. A card becomes a node only by `prog promote` with a reason (04.1).
2. **Two trees, one rule.** Concepts live in the topic tree; language features live in the language trees; every feature links to the concept it realizes (E10), and every concept page shows every tier-1 and tier-2 language (05.1).
3. **The matrix is decided.** Each concept × language cell is `native`, `stdlib`, `library`, `idiom`, `absent` or `n/a`, with an example or a reason (04.8, E11).
4. **Every example runs.** The harness `prog run` executes every block in OCI images pinned by digest, under gVisor, with no network, for every version the block claims. Output blocks come only from `prog run -update` (07).
5. **Versions are computed.** Compat entries follow MDN BCD's model, with implementation columns. A probe proves each `version_added` (K2). Status labels are computed by status-v1, never written (09).
6. **Cite the normative source.** References and standards at the stated version, checked against stored snapshots (13.6). Secondary sources never prove behaviour (G06).
7. **Tiers.** Tier 1 is Python, C, C++, Java, C#, JavaScript, TypeScript, Go, Rust, Kotlin, plus SQL, Bash, regular expressions, HTML/CSS and PowerShell. Tier 2 is about 34 languages. Tier 3 is any language with a runnable image. Tier 4 is cards only (08.2).
8. **Harvest first.** PLDB (public domain) and Wikidata (CC0) are the backbone; release feeds, proposal indexes, stdlib inventories and diagnostic indexes make most cards without a model (13).
9. **Algorithm's pipeline, reused.** Harvest, universe, catalog, stub and promote, write ladder, examples, gate, audit and verifier, one PR per sub-area, merge, deploy, live check (14). The packages are copied from algorithm-agent, not imported (20.2).
10. **Labels for formal claims.** "Proved", "executable" or "normative", never more than the evidence (11.1). The Lean package `Prog` shares algorithm's pin (Lean v4.35.0) and builds on CSLib.
11. **Agents are readers.** Every page has `.md` and `.json` twins; the site has `llms.txt` files, datasets (compat, matrix, cards) and an MCP server (16.10).
12. **Workers and R2 from day one.** A Cloudflare Worker `programming` with static assets for the hot set and R2 `programming-site` for the rest; raw snapshots in R2 `programming-raw` (16.20). Not Pages: algorithm's Pages deploy failed at about 3,500 nodes.
13. **The playground runs in the browser.** WebAssembly runtimes only in v1 (Pyodide, MicroPython, sqlite-wasm, PGlite, wasmoon, ruby.wasm, php-wasm, wasm_of_ocaml, JS and TS). No server execution in v1 (Q-07).
14. **"No missing" is measured.** Every row of every source has an outcome; the scope rule for libraries is published; the claims are per scope and computed (18).

## Conventions

- Section references look like 07.4 (document 07, section 4). Rules look like 15.3 X03, or just X03.
- Tags have 6 characters of Crockford base32, as in algorithm, for example `7K2QPM`. URIs look like `pg://7K2QPM` (Q-02). Slugs look like `closures`; language-owned slugs look like `python/match-statement` (04.3).
- "Estimate" marks a number from research that M0 to M1 must measure. "Default" marks a decision of document 22 that can still change.
- Plain English, short sentences, one sentence per line in the Markdown source of the content repository, no em or en dashes (S01).
- Rule families: T tags and schema, E edges and graph, S style, G grounding, X execution, K compat, F formal, V visualization, C copyright, R cards, L links (15.3).
