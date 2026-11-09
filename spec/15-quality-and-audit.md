# 15. Quality and audit

The site claims that its pages are correct for stated versions. This document defines the checks that support that claim: the rule families that `prog check` runs, the reviews, the audit by models, the claim check, the hallucination guards, and the metrics that show whether quality holds as the site grows.

## 15.1 The threat model of a generated library

An AI-written programming reference fails in known ways. Each failure has a guard:

| Failure | Example | Guard |
|---|---|---|
| invented API | `str.splitlines(keep=True)` (the parameter is `keepends`) | X09: every code span with a language prefix and every call in a run block resolves to an api card of that language and version; and the block runs |
| wrong version | "f-strings were added in Python 3.5" (3.6) | K04: no version claims in prose; the compat entry has a probe (K2) |
| wrong output | an output block that the model "expected" | X03: output only from `prog run -update` |
| stale fact | "the latest Java LTS is 21" | S07 bans "latest" and "currently"; the facts box is computed; the window moves on its own (9.3) |
| invented citation | a clause number that does not exist | G05: the citation check against the raw snapshot (13.6) |
| plausible but wrong rule | "Go maps iterate in insertion order" | the verifier (15.4) checks each normative claim against its cited source; the X3 matrix runs catch behaviour that differs from the claim when the page shows it in code |
| close paraphrase | a section that follows the Python docs sentence by sentence | the copy check C01 against the snapshots of the cited sources (19.4) |
| bias and filler | "Rust is the best choice for safe code" | S02 banned words; the critic's neutrality check (6.1, rule 7) |
| copied code under a bad licence | a Rosetta Code solution pasted in | C01: the code similarity check against the code corpora of the sources (19.4); C06: no code in the MIT zone from a class B or D source |
| unsafe example | an example that deletes files outside its directory | the sandbox (07.4); no network; review of `vulnerable` blocks (10.8) |

## 15.2 Reviews

As in algorithm: a review is a signed record in `reviews/<tag>.jsonl` with the reviewer, the role, the date, the node hash, the scope (`prose`, `examples:<lang>`, `compat`, `formal`, `all`) and the verdict. A review counts only for the hash it reviewed.

| Ladder level | Review needed |
|---|---|
| W5 | a reviewer of the node's subject, scope `prose` or `all` |
| X4 | a reviewer of the language, scope `examples:<lang>` |
| K3 | a reviewer of the language, scope `compat` |
| F1, F3 | a PL reviewer, scope `formal` |
| V3 | any reviewer, scope `viz` |

Review sampling: when there are fewer reviewers than nodes (always), the review queue samples nodes by risk: the core first, then nodes with many readers (16.15 analytics), then nodes whose audit score was close to the line, then a random 2% of the rest, so the measured error rate of 15.6 is unbiased.

## 15.3 Rule families

`prog check` runs the rules. Each rule has an id, a level (error blocks the gate; warning shows in the report), and the ladder level at which it starts to block. The families:

