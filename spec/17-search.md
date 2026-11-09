# 17. Search

Search is the main door into the site. A programmer comes with a task ("how do I read a file in Go"), an error message copied from a terminal, a symbol (`str.split`), an id (`E0382`, `CWE-89`, `PEP 634`), a version ("python 3.12 new") or a choice ("python vs go"). Search must take each of these to the one page that answers it, in under 50 ms per keystroke, on a site of about 1M entries. This document extends the search of algorithm (2150/15: the palette, the sharded static index, the matcher, the priors, the evaluation set, the stages). Where this document is silent, 2150/15 applies.

## 17.1 What people search for

| Query type | Examples | Resolves to | How (section) |
|---|---|---|---|
| name | `closure`, `ruff`, `json`, `borrow checker` | the node with that name or alias | prefix and fuzzy match on names (17.7) |
| task in a language | `how do I read a file in go`, `reverse a list python`, `go read file` | the task card for that task × language (04.2), or the task page with the language tab open | task intent parse (17.6.1) |
| error message | `borrow of moved value: \`x\``, `TypeError: 'NoneType' object is not subscriptable`, `undefined reference to \`main'` | the diagnostic card (04.5) | template match (17.6.2) |
| api symbol | `str.split`, `Array.prototype.map`, `std::vector::push_back`, `os.ReadFile`, `HashMap::entry` | the API card | symbol index (17.6.3) |
| id alias | `E0382`, `cwe-89`, `pep 634`, `pep-634`, `CVE-2021-44228`, `RFC 9110`, `E501`, `TS2322`, `SC2086` | the node that owns the id | id index (17.6.4) |
| version | `python 3.12 new`, `java 21 features`, `c++23`, `es2024` | the release page or the "what is new" view of 09 | version parse (17.6.5) |
| comparison | `python vs go`, `rust or c++`, `tabs vs spaces`, `npm vs pnpm` | the comparison page `/compare/a-vs-b/` (08.8), or a generated table if no page exists | comparison parse (17.6.6) |
| concept in a language | `closures in rust`, `rust closures`, `generics java` | the feature card that `realizes` the concept in that language (04.6) | concept × language parse (17.6.1) |
| command | `> random concept`, `> toggle theme` | a palette command | prefix `>` (as in algorithm) |
| tag | `#memory-safety`, `#async` | the tag page | prefix `#` |
| full text | a phrase from a page body | any page | full-text stage (17.10) |

The site is search first, not search only. Every result is a page with a URL, so search never shows content that is not on a page.

## 17.2 What is indexed

Each node of the graph (04) and each card (04.2) is one entry. An entry holds only the fields that search needs. The page body is not in the small index; it is in the full-text index (17.10).

| Field | From | Weight in BM25F (Default) |
|---|---|---|
| `title` | the node's name; for a card, "`<item>` in `<language>`" | 3.0 |
| `aliases` | the `aliases` and `also_known_as` of the node (04.3), the translated names (20) | 2.5 |
| `ids` | external ids: `pep`, `cwe`, `cve`, `rfc`, `ecma`, diagnostic codes, lint codes | exact match only (17.6.4) |
| `symbol` | for API cards, the qualified name and its parts (`str.split`, `split`) | exact and prefix (17.6.3) |
| `templates` | for diagnostic cards, the `message_templates` (04.5) | template match only (17.6.2) |
| `summary` | the one-line summary (06.2) | 1.0 |
| `tags` | tags and the area (05) | 0.8 |
| `lang` | the language of a card, as a facet | facet and parse |
| `kind` | the node kind (04.1) | facet and intent |
| priors | `level`, `core`, `coverage` (C0 to C6, 18.2), `tier` of the language (08.2), views in 30 days, `status` | ranking (17.7) |

Entries by kind (Estimate at full coverage):

| Kind | Entries | Note |
|---|---|---|
| languages | about 9,000 | all known languages with a page (08) |
| concepts, patterns, paradigms | about 4,000 | |
| features (concept × language cards) | about 60,000 | 15 tier-1 languages × the core concepts, plus tiers 2 and 3 |
| tasks and task cards | about 120,000 | about 3,000 tasks × the languages that have a card |
| API cards | about 500,000 | the standard libraries of tiers 1 and 2, and the core packages |
| diagnostics and lints | about 150,000 | compiler errors, warnings, linter rules |
| tools, libraries, formats, protocols, standards | about 80,000 | |
| weaknesses, proposals, releases, people, organizations | about 60,000 | |
| comparisons | about 20,000 | only the pairs that have a page |
| **total** | **about 1M** | |

