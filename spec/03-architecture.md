# 03. Architecture

The shape is algorithm's (2150/03): one content repository, one toolchain repository with one Go binary, several stores, a few flows, and a deploy to a Cloudflare Worker with R2. Three things are new. There are **cards** next to nodes, about ten cards for each node (04.1). There are **runner images** for many language versions, about 400 at M5, because execution is the main evidence (07). And there is a **harvest** that makes most of the cards without a model (13).

## 3.1 Two repositories

| Repository | Holds | Licence (19) |
|---|---|---|
| `tamnd/programming` | the content: node directories (prose, examples, compat entries, scenes), the card registers (JSONL, sharded), the `Prog` Lean package, the registers (languages, tags, universe, canon, scope rules), the reviews, the reports | prose CC BY-SA 4.0; code in examples and tasks MIT; Lean Apache-2.0; cards, compat data and metadata CC0 1.0 |
| `tamnd/programming-agent` | the toolchain: the Go binary `prog` (check, run, compat, harvest, build, site, search, agents, mcp), the runner image definitions, the harvest adapters, the visualization runtime, the prompts and `models.yaml` | Apache-2.0; the visualization runtime MIT |

As in algorithm, content changes are pull requests on `tamnd/programming` and the checks are the gate. Toolchain changes are releases of `tamnd/programming-agent`. `manifests/releases/next.yaml` in the content repository pins the version of `prog` and the digests of the runner images that CI uses. Names are Q-01 (22).

## 3.2 The boxes

```
                    ┌────────────────────────────────────────────┐
 sources (13) ────▶ │ harvest: adapters → raw snapshots (R2)     │
 release feeds ───▶ │   → cards, compat imports, universe rows   │
                    └───────────────────┬────────────────────────┘
                                        │ gaps (18), new releases (9.7)
                                        ▼
 people ──PR──▶ ┌───────────────────────────────────────────────┐ ◀──PR── agents (14)
                │ tamnd/programming: nodes/, cards/, lean/,      │
                │ registers/, compat/, reviews/                  │
                └───┬──────────────┬───────────────┬────────────┘
                    │              │               │
                    ▼              ▼               ▼
          ┌──────────────┐ ┌───────────────┐ ┌────────────────────────┐
          │ check: rules │ │ run: examples │ │ compat probes (9.4)    │
          │ T E S G X K  │ │ in pinned OCI │ │ lean build (11)        │
          │ F V C R L    │ │ images, cache │ │ tasks and tests (10)   │
          │ audit (15)   │ │ (07)          │ │                        │
          └──────┬───────┘ └──────┬────────┘ └───────────┬────────────┘
                 └────────────────┼──────────────────────┘
                                  ▼
                 ┌────────────────────────────────────┐
                 │ build: graph, matrix, ladders,     │
                 │ status labels, coverage, search    │
                 │ shards, datasets, scene bundles    │
                 └─────────────────┬──────────────────┘
                                   ▼
                 ┌────────────────────────────────────┐
                 │ site: HTML, .md, .json, llms.txt,  │──▶ Worker `programming`
                 │ datasets, feeds, MCP               │    + static assets (hot set)
                 └────────────────────────────────────┘    + R2 `programming-site`
```

## 3.3 The stores

| Store | Where | What | Size at M5 (Estimate) |
|---|---|---|---|
| content | git, `tamnd/programming` | nodes, cards (JSONL sharded by kind and tag prefix), compat entries, Lean, registers, reviews | 3 to 6 GB; most of it is cards and examples. If git slows down, cards move to a separate repository `tamnd/programming-cards` (Q-09) |
| raw | R2 `programming-raw`, keyed by BLAKE3 | every harvested response and every cited page snapshot (13.4) | 50 to 200 GB |
| images | GHCR `ghcr.io/tamnd/prog-runner/<lang>:<version>`, pinned by digest (20) | about 400 runner images (07.4) | 200 to 400 GB in the registry |
| run cache | R2 `programming-build` (Default), keyed by block hash × image digest | the outputs of `prog run` | 5 to 20 GB |
| build | `build/` in CI, cached by input hash | the graph, the matrix, the ladders, run results, Lean output, scene bundles | 10 to 30 GB, not committed |
| reports | `reports/` in git | coverage, quality, compat, the matrix, cost | small, JSON and Parquet |
| site | Worker static assets (hot set) and R2 `programming-site` (the rest) | the rendered site, datasets, search shards (16.20) | about 15 GB |

The content store is the only source of truth. Every other store can be rebuilt from it, the raw store and the image definitions. The raw store cannot be rebuilt (the web changes), so it has a monthly backup to a second bucket.

## 3.4 The packages of `prog`

