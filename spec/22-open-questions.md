# 22. Open questions

Each question has a Default, which is what the spec assumes until someone decides. The ids are stable; other documents refer to them.

| Id | Question | Options | Default | Decide by |
|---|---|---|---|---|
| Q-01 | Repository and binary names | `programming` and `programming-agent` with binary `prog`; or `prog` and `prog-agent` | `tamnd/programming`, `tamnd/programming-agent`, binary `prog` | M0 |
| Q-02 | The URI scheme for tags | `pg://`, `prog://`, `pl://` | `pg://` (short, and distinct from `al://` and `m2://`) | M0 |
| Q-03 | One site or a sub-site per language? | `programming.tamnd.com/python/`; or `python.programming.tamnd.com` | one site, paths (16.2) | M0 |
| Q-04 | The anchor language of concept pages | the first tier-1 language where the concept is native (by register order); or the reader's chosen language always | register order, and the switcher remembers the reader's language (16.4) | M0 |
| Q-05 | Is TypeScript a separate language tree, or a layer on JavaScript? | separate tree; or a sub-tree `javascript/typescript` | separate tree, with `variant_of: javascript` and the features that are only about types | M1 |
| Q-06 | C and C++: one tree or two? | two trees; or one "C family" tree | two trees; the shared preprocessor and the C subset of C++ are features of C linked from C++ | M1 |
| Q-07 | Server-side "Run" for compiled languages | none (v1); self-hosted Judge0 CE or Piston; Cloudflare Sandbox SDK 1.0 with our images | none in v1; revisit at M3 with the cost per run and an abuse plan | M3 |
| Q-08 | How many tier-2 languages | algorithm's 13; or about 34 (8.2) | about 34 | M1 |
| Q-09 | Store cards in git or only in R2? | JSONL in git, sharded; or a database exported to R2 | JSONL in git, sharded by kind and prefix; raw snapshots in R2 only (20) | M0 |
| Q-10 | The library scope threshold (18.5) | dependents ≥ 100; ≥ 50; a top-N per registry | ≥ 100 for npm and PyPI, ≥ 50 elsewhere, at most 5,000 per registry | M1 |
| Q-11 | MDN prose: adapt (CC BY-SA 2.5) or write our own? | adapt with attribution; write our own and link | write our own, link MDN, use BCD (CC0) for data; adapt only when a reviewer asks | M1 |
| Q-12 | Do we show popularity ranks on language pages? | yes, computed with dates; no | yes, from the build, with source and date, in the facts box only | M1 |
| Q-13 | Lean pin: shared with algorithm or own? | shared; own | shared (one pin for the library, 11.2) | M0 |
| Q-14 | Translations | English only; machine-translated with review; community | English only until M5 | M5 |
| Q-15 | Reader accounts (saved language, progress) | none; local storage only; accounts | local storage only | M2 |
| Q-16 | Do we accept library authors' own pages? | yes, as PRs under the rules; no | yes, as PRs under the same rules, reviewed, with a disclosure of authorship | M2 |
| Q-17 | AI disclosure wording on pages | per page provenance line; site-wide statement only | both: a site statement and a per-page line that names the writer roles and the review status (19) | M0 |
| Q-18 | How deep for domain frameworks (React, Django, Spring)? | cards only; nodes for the core; full guides like a language | node for the core frameworks with D2 depth; full guides after M5 | M2 |
| Q-19 | The second model family for audits | a fixed vendor; rotate | a fixed vendor pinned in models.yaml, revisited each milestone | M0 |
| Q-20 | Windows and macOS examples on every PR? | weekly scheduled job; every PR | weekly scheduled job (07.4) | M0 |
| Q-21 | C++26 and other drafts: show as standards or as drafts? | as published; as draft | as `draft` until ISO publishes; show the working draft number (9.2) | now |
| Q-22 | Search back end | static sharded index; Worker with D1 or Vectorize | see 17 | M1 |
