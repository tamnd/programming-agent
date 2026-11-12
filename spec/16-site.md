# 16. The site: programming.tamnd.com

The site is where people and agents use the corpus. `prog site build` builds it from the content tree and `build/`. The renderer is the `site` package of algorithm-agent, copied into programming-agent at M0 (20.6). It keeps the following from algorithm: Go `html/template`, KaTeX at build time, Chroma highlighting at build time, the math2 visual design (one reading column, Inter, a hue for each area, light and dark themes, no rounded corners), pages that work with no script, and the commands `prog site serve` (local, on `:8151`), `prog site check` and `prog site deploy`. This document covers what programming adds: owned URLs, card collections, language tabs, version facts, the playground, the data and agent surface, and hosting at a size about 2 to 3 times that of algorithm.

Hosting evidence (checked 2026-11-09 at https://developers.cloudflare.com/workers/platform/limits/):

- Workers static assets take at most 20,000 files per Worker version on Free and 100,000 on Paid.
- Each asset file can be at most 25 MiB.
- A Worker script can be at most 64 MiB uncompressed.
- On Paid, CPU time per HTTP request defaults to 30 s and can be set up to 5 min.

Algorithm's spec (2150, 14.9) gives the other figures: GitHub Pages takes 1 GB per site, and R2 costs $0.015/GB-month with no egress fee.

## 16.1 Design principles

Algorithm's eight principles apply as they are (2150, 14.1): reading first, the answer at the top, one page with many depths, honest evidence, the same look for the same idea, no script needed to read, fast on a cheap phone, and stable links. Programming adds four:

1. **The reader's language first.** On every concept, task and pattern page, the language the reader chose comes first: in the tabs, in the comparison rows and in the examples. A reader who chose Go does not have to scroll past Python to find Go.
2. **Version is always visible.** Each fact box, example and compat badge says which versions it holds for and when it was verified. A reader never has to guess whether a page is old.
3. **Every fact is clickable to its evidence.** A badge links to its run, probe or source (6.6, 9.4). A citation links to the anchor it cites. An executed example links to its run record.
4. **Cards are first-class.** Most items are cards (04.1). A card has a URL, an anchor, a JSON record and a Markdown record, and search finds it. Only its prose page is missing, and the card says so with a "Help write this page" link (16.19).

## 16.2 URLs and reserved words

There are two namespaces (04.3):

- **Global slugs** sit at the root: `/<slug>/`.
- **Owned slugs** sit under their owner: `/<owner>/<slug>/`.

A URL never changes when an item moves in the topic tree or the language tree (05.1). It also never changes when a card is promoted to a node (04.1): the card's URL becomes the node page's URL.

| URL | Page | Kind |
|---|---|---|
| `/` | home (16.6) | |
| `/<slug>/` | the page of any global-slug kind: `/closures/`, `/observer-pattern/`, `/read-a-file-line-by-line/`, `/git/`, `/json/`, `/grace-hopper/` | concept, definition, task, pattern, technique, practice, tool, library, format, standard, weakness, work, person, organization, area |
| `/<lang>/` | the language page, which is the front page of the guide: `/python/`, `/rust/`, `/sql/` | language |
| `/<lang>/<area>/` | a sub-area page of the language template (05.3): `/python/control-flow/` | area |
| `/<owner>/<slug>/` | a feature page: `/python/match-statement/`, `/rust/question-mark-operator/`, `/react/hooks/` | feature |
| `/<owner>/<version>/` | a release page: `/python/3.15/`, `/java/27/`, `/cpp/23/`, `/ecmascript/2026/`, `/node/26/` | release |
| `/<owner>/api/` and `/<owner>/api/<group>/` | the api index, and an api collection page (one module, package, class or header): `/python/api/str/` | collection of api cards |
| `/<owner>/api/<path>/` | an api card: `/python/api/str.split/`. A node when promoted, otherwise a redirect to its anchor on the collection page (16.7) | api |
| `/<owner>/errors/` and `/<owner>/errors/<code>/` | the diagnostics of a language or a compiler: `/rust/errors/E0382/`, `/typescript/errors/TS2339/`, `/python/errors/KeyError/` | diagnostic |
| `/<tool>/lints/` and `/<tool>/lints/<code>/` | the lint rules of a linter: `/ruff/lints/E501/`, `/clippy/lints/needless_range_loop/` | diagnostic |
| `/<owner>/proposals/` and `/<owner>/proposals/<id>/` | proposals: `/python/proposals/pep-634/`, `/java/proposals/jep-444/`, `/cpp/proposals/p2300/` | proposal |
| `/exercises/<site>/<id>/` | an exercise card | exercise |
| `/<slug>/<lang>/` | for a concept, task or pattern: a redirect to `/<slug>/?lang=<lang>#<lang>`, the tab of that language (16.4). `/closures/rust/` and `/read-a-file-line-by-line/go/` are examples | redirect |
| `/cwe-89/`, `/pep-634/`, `/E0382/` | well-known ids as aliases: a redirect to the item (04.3). An ambiguous id (`E0001` exists in several tools) goes to a disambiguation page | redirect |
| `/t/<tag>` | the item by tag, node or card | redirect |
| `/t/<tag>@<release>` | the item in a site release (16.21) | redirect |
| `/compare/` and `/compare/<a>-vs-<b>/` | comparisons (16.5): `/compare/python-vs-go/` | comparison |
| `/compat/`, `/compat/<lang>/` and `/compat/<lang>/<area>/` | the compat tables (16.7): `/compat/python/` | compat |
| `/tasks/` | the task index, by group (10.6) | index |
| `/concepts/`, `/patterns/`, `/languages/`, `/tools/`, `/libraries/`, `/formats/`, `/standards/`, `/weaknesses/`, `/people/`, `/works/` | the index of each kind | index |
| `/languages/all/` and `/languages/family/<family>/` | the catalogue of all languages, and the family pages (08.7) | collection of language cards |
| `/paths/` and `/paths/<path>/` | reading and learning paths (16.12) | path |
| `/play/` | the full-window playground (16.8) | playground |
| `/search/` | search (17) | |
| `/coverage/`, `/matrix/`, `/gaps/`, `/universe/`, `/latest/` | the dashboards (16.14) | |
| `/about/`, `/about/status/`, `/about/badges/`, `/contribute/`, `/spec/`, `/sources/` | the project, the computed status definitions (9.5), the badge meanings, and the source register | |
| `<any page path>.md`, `<any page path>.json` | the Markdown twin and the JSON record (16.10). For example `/closures.md`, `/python/match-statement.json` and `/python/api/str.split.json` | data |
| `/llms.txt`, `/<lang>/llms.txt`, `/<area>/llms.txt`, `/data/`, `/mcp`, `/.well-known/mcp`, `/sitemap.xml`, `/robots.txt` | the machine surface (16.10) | data |

**How a path is read.** The Worker and the build use the same function (`site.Resolve`, with golden tests):

1. If the first segment is a reserved word, it is a site section.
2. If the first segment is an owner (a language, tool, library, standard or organization with owned items), the second segment is a sub-area id, a release version, `api`, `errors`, `lints`, `proposals`, or an owned slug, in that order.
3. Otherwise the first segment is a global slug. A second segment must be a language id (a tab redirect), or the path is a 404.

**Reserved words.** `registers/reserved.yaml` holds the list. Rule T04 rejects a global slug on the list. It also rejects an owned slug that is a reserved second segment or a language id.

```
about all api areas assets badges compare compat concepts contribute coverage data
embed errors exercises feed formats gaps help index kind languages latest libraries
lints llms mcp matrix patterns paths people play proposals random releases search
sitemap sources spec static status t tasks tools universe weaknesses works x _cards 404
```

**Canonical form.**

- Trailing slashes are canonical, and lower case is canonical for slugs.
- Codes keep the case their owner uses (`E0382`, `TS2339`, `KeyError`). Any other case redirects to the canonical case through the code index of the owner.
- Versions keep the owner's form (`3.15`, `27`, `2026`).
- An alias in `registers/aliases.reg` redirects with 301.
- `/rust/errors/e0382/` redirects to `/rust/errors/E0382/`.

## 16.3 Navigation

- **Header.** The site name, the search field (`/` or `⌘K`), the language picker (16.4), the theme switch.
- **Breadcrumbs.**
  - A topic tree page shows the area path: Functions › Closures and scope › Closure.
  - A language tree page shows the guide path: Python › Control flow › The match statement.
  - A card shows its collection: Python › API › str › split.
- **Sidebar** on wide screens:
  - On a topic page, the sub-area list in reading order (05.5).
  - On a language tree page, the 25 template sub-areas of the guide with their counts.
  - It folds under the head on narrow screens.
- **Previous and next** in the reading order of the sub-area, or of the path the reader follows (16.12).
- **"In other languages" chips** on a feature page: the same concept in other languages (from `realizes`). One click goes to the feature page of that language when it exists, otherwise to the tab of the concept page.
- **"Elsewhere" links** to algorithm and math2 (5.8). They have a small label that names the other site, so the reader knows they are leaving.

## 16.4 Page templates

Each kind has one template. Each template follows the sections of 6.3 and adds the parts the build computes. The parts are the same across kinds, so the reader learns them once.

**The facts box** is at the top of every node page, in the first screen. It is computed from the front matter, the registers and the ladders. It is never written by hand (T05).

| Kind | Facts box lines |
|---|---|
| concept | kind and level; the matrix row for tier 1, as chips (native, stdlib, library, idiom, absent, n/a); the prereqs; "Formal model: proved / executable / none" (11.1) |
| language | class, paradigms, typing, memory, execution, appeared, designers, latest release with its date, the support window, implementations, file extensions, popularity with its date (04.5) |
| feature | the owner, the synopsis in one line, the computed status (9.5) with the version where it was added, the proposals, the concepts it realizes |
| release | the release date, the end of life, the status, the counts of added, changed, deprecated and removed items |
| task | the contract type and signature, the number of languages that pass, the edge cases count |
| pattern | the family, the polarity, the intent in one line, the languages where it is built in |
| tool, library | role, licence (SPDX), latest version with its date, languages, repository |
| weakness | CWE id, abstraction, OWASP Top 10 entry, CWE Top 25 rank, languages affected |
| card | the card's fields, and "No page yet. Help write it." (16.19) |

**The verified line.** Each facts box ends with a line such as "Verified on 2026-11-20 against Python 3.11 to 3.15, CPython and PyPy". The date is the last run of the examples (07) and probes (09.4) on the current images. A click opens the run record: the image digests, the commands, and the result of each block. When the run is older than 30 days, or the language has a release that the page has not been run on, the line shows "Not yet verified against <version>" in the warning colour (Default).

**The badges.** The ladders of 04.9 are shown as small badges: W for writing, G for grounding, X for execution, K for compat, F for formal, V for visualization.

- Each badge shows the level digit and has a tooltip that explains it.
- Each badge links to its evidence:
  - W links to the audit and reviews.
  - G links to the citation check report.
  - X links to the run record, by language.
  - K links to the probes.
  - F links to the Lean declarations.
  - V links to the scene.
- A badge appears only when its ladder applies to the kind (F only for `formal: expected` or `optional`).
- Text that an agent wrote and no person reviewed carries the label "Written by an agent, checked by the build, not yet reviewed by a person", as in algorithm (2150, 14.4).

**The language switcher.** Concept, task and pattern pages have one tab per language section (6.3).

- The tab row is ordered as follows: the reader's chosen language first, then the reader's other recent languages, then tier 1 in register order, then tier 2. A "More" menu holds the rest.
- The choice is kept in `localStorage` under `pg-lang`, with no cookie. The `?lang=` query sets it for one visit, which is how `/<slug>/<lang>/` works.
- With no script, all sections are stacked, in register order, each with its heading. The anchor `#<lang>` works in both modes.
- A tab shows the matrix verdict before the code ("Native since Java 8", "Library: `regex` crate", "Absent: reason").
- A cell that is `absent` or `n/a` has a tab with one line and no code. The tab is not hidden, because "absent" is an answer.
- The header picker sets the same preference. Feature pages are already in one language. They use the preference only for the comparison section and the "In other languages" chips.

**Code blocks** (from 07):

- A header row shows the language, the version range of the block (`since` and `until`), and the result of the last run for each version as a compact strip of dots. A click on a dot opens the run record.
- The **Copy** button copies the prelude, the wrapper and the block, which is the full program the build ran (X06).
- **Run** appears when the block is marked `browser` (16.8). **Open in** appears for languages with an official playground (07.10).
- Output blocks are joined to their code block, shown in the muted surface, and labelled "Output (recorded on <date>, <version>)".
- `fails` blocks have a red left bar and the label "This code does not compile" or "This code raises `<Type>`", with the diagnostic linked (X05). `wrong` blocks have a red bar and "This code runs but gives a wrong result".
- `norun` blocks show their reason in words (07.8).

**Version-dependent text.** `{since:...}` renders as a small version chip that links to the compat entry. `> [!since v="3.10"]` callouts render with the same chip. On a page with side-by-side behaviour (6.7), the versions are tabs inside the example, labelled by version.

## 16.5 Comparison pages

Comparison pages are generated from facts that the nodes already hold (08.8). They never contain prose that disagrees with a node.

| Page | Built from | Content |
|---|---|---|
| `/compare/<a>-vs-<b>/` | the two language nodes, their `concepts.yaml` rows, the concept examples, the tasks, the `compare/<a>-to-<b>` and `compare/<b>-to-<a>` nodes | the facts boxes side by side; the "coming from" notes for each direction; a table of every core concept with both cells, filtered by default to the cells that differ; for each differing concept, the two examples side by side; a selection of 20 tasks side by side (the core tasks with the most search demand, 17); the links to both guides |
| `/compare/` | the matrix and the examples dataset, loaded as JSON in the browser | a picker for up to 6 languages and any concepts, tasks or areas; a table with one row per concept or task and one column per language; each cell shows the verdict and, when opened, the example with its output. With no script, it shows the list of generated pair pages |
| `/compare/<a>-vs-<b>/<concept>/` | none (redirect) | `/<concept>/?lang=<a>,<b>`, the concept page with both tabs side by side |

Rules:

- Pair pages are written to files only for the pairs in `registers/compare-pairs.yaml` (about 200 at M3, Estimate). Any other pair is served from `/compare/?langs=a,b` in the browser.
- The order of the two names in a URL is the order in `registers/compare-pairs.yaml` (`/compare/python-vs-go/`). The other order (`/compare/go-vs-python/`) redirects to it. The page shows the reader's chosen language on the left, whatever the URL order.
- No ranking, no scores, no "winner". Benchmarks appear only from benchmark nodes (07.9).
- The page title is "Python and Go compared", not "Python vs Go: which is better".

## 16.6 The home page and index pages

- **The home page** has:
  - The search field, at the top.
  - "Pick your language": the tier-1 and tier-2 languages as tiles. Each tile goes to the guide and sets the language preference.
  - The six groups of the topic tree (05.2) as tiles in the area hues, each with its top areas.
  - "Common tasks": the 12 core tasks with the most demand.
  - "New in languages": the newest release nodes, with their dates.
  - "Recently verified": the latest pages to reach X3 or K2.
  - The coverage summary, as a single line with the link to `/coverage/`.
- **Area pages** list sub-areas and items in reading order, with the kind word, the level and the ladder digits, as in algorithm (2150, 14.7.2). Cards in an area appear in a separate "Not yet written" list, so a reader sees the whole area and its gaps.
- **Kind indexes** (`/concepts/`, `/patterns/`, ...) can be filtered by area, level, language and coverage. With no script they are sorted by area.

## 16.7 Language guides, releases, compat and card collections

**The language page** (`/python/`) is the front page of the guide (6.3, language). Its "The guide" section is the 25 template sub-areas (05.3) as a grid. Each cell shows the sub-area name, the count of features, and a small coverage bar. The release timeline (12.6) and the support window come next.

**Release pages** (`/python/3.15/`) list what is new, changed, deprecated and removed, from compat. Each line links to its feature. Release cards that are not yet nodes show on `/python/releases/`, which holds the full table of every release since the first.

**Compat pages.**

- `/compat/<lang>/` has one table per sub-area. Each row is a feature, with one column per implementation (9.2). The cells show `version_added` and are coloured by computed status, with a word in each cell as well as the colour. Notes and flags are in footnotes.
- Each cell links to its probe and its last result.
- Web features (HTML, CSS, JavaScript web APIs) show the BCD browser columns and the Baseline status, with credit to MDN and the W3C WebDX CG (9.6).
- The same tables appear in the Versions section of each feature page.

**Card collection pages** hold cards of one kind with one owner, grouped:

| Collection | Grouping | Page size |
|---|---|---|
| api | one page per module, package, class, header or interface; huge modules split by first letter | at most 200 cards per page |
| diagnostics | one page per owner and code prefix (`/rust/errors/` lists E0001 to E0799 in pages of 200) | 200 |
| lints | one page per linter and rule category | 200 |
| proposals | one page per owner and status (final, accepted, draft, rejected, withdrawn) | 200 |
| language catalogue | one page per family, and `/languages/all/` paged by first letter | 300 |
| exercises | one page per site and per concept | 200 |

Each card on a collection page is a block with an anchor (`#str.split`, `#E0382`). The block holds the card's fields, its one-sentence summary, its links, its sources, and the "Help write this page" link. A card's own URL (`/python/api/str.split/`) redirects with 302 to `/python/api/str/#str.split`. The redirect is 302, not 301, because promotion turns the URL into a real page later (Default).

**Card redirects at scale.** There are hundreds of thousands of cards (Estimate, 02), too many for one redirect map. The build writes one small index per collection, `/_cards/<owner>/<kind>/<group>.json`, mapping each card path to its collection page and anchor. The Worker reads the index for the path's owner and kind from R2, caches it at the edge, and answers the redirect. The same index answers `<card path>.json` (the record) and `<card path>.md` (16.10).

## 16.8 The playground

The playground runs examples in the reader's browser for the languages of 07.10, and in nothing else.

**Inline run.** A block marked `browser` has a **Run** button and an **Edit** button.

- **Edit** turns the block into an editor: a plain `<textarea>` with highlighting from the same Chroma themes, line numbers and tab handling. There is no heavy editor framework (Default). The editor is loaded on first use.
- **Run** starts the language's runtime in a Web Worker on the first click, never on page load. Then it runs the code, with the prelude and wrapper, and shows stdout, stderr and the exit status under the recorded output.
- A size notice appears before the first download of a large runtime ("Loads Python (6 MB) once"). After that the browser caches it, because runtime files have hashed names and an immutable cache header.
- When the reader has not edited the code and the output differs from the recorded output, the notice of X11 appears (07.10).
- **Reset** restores the page's code. Edits are kept per block in `sessionStorage` only.
- A run is limited to 10 s of wall time by default, after which the Worker is ended. Memory is bounded by the runtime.

**The full-window playground** (`/play/`) has a language picker, a version label (the version of the browser runtime, which can differ from the CI images), a stdin field, and a share action. Share puts the code in the URL fragment (`#code=<compressed>`). Fragments do not reach the server, so shared code never touches our servers (07.10).

**Open in.** For languages with no browser runtime, the block has "Open in <playground>": go.dev/play, play.rust-lang.org, play.kotlinlang.org, Compiler Explorer, SwiftFiddle, and others in `registers/playgrounds.yaml`. The link passes the code in the URL where the playground supports it. A notice says the code goes to a third party. The link is a plain link, so it works with no script.

**Cross-origin isolation.** Some runtimes need `SharedArrayBuffer` (threads in Pyodide builds, Runno's clang, some wasm builds with pthreads), and that needs cross-origin isolation.

- The Worker sends `Cross-Origin-Opener-Policy: same-origin` and `Cross-Origin-Embedder-Policy: require-corp` **only** on `/play/` and on the runtime files under `/static/runtimes/`. It never sends them on content pages, because COEP would block cross-origin images and embeds that a page may have, and would break "Open in" pop-ups.
- Inline runs on content pages use the runtimes that work without isolation. A runtime that needs isolation opens the block in `/play/` instead, with the code passed in the fragment.
- Runtime files are served from R2 by the Worker, with `Cross-Origin-Resource-Policy: same-origin` and the correct `Content-Type: application/wasm`, so streaming compilation works.

**Security.** The playground runs only in the reader's browser, in a Web Worker or a sandboxed iframe with no access to the page's origin storage (`sandbox="allow-scripts"`, no `allow-same-origin`). The Content Security Policy permits `'wasm-unsafe-eval'` on `/play/` and on pages with runnable blocks, and no `unsafe-eval` anywhere. A server-side run for compiled languages is not in v1 (07.10, 22 Q-07).

## 16.9 Visual design, interaction and mobile

These are the same as algorithm (2150, 14.10 to 14.12), with these changes:

- **Hues.** One hue per top area of the topic tree (36 areas, 05.2). Language trees use a neutral hue, and each language has a small identity colour taken from GitHub Linguist's `color` (MIT). It is used only in tab underlines and chips, never as the only signal. There are no logos (19.6).
- **Code highlighting** uses Chroma at build time, and the build fails if a language has no Chroma lexer. Shiki is the fallback through a build-time Node step only for languages Chroma lacks (Default). Rare tier-3 and tier-4 languages are shown with no highlighting.
- **Keyboard.** Algorithm's shortcuts, plus:
  - `l` and `Shift L` move to the next and previous language tab.
  - `r` runs the focused block.
  - `e` edits the focused block.
  - `g c` goes to `/compare/`.
  - `g t` goes to `/tasks/`.
  - `v` cycles the version tabs of an example.
- **Preferences** are stored in `localStorage` under keys that start with `pg-`: theme, language, recent languages, version display (the newest version only, or all), shortcuts on or off, followed path, and progress.
- **Mobile.**
  - The language tab row scrolls sideways, with the chosen tab in view.
  - Compat tables keep the feature column fixed and scroll the implementation columns.
  - The playground on narrow screens is a stacked editor, then output, with large touch targets (at least 44 × 44 px).

## 16.10 The data and agent surface

Agents are a main audience, not an afterthought (01.7). Every page and card has a machine form. The site publishes its data as datasets, and an MCP server answers questions with the same data.

**Markdown twins.**

- Every page has a twin at `<path>.md` (`/closures.md`, `/python/match-statement.md`). The build writes it from the source, not by converting the HTML. It has front matter (tag, kind, URL, versions verified, ladders, licence) and the body with all language sections. Executed examples carry their recorded outputs and version ranges.
- Card twins are generated from the card record by one fixed function. The build and the Worker share that function's golden tests, so the Worker's output for a card is byte-identical to what the build would write.
- **Content negotiation.** A request with `Accept: text/markdown` that prefers Markdown over HTML gets the twin with `Content-Type: text/markdown; charset=utf-8` and `Vary: Accept`. We serve our own twins rather than Cloudflare's "Markdown for Agents" conversion (beta, Feb 2026), because our twins come from the source and keep the code exact (Default).
- Twins are `X-Robots-Tag: noindex` and carry a `Link: <html url>; rel="canonical"` header.

**JSON records.** `<path>.json` is the record of any item, node or card. It holds:

- the front matter fields;
- the computed ladders and status;
- the edges in both directions;
- the compat entries;
- the examples, as references to the examples dataset;
- the sources.

The schema is `schema/record.v1.json`, published at `/data/schema/`. A record never contains the prose. The prose is in the twin, because the two have different licences (19.2).

**llms.txt.**

- `/llms.txt` follows the llmstxt.org format: an H1, a summary blockquote, then H2 lists of links, ending with an "Optional" section. It lists the top areas, the language guides, the indexes, the datasets and the MCP endpoint.
- Each top area and each language has its own `/<area>/llms.txt` and `/<lang>/llms.txt`, listing that subtree's twins.
- There is no single `llms-full.txt`. At this size it would be several GB. Instead, each language and area has `/<lang>/llms-full.txt` with the twins of its core pages, at most 5 MB each (Default).

**Datasets** at `/data/`. Each is versioned, has a changelog, and is mirrored to Hugging Face as the dataset repository `tamnd/programming`. Each monthly site release (16.21) is a tag.

| Dataset | Content | Licence | Format |
|---|---|---|---|
| `items` | every node and card record | CC0 (metadata) | JSONL, sharded by kind |
| `examples` | every code block with its recorded outputs and run status per version (07.11) | MIT (code), CC0 (metadata) | JSONL |
| `compat` | the compat data per language, in BCD's shape (9.9) | CC0 | JSON per language |
| `matrix` | the concept × language matrix with verdicts and evidence links (4.8) | CC0 | CSV and JSON |
| `tasks` | contracts and cases (10.2) | MIT | JSONL |
| `pages` | the Markdown twins of all nodes | CC BY-SA 4.0 (prose), MIT (code blocks) | JSONL, sharded by area |
| `graph` | all edges | CC0 | JSONL and GraphML |

**The MCP server** is at `https://programming.tamnd.com/mcp`, with discovery at `/.well-known/mcp`.

- It implements MCP specification version 2026-07-28: protocol version per request in `_meta`, and the `server/discover` method. It also accepts the 2025-11-25 handshake for older clients (Default).
- It runs in the `programming` Worker, is read-only, needs no account, and is rate limited per IP (Default: 60 requests per minute).

| Tool | Input | Output |
|---|---|---|
| `search` | a query; optional kind, language, version, area, limit | ranked items with URL, kind, summary and ladders, from the search index of 17 |
| `get_item` | a URL, a slug, a tag, or a well-known id (`cwe-89`, `pep-634`, `E0382`); optional `format` of `markdown` or `json`; optional language filter | the twin or the record; a language filter returns only that language's section of a concept or task page |
| `get_example` | an item and a language, or an example id; optional version | the code as the reader would copy it, the recorded output, the version range and the run status per version |
| `compat_lookup` | a feature, api or diagnostic, and optionally a language, implementation and version | whether it is supported, `version_added` and `removed`, flags, status, and the probe result |
| `run_example` | code and a language | **only if the server-run option of 22 Q-07 is chosen**; otherwise it is not listed |

Each tool result includes the item URL and the "verified on" date, so an agent can cite the page and say how fresh the fact is.

## 16.11 Search

Search is document 17. The site part:

- the palette (`/`, `⌘K`);
- the `/search/` page;
- static search shards in the Worker assets, loaded on demand;
- the "paste an error" field on `/search/` and on each `/<owner>/errors/` page, which matches the message against diagnostic templates (17).

## 16.12 Reading paths and learning paths

| Path type | URL | How it is made |
|---|---|---|
| **Learn L** | `/paths/learn-python/` | computed from the language tree: the `intro` and `core` features in template order, then prereqs order (08.9) |
| **L for programmers of M** | `/paths/python-for-go-programmers/` | computed: the features of L whose matrix cells differ from M's, ordered by the template, with each step linking the `compare/m-to-l` note; generated for every pair of tier-1 languages and for registered tier-2 pairs |
| **Topic paths** | `/paths/concurrency/` | computed from the topic tree: an area's items at `intro` and `core`, in prereqs order, with the reader's language chosen for every example |
| **Curriculum paths** | `/paths/cs2023-sdf/`, `/paths/swebok-construction/` | from the crosswalks (5.7): the items mapped to each knowledge unit of CS2023 or each topic of SWEBOK v4, in the outside system's order |
| **Edited paths** | `/paths/<slug>/` | written by an editor as a `path` file in `paths/`: a list of items with a sentence before each; reviewed like a node |

A path page shows its steps and the reader's progress, kept in `localStorage`, with no account. The next-step link appears on each page of the followed path. With no script, a path page is an ordered list. Each step shows the ladder digits, so a reader sees which steps are complete.

## 16.13 Performance budgets

`prog site check` measures these budgets on a fixed sample of 100 pages: the 30 largest concept and task pages, 20 language guide pages, 20 collection pages, and 30 random pages. A page over budget fails the build.

| Budget | Limit |
|---|---|
| HTML of a node page, gzip | ≤ 80 KB with all tier-1 tabs in the HTML; tier-2 and later sections are in the HTML only as `<details>` content (Estimate; algorithm's limit is 60 KB) |
| HTML of a collection page, gzip | ≤ 120 KB (200 cards) |
| HTML of any page, uncompressed | ≤ 800 KB |
| CSS, gzip | ≤ 20 KB, one file |
| site script, gzip | ≤ 18 KB, deferred (tabs, copy, palette loader, theme, version tabs) |
| editor and playground shell, gzip | ≤ 25 KB, loaded on first Edit or Run |
| language runtimes | never before a Run click; size shown before the first download |
| fonts before first paint | Inter text subset only, ≤ 110 KB; code uses the system monospace stack |
| Largest Contentful Paint | ≤ 1.8 s on a mid-range phone on Fast 4G (Lighthouse) |
| Cumulative Layout Shift | ≤ 0.05; tab rows and output blocks reserve their height |
| Interaction to Next Paint | ≤ 200 ms, including a tab switch |
| Lighthouse performance score | ≥ 95 on the sample |

## 16.14 Accessibility, SEO, dashboards

**Accessibility.**

- The target is WCAG 2.2 level AA on every page, with algorithm's rules (2150, 14.14), and `prog site check` runs axe-core on the sample (rule W08).
- Language tabs and version tabs follow the ARIA tab pattern.
- Compat cells always have a word, never colour alone.
- Output blocks are labelled as output for screen readers.
- `fails` and `wrong` blocks announce their label before the code.
- The playground editor is a real `<textarea>`, so assistive technology works. Run results go to an `aria-live="polite"` region.
- A person tests the core page types with VoiceOver and NVDA at each milestone.

**SEO.**

- Titles have these forms:
  - "Closures: concept, examples in 15 languages · programming"
  - "The match statement in Python · programming"
  - "Rust error E0382: borrow of moved value · programming"
- Structured data:
  - `TechArticle` for nodes.
  - `SoftwareSourceCode` for task implementations and examples, with `programmingLanguage` and `runtimePlatform`.
  - `DefinedTerm` for definitions and diagnostics.
  - `ComputerLanguage` for languages.
  - `Dataset` for `/data/`.
- Sitemaps form an index, split by kind and by language, each under 50,000 URLs. They list node pages and collection pages, never card redirect URLs.
- Social images are made only for node pages (not cards), at 1200 × 630.

**Dashboards.**

- `/coverage/` shows the C0 to C6 levels per area and per language (18).
- `/matrix/` shows the concept × language matrix with verdicts and X levels.
- `/gaps/` is the ranked gap list with "Help write it" links.
- `/universe/` is the searchable universe register.
- `/latest/` lists new and changed pages, new releases covered, and new verified levels, with an Atom feed.
- Each dashboard has a JSON copy under `/data/`.

## 16.15 Languages of the site and analytics

**English first.** The site is English only through M4. Aliases in other languages are in front matter (04.5), so search finds pages by them. After M4, translations come in as a layer `i18n/<lang>/nodes/...` with translated prose only. Code, outputs, compat and cards are shared. Translated pages live at `/<lang-code>/<path>/` with `hreflang` links and a "Translated from the English page of <date>" line. The site language code uses a `~` prefix (`/~ja/closures/`), so it can never clash with a programming language id such as `/d/` or `/r/` (Default; 22).

**Analytics.**

- No cookies, no trackers, no third-party scripts.
- Counts come from Cloudflare's own request logs (Workers Analytics Engine, written by the Worker, with the page path and the country, and no IP address and no user agent stored).
- Search terms are logged without any client data. Terms with no result feed the gap list (18.6).
- The site has no cookie banner, because it sets no cookies.

## 16.16 Rendering

- The renderer is the `site` package in Go, with `html/template`, writing plain HTML, as in algorithm and math2. There is no client framework.
- Markdown goes through goldmark with the `markup` extensions of 6.2, plus the code block attributes of 7.2, the `{since:}`, `{std:}` and `{rule:}` references, the `:::rule` inference-rule blocks of 11.4, and the owned-link forms.
- KaTeX renders to MathML at build time through goja, as in math2 and algorithm.
- Each template has a golden test, with one fixed node directory per kind.
- The function that resolves a path to its page (16.2) and the function that renders a card twin are shared Go code. The Worker's copies are tested against them: either the Go is compiled to wasm for the Worker, or a TypeScript port passes the same golden files (Default: the TypeScript port, for Worker size and start time; check at M0).

## 16.17 Build

| Stage | Input | Output | Cache key |
|---|---|---|---|
| load | content tree, registers, cards | the in-memory graph, with reverse edges | none (always) |
| ladders | graph, run records, probe records, audits, reviews | the computed levels and statuses | hash of inputs per item |
| render | graph, templates | HTML, twins, records | hash of the item's files, the titles and ladders of items it links, the template version |
| assets | CSS, JS, fonts, scene bundles | hashed files | content hash |
| search | graph | search shards (17) | hash of the indexed fields |
| datasets | graph, run records | `/data/` files | hash of inputs |
| social | node pages | JPEGs | hash of title, kind, hue |

**Incremental builds.**

- The build keeps a dependency index: which outputs depend on which inputs.
- A change to one node re-renders that node, the pages that link it (for titles and ladder digits), its area page, its collection page if it was a card, the comparison pages that use its examples, and the affected search shards.
- A change to a template or the CSS re-renders every page of that template.
- A change to `registers/languages.yaml` re-renders every page that shows that language.

**Time budget** (Estimate, on the build server, 16 cores):

| Build | Limit |
|---|---|
| a one-node change | ≤ 3 min from merge to an uploaded diff |
| a template change | ≤ 20 min |
| a full build from an empty cache | ≤ 60 min (algorithm's budget is 30 min at about 175,000 files; this site is about 2 to 3 times larger) |
| a full example run of all languages and versions | separate, nightly, sharded (07.4); never part of the site build |

`prog site check` runs after every build:

- broken internal links;
- reserved words;
- budgets;
- axe-core;
- the hot-set file count against the assets limit;
- file size against 25 MiB;
- the twin and record for every page;
- the card index for every card.

## 16.18 Security headers

These are the same as algorithm (2150, 14.9). Hashed assets have `Cache-Control: public, max-age=31536000, immutable`. HTML has `max-age=300, stale-while-revalidate=86400`. Responses carry a strict Content Security Policy (scripts from the site only, no inline script except the theme line by hash), `X-Content-Type-Options: nosniff` and `Referrer-Policy: strict-origin-when-cross-origin`, and the site makes no third-party requests. Changes for this site:

- COOP and COEP on the playground routes only (16.8).
- `'wasm-unsafe-eval'` where runnable blocks exist.
- `Access-Control-Allow-Origin: *` on `/data/`, `.json` and `.md` paths, so tools in other origins can read the data.

## 16.19 Help write this page

Every card and every node below W5 has a "Help write this page" or "Improve this page" link. The link opens one of three things, chosen by the reader:

| Action | Opens | Prefilled with |
|---|---|---|
| **Write this page** (cards) | a GitHub issue form `write-page.yml` in `tamnd/programming` | the card's tag, URL, kind, owner and sources, and the `prog promote <tag>` command for people who want to send a PR |
| **Edit on GitHub** (nodes) | the `node.md` or the `examples/<lang>.md` file of the open tab, in GitHub's editor | none; the DCO note is shown in the PR template (19.4) |
| **Report a problem** (any) | an issue form `report.yml` | the URL, the tag, the commit of the build, the section the reader was reading (from the nearest heading), the reader's language and the versions shown; the reader picks a type: wrong fact, wrong output, out of date, broken example, missing language, unclear text |

Reports get the label `reader-report`. The work queue (14) takes them in as gap items with a high priority, and the verifier agent checks a "wrong fact" report before any writer changes the page. The links are plain links to GitHub, with no script and no account on our side.

## 16.20 Hosting

**Decision.** The site starts on Cloudflare Workers with static assets and R2 from M0. It does not start on Cloudflare Pages.

**State at M0.** The Worker `programming` serves the static assets and has no script yet. The host reads `_redirects` and `_headers` from the asset tree, adds the trailing slash to a directory URL and serves `404.html` for a missing path. `prog site deploy` uploads a new version and makes it live on the custom domain. With `-preview`, it uploads a version with the preview alias `pr-<n>` and does not change production. The limits at this stage are 20,000 files for each version and 100 `_redirects` rules with a placeholder, and `prog site check` fails before the deploy when the tree breaks one. The script of the Worker and the R2 bucket come when the hot set passes these limits.

**Record packs.** The JSON records of the cards of a collection are one JSON array. The site does not write a file for each collection. `prog site build` puts the arrays into packs in `/api/records/`, in the order of the collection URLs, so that the collections of one owner share their packs. A pack is a file of JSON lines, with one line for each collection, and it closes at about 2 MiB. A collection that is larger has a pack of its own. The name of a pack has the hash of its bytes. `/api/records.json` gives the list of the packs and, for the URL of each collection, the number of its pack, the byte offset of its line and the length of the line. The link "JSON record" of a card opens `/api/record/#<collection URL>`, and the script of that page reads the line. It asks for the range of the line first. The static assets of the Worker do not answer a range request (checked 2026-11-12: the host sent the whole file with status 200), so the browser gets the whole pack, keeps it in the Cache API and cuts the line from it. The next record of the same pack needs no request. A host that answers a range request, such as a server mirror or R2, sends only the line. The 404 page sends a request for an old `<collection>/cards.json` to the record page. The packs are the reason why a collection adds one file to the tree and not two: at 2,688 collections, the packs took the tree from 8,809 files to about 6,200.

**Caching.** Algorithm put its packs behind a `_worker.js` and then removed it, because each request to a worker counts against the daily limit of the free plan, and the packed files fail when the limit is used up. This site has no worker script, so each request is a request for a static asset, and the cache is in the browser:

- A pack is immutable for a year (`_headers`), because its name changes when its bytes change.
- The links to the CSS and the scripts have the hash of the file in the query (`?v=`), so that a new deployment loads the new file at once and the old one is never used with a new page.
- `search.json` is kept in `localStorage` with the commit of the deployment (`data-build` on `<html>`). The search opens with no request until the next deployment.
- The service worker `/sw.js` gets a page from the network first, so that a reader sees a new deployment at once. When the network fails or is slower than 3 seconds, it uses the copy of the page that the reader opened last time, and it keeps 200 pages. An asset with a hash in its URL, a font or a KaTeX file comes from the cache first. Another asset comes from the cache and is fetched again in the background.
- The script fetches a page when the pointer rests on its link for 65 ms or a finger touches it, so that the next page opens at once. It does not do this when the browser asks to save data.

**Build server.** A cron job on a server of the project runs `scripts/deploy-cloudflare.sh` of the content repository every 5 minutes, as math2 does. The script fetches main, and it does nothing when main is the same as at the last deploy. Otherwise it gets the `prog` release that `.github/workflows/ci.yml` pins, builds the site, deploys it and removes the tree. So a merge reaches the site in minutes, and it does not wait for the test jobs of the CI.

- Algorithm started on Pages and passed the 20,000-file limit at about 3,500 written nodes. It then had to turn language and proof pages into redirects (2150, 14.9).
- This site has more items than algorithm's whole plan by M1, and card collections and twins from day one.

**The page and file budget** at the scale of 02 (all Estimate; recompute at M1 from the universe):

| Files | Count | Size | Where |
|---|---|---|---|
| node pages (concepts, features, tasks, patterns, tools, ...) | 60,000 | 3.0 GB | Worker assets for the core and the tier-1 guides (about 35,000); R2 for the rest |
| area pages, language template sub-area pages, path pages | 4,000 | 80 MB | Worker assets |
| card collection pages | 5,000 | 500 MB | Worker assets for the collections of tier-1 owners (about 2,000); R2 for the rest |
| card indexes (`/_cards/...`) | 5,000 | 700 MB | R2 |
| node twins (`.md`) and records (`.json`) | 120,000 | 1.5 GB | R2 |
| comparison pages (registered pairs) | 250 | 25 MB | Worker assets |
| compat pages | 1,500 | 150 MB | Worker assets for tier 1; R2 for the rest |
| `llms.txt` and `llms-full.txt` files | 400 | 1 GB | R2 |
| scene bundles and figures | 20,000 | 300 MB | Worker assets for bundles; R2 for frames |
| social images | 60,000 | 3.0 GB | R2 |
| browser runtimes (Pyodide, ruby.wasm, sqlite-wasm, PGlite, wasmoon, wasm_of_ocaml, ...) | 500 | 800 MB | R2 |
| datasets (`/data/`) | 100 | 3 GB | R2, and Hugging Face |
| search shards, CSS, JS, fonts, icons | 6,000 | 150 MB | Worker assets |
| **total** | **about 280,000** | **about 14 GB** | |

**The outlets.**

| Outlet | Files | Serves |
|---|---|---|
| **Workers static assets**, Worker `programming` on `programming.tamnd.com` (Workers Paid, $5/month) | the **hot set**, at most 95,000 files: the shell, all node pages of the core and of the tier-1 guides, area and path pages, tier-1 collection pages, compat pages for tier 1, comparisons, scene bundles and search shards | the fastest path; the limit is 100,000 files per Worker version on Paid and 25 MiB per file (checked 2026-11-09), so 95,000 leaves a margin of 5,000 |
| **R2 bucket `programming-site`** behind the same Worker | every other file, at its URL path | a request that misses the assets reaches the Worker, which reads `r2://programming-site/<path>`, sets the headers and caches the response at the edge |
| **R2 bucket `programming-raw`** | not served; the raw source snapshots of 13 | the citation check (13.6) |
| **GitHub Pages mirror** (`tamnd.github.io/programming`) | the core only: the core node pages, the tier-1 guides' `intro` and `core` features, the indexes and search, under 900 MB | for readers who cannot reach Cloudflare; a banner says which pages it leaves out and links to the main site. GitHub Pages takes at most 1 GB per site. If the core passes 900 MB, the mirror drops the tier-1 guides first and keeps the concepts and tasks. If it cannot hold the core concepts and tasks, it is dropped and the offline archive (16.21) takes its role (Default: keep it to M3, review at M3) |

`registers/site-tiers.yaml` decides which files go to the assets, by path pattern and by the core list. `prog site check` fails when the hot set passes 95,000 files or any file passes 25 MiB. Files over 25 MiB (some runtimes and datasets) always go to R2.

**The Worker `programming`** lives in `programming-agent/worker/`, in TypeScript, about 1,000 lines (Estimate). On each request it tries these steps in order:

1. **Static assets.** Serve the file if it is there. Assets are matched first, with the `run_worker_first` setting off for asset paths, so most requests cost no Worker CPU (check the binding setting names at M0).
2. **Content negotiation.** `Accept: text/markdown` on a page path serves the `.md` twin (16.10).
3. **Special routes:**
   - `/mcp` and `/.well-known/mcp` go to the MCP server (16.10).
   - The playground routes get the COOP and COEP headers (16.8).
   - `/data/` is served from R2 with CORS.
4. **Redirects.** The Worker answers these:
   - `/t/<tag>`;
   - the aliases of `registers/aliases.reg`;
   - well-known ids (`/cwe-89/`);
   - case folding and missing slashes;
   - `/<slug>/<lang>/` tabs;
   - comparison order.

   They come from a redirect map that the build writes (about 50,000 entries, Estimate, loaded once per isolate from assets). Card URLs are answered from the card indexes (16.7).
5. **R2.** Read `programming-site/<path>` and cache it at the edge.
6. **404.** The 404 page offers a search for the words of the path. If the path looks like a card (`/<owner>/api/...`) that is not in the universe, it also offers "Report a missing item" (16.19).

**Cost** (Estimate; check at M0):

- Workers Paid at $5/month includes 10 million requests per month.
- R2 storage of about 14 GB is about $0.21/month.
- R2 reads are Class B operations, which are cheap and mostly absorbed by the edge cache.
- R2 has no egress fee.
- The expected total is under $15/month at 10 million page views per month.

**DNS.** The custom domain `programming.tamnd.com` is a Workers custom domain on the `tamnd.com` zone. Releases are served at `<release>.programming.tamnd.com` by the same Worker, from the prefix `releases/<release>/` in R2 (16.21).

## 16.21 Deploy, rollback and releases

**Deploy flow** (`prog site deploy -out build/site -commit <sha>`). It follows algorithm's `deploy.sh`.

1. **Manifest.** The build writes `build/site/manifest.json`: each file with its path, BLAKE3 hash, size and tier (assets or R2). Deploy reads the last deployed manifest from `programming-site/_deploys/current.json` and computes the diff.
2. **Stamp.** Every HTML page has a footer link to `github.com/tamnd/programming/commit/<sha>`. The deploy refuses to run if the build's commit differs from `-commit`.
3. **Back up.** Before it overwrites or deletes an R2 object, deploy copies the old object to `_deploys/<previous-sha>/<path>`. Only changed paths are copied.
4. **Upload R2 first.** Deploy uploads changed and new R2 objects, in parallel with retries. Pages in R2 can link to new assets, so R2 goes first and assets second. An asset page never links to an R2 file that does not exist yet.
5. **Deploy the Worker.** If any hot file or the Worker code changed, deploy uploads a new Worker version with its assets. Workers upload only changed assets by hash.
6. **Purge.** Deploy purges the edge cache for the changed R2 paths, in batches. It purges everything when a template, CSS or JS file changed.
7. **Live check.** Like `deploy.sh`, deploy fetches fixed probe pages and checks that they show `commit/<sha>`:
   - `/closures/` (assets);
   - one R2 page;
   - one card redirect;
   - `/llms.txt`;
   - one `.md` twin.

   It retries every 30 s for up to 10 minutes.
8. **Commit or roll back.**
   - On success, deploy writes `_deploys/current.json` and keeps the old manifest as `_deploys/<previous-sha>.json`.
   - On failure, deploy rolls back automatically: it restores the backed-up R2 objects, rolls the Worker back to its previous version, and purges again. Then it opens an issue labelled `deploy-failed` with the last lines of the log.
   - Backups older than 30 days are deleted.

`prog site deploy -rollback <sha>` does the same restore on request.

**Who runs it.** The build server runs the deploy loop every 5 minutes, as for algorithm and math2. It does nothing when the commits of `tamnd/programming` and `tamnd/programming-agent` did not change. A merge to `main` is live in 10 to 20 minutes (Estimate): the tests of the changed nodes run first (07, 09).

**Secrets.**

- The deploy reads `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` from the environment. The wrapper script loads them with `set -a; source "$HOME/data/.env"; set +a`, as `deploy.sh` does.
- The token is scoped to Workers Scripts edit, R2 edit for the two buckets, and Cache Purge for the zone.
- No command prints the environment. Deploy logs print only paths, counts and the commit.

**Site releases.**

- Each monthly release (21) builds the site of its commit into `releases/<YYYY-MM>/` in R2. It is served at `https://<YYYY-MM>.programming.tamnd.com`, with the same paths. `/t/<tag>@<release>` redirects there.
- Each release also produces an offline archive, `programming-<release>-core.zip`, attached to the GitHub release: the core pages, their tier-1 code, search with an embedded index, and the Markdown twins. It opens from disk.
- The datasets of the release are tagged on Hugging Face with the same name.
