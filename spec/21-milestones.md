# 21. Milestones

Status: **Draft**. Dates are targets, not promises. Algorithm wrote 8,396 stubs and 3,617 pages in 7 days (01.6), so writing speed is not the limit; the gates, the images, the harvest, and review are. Each milestone has exit criteria that the coverage report (18.7) can show, and a budget line (14.10).

GitHub tracking, as in algorithm: milestones `M0` to `M5`, plus one milestone per top area (`Area: Concurrency`), one tracking issue per sub-area, one PR per sub-area.

## M0. Foundations (2026-11 to 2026-12)

Goal: the system works end to end on a small gold set, and the site is live.

- Create `tamnd/programming` and `tamnd/programming-agent`. Copy the algorithm-agent packages (node, check, lang, markup, katex, viz, site, prov, mint) into `prog`, as algorithm did from math2 (20.2).
- The schema for all kinds (04.5), the card registers (04.1), the rule families T, E, S, G, X, K (15.3).
- The harness `prog run` (07) with the attributes of 7.2, the cache, the normalizers.
- Runner images for tier 1 with the full version matrix (07.4), pinned by digest, built in CI.
- `prog compat probe` (09.4).
- The site: templates for all kinds, the facts box, the language switcher, the playground for Python (Pyodide), JavaScript, SQL (sqlite-wasm), Lua.
- Hosting on Workers static assets + R2 + the Worker from day one (16.20), with `prog site deploy` and the live check.
- **The gold set**: 15 to 20 nodes written to C5 by hand and agent together, one or two of each main kind, used as the examples in every prompt:

| Gold node | Kind |
|---|---|
| closures | concept |
| async/await | concept |
| garbage collection | concept |
| Python: the match statement | feature |
| Rust: ownership | feature |
| Python | language |
| Python 3.15 | release |
| read a file line by line | task |
| observer pattern | pattern |
| extract function | technique |
| unit testing | concept + practice |
| git | tool (with its first sub-nodes) |
| JSON | format |
| SQL injection (CWE-89) | weakness |
| E0382 (borrow of moved value) | diagnostic, promoted to a node |
| PEP 634 | proposal card |
| str.split | api card |

**Exit**: the gold set at C5 live on programming.tamnd.com; CI runs the matrix; a deploy is one command.

## M1. The universe and the cards (2027-01 to 2027-02)

- Harvest adapters for the sources of 13.2. All sources of 18.4 snapshot into R2 `programming-raw`.
- The universe register with outcomes for every row; `pending` below 5% overall, 0% for languages, releases and CWE.
- All cards: languages (about 12,000), releases, proposals, diagnostics, api for tier 1, CWE, SPDX, tools, libraries by the scope rule.
- The topic tree and the language template fixed (05); the core list approved (5.4).
- Search v1 with the card index (17).
- Site v1 at full size: collection pages for every card kind.

**Exit**: every card scope at C1; the core list approved; the site serves about 500,000 card entries.

## M2. The core concepts and the tier-1 guides (2027-03 to 2027-06)

- The 400 core concepts at C3, with all tier-1 matrix cells decided.
- The tier-1 language guides at D3 (8.3): the tutorial features of each of the 15 tier-1 languages at C3.
- Compat at K2 for every tier-1 feature in the guides.
- The core tasks (300) with tier-1 implementations.
- Comparison pages for the top 50 pairs.

**Exit**: core concepts at C3; tier-1 at D3; K2 share for tier 1 above 80%.

## M3. Breadth (2027-07 to 2027-12)

- Core tasks in all 49 tier-1 and tier-2 languages.
- Patterns, refactorings, practices of the core at C3.
- The security area: CWE Top 25 and the OWASP Top 10 as nodes at C3, with the security harness (10.8).
- Tools of the core at C3; Git in depth.
- Tier-2 languages at D3.
- Visualization: the scenes for the concepts with `viz: expected` at V2.
- The core at C4.

**Exit**: the core at C4; tier 2 at D3; the error rate of W4 pages measured (15.6).

## M4. Depth (2028-01 to 2028-06)

- Tier-1 languages at D4: every section of each reference mapped (8.3); the full stdlib as api cards; the full compat window.
- The formal developments of 11.2 at F2; F3 for the interpreters.
- The non-core concepts at C3.
- Tier-3 languages at D2.
- Review drive: the core at C6 as far as reviewers allow.

**Exit**: tier 1 at D4 and C4; the core at C5.

## M5. No missing (2028-07 to 2028-12)

- Every scope of 18.4 with no `pending` rows.
- Every `node` outcome at C3 or more; the core at C5; the reviewed share published.
- The claims of 18.1 computed and shown for every scope.

**Exit**: the "complete" claims hold for every scope, with the date and source versions. After M5 the work is maintenance (14.9), the review drive, and the new knowledge that each month brings.

## Budget lines

Each milestone gets a model-cost line and a CI-minutes line, set at the start of the milestone from the measured cost per node of the previous one (14.10). The first lines are set at the end of M0 from the gold set's measured cost. No number is fixed here.

## Risks

| Risk | Effect | Response |
|---|---|---|
| CI time of the version matrix | thousands of images × hundreds of thousands of blocks | the cache (07.4); run the full matrix only on image changes and at W3; shard on free runners; a paid runner pool if needed |
| reviewer shortage | W5, X4, K3 stay low | sampling (15.2); show computed evidence first; invite language communities per language |
| source access | Rosetta Code blocks automated access (403 in Nov 2026); HOPL has no dump; Stack Overflow data terms | names-only sources are not critical; use PLDB and Wikidata as the backbone (13.2) |
| licence mistakes | a takedown, a forced rewrite | C rules (15.3), the source table (19.3), provenance on every change |
| hosting limits at scale | the algorithm Pages failure at 3,500 nodes | Workers + R2 from day one (16.20) |
| standards in flux | C++26 not published (ISO ballot failed Oct 2026); C2y in draft | the `draft` marker on standard columns (9.2); pages say "draft" until publication |
| the model writes plausible errors | wrong facts at scale | the verifier, the citation check, execution (15) |
