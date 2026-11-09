# 20. Repository layout

This document gives the file trees of the two repositories, the storage rules (what lives in git and what lives in R2), the size of each part, the `prog` command tree with its exit codes, the releases and workflows, and the rules that keep programming and algorithm in step. It follows 2150/18 (algorithm). The licence of each path is in 19.1.

| Repository | Holds | Licence |
|---|---|---|
| `tamnd/programming` | the content: nodes, cards, compat entries, registers, crosswalks, the `Prog` Lean package, provenance, reviews, reports | per path (19.1) |
| `tamnd/programming-agent` | the toolchain: the Go module with the `prog` binary, the runner image definitions, the agent prompts, the Worker, the spec copy | Apache-2.0; the viz package and the Worker MIT |

## 20.1 `tamnd/programming`, the content

```
programming/
  README.md  CONTRIBUTING.md  CITATION.cff  LICENSE  REUSE.toml  SECURITY.md
  LICENSES/                 CC-BY-SA-4.0.txt  MIT.txt  Apache-2.0.txt  CC0-1.0.txt  README.md (the split of mixed files, 19.1)
  prog.yaml                 pins: the prog version, Lean v4.35.0, CSLib and Mathlib revisions, the image register digest, the model register version
  nodes/
    <top>/<sub>/<slug>/     one node per directory (04.2); <top> is an area of topics/areas.yaml
      node.md               front matter (CC0) and body (CC BY-SA 4.0)
      notes.md              optional long notes, history, open questions
      examples/<lang>.md    the language sections of a concept: prose and fenced blocks with attributes (07.2)
      matrix.yaml           the concept x language cells: status, reason, the example ids (04.8)
      compat.yaml           compat entries of this node per language and version (09)
      probes/<lang>/        the probe programs that decide each compat entry (09.4)
      contract.yaml         a task's contract: signature, input and output types, rules (10.2)
      cases/*.jsonl         shared test cases of a task (10.3)
      fixtures/             input files the cases and examples read
      impl/<lang>/          a task's implementation and its generated test file (10.4)
      vuln/<lang>/vulnerable/  vuln/<lang>/fixed/   the two programs of a weakness node (10.8)
      vuln/exploit.jsonl    inputs that break the vulnerable program and not the fixed one, against the local fixture only
      viz/                  scene.yaml, scene code, trace hooks (12)
      figures/              SVG figures made from scenes or by hand
      tests/                extra tests: the property and differential tests of the node
    languages/<lang>/       the language tree (05.4, 08)
      node.md               the language page
      releases.yaml         every release: version, date, end of life, the release card ids (09.3)
      implementations.yaml  compilers, interpreters, runtimes, with the versions in the window
      concepts.yaml         the language's position on each core concept, which drives the matrix (04.8)
      style.yaml            naming, formatting, the formatter and linter pins for this language's pages
      areas.yaml            the language's own sub-areas (stdlib, tooling, ecosystem, idioms)
      <sub>/<slug>/         feature, release, api, diagnostic and idiom nodes, with the same files as above
  cards/
    <kind>/<owner>/<pp>.jsonl   card rows sharded by kind, owner and prefix (20.3)
    <kind>/README.md        the schema version and the shard rule of that kind
  registers/
    tags.reg  aliases.reg  reserved.yaml  disambiguation.yaml
    languages.yaml          every language with tier, depth band, comment syntax, file extension, test framework, wrappers; the source of truth for algorithm (20.6)
    implementations.yaml    every implementation with its language and image family
    support-policy.yaml     the version window rule per language (09.3)
    images.yaml             image name, language, version, implementation, digest, build date, smoke test result (07.4)
    terms.yaml  names.yaml  banned-words.yaml  notations.yaml   the style registers (06)
    sources.yaml            every source with licence, class, date checked (19.3)
    sources-deny.yaml       class E sources and domains we must not fetch
    tasks.yaml              the task list with ids and the source of each name (10.1)
    core.yaml               the core concept list (05.6)
    compare-pairs.yaml      the language pairs that get comparison pages (16)
    canon.yaml              the canon of books, standards and references (13.3)
    universe/<source>.jsonl the universe rows and their outcomes (18.3)
    scope.yaml              the scope rules for libraries, tools and people (18.5)
  topics/
    areas.yaml              the area tree: top areas, sub-areas, owners, hues (05.2)
    <area>.md               the area introductions
  languages/<lang>/         per-language configuration of the harness, not pages
    wrappers/               program, main, expr, repl templates (07.3)
    normalize.yaml          the output normalizers beyond the defaults (07.6)
    style.yaml              the formatter and linter command lines with their pins
  compat/
    imports/<source>/       pinned imports, for example bcd.json and web-features.json, each with its own licence file (19.1)
    status-v1.yaml          the status vocabulary of compat entries (09.2)
  crosswalk/
    cs2023.yaml  swebok-v4.yaml  acm-ccs-2012.yaml  cwe.yaml  wikidata.yaml  pldb.yaml
    linguist.yaml  so-tags.yaml  rosetta.yaml  spdx.yaml  algorithm.yaml
  lean/
    lakefile.toml  lean-toolchain  lake-manifest.json
    Prog/                   the Prog package: Core, Semantics/<lang>/, Types/, the node bindings (11)
  harvest/
    <source>.yaml           adapter config: URL roots, schedule, rate limit, licence class, outputs (13.2)
    manifests/<source>/<date>.jsonl   one line per fetched response: URL, time, status, BLAKE3 hash, size; the bytes are in R2 only
  prov/
    <tag>.jsonl             provenance of node files (04.7)
    cards/<kind>.jsonl      provenance of card rows per harvest run
    runs/<run-id>.yaml      agent run records: model id, prompt version, items, cost (14)
  reviews/<tag>/<part>.yaml person review records with hashes (15.5)
  schema/
    node.v1.json  card.v1.json  compat.v1.json  contract.v1.json  register.v1.json  scene.v1.json
  manifests/releases/<YYYY.MM>.yaml   the release manifests (20.5)
  reports/                  coverage, quality, compat, matrix, cost; JSON and Parquet, written by prog
  vendor/<source>/          data imported as-is under the source's licence (19.1)
  build/                    ignored by git; made by prog
    runs/                   run outputs keyed by block hash and image digest (mirrored to R2 programming-build)
    data/examples.jsonl     every executed block with its result
    data/compat/<lang>.json the compat dataset per language
    graph/  ladders/  lean/  viz/  site/  search/
  .github/workflows/        20.5
```