## 17.3 The palette

As in algorithm (2150/15.2): it opens on `/`, `⌘K` or `Ctrl K` from any page. Results show as you type. They are grouped by kind with at most 5 per group and a "more" link. There are recent pages from `localStorage`, `>` for commands, `#` for tags, facet chips with `Tab`, and "Tell us it is missing" when nothing matches. These parts are new:

- **The context language.** On a page in a language's tree (`/rust/...`), the palette has a removable chip `lang:rust`. A query with no language in it is ranked as if the reader typed "in rust". The chip shows, so the reader sees why Rust results come first, and `Backspace` on an empty query removes it.
- **The interpreted query.** When the parser understands the query (17.6), the first row says what it understood: "Task: read a file · Language: Go", "Error: E0382 (rustc)", "Compare: Python and Go". The reader can press `Tab` on the row to change the language or the kind.
- **The answer preview.** For a task card, an API card or a diagnostic card, the first row shows the first code block or the one-line fix of the card (06.3). Many queries end there. `Enter` opens the page, and `⌘C` or `Ctrl C` copies the code.
- **Paste of a long error.** If the reader pastes more than one line, the palette keeps the first line that matches a template (17.6.2) and shows "Matched line 3 of 12". The full paste is not stored and not sent (17.9).
- **Keyboard.** `↑` and `↓` move, `Enter` opens, `⌘Enter` or `Ctrl Enter` opens in a new tab, `Esc` closes, `Tab` cycles the chips. The list is an ARIA combobox with `aria-activedescendant`, and the number of results is read by a live region.

## 17.4 The default view by level

When the query is empty, or when the reader opens `/search/` with no query, the palette shows a start view. The view depends on the level that the reader chose on the home page (04.9, 06.5), kept in `localStorage`. With no choice it uses `intro`.

| Level | Empty-query view | Ranking change for the same query |
|---|---|---|
| `intro` | "Start here" (the first 10 concepts of the learning path), the tasks that beginners ask most, the tier-1 languages | intro and core entries ×1.3; research entries ×0.6; API cards and diagnostics stay (beginners paste errors too) |
| `core` | recent pages, the most viewed task cards of the context language, the newest release notes of tier-1 languages | no change (the default weights of 17.7) |
| `advanced` | recent pages, the context language's advanced features, the recent proposals | advanced entries ×1.2 |
| `research` | recent pages, formal semantics (11), papers, open proposals | research and formal entries ×1.3; intro entries ×0.8 |

The level only changes the order. It never hides a result.

## 17.5 The index

The index is static (Default). The build writes it, and the browser loads only what it needs. There is no server in the path of a keystroke.

| File | Holds | Size, gzip (Estimate) | Loaded |
|---|---|---|---|
| `head.json` | the core (about 3,000 items, 05.4), the tier-1 language nodes, the top 10,000 entries by views, the alias table, the stop words, the language-name table (17.6.7) | about 400 KB | on the first open of the palette, then cached by the service worker |
| `names-<k>.bin` | a prefix trie of titles and aliases, split by the first two characters (about 400 shards) | 5 to 40 KB each | on the first keystroke of that prefix |
| `lang-<lang>-<k>.json` | the cards of one language: features, tasks, API cards, diagnostics; split into shards of at most 2,000 entries | 20 to 60 KB each | when the query names the language, or the context chip is set |
| `ids.json` | external ids to node ids: PEP, CWE, CVE, RFC, diagnostic and lint codes | about 900 KB, split into 16 shards by id prefix | when the query looks like an id |
| `templates-<lang>.bin` | the compiled message templates of one language's diagnostics (17.6.2) | 30 to 300 KB each | on paste, or when the query looks like an error |
| `symbols-<lang>-<k>.bin` | the API symbol trie of one language | 20 to 80 KB each | when the query looks like a symbol |
| `facets.json` | counts per kind, language, level and area | about 30 KB | with `head.json` |
| `manifest.json` | the list of shards with their content hashes | about 60 KB | with `head.json` |

The full index is about 120 MB gzip (Estimate). A typical session loads less than 1 MB. All shard file names hold a content hash, so the CDN can cache them forever, and a new build changes only the shards whose content changed.