| Package | Job | Spec |
|---|---|---|
| `node` | parse and write node directories, front matter, sections | 04 |
| `card` | read, write, shard, merge and promote cards; the matcher | 04.1, 13 |
| `schema` | JSON schemas of front matter, cards, compat entries, contracts, registers | 04 |
| `tag` | mint, check and retire tags (copied from algorithm-agent) | 04.3 |
| `graph` | edges, closures, the prerequisite DAG, paths | 04.6, 05 |
| `taxonomy` | the topic tree, the language template, the core list | 05 |
| `matrix` | the concept × language matrix and E11 | 04.8 |
| `markup` | Markdown with the block attributes of 7.2, cross-links, math (copied) | 06, 07 |
| `lang` | the language register, tiers, depth bands, style per language | 08 |
| `run` | run blocks in pinned OCI images with gVisor, no network, limits, normalizers and the cache | 07 |
| `images` | build, pin and test the runner images from their definitions | 07.4 |
| `compat` | compat entries, probes, imports (BCD, web-features), status-v1 | 09 |
| `release` | the release register, the window, freshness | 09.3, 09.7 |
| `testkit` | task contracts, shared cases, the security harness | 10 |
| `lean` | the `Prog` Lean build, decl-node bindings, `sorry` scan (copied) | 11 |
| `viz` | scenes, tracers, diagrams, frames (copied and extended) | 12 |
| `harvest` | the adapters, the raw store, the manifests | 13 |
| `cite` | the citation check against snapshots | 13.6 |
| `universe` | the universe register, outcomes, scope rules | 18 |
| `coverage` | levels C0 to C6, claims, the report | 18 |
| `status` | the ladders W, G, X, K, F, V and badges | 04.9 |
| `check` | the rule families (15.3) | 15 |
| `audit` | the audit and verifier prompts and their JSON | 15.4 |
| `copycheck` | similarity of prose and code against the snapshots (copied) | 15.3 C01, 19.4 |
| `agent` | the roles, the loop, models, budgets, run records | 14 |
| `queue` | work items from gaps, releases, reports | 14.8 |
| `prov` | provenance records (copied) | 04.7 |
| `site` | renderer, templates, assets, deploy, live check | 16 |
| `search` | the sharded index | 17 |
| `mcp` | the MCP server | 16.10 |
| `dataset` | the published datasets (compat, matrix, cards) | 9.9, 16 |

Copied packages start as copies of algorithm-agent and do not import it (20.2), as algorithm did with math2.

## 3.5 The flows

**F1. Harvest (scheduled per source, 13.2).** `prog harvest <source>` fetches through the adapter, writes raw snapshots, and writes cards, compat imports and universe rows. The matcher proposes outcomes. A pull request per source carries the change; the R rules gate it.

**F2. Release watch (daily).** `prog release watch` reads the release feeds. A new release creates its release card, a queue item for the release node (9.7), new image definitions for the runner, and a queue item to rerun the matrix on the new image. A release enters the window only when its image passes the smoke tests.

**F3. Write (each change).** A person or an agent changes a node, a set of cards, or a compat entry. `prog check -changed` runs the rules of the changed items, `prog run` for the changed blocks in the images their attributes name, the probes of changed compat entries, and the Lean build of changed modules. The pull request shows the ladders before and after.

**F4. Build (each push to `main`, and nightly).** `prog build` runs over the whole tree with the input-hash cache: the matrix (07.4), the probes, Lean, the scenes, the graph, the ladders, coverage and the datasets. The full version matrix runs only for changed blocks, changed images, and nodes going to W3 (Default; 21 risks).

**F5. Publish (after F4).** `prog site build` renders; `prog site check` checks links, sizes, accessibility and the file budget; `prog site deploy` uploads the changed files (manifest diff) to the Worker assets and R2, then runs the live check, with rollback on failure (16.20).

**F6. Agents (continuous, with budgets).** `prog queue` turns gaps, new releases and failed checks into items. `prog agent run` takes items, works in a sandbox clone, runs F3 until the checks pass, and opens a pull request per sub-area with its run record (14).

**F7. Maintenance (scheduled).** Link checks, image rebuilds for security fixes, the reruns that the window moves force, the review sampling, and the monthly backup of the raw store (14.9).

## 3.6 Cross-cutting decisions

- **Determinism.** Every step is a function of its inputs and the pins: `prog`, image digests, Lean, model ids and prompt versions. Run outputs are normalized (7.6) so they are stable.
- **Sandbox.** Every example, probe and test runs in its image under gVisor with no network, a read-only tree, a writable temp directory, and limits on time, memory, processes and output (7.5). Code from a pull request never runs outside it.
- **Incremental.** A change to a block runs that block. A change to an image reruns every block that uses it. A change to a task contract reruns all its implementations.
- **Cards are cheap, nodes are earned.** A card costs no model call when harvested. A card becomes a node only through `prog promote`, with a reason (demand, the core list, or a tier rule).
- **Machine copies.** Every page has `.md` and `.json` twins, the site has `llms.txt` files and an MCP server, and the data are published as datasets (16).
- **No secrets in content.** Deploy tokens live in CI secrets only; the content repository has none.

## 3.7 Deployment

| Part | Where | Cost |
|---|---|---|
| main site | Cloudflare Worker `programming` with static assets (hot set, at most 100,000 files on the paid plan) and R2 `programming-site` for the rest | Workers Paid plus R2 storage; check at M0 (16.20) |
| mirror | GitHub Pages, the core only, under 1 GB, or none (16.20) | free |
| CI | GitHub Actions on free runners (Linux x64 and arm64 for public repositories; Windows and macOS weekly, Q-20); a self-hosted build server for the nightly full build | the server exists |
| images | GHCR, public | free for public images |
| raw and cache | R2 `programming-raw`, `programming-build` | storage only; R2 has no egress fees |
| agents | model API spend under the budgets of 14.10 | 21 |
| playground | in the browser only, WebAssembly runtimes from the site (7.10); no server runs in v1 (Q-07) | none beyond hosting |