**Rules of the tree.**

- A path is the address of a node in the tree and is not its identity; the tag is (04.3). `prog mv` moves a directory and writes an alias.
- Language-specific nodes live under `nodes/languages/<lang>/`. Concepts that span languages live under `nodes/<top>/<sub>/` and hold one `examples/<lang>.md` per language.
- The `languages/<lang>/` tree at the root is configuration for the harness. The `nodes/languages/<lang>/` tree is content. The check rule R-L1 fails if a language has one without the other.
- `build/` is never committed. Everything in it can be made again from the tree, the pins and the raw store.
- No file holds a secret. Deploy tokens live in CI secrets only.

## 20.2 `tamnd/programming-agent`, the toolchain

```
programming-agent/
  go.mod                    module github.com/tamnd/programming-agent; Go 1.26 (Default, pinned at M0)
  Makefile  README.md  CONTRIBUTING.md  SECURITY.md  THIRD_PARTY.md  LICENSE  LICENSES/
  cmd/prog/                 main.go and the command tree (20.4)
  # copied from algorithm-agent at M0, then owned here (20.6)
  node/                     node directories, front matter, sections
  check/                    the rule families and the rule runner (15)
  lang/                     the language register, tiers, style per language (08)
  markup/                   Markdown with block attributes, cross-links, math (06, 07.2)
  katex/                    math at build time
  viz/                      scenes, tracers, frames; published as @prog/viz (MIT) (12)
  site/                     renderer, templates, assets, deploy, live check (16)
  prov/                     provenance records (04.7)
  mint/                     tag mint, check, retire (04.3)
  # new in programming-agent
  card/                     read, write, shard, merge and promote cards; the matcher (04.1, 13)
  compat/                   compat entries, probes, imports, status-v1 (09)
  run/                      the harness: blocks, wrappers, normalizers, the gVisor sandbox, the run cache (07); the runner package of 03.4
  images/                   the Go code that builds, pins and tests runner images (07.4)
  harvest/
    <source>/               one adapter per source: mdn, bcd, web-features, pldb, wikidata, linguist, python-docs, rustdoc, javadoc, cwe, osv, peps, jeps, ... (13)
    raw/                    the R2 raw store client, BLAKE3 keys, manifests
  matrix/                   the concept x language matrix and rule E11 (04.8)
  search/                   the sharded index, the query parser, the eval set (17)
  mcp/                      the MCP server (16.10)
  graph/  taxonomy/  schema/  release/  testkit/  lean/  cite/  copycheck/  universe/  coverage/  status/  audit/  agent/  queue/  dataset/
  images/                   (data, beside the Go package above) runner image definitions
    <lang>/<version>/Dockerfile      one per language version, for example images/python/3.14/Dockerfile
    <lang>/<version>/smoke/          smoke programs the image must pass before it enters the window
    <lang>/<impl>/<version>/Dockerfile  where an implementation needs its own image (pypy, graalvm, bun, deno)
    base/                   the shared base layers (Debian slim, gVisor-friendly users and limits)
  agents/
    models.yaml             the pinned models by class (14.2)
    prompts/<role>/v<n>.md  versioned prompts: harvester, compat-prober, example-writer, task-implementer, verifier, catalog, curate, write, audit, critic (14)
    roles.yaml              the roles, their tools, budgets and gates
  worker/                   the Cloudflare Worker programming, TypeScript (Default), about 1,000 lines (16.20)
    src/  test/  wrangler.toml (no secrets; the account id comes from CI secrets)  package.json
  testdata/                 small trees used by the Go tests
  spec/                     a copy of these notes, checked by spec.yml
  scripts/                  one-off maintenance scripts
  .github/workflows/        20.5
```

