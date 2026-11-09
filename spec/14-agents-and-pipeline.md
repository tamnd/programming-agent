# 14. Agents and the pipeline

Most of the writing is done by AI agents, under the checks of 15 and the review of people. This document defines the roles, the models, the pipeline from an empty sub-area to a merged PR and a live page, the work queue, the budgets, and what we keep from algorithm's practice. Algorithm wrote 8,396 stubs and 3,617 pages in its first 7 days (2026-11-03 to 2026-11-09) with this pipeline (01.6); the bottleneck was the gates, not the writing. Programming is about 10 times larger, so the design puts the effort in the gates, the harvest and the queue.

## 14.1 What we keep from algorithm

The `algorithm-work` scripts were the real pipeline, more than the `algo` agent framework. We keep their shape and move them into `prog` as commands (20.4):

| algorithm-work | What it did | In programming |
|---|---|---|
| `catalog.py` | one model call per sub-area lists every item that should be there | `prog agent catalog <area>`; now it starts from the harvested universe and asks the model only for what the sources miss (14.4) |
| `curate.py` | dedupe; a model decides groups and merges | `prog agent curate` with the universe outcomes (18.3) |
| canon match | match the list against the canon books | `prog agent canon`, against the canon of 13.3 and the crosswalks (5.7) |
| `stub.py` | write node stubs | `prog new` in bulk, and `prog promote` for cards (04.1) |
| `writer.py` | a ladder of writer models with gates | `prog agent write`, the same ladder (14.6) |
| `gate.py` | `algo fmt` and `algo check` until PASS | `prog fmt`, `prog check`, `prog run`, `prog compat probe` until PASS (14.7) |
| audit prompt | a cross-model audit returns JSON: verdict, score, issues by severity | the same JSON, plus the claim check of the verifier (15.4) |
| `pr.py` | one PR per sub-area, labels, milestone, `--merge` | the same; the PR body carries the gate report |
| `ship.sh` | writes provenance | `prog prov` |
| `deploy.sh` | `algo site deploy -out build/site -commit $H`, retries until the live page shows the commit; reads the Cloudflare token from `$HOME/data/.env` with `set -a`, never prints it | `prog site deploy` (16.20) with the same live check |
| `report.py` | ticks the tracking issues | `prog report` |

What we change: the framework of roles in algorithm-agent (12-agents.md) is kept as the definition of roles and prompts, but the driver is the pipeline of 14.3, run as batch jobs, not as a long-running agent server.

## 14.2 Roles

Algorithm's roles stay: scout, curator, writer, proof-writer, formalizer, prover, implementer:<lang>, tester, visualizer, linker, critic, fixer. New roles:

| Role | Input | Output | Model class (14.5) |
|---|---|---|---|
| **harvester** | a source adapter's raw snapshot (13.2) | cards and K1 compat entries; no model, or a small model for field extraction only | none or small |
| **matcher** | cards from two sources | `same_as` links, merge proposals with scores | small, then medium for the doubtful |
| **example-writer:<lang>** | a concept or task node and the language's `style.yaml` | the `examples/<lang>.md` section with run blocks; runs `prog run -update` | medium |
| **task-implementer:<lang>** | a task contract | `impl/<lang>/` passing the cases | medium |
| **compat-prober** | a K1 entry or a disagreement between sources | a probe file and the corrected entry, run in the matrix | medium |
| **verifier** | a node at W3 | a claim list: each normative sentence with its citation and a verdict (supported, unsupported, contradicted) from the raw source snapshot | large, different family from the writer (15.4) |
| **matrix-decider** | a concept and a language with an empty cell | the cell value with evidence or reason (4.8) | medium |
| **diagnostic-writer** | a diagnostic card | a trigger program and a fix program, both run (`fails expect=error:<code>`, then `run`) | small or medium |
| **release-watcher** | the release feeds (13.2) | a release card, the pre-release cards, the gap items for new features | none |

## 14.3 The pipeline

One unit of work is a **sub-area** (or a language sub-area for language trees), as in algorithm. The stages:

```
harvest ─▶ universe ─▶ catalog ─▶ curate ─▶ stub/promote ─▶ write ─▶ examples ─▶ gate ─▶ audit ─▶ PR ─▶ merge ─▶ deploy ─▶ live check ─▶ report
   │            │                                            ▲                      │
   └── cards ───┘                                            └──── fixer ◀──────────┘
```

| Stage | Command | Done when |
|---|---|---|
| harvest | `prog harvest <source>` | the snapshot is in R2 `programming-raw`, the cards are written, provenance recorded |
| universe | `prog universe update` | every row of every source has an outcome or is `pending` (18.3) |
| catalog | `prog agent catalog <area>` | the sub-area list is complete against the universe and the model's list (14.4) |
| curate | `prog agent curate <area>` | duplicates merged, groups set, levels set, primary areas set |
| stub, promote | `prog new`, `prog promote` | directories exist, tags minted |
| write | `prog agent write <area>` | the node reaches W2 (writer ladder, 14.6) |
| examples | `prog agent examples <area> --lang L` | each matrix cell of the sub-area's concepts has its section, X1 |
| gate | `prog gate <area>` | fmt, check, run, compat probes, links: all PASS (14.7) |
| audit | `prog agent audit <area>` | no blocker, no major (15.4) |
| PR | `prog pr <area>` | one PR per sub-area, labels `area:<top>`, `kind:<k>`, milestone, the gate and audit report in the body |
| merge | `prog pr <area> --merge` | CI green, then squash merge |
| deploy | `prog site deploy` | the live page shows the merge commit |
| report | `prog report` | the tracking issue is ticked, the coverage report is updated |

Language trees are worked one sub-area at a time across all tier-1 languages ("functions" for all 15 languages, then "collections" for all 15), so the matrix is filled evenly and the comparisons appear early (Default).

## 14.4 The catalog step, harvest first

Algorithm's catalog asked a model to list every item of a sub-area. That works for a closed field with a canon. Programming has structured sources for most of its items, so:

1. The harvester writes cards from the sources (13.2): languages, releases, proposals, diagnostics, api entities, CWE, SPDX, tools, libraries, formats, standards.
2. The universe register (18.3) holds every row with an outcome.
3. The catalog step runs per sub-area with the universe rows that map to it as input, and asks the model (14.5, large class) only: "What items of this sub-area are missing from this list? Name sources." The model's additions go in as `pending` with the source it named. The curator checks each source exists (the link check of 13.6).
4. Two catalog runs with different model families must agree, or the difference goes to the curator (the "two-model catalog" rule, Default).

## 14.5 Models

`agents/models.yaml`, as in algorithm, pins models by class. In Nov 2026:

| Class | Models | Use |
|---|---|---|
| large | `claude-opus-5-5`, `claude-fable-5-1` | catalog, verifier, audit, hard writing (concept pages at the advanced level, formal sections) |
| medium | `claude-sonnet-5-5` | writing, examples, tasks, probes |
| small | `claude-haiku-5-5` | extraction, diagnostics, cards, matching |
| other families | a model from another vendor (as `codex:gpt-6-luna` was used in algorithm), pinned by name in models.yaml | the cross-family jury for the audit and the verifier, and the first rung of the writer ladder |

Rules:

- The audit and the verifier never use the same model family as the writer of the node (15.4).
- Prompts are versioned files in `agents/prompts/<role>/v<n>.md`. A node's provenance records the prompt version and the model id of every change (04.7).
- Every prompt carries the style rules of 6.1 and the prohibitions of 6.10, as algorithm's prompts did ("plain English, no em or en dashes, never invent page or theorem numbers, link only to slugs that exist"), and these new ones: never write a version claim in prose (6.7), never write output by hand (7.2), never name an API that is not in the api cards given in the context (the writer gets the api card list of the language as context).

## 14.6 The writer ladder

As in algorithm's `writer.py`: try the cheapest writer first, run the gates, and climb the ladder on failure:

1. a model of the other family (the first rung, cheap and fast);
2. `claude-haiku-5-5`;
3. `claude-sonnet-5-5`;
4. `claude-opus-5-5`.

Each rung gets the previous failure report. A node that fails at the top rung goes to the human queue with its report. Algorithm's experience: most of the pages were written in unmetered interactive sessions, and about $73 was metered over 957 API events (01.6). Programming plans for metered batch use at scale (14.10).

## 14.7 The gates

A gate is a set of checks that must pass before a node moves on. The checks are the rules of 15. The gate runs in this order, so cheap checks fail first:

1. `prog fmt` (front matter, markup, normalization)
2. `prog check` static rules: schema, links, words (S), grounding form (G01 to G04), version language (K04), symbols (X09)
3. `prog run` for the node's blocks in the default image (X1, X2)
4. `prog test` for tasks
5. `prog compat probe` for new compat entries (K2)
6. `prog run --matrix` (X3), only at W3
7. the citation check against the raw snapshots (G05, 13.6)
8. the audit (15.4), only at W3

A gate failure returns a JSON report that the fixer reads. The fixer has 3 attempts, then the node climbs the writer ladder.

## 14.8 The work queue

The queue is computed, as in algorithm (gaps feeding a work queue), from:

- universe rows with outcome `node` and no node (18.3);
- nodes below their target level for the current milestone (21);
- empty matrix cells (4.8);
- compat entries at K1 without a probe;
- release features with no compat entry (K06);
- failing blocks after an image update (X3 regressions);
- search queries with no result (17.9);
- reader reports (16.19).

Priority: core first (5.4), then by the milestone, then by demand (the search and survey signals, 18.6). The queue is published at `/about/queue/` so that people can take items.

## 14.9 Maintenance jobs

Scheduled jobs that keep the site true:

| Job | When | What |
|---|---|---|
| release watch | daily | new releases and pre-releases → release cards, gap items (9.7) |
| harvest refresh | weekly per source, or by the source's own feed | re-harvest; diff the cards; changed facts go to review |
| image refresh | when a language publishes a patch release | rebuild and pin the image; re-run all blocks of that language (07.4); new failures go to the queue |
| window move | on a release or an end of life | move the support window (9.3); remove redundant `since` attributes; re-run the matrix |
| link check | weekly | every external URL; dead links go to the queue with an archive.org fallback |
| citation recheck | monthly | re-run the citation check against fresh snapshots; changed anchors go to the queue |
| platform jobs | weekly | Windows, macOS, network, GUI blocks (07.8) |

## 14.10 Budget

Estimate, for planning, with prices from the Claude API price list at the time of the run (the `claude-api` reference is the source; no price is fixed in this spec):

| Work | Count (Estimate) | Model class | Notes |
|---|---|---|---|
| node writing to W3 | 60,000 nodes by M5 | mixed, through the ladder | the main cost |
| example sections | 400 core concepts × 15 tier-1 + 400 × 34 tier-2 (decided cells only) ≈ 15,000 sections | medium | |
| task implementations | 300 core × 49 + 1,200 × 15 ≈ 33,000 | medium | |
| diagnostics trigger and fix | 20,000 | small | |
| verifier runs | one per node at W3 and per change | large, other family | |
| harvest extraction | millions of rows | none or small | most sources are structured |

Rules: every run logs its tokens and cost in `prov/` (04.7); `prog report cost` sums them per milestone; a milestone has a budget line in 21, and the queue stops at the line until a person raises it.

## 14.11 People

People are editors, reviewers and contributors, as in algorithm. Roles: **editor** (owns a top area: approves the area tree, the core list, the tier changes), **reviewer** (owns a language or a subject: reviews nodes to W5, compat to K3, examples to X4), **contributor** (anyone, through PRs, with DCO). A W5 review needs a reviewer of the right subject: a language reviewer for a feature, a security reviewer for a weakness, a PL reviewer for a formal section. The reviewer list and subjects are in `registers/reviewers.yaml`.