**Matcher.** The algorithm matcher (2150/15.4: prefix, trigram fuzzy with edit distance 1 or 2, BM25F, priors) is extended with the parsers of 17.6. It is one TypeScript module of about 12 KB gzip (Estimate) that runs in a Web Worker so that typing never blocks the page. The same matcher is compiled for the agent API (17.11).

**Normalization.** As in algorithm: NFKC, case folding, punctuation folding, stop words with low weight, light stemming, aliases in the translated languages (20). These rules are new for programming:

- Code punctuation is kept inside tokens that look like code: `str.split`, `std::vector`, `C++`, `C#`, `F#`, `.NET`, `node.js`, `E0382`, `->`, `?.`. The tokenizer has a list of these (Default, `registers/search-tokens.yaml`).
- `cpp`, `c plus plus` and `cxx` map to `C++`; `csharp` and `c sharp` map to `C#`; `js` maps to `JavaScript`; `ts` maps to `TypeScript`; `golang` maps to `Go`; `py` maps to `Python`. The list is in the language-name table (17.6.7).
- Snake case, camel case and kebab case are split for matching (`readFile` matches `read file`) and kept whole for exact match.

## 17.6 Query parsing

The matcher first tries to parse the query as one of the structured types. A parse that succeeds adds a strong candidate at the top. The plain name match still runs, so a wrong parse never hides the plain results.

### 17.6.1 Task × language and concept × language

The parser looks for a language name (17.6.7) and removes it, together with the joining words ("in", "with", "using", "for"). It removes "how do I", "how to", "how can I" and similar openers from a list. It then matches the rest against task titles and aliases and against concept names.

- A task match plus a language gives the task card for that language. If the card does not exist, the result is the task page with the note "No Go card yet", and the query is logged as a gap (17.9).
- A concept match plus a language gives the feature card that `realizes` the concept in that language. If the language has no such feature (the concept matrix says `absent`, 04.8), the result is the concept page with the cell "Go has no exceptions; see error values".

### 17.6.2 Error messages

Diagnostic cards carry `message_templates` (04.5), for example `borrow of moved value: \`{}\``. The build compiles the templates of each language into a matcher:

1. Each template is split into literal parts and holes. A hole matches any text of at most 200 characters on one line.
2. The query is normalized: paths, line and column numbers, addresses, quoted identifiers and numbers are replaced by placeholders, with the normalizers of 07.6.
3. The matcher finds the templates whose literal parts all occur in order. The score is the length of the matched literal text divided by the length of the query line. A score of at least 0.6 is a match (Default).
4. A code in the query (`error[E0382]`, `TS2322`, `SC2086`) is exact and wins over a template match.
5. If two languages have templates that match (for example a message that both GCC and Clang print), the palette shows both, with the compiler name, and the context language decides the order.

The diagnostic card is the result. The answer preview shows the one-line cause and the fix (06.3). Messages change between compiler versions; each template records the versions it applies to (09), and an older message still matches its card.

### 17.6.3 API symbols

A query that looks like a qualified name (it holds `.`, `::`, `#` or `->` between identifier parts, or it is a known unqualified name of an API card) is looked up in the symbol trie. `str.split` finds Python's `str.split`. `split` alone finds the `split` entries of all languages, with the context language first. The symbol forms of each language are in `languages/<lang>/symbols.yaml`: `Array.prototype.map` and `Array#map` are the same entry for JavaScript.

### 17.6.4 Id aliases

The id parser knows the id forms of 04.3. It matches the form with or without a separator, space or case: `pep634`, `PEP 634`, `pep-634` are the same; `cwe89` and `CWE-89` are the same. The result is the node that owns the id, at its owned URL (`/python/proposals/pep-634/`, `/rust/errors/E0382/`). An id that we do not have yet gives "We have no page for PEP 9999" and a link to the source (python.org for PEPs), and it is logged as a gap.