The Go package `images/` and the image definitions share the directory name. Default: the Go code lives in `internal/imagebuild/` if the Go tests or the Docker build context make the shared directory awkward; the decision is made at M0.

Copied packages start as copies of algorithm-agent and do not import it, so the two toolchains can change at their own pace (20.6).

## 20.3 Storage and size

**What goes where.**

| Data | Store | Rule |
|---|---|---|
| nodes, examples, compat entries, probes, contracts, cases, implementations, Lean, registers, crosswalks, provenance, reviews, reports | git, `tamnd/programming` | plain text, reviewed in pull requests |
| cards | git, `cards/<kind>/<owner>/<pp>.jsonl` | JSONL, one card per line, sorted by tag, keys in schema order, so a diff shows one line per changed card |
| raw source snapshots | R2 `programming-raw/<source>/<date>/<blake3>` | never in git; git holds only the manifests with the hashes (Default) |
| run cache | R2 `programming-build`, keyed by block hash x image digest | rebuildable, so no backup |
| runner images | GHCR `ghcr.io/tamnd/prog-runner/<lang>:<version>`, used by digest | the digest is in `registers/images.yaml` |
| the site and datasets | Worker static assets (hot set) and R2 `programming-site` | made by `prog site build` (16.20) |
| release archives | GitHub releases and R2 `programming-site/releases/<YYYY-MM>/` | per release (20.5) |

**No Git LFS.** LFS bandwidth and storage on GitHub are metered, and LFS files do not diff. Large binary inputs belong in R2 with a hash in git. Rule R-G1 fails any file over 1 MB in git except card shards and the imported compat data; rule R-G2 fails any binary file outside `figures/` and `fixtures/`.

**Card shards.** A card file is `cards/<kind>/<owner>/<pp>.jsonl`: `<owner>` is the language slug for language-owned kinds (api, diagnostic, feature, release, library) and `_` for the others; `<pp>` is the first two characters of the card's slug. A shard splits to three characters when it passes 5,000 rows or 2 MB (Default). `prog fmt` writes shards in canonical order, so two runs that make the same cards make the same bytes.

**Raw snapshots.** The bytes stay in R2 because they are large, they change often, and many of them are under licences that do not allow redistribution in a public repository (19.3); only the hash and the URL are public. The raw store cannot be rebuilt (the web changes), so it has a monthly copy to a second bucket (03.3).

**Size at M5 (Estimate).**

| Part | Files | Size |
|---|---|---|
| `nodes/` (about 60,000 nodes) | about 500,000 (node.md, examples, compat, probes, cases, impl) | 2 to 4 GB in the working tree |
| `cards/` (about 550,000 cards) | about 1,500 shards | 400 to 700 MB |
| `prov/` | about 62,000 | 200 to 400 MB |
| `registers/`, `crosswalk/`, `topics/`, `schema/`, `languages/` | about 3,000 | about 150 MB, most of it `registers/universe/` |
| `lean/` | about 3,000 | about 30 MB |
| `harvest/manifests/` | about 20,000 | about 200 MB |
| all of `tamnd/programming` | about 600,000 | 3 to 6 GB packed, as 03.3 |
| `tamnd/programming-agent` | about 3,000 (with about 400 image definitions) | under 100 MB |
| R2 `programming-raw` | about 10 million objects | 50 to 200 GB |
| GHCR runner images | about 400 | 200 to 400 GB |