| Family | Subject | Examples (ids used in other docs) |
|---|---|---|
| **T** | tags, slugs, front matter, the schema (as in algorithm) | T01 schema valid; T02 tag minted; T03 slug unique in its namespace (04.3); T04 reserved words; T05 no status, level, badge or popularity written by hand; T07 no two nodes with the same name and kind without a disambiguation entry (17.6) |
| **E** | edges and the graph | E01 no cycle in `prereqs`, `variant_of`; E02 edge targets have the right kind; E10 every feature realizes a concept or says why not; E11 matrix cells decided (4.8); E12 template sub-areas exist (5.3) |
| **S** | style and words | S01 no em or en dash, no line break inside a sentence; S02 banned words; S03 sentence length (warning above 25 words); S04 headings match the kind's section list (6.3); S05 preferred terms (6.5); S06 names of languages and products; S07 time words that rot (6.7); S08 length limits (6.9) |
| **G** | grounding | G01 locators; G02 footnotes used and resolved; G03 deny-listed sources; G04 quotes at most 50 words with a licence; G05 citation check; G06 secondary sources not used as proof of behaviour; G07 every number and date in front matter has a source (04.5) |
| **X** | examples and execution (07) | X01 run mode; X02 output block; X03 output from the runner; X04 failure signature; X05 diagnostic codes have cards; X06 copied text runs; X07 stdlib only; X08 `norun` share; X09 symbols resolve; X10 benchmarks; X11 browser parity |
| **K** | compat (09) | K01 schema of compat entries; K02 every `version_added` is a real release in the register; K03 front matter `since` agrees with compat; K04 no version claims in prose; K05 probe failure signatures; K06 every feature in a release's notes has a compat entry |
| **F** | formal (11) | F01 "proved" only with a theorem; F02 model versus language; F03 executable semantics with version and date |
| **V** | visualization (12) | defined in 12 |
| **C** | copyright and licence (19.11) | C01 the copy check (thresholds in 19.4); C02 `reuse lint` passes; C03 adapted sources are class A, B or C and carry their attribution; C04 person photos; C05 exercise statements not copied; C06 no class B or D origin in the MIT zone; C07 no class B, C or D text in the CC0 zone; C08 every cited source is in `registers/sources.yaml`; C09 no links to exploit code; C10 living persons. 19.11 is the authority |
| **R** | cards (04.1) | R01 card schema; R02 every card has a source with a snapshot hash; R03 `same_as` ids resolve; R04 no duplicate cards (the matcher's threshold); R05 owned cards have an owner that exists |
| **L** | links | L01 internal links resolve (`[[...]]`, `al:`, `m2:` against the other sites' published tag lists); L02 external links answer (the weekly job, 14.9) |

The full list with messages is generated into `/about/rules/` from the code, as in algorithm.

## 15.4 The audit and the verifier

At W3 every node gets two model checks, each by a model of a different family from the node's writer (14.5).

**The audit** (algorithm's audit prompt, extended). Input: the node as rendered Markdown, its kind's section list, the style rules, the run results. Output JSON:

```json
{"verdict": "pass|revise|reject", "score": 0.0,
 "issues": [{"severity": "blocker|major|minor", "section": "Semantics", "rule": "r3",
             "text": "Claims that dict iteration order is insertion order since 3.6; it is guaranteed since 3.7 (3.6 was an implementation detail of CPython).",
             "fix": "State 3.7 as the language guarantee, link the compat entry."}]}
```

A blocker or a major issue sends the node to the fixer. The audit also checks: neutrality, level fit, the Diátaxis form (6.1), close paraphrase by reading, and the pitfalls a reader would expect.

**The verifier** (the claim check). It splits the node into claims (each normative sentence, each rule `rN`, each table row with a fact). For each claim it reads the cited source from the raw snapshot (not from memory, and not from the web at audit time), and returns `supported`, `unsupported` (the source does not say it), or `contradicted`. Rules:

- G2 needs every normative claim `supported`.
- An `unsupported` claim gets a better citation, or it is removed, or it is turned into an executed example that shows it (behaviour that an example shows needs no citation for G2, because the run is the evidence for the stated versions).
- A `contradicted` claim is a blocker.

## 15.5 Human-facing signals of quality

Each page shows its ladder levels as badges (16.4), the "verified on <date> against <versions>" line, the review status, and a "Report a problem" link (16.19). A page below W4 shows a notice that it has not been checked by the audit yet. A page with a reader report that a reviewer confirmed shows the report until it is fixed.

## 15.6 Metrics

`prog report quality` publishes, per area and per milestone:

| Metric | How it is measured | Target at M5 (Default) |
|---|---|---|
| error rate of W4 pages | reviewers' sample (15.2): share of sampled nodes with at least one major issue | below 2% |
| error rate of cards | a sample of 500 cards per quarter checked against sources | below 0.5% |
| X3 share | blocks passing the full matrix / all `run` blocks of tier 1 | above 98% |
| K2 share | tier-1 compat entries with a probe that passes both checks | above 90% |
| G2 share | nodes at W3 or more with every normative claim supported | 100% |
| time to fix a confirmed reader report | median days | below 7 |
| regressions per image update | blocks that fail after an image update / blocks run | tracked, no target |

The audit itself is audited: each quarter, a person checks 100 random audit verdicts, and the rate at which the audit missed a major issue is published.