Short codes that more than one tool uses (`E501` is a lint code of pycodestyle and Ruff, and also other tools' codes) give a list, with the context language and the most viewed tool first.

### 17.6.5 Versions

A language name followed by a version (`python 3.12`, `java 21`, `c++23`, `es2024`, `go 1.23`) with an optional word from a list ("new", "features", "changes", "what's new", "release") gives the release page of that version (09.1). A version that is not released yet gives the planned release page if we have one. A language name with "latest" gives the newest release.

### 17.6.6 Comparisons

Two names joined by "vs", "versus", "or", "compared to" or "and" (with "difference" or "compare") give the comparison page `/compare/a-vs-b/` with the names in the canonical order (08.8). If no page exists and both names are languages of tier 1 or 2, the result is the generated comparison table from the matrix (08.8), marked "generated". Otherwise the result is the two pages, and the pair is logged as a gap.

### 17.6.7 Short and ambiguous language names

Some language names are one letter (`C`, `R`, `D`, `J`, `K`, `Q`) or are common English words (`Go`, `Rust`, `Swift`, `Ruby`, `Dart`, `Julia`, `Elixir`, `Crystal`, `Pony`, `Io`, `Red`, `Lean`, `Mercury`, `Logo`, `Basic`, `Elm`, `Nim`, `Zig`, `Odin`, `Pure`, `Clean`, `Factor`, `Joy`). The parser treats these names as a language only in these cases:

| Case | Example | Result |
|---|---|---|
| the name is the whole query | `go`, `r` | the language node is the first result, and other nodes with that name follow |
| the name is joined to a task or concept by "in", "with", "using" | `read a file in go` | language |
| the query holds a context word for that language from `registers/lang-context.yaml` | `go goroutine`, `r dataframe`, `c pointer`, `rust cargo` | language |
| the context chip is set to that language | on a `/go/` page, `read file` | language |
| a symbol form of that language | `fmt.Println` | language |
| another word in the query is a code term and the name is next to it | `c struct padding` | language |

In other cases the word is a plain word. "go to statement" finds the `goto` concept. "swift review" does not assume Swift. For a one-letter name, the parser needs one of the cases above. A bare word that is in the last row of the table with no other code term gives both readings, with a disambiguation row first ("Did you mean the R language?").

The context-word list has about 20 to 50 words per ambiguous language (Estimate). The curator builds it from the language's keywords, its tool names and its standard library names, and the search evaluation (17.8) tests it.

## 17.7 Ranking

The score of an entry is the text score (BM25F over the fields of 17.2, or the parse score of 17.6) times the priors. The priors of algorithm (2150/15.5) apply, with these changes for programming (all Default, tuned by 17.8):

| Prior | Factor |
|---|---|
| structured parse matched this entry (17.6) | ×3.0 |
| kind matches the query intent ("error" words, symbol form, "vs", "how do I") | ×2.0 |
| `core: true` | ×1.3 |
| level: by the reader's level (17.4) | ×0.6 to ×1.3 |
| coverage: C0 ×0.5, C1 ×0.7, C2 ×0.85, C3 ×1.0, C4 ×1.05, C5 and C6 ×1.1 (18.2) | as listed |
| language tier: tier 1 ×1.2, tier 2 ×1.1, tier 3 ×1.0, tier 4 ×0.9 (08.2) | as listed |
| context language match | ×1.5 |
| popularity: `1 + 0.3 × log10(1 + views30)` | about ×1.0 to ×2.5 |
| status: deprecated or removed API, retired language | ×0.5; ×0.1 if a replacement exists (the replacement is shown in the row) |
| evidence: page runs at X2 or better (07) | ×1.1 |

Ties are broken by the shorter title, then by the node id. A pin list (`registers/search-pins.yaml`, at most 500 pins, Default) can put one entry first for one query; each pin has a reason and is checked by 17.8.

## 17.8 Evaluation

The evaluation set is `registers/search-eval.yaml`: queries with the expected node. It starts with about 3,000 queries (Estimate) at M2 and grows from the logs (17.9). Each query type of 17.1 has at least 100 queries. The error-message part uses real messages from the `fails` blocks of 07.2, so the message and the card come from the same run.

| Measure | Gate (Default) |
|---|---|
| success@1 on id, symbol and error queries | at least 95% |
| success@5 on all queries | at least 95% |
| MRR on all queries | at least 0.85 |
| false language reading on the ambiguous-name set (17.6.7) | at most 3% |
| p95 keystroke latency, mid-range phone, warm cache | at most 50 ms |
| first result after the palette opens, cold cache, fast 4G | at most 500 ms |

`prog search eval` runs the set in CI on every change to the matcher, the priors, the tokens or the context words. A drop below a gate fails the build. Rule T07 (no two nodes with the same name and kind without a disambiguation entry) holds as in algorithm.

## 17.9 Query logs and gaps

The search log tells us what is missing. Each query that finds nothing, or that finds a gap (a task with no card for the language, an id with no page, a comparison with no page), is a request for content.

**What is logged.** The browser sends one event per closed palette session. The event holds:

| Field | Note |
|---|---|
| the normalized query | after 17.5 normalization; error messages after the placeholder step of 17.6.2, so paths, names and values from the reader's code are already removed |
| the parse type | task, error, symbol, id, version, comparison, name, none |
| the result position that was opened, or none | |
| the gap type, if any | no result, missing card, missing id, missing comparison |
| the context language and the level | |
| the site language | |

**What is never logged.** No IP address is stored, no user id, no cookie, no session id across visits, no user agent string, and no full pasted text. A query longer than 200 characters after normalization is dropped. A query that holds something that looks like a secret (a token pattern, an email address, a key) is dropped in the browser before it is sent. The reader can turn logging off in the settings, and the site honours the Global Privacy Control signal.

**Aggregation.** The collector (a Worker, 13) adds each event to daily counters by (normalized query, parse type, gap type). It does not keep the events. A query is shown in reports only if at least 5 sessions sent it in the window (Default), so a rare query that might point to one person never shows. The counters are kept for 13 months.

**From gap to work.** Each week the `curator` agent (14) reads the gap report:

| Gap | Action |
|---|---|
| a task with no card for the language, asked at least 20 times | a work item for the `writer` with that task × language |
| an error message that matches no template | the message is grouped by similarity; a group of at least 10 sessions is a work item to add the diagnostic card or a template |
| an id with no page | a work item to add the node, if the id is in the scope of 05 |
| a comparison with no page | a work item when both are tier 1 or 2 and the pair is asked at least 50 times |
| a query that opens the 4th or lower result most of the time | a ranking item for 17.8, and the query is added to the evaluation set |
| a word that the parser read as a language and the reader then changed | a context-word item for 17.6.7 |

The gap report is public, in aggregate, at `/about/gaps/` (counts only, with the 5-session rule). It is also an input to the coverage plan of 18.

## 17.10 Stages and the server option

| Stage | What | When |
|---|---|---|
| S1 | the static index of 17.5 and the parsers of 17.6, in the browser | M1 (Default) |
| S2 | full text with Pagefind over the page bodies, one index per site language, loaded only for "Search all text" | M2 |
| S3 | a Worker at `/api/search` with D1 (SQLite FTS5) for full text and for agents, when the static first-load files pass 2 MB or the S2 index becomes too large to load in parts | Open question (22): at about 1M entries we expect to need S3 for full text, but the static name, id, symbol and template shards should stay static |
| S4 | semantic search with Vectorize: embeddings of the summary and the first section of each page, used for queries that S1 cannot parse and that have few text results; the semantic results are mixed in below the text results and are marked "related" | Open question (22): it needs a choice of embedding model and a cost check |

**Default.** The keystroke path stays static at all stages. The Worker serves full text, agents and the log collector. If the Worker is down, the palette still works for names, tasks, ids, symbols and errors.

**Pagefind.** Pagefind builds a static, chunked full-text index at build time and loads only the chunks a query needs. We use it for page bodies in S2. We do not use it for the structured parsers, which need our own shards. Its limits at our size are measured in M2 (Estimate: about 1M pages is large for one Pagefind index, so we split by site section).

**Sizes and costs** (Estimate): the static index is about 120 MB gzip, inside the static asset limits of the hosting plan (13). D1 for S3 holds about 1M rows of summaries and section text, a few GB. Vectorize for S4 holds about 1M vectors. Exact limits and prices are checked in 13 before S3 or S4 starts.

## 17.11 Search for agents

Agents use the same matcher and the same index as the palette, so a query gives the same results to a person and to an agent.

- **HTTP.** `GET /api/search?q=<query>&lang=<lang>&kind=<kind>&level=<level>&limit=<n>` returns JSON: the interpreted query (17.6), and the results with `id`, `uri` (`pg://`), `url`, `kind`, `title`, `summary`, `lang`, `score` and, for cards, the answer preview. It is served by the Worker (S3), or, before S3, by a static fallback that serves only the parsed types (id, symbol, task, error) from the shards. The full contract is in 16.
- **MCP.** The `search(query, filters)` tool of the programming MCP server runs the same matcher (a Go build of the same rules, tested against the TypeScript one on the evaluation set). It adds `resolve_error(message, lang?)` for error text. The tools are in 16.
- **Logs.** Agent queries are counted apart from reader queries and are never part of the 5-session reports, so that an agent loop cannot fill the gap report.
- **Limits.** `/api/search` has a rate limit per IP that is applied at the edge and not stored (16).