**Working with a large tree.** `prog` works on a sparse checkout: `prog sparse <area>...` sets the sparse patterns to the named areas plus `registers/`, `schema/`, `topics/` and `languages/`. CI uses partial clone (`--filter=blob:none`) and the changed paths only. Default threshold: if a fresh clone passes 8 GB or `git status` on the build server passes 2 seconds, cards move to `tamnd/programming-cards` (Q-09 in 22).

## 20.4 The `prog` command tree

All commands share one exit code convention, the same as `algo`:

| Code | Means |
|---|---|
| 0 | the command did its work and found nothing wrong |
| 1 | the command ran and found failures: rule violations, failed runs, output drift, failed probes, a broken link |
| 2 | usage error: an unknown command or flag, a slug or tag that does not exist, a bad path |
| 3 | environment error: a missing tool or image, the sandbox is not available, R2 or the network failed, a budget is exhausted, a lock is held |

CI treats 1 as a red check and 3 as a retry. Every command takes `-json` to print its report as JSON for agents and CI.

| Command | Does | Spec |
|---|---|---|
| `prog new <kind> <slug> [-area A] [-lang L]` | make a node directory with a provisional id and the sections of its kind | 04, 06 |
| `prog new viz <node> [-from-trace ID]` | make a scene, recorded from an example's trace if asked | 12 |
| `prog fmt [paths]` | format front matter, card shards, compat files and each code block with its language's formatter | 04, 07 |
| `prog mint` | give tags to provisional ids before the merge | 04.3 |
| `prog promote <card> [-reason R]` | turn a card into a node, keeping the tag | 04.1 |
| `prog mv <slug> <area>` | move a node and write an alias | 04.3 |
| `prog check [-changed] [-rule ID]` | the rule families of 15 over the changed or named items | 15 |
| `prog status <slug>` | the ladders W, G, X, K, F, V of a node and the reasons | 04.9 |
| `prog gen [tests\|wrappers] [lang]` | make generated files: test files from cases, build trees | 10.4 |
| `prog run [slug] [-lang L] [-update] [-matrix]` | run the blocks of a node in their images and compare the output; `-update` writes the new output into the block; `-matrix` runs every version in the window | 07 |
| `prog test [slug] [-lang L] [-kind cases\|prop\|diff\|security]` | run task tests and the security harness | 10 |
| `prog compat probe [slug] [-lang L] [-version V]` | run the probes and write compat entries | 09.4 |
| `prog compat import <source>` | import BCD, web-features and other compat data | 09.6 |
| `prog harvest <source> [-since D]` | run one adapter: raw snapshots to R2, cards, universe rows, provenance | 13 |
| `prog universe update` | join harvest rows into the universe register and set outcomes | 18.3 |
| `prog release watch` | read release feeds; make release cards, image definitions and queue items | 09.7 |
| `prog images build [lang] [version]` | build runner images from `images/` and run their smoke tests | 07.4 |
| `prog images pin [lang]` | write the digests of passing images into `registers/images.yaml` | 07.4 |
| `prog lean build \| status` | build the `Prog` package, report bindings and the F level | 11 |
| `prog viz check [slug]` / `prog viz serve <slug>` | check scene traces; open the player with reload | 12 |
| `prog build` | the full build with the cache | 03.5 F4 |
| `prog search build \| eval` | build the index; run the eval set | 17 |
| `prog site build \| serve \| check \| deploy [-rollback]` | render, preview on :8151, check and publish | 16 |
| `prog mcp serve` | run the MCP server | 16.10 |
| `prog queue [list\|add\|lease]` | the work queue | 14.8 |
| `prog agent run <role> [-item ID] [-budget USD]` | run an agent role on queue items | 14 |
| `prog agent catalog \| curate \| canon \| write \| examples \| audit <area>` | the steps of the pipeline | 14.3 |
| `prog gate <area>` | fmt, check, run, probes, links for a sub-area | 14.7 |
| `prog pr <area> [-merge]` | open or merge the pull request of a sub-area | 14.3 |
| `prog prov <slug>` | show or write provenance | 04.7 |
| `prog review <slug> <part>` / `prog quarantine <slug> <part> -reason R` | person review records; withdraw a part (19.10) | 15 |
| `prog report [coverage\|quality\|compat\|matrix\|cost]` | write the reports into `reports/` | 15, 18 |
| `prog sparse <area>...` | set a sparse checkout | 20.3 |
| `prog sync algorithm` | write the languages export for algorithm (20.6) | 20.6 |
| `prog release cut \| verify \| tag` | the monthly release | 20.5 |

## 20.5 Releases and workflows

**Releases.** Calendar versions `YYYY.MM`, one each month, as in algorithm. The first is 2027.03. `prog release cut` builds `main` and writes `manifests/releases/<YYYY.MM>.yaml` (the commit, the `prog` version, the Lean and CSLib pins, the image register digest, the model ids and prompt versions, the ladder counts, the coverage), a second runner checks it with `prog release verify`, a person writes the summary and merges, and `prog release tag` tags, makes the GitHub release with the offline archive, and builds the release site at `https://<YYYY-MM>.programming.tamnd.com` (16). The tag `2027.03` and the URI form `pg://<tag>@2027.03` use the dot; the host name uses a hyphen, because a dot would make a new DNS label.

**Workflows of `tamnd/programming`.**

| Workflow | Runs on | Does |
|---|---|---|
| `check.yml` | each pull request | `prog check -changed`, `prog run` for changed blocks (sharded by image), probes of changed compat entries, Lean of changed modules, `reuse lint`; posts the ladder table |
| `tags.yml` | each pull request | fails while a provisional id exists; DCO sign-off check |
| `build.yml` | each push to `main`, and nightly on the build server | `prog build` with the cache |
| `pages.yml` | after `build.yml` | `prog site build`, `site check`, `site deploy` to the Worker and R2; the GitHub Pages mirror |
| `harvest.yml` | per source schedule (13.2) | `prog harvest <source>`, `universe update`; a pull request per source |
| `release-watch.yml` | daily | `prog release watch`; pull requests for new releases and images |
| `matrix.yml` | nightly, and weekly for Windows and macOS | `prog run -matrix` over the changed images and the nodes going to W3 |
| `queue.yml` | daily | `prog queue` from gaps; agents lease items on the build server |
| `audit.yml` | daily | the full audit; a pull request when the report changes |
| `release.yml` | the first day of each month | `prog release cut` when the variable `PROG_RELEASES` is `true` |
| `codeql.yml`, `scorecard.yml` | weekly | security checks |

**Workflows of `tamnd/programming-agent`.**

| Workflow | Does |
|---|---|
| `ci.yml` | Go tests, lint, the `@prog/viz` tests, builds of `testdata/` trees |
| `images.yml` | build and smoke-test runner images when `images/` changes; push to GHCR; open a pull request in `tamnd/programming` with the new digests |
| `worker.yml` | type check and test the Worker; deploy it on a tag (secrets from CI only) |
| `release.yml` | build `prog` for Linux, macOS and Windows on a tag; publish `@prog/viz` |
| `spec.yml` | check the links and section references of `spec/` |

## 20.6 Sync with algorithm

Programming and algorithm share the language facts, parts of the toolchain, and links between their nodes. The rules:

**The languages register.** `tamnd/programming` is the source of truth for languages: names, slugs, the toolchain versions, the image digests, the test frameworks, the formatters and the comment syntax. Algorithm keeps its own `registers/languages.yaml` with its 23 language rows (05.8), but it does not edit their toolchain fields by hand.

1. At each programming release, `prog sync algorithm` writes `build/export/algorithm-languages.yaml`: the 23 rows that algorithm uses, with the fields that algorithm's register has, from programming's `registers/languages.yaml` and `registers/images.yaml`.
2. A workflow in `tamnd/algorithm` (`bump.yml`) reads the export from the programming release and opens a pull request that replaces those fields. A person merges it after algorithm's tests pass on the new images.
3. Fields that only algorithm has (for example its benchmark settings) stay in algorithm and are not touched.
4. A new language in algorithm needs a row in programming first. Rule R-S1 in programming fails if a language slug used in `crosswalk/algorithm.yaml` has no row.

Algorithm may stay one release behind programming. It must not be more than two releases behind (Default); the programming audit reports the gap.

**The images.** Algorithm's runner images and programming's runner images converge on one set, `ghcr.io/tamnd/prog-runner/<lang>:<version>`, at M2 (Default). Until then each repository pins its own digests.

**The copied packages.** The packages node, check, lang, markup, katex, viz, site, prov and mint are copied from algorithm-agent at M0 with the commit recorded in `THIRD_PARTY.md`. After the copy they are owned by programming-agent and do not import algorithm-agent. A fix in one is ported to the other by hand, with a pull request that links the original.

**Cross links.** `crosswalk/algorithm.yaml` maps programming nodes to algorithm tags (for example a sorting task to the algorithm's node). Pages link across the two sites with the target's tag URL, so a move on either side does not break the link. `prog check` reads algorithm's released tag register to check these links.
