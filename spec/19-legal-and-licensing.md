# 19. Legal and licensing

This document says under which licence each part of the corpus is published, what the project may take from each source, how quotes and paraphrase are limited, how agent-written content is labelled, and how the project handles trademarks, security content, privacy, export rules and takedowns. It follows 2150/17 (algorithm) and changes it where programming forces a change: far more of the sources are documentation sites with their own licences, the code examples sit inside Markdown files, and the weakness pages show vulnerable code.

This is a design, not legal advice. The points marked "Default" are decisions that a lawyer's check may change (22). Every licence fact in 19.3 says whether it was checked, and when.

**Principles.**

1. **Write fresh.** As in algorithm 17.2, prose and code are written fresh. Sources are read, closed, and then the writer writes. The licences of classes A to C (19.2) would allow more; the project does not use that freedom by default, so that one rule covers every writer.
2. **Facts are free.** Names, version numbers, dates, error codes, API signatures, rule ids, clause numbers and the structure of a field are facts. The project takes facts from any source, with a citation, unless the source's terms forbid automated access (class E).
3. **Code must stay MIT.** Readers copy examples into closed programs. No example, probe or task implementation may come from a share-alike or GFDL source (rule C06).
4. **Data must stay CC0.** Cards, compat entries and registers are CC0, so text from share-alike or attribution-only sources must not enter them (rule C07).
5. **Show the evidence.** Every page lists its sources with their licences, and the site's `/sources/` page lists every source with its class and what the project took.

## 19.1 Licence zones

Each file of `tamnd/programming` is in one zone, set by its path in `REUSE.toml`. The zones are those of algorithm, with three new paths (cards, compat, probes) and one new mixed zone (example files).

| Zone | Paths | Licence | Why |
|---|---|---|---|
| prose | `nodes/**/node.md` (body), `nodes/**/notes.md`, `topics/**/*.md`, `i18n/**` | CC BY-SA 4.0 | the licence of Wikipedia, so class B text can be adapted with attribution (19.2); as algorithm and math2 |
| example files | `nodes/**/examples/*.md` | prose parts CC BY-SA 4.0; fenced code blocks and their `output` blocks MIT | the file holds both; see "Mixed files" below |
| figures | `nodes/**/figures/**`, frames and images made from scenes | CC BY-SA 4.0 | they explain the prose and travel with it |
| code | `nodes/**/impl/**`, `nodes/**/probes/**`, `nodes/**/vuln/**`, `nodes/**/tests/**`, `nodes/**/cases/**`, `nodes/**/fixtures/**`, `nodes/**/contract.yaml`, `nodes/**/viz/**`, `languages/**/wrappers/**` | MIT | readers copy code into their own programs, often closed ones |
| Lean | `lean/**` | Apache-2.0 | the licence of Mathlib and CSLib, so code can move in both directions (11.3) |
| data | front matter (the YAML block of each `node.md`), `cards/**`, `nodes/**/compat.yaml`, `nodes/**/matrix.yaml`, `nodes/**/releases.yaml`, `nodes/**/implementations.yaml`, `nodes/**/concepts.yaml`, `registers/**`, `topics/*.yaml`, `crosswalk/**`, `prov/**`, `reviews/**`, `schema/**`, `manifests/**`, `reports/**` | CC0 1.0 | facts about programming should be free to reuse with no conditions, as Wikidata's and MDN BCD's are |
| toolchain | all of `tamnd/programming-agent` | Apache-2.0; the viz package MIT; the Worker MIT | as algorithm-agent |

**Mixed files.** REUSE allows one licence expression per file. `node.md` is marked `CC-BY-SA-4.0 AND CC0-1.0`, and `examples/*.md` is marked `CC-BY-SA-4.0 AND MIT`. `LICENSES/README.md` states the split in plain words: in `node.md`, the YAML block between the first two `---` lines is CC0 and the rest is CC BY-SA 4.0; in every Markdown file, the content of each fenced code block (including `output` blocks and hidden preludes) is MIT and the rest is CC BY-SA 4.0. The `.json` and `.md` copies of pages on the site (16.10) carry a `license` object with the same split.

**REUSE.** The repository follows REUSE 3.3, with licences set by path in `REUSE.toml`, not by a header in each file, as in algorithm:

```toml
version = 1

[[annotations]]
path = ["nodes/**/node.md", "nodes/**/notes.md", "topics/**.md", "i18n/**", "README.md", "CONTRIBUTING.md"]
SPDX-FileCopyrightText = "The programming contributors <https://programming.tamnd.com/about/>"
SPDX-License-Identifier = "CC-BY-SA-4.0 AND CC0-1.0"

[[annotations]]
path = ["nodes/**/examples/*.md"]
SPDX-FileCopyrightText = "The programming contributors <https://programming.tamnd.com/about/>"
SPDX-License-Identifier = "CC-BY-SA-4.0 AND MIT"

[[annotations]]
path = ["nodes/**/impl/**", "nodes/**/probes/**", "nodes/**/vuln/**", "nodes/**/tests/**",
        "nodes/**/cases/**", "nodes/**/fixtures/**", "nodes/**/contract.yaml", "nodes/**/viz/**",
        "languages/**/wrappers/**"]
SPDX-FileCopyrightText = "The programming contributors <https://programming.tamnd.com/about/>"
SPDX-License-Identifier = "MIT"

[[annotations]]
path = ["lean/**"]
SPDX-FileCopyrightText = "The programming contributors"
SPDX-License-Identifier = "Apache-2.0"

[[annotations]]
path = ["cards/**", "nodes/**/compat.yaml", "nodes/**/matrix.yaml", "nodes/**/releases.yaml",
        "nodes/**/implementations.yaml", "nodes/**/concepts.yaml", "registers/**", "topics/*.yaml",
        "crosswalk/**", "prov/**", "reviews/**", "schema/**", "manifests/**", "reports/**", "CITATION.cff"]
precedence = "aggregate"
SPDX-FileCopyrightText = "The programming contributors"
SPDX-License-Identifier = "CC0-1.0"
```

Rule C02 runs `reuse lint` and fails on any file with no licence.

**Imported data keeps its own licence.** Data that the project imports as-is is not relicensed. It sits under `vendor/<source>/` with the source's own licence file and its own REUSE annotation: MDN BCD (CC0), web-features (Apache-2.0), SPDX license list data (CC0), endoflife.date snapshots (MIT), Exercism canonical data (MIT). The build may merge these into published datasets only when the result's licence allows it, and the dataset's README names each input and its licence (16.10).

**On the site.** Each page footer shows "Text CC BY-SA 4.0 · Code MIT · Data CC0", and "Lean Apache-2.0" where it applies. The copy button copies the code alone. The "Raw" file of an example starts with one SPDX comment line in the language's comment syntax (from `registers/languages.yaml`), added by the build, with the page URL, so a copied file keeps its notice.

**Inbound equals outbound.** A contribution is under the licence of the zone it lands in. Each commit is signed off under the Developer Certificate of Origin (`Signed-off-by:`), checked by the DCO app. There is no CLA (Default, as algorithm).

**Third-party parts** of the site and the toolchain are listed in `programming-agent/THIRD_PARTY.md`, as in algorithm, with the browser runtimes of 07.10 added: Pyodide (MPL-2.0), MicroPython (MIT), ruby.wasm (MIT), php-wasm (GPL-2.0-or-later, served as a separate program in a Worker and not linked into our code; Default, checked at M1), wasmoon (MIT), sqlite-wasm (public domain), PGlite (Apache-2.0 or the PostgreSQL licence), wasm_of_ocaml (LGPL-2.1 with the OCaml linking exception; checked at M1). A runtime whose licence needs a notice gets it on `/about/legal/`.

## 19.2 Source classes

Each source in `registers/sources.yaml` has a class. The class says what the project may take from it. The classes are those of algorithm 17.2.

| Class | Licence of the source | May take |
|---|---|---|
| A | public domain, CC0 | anything, with a credit line as courtesy; data may enter the CC0 zone |
| B | share-alike, and the project's prose zone can accept it (CC BY-SA 1.0 to 4.0, since CC BY-SA 3.0 and earlier allow later versions) | prose may be adapted into the prose zone with attribution; never into code (C06) or data (C07) |
| C | permissive: CC BY, MIT, Apache-2.0, BSD, 0BSD | facts, structure, short quotes with attribution, adapted prose with attribution; code ideas with a citation; our code is still written fresh |
| D | all rights reserved, a licence for personal or internal use only, GFDL, NC licences, or terms that are unclear | facts, names, the structure of a field, clause and section numbers for citations, and quotes under G04 (19.4) only. Never text, figures, exercises or code |
| E | forbidden to harvest: terms that forbid automated access, or the owner asked us to stop | nothing by machine; a person may cite it |

GFDL sources are class D. GFDL text cannot become CC BY-SA 4.0 (the Wikipedia relicensing window of 2009 does not cover these sources), and its invariant-section and cover-text conditions do not fit code. Non-commercial licences are class D because the project's licences allow commercial use.

## 19.3 Per-source rules

Other documents (05, 06.6, 08, 09, 10, 13) refer to this table. "Checked" gives how the fact was confirmed on 2026-11-09: "repo" means the LICENSE file or the SPDX id of the source's GitHub repository, "site" means the source's own licence page, "research" means the research agents' reports (01). "Unverified" means not confirmed, to be checked at M1 before any use beyond facts.

**Catalogues and data**

| Source | Licence | Class | What we take | Checked |
|---|---|---|---|---|
| PLDB | public domain | A | the language cards' facts, ids, years, creators, file extensions (08.7) | research (readme) |
| Wikidata | CC0 | A | QIDs, labels, dates, people, influences | research |
| Esolang wiki | CC0 | A | names and facts of esoteric languages; text may be adapted, but we write fresh | research (site rights info) |
| MDN browser-compat-data | CC0 | A | imported as-is for web features (09.6); v8.1.8 at the time of writing | repo |
| SPDX license list data | CC0 for the list data; each licence text keeps its own terms | A | the licence cards (05.2, `open-source-and-licensing`) | research; the per-text terms unverified, check at M1 |
| Jargon File 4.4.7 | public domain (ESR asks that it not be printed) | A | terms and history, with credit | research |
| GitHub Linguist | MIT | C | names, aliases, extensions, colours (`crosswalk/linguist.yaml`) | repo |
| Exercism problem-specifications and tracks | MIT | C | canonical-data cases as task test data with credit (10.1); concept names | repo |
| web-features | Apache-2.0 | C | Baseline status for web features (09.5) | repo |
| endoflife.date | MIT | C | release and end-of-life dates as a cross-check (09.3) | repo |
| node.green (node-compat-table) | MIT | C | a cross-check of the Node.js column | repo |
| caniuse | CC BY 4.0 | C | a cross-check only, with credit | repo |
| tldr-pages | CC BY 4.0 for pages, MIT for scripts | C | command names and option facts; a page may link the tldr page | repo |
| free-programming-books | CC BY 4.0 | C | links to free books, with credit | research |
| Go by Example | CC BY 3.0 | C | titles as a check of the task list (10.1) | research |

**Language and platform documentation**

| Source | Licence | Class | What we take | Checked |
|---|---|---|---|---|
| Python documentation | PSF License v2; code in the docs is 0BSD (since 3.8.6) | C (with care: the PSF licence is written for software) | facts, anchors, `versionadded` and related directives as compat facts (09.6) | site (license page) and research |
| Go documentation (go.dev) | content CC BY 4.0, code BSD-3-Clause | C | facts, anchors, the spec's grammar with credit (11.5) | research; the website repo is BSD-3-Clause (repo) |
| The Rust Reference, the FLS, The Rust Book | MIT or Apache-2.0 | C | facts, rule ids, the grammar with credit | repo (Reference Apache-2.0 SPDX; Book has LICENSE-MIT and LICENSE-APACHE) and research |
| Kotlin documentation and the Kotlin spec | Apache-2.0 | C | facts, anchors | repo (kotlin-web-site) and research |
| Swift book (The Swift Programming Language) | Apache-2.0 | C | facts, anchors | repo |
| Dart documentation (dart.dev) | content CC BY 4.0, code BSD-3-Clause | C | facts, anchors | repo (LICENSE text) |
| TypeScript website and handbook | CC BY 4.0 | C | facts, anchors | repo (SPDX) |
| Microsoft Learn (.NET docs, PowerShell docs) | docs CC BY 4.0, code MIT | C | facts, anchors, the cmdlet reference as api card facts | repo (dotnet/docs CC-BY-4.0; PowerShell-Docs LICENSE.md CC BY 4.0 and LICENSE-CODE.md MIT) |
| C# standard (ECMA-334) | the dotnet/csharpstandard repo is CC BY 4.0 | C | facts, clause numbers | research |
| PHP manual | CC BY 3.0 or later | C | facts, anchors | repo (bookinfo.xml) |
| MDN prose | CC BY-SA 2.5 (or later, per research) | B | prose may be adapted with attribution; we write fresh by default | repo (LICENSE.md) |
| MDN code samples | CC0 if added on or after 2010-08-20, MIT before; MDN cannot tell which applies to a given sample | treat as C | ideas only; our examples are written fresh, so the dating question does not arise | repo (LICENSE.md) |
| Java SE documentation and the JLS | Oracle "Limited License Grant", internal evaluation only | D | facts, names, section numbers | site (JLS copyright page) and research |
| cppreference | CC BY-SA 3.0 and GFDL (dual) | D by default | facts and cross-checks only (09.6); never text | research |
| GNU manuals, including the Bash manual | GFDL (the version per manual is unverified, check at M1) | D | facts, option names, section names for locators | unverified |
| Linux man-pages | per-page SPDX tags: Linux-man-pages-copyleft (941 pages), GPL variants, BSD, MIT | D by default; per page at M1 | facts, names, section numbers | research |
| Arch Wiki | GFDL 1.3 or later | D | facts and links | research |
| Ruby documentation | unverified, check at M1 | D until checked | facts | unverified |
| Other official docs (Scala, Haskell, OCaml, Elixir, Julia, R, Zig, Lua, Erlang, Perl, ...) | per language; recorded in `registers/sources.yaml` at M1 | D until checked | facts | unverified, check at M1 |

**Standards**

| Source | Licence | Class | What we take | Checked |
|---|---|---|---|---|
| ISO, IEC and ISO/IEC JTC 1 standards (C, C++, SQL, Ada, Fortran, COBOL, ...) | all rights reserved | D | clause numbers and titles only, never text; our own simplified grammar forms, marked as simplified (11.5) | general knowledge of ISO terms |
| WG14 and WG21 working papers and drafts (open-std.org) | the papers' own notices; drafts are public | D | facts, paper numbers, section numbers | unverified per paper |
| IEEE standards | all rights reserved | D | clause numbers only | general |
| POSIX.1-2024 (The Open Group) | the free HTML is for personal use only | D | names, facts, section numbers | research |
| ECMA-262 (ECMAScript) and other Ecma standards | Ecma "alternative copyright": copy, modify and distribute with the notice | C | facts, section numbers, the grammar with the notice | research |
| IETF RFCs | the IETF Trust Legal Provisions under BCP 78 and BCP 79; code components are BSD | C for code components, D for text by default | facts, section numbers; code components with the BSD notice | general knowledge; unverified for the current TLP version, check at M1 |
| W3C specifications | W3C Software and Document License | C | facts, section numbers | research |
| WHATWG standards (HTML, DOM, URL, Fetch, ...) | CC BY 4.0 | C | facts, anchors | general knowledge; unverified, check at M1 |
| Unicode standard and data | Unicode License v3 (permissive) | C | data files with the notice (normalization, properties) | research |

**Security**

| Source | Licence | Class | What we take | Checked |
|---|---|---|---|---|
| CWE 5.0 (2026-10-28; "Category" renamed "Grouping") | custom terms: non-exclusive, royalty-free, any copy must reproduce MITRE's copyright designation and the licence, CWE is a trademark, unclear whether modified versions may be published | C for ids and names, D for text | ids, names, the hierarchy, the abstraction level; our own descriptions (06.3); the notice on `/about/legal/` | site (Terms of Use) and research |
| CAPEC | assumed the same terms as CWE | D until checked | ids and names | unverified, check at M1 |
| OWASP (Top 10, ASVS 5.0, Cheat Sheets) | CC BY-SA 4.0 | B | prose may be adapted with attribution; ids and names as facts | repo (Top10 LICENSE) and research |
| NVD and CVE records | NVD data is a US government work; CVE records under the CVE terms of use | A for NVD facts; C for CVE | ids, dates, CVSS scores, affected versions, as card facts with the advisory link | unverified for the current CVE terms, check at M1 |
| OSV | CC BY 4.0 for the aggregated data (per research of algorithm) | C | advisory facts | unverified, check at M1 |

**Reference sites, books and curricula**

| Source | Licence | Class | What we take | Checked |
|---|---|---|---|---|
| Wikipedia | CC BY-SA 4.0 | B | prose may be adapted with attribution; facts with the article's own citations followed to the primary source | research |
| Stack Overflow | CC BY-SA 4.0 since May 2018 (earlier posts CC BY-SA 3.0 or 2.5) | B | question titles and tag counts as facts (`crosswalk/so-tags.yaml`); prose may be adapted with attribution, but G06 limits citations; the official dumps now need a login and an agreement not to train an LLM, so the project uses the API for tag data only | research |
| Learn X in Y minutes | CC BY-SA 3.0 | B | titles and structure; we write fresh | repo and research |
| Hyperpolyglot | CC BY-SA 3.0 | B | the comparison rows as a check of the matrix | research |
| Programming-Idioms | content CC BY-SA 3.0, code Apache-2.0 | B for text, C for code (but our code is written fresh) | idiom names and ids (`crosswalk`) | research |
| Azure Architecture Center patterns | CC BY 4.0 | C | pattern names and facts, with credit | research |
| Rosetta Code | GFDL 1.2 | D, and E for automated access while the site returns 403 to bots | task names only (10.1) | research |
| FOLDOC | GFDL 1.1 or later (15,284 entries) | D | terms as a check of the definition list | research |
| SEVOCAB (ISO/IEC/IEEE 24765) | definitions may be copied with "Copyright ©, 2026, IEEE. Used by permission." | C for single definitions with that notice | a definition may be quoted on a `definition` page, with the exact notice | research |
| SWEBOK v4 | no redistribution or alteration | D | area and topic names (`crosswalk/swebok-v4.yaml`) | research |
| CS2023 | all rights reserved | D | knowledge area and unit names (`crosswalk/cs2023.yaml`) | research |
| ACM CCS 2012 | free for educational and research use with credit | D | ids and names with credit | research |
| Design Patterns (GoF), POSA, PoEAA, Enterprise Integration Patterns, Fowler's refactoring catalogue, refactoring.guru, microservices.io, the c2 wiki | all rights reserved, or unclear | D | names, the existence of a pattern or refactoring, the structure of the catalogue, citations with page or section; never text, figures or code | research |
| Pro Git | CC BY-NC-SA 3.0 | D | facts, chapter citations | research |
| PLAI | CC BY-NC-SA 4.0 | D | facts, chapter citations | research |
| MIT Missing Semester | CC BY-NC-SA 4.0 | D | topic names, links | research |
| roadmap.sh | personal use only; redistribution needs consent | D | nothing beyond facts and a link | research |
| freeCodeCamp curriculum | no named licence (the software is BSD-3-Clause) | D | links only | research |
| textbooks, papers, talks, blogs | all rights reserved unless stated | D | facts, citations with locators, quotes under G04 | general |

`registers/sources.yaml` holds this table as data: the source id, the licence (an SPDX id or `LicenseRef-<name>`), the class, the date checked, the URL of the licence statement, and what was taken. Rule C08 fails a page that cites a source with no row.

## 19.4 Quotes and paraphrase

- **Quotes.** G04 (06.6) applies: a quote is at most 50 words, marked as a quote, with the source, the locator, and a source class that allows it (A, B or C, or D under fair use or fair dealing for a short quote with a clear purpose such as the exact wording of a rule). A page has at most 3 quotes from one class D source. A definition from SEVOCAB carries its notice.
- **Normative wording.** The exact words of a standard matter in some places (for example a clause that defines undefined behaviour). The page may quote the sentence under G04 and must give the clause. It may not reproduce a list or a table of a class D standard.
- **No close paraphrase.** A section that follows one source sentence by sentence, in the same order with the words changed, is a copy even if no sentence matches. The critic checks for it (15.4), and the copy check measures it: shingled n-gram overlap for prose against the raw snapshot of each cited source, and token-level winnowing for code against the code corpora of the sources (Rosetta Code, Programming-Idioms, Learn X in Y minutes, MDN samples, Go by Example, the official tutorials). A change over the threshold of C01 (Default: shingled 5-gram overlap above 15% of a section for prose; a winnowing match of more than 10 consecutive non-trivial lines for code) fails rule C01 and goes to a person.
- **Structure.** The order of topics in a book or a reference is a fact about that source and may inform our taxonomy, but our area tree and our page sections follow our own templates (05, 06.3), not one book's table of contents.
- **Code.** Common idioms (`for i in range(n)`, `if err != nil { return err }`) are not protected and appear in every language's code. What counts as copying is the same sequence of non-obvious choices in a whole example, the same worked example with the same names and data, or the same comments.

## 19.5 Content written by agents

- **Provenance.** Each file change has a provenance record (04.7) with the author: a GitHub login or an agent run id. The run record names the model id, the prompt version and the sources the agent was given.
- **Disclosure on pages.** The site labels each page by how it was made, from the provenance and the review records: "Written by AI, checked by the build" (W4 without review), "Reviewed by <name>" (W5), or "Written by <name>". The labels are about trust, and they are also the honest statement that the law on the authorship of machine output is not settled. The About page explains the pipeline (14) and its checks.
- **Licence of agent output.** Agent output is published under the licence of its zone, with "The programming contributors" as the copyright holder, as in algorithm 17.4. Where the output has no author in law, the notice does no harm, and the human selection, review and arrangement of the corpus is still the contributors' work.
- **Upstream policies on AI text.** Some sources forbid AI-written contributions to their own sites (Stack Overflow's policy on generated answers, Wikipedia's guidance on LLM-generated text, several open-source projects' contribution rules). The project does not post AI-written text to any upstream site. When a writer finds an error in an upstream source, a person reports it (14.9), not an agent.
- **No training claims.** The site does not say a page is "original" beyond what the copy check measured.
- **Model terms.** The project uses model providers whose terms give the user the rights in the output (14.2, `models.yaml`).

## 19.6 Trademarks

- **Nominative use only.** Names of languages, tools, companies and products are used to identify them, in the form their owners write them (06.5). The site does not use logos, mascots or brand colours of languages or products (Default: no logos at all; a language page may show the language's own name in plain text only).
- **No sign of endorsement.** The About page and the footer say that the site is independent and not endorsed by any language's steward.
- **Policies of the main stewards** (each checked at M1; the summaries are from general knowledge and unverified in detail):

| Mark | Steward | What it means for us |
|---|---|---|
| Python | Python Software Foundation | the word mark may be used to refer to the language; no logo use without permission; no names that suggest an official PSF site |
| Java | Oracle | "Java" is a trademark of Oracle; refer to it as such on `/about/legal/`; do not use "Java" in a product or domain name |
| Rust | Rust Foundation (trademark policy revised after the 2023 consultation) | the word may be used to refer to the language; the logo is not used |
| Go | Google (Go brand and trademark guidelines) | the word may be used to refer to the language; the Gopher artwork has its own licence (CC BY 4.0 per its author) but we do not use it |
| JavaScript | Oracle holds the US mark | refer to the language as JavaScript or ECMAScript in the normal way |
| CWE, CVE | The MITRE Corporation | CWE and CVE are trademarks; the notice is on `/about/legal/` |

- The site's own name and the URL do not contain any language's mark (`programming.tamnd.com`).

## 19.7 Security content

The `security` area (05.2) and the weakness nodes (04.2) teach defects so that programmers can avoid them. The policy:

- **Teaching examples only.** Vulnerable code is the smallest program that shows the weakness, with a local fixture (10.8). It runs only in the security harness: no network, local fixtures only, no real targets.
- **No weaponized exploits.** The project does not publish exploits for real product versions, exploit chains, payloads meant for evasion, or scanners meant to find targets. An exploit input in `vuln/exploit.jsonl` must only work against the node's own vulnerable fixture.
- **CVE items.** A CVE card holds the id, the dates, the affected versions, the weakness it belongs to, and the link to the vendor or project advisory and the NVD record. It does not link to exploit code or proof-of-concept repositories.
- **Responsible disclosure.** If a writer or an agent finds what looks like a real vulnerability in a real project while working (for example in a library's documented example), the run stops that item and a person reports it privately through the project's security policy, `SECURITY.md`, or the vendor's channel. Nothing about it is published until the project has fixed it or 90 days have passed (Default), and then only as a CVE card linking the advisory.
- **Our own security policy.** `SECURITY.md` in both repositories says how to report a security problem in the site or the toolchain (for example, a sandbox escape in the runner).
- **Cryptography pages** say that their code is for learning and is not hardened against side channels, and they point to the vetted libraries of each language.

## 19.8 Privacy

- **Readers.** The site sets no cookies and loads nothing from a third party (16). Preferences and playground state live in the reader's browser only. The playground never sends code to our servers (07.10). Page view counts come from Cloudflare's aggregated request logs, by path, with no identifiers kept by the project. The search "no results" counts are off by default (17).
- **"Open in" links** send code to a third-party playground only on the reader's click, with a notice naming the third party (07.10).
- **Contributors.** Names and e-mail addresses are in git history, as in any public repository. Provenance logs store the GitHub login only.
- **People pages (living persons).** The rules follow the spirit of Wikipedia's policy on biographies of living persons:
  - only public facts about published work and public roles, each with a reliable source: names, years, affiliations, the items they made or named, and the dates;
  - no contact details, no home location, no family, no health, no private life, no social-media quotes that are not about their published work;
  - no criticism of a living person except criticism of their published work, from a reliable source, in neutral words;
  - a person who is the subject of a page can ask for a correction, or for the page to be cut to the list of their items (19.10); the default answer to such a request is yes;
  - a person node about a living person needs a person review (W5) before it shows anything beyond the card (Default).
- `/about/privacy/` says all of this in plain words.

## 19.9 Export rules for cryptography code

As in algorithm 17.5: code examples that use or implement cryptography are public, open source code. Under the US Export Administration Regulations, publicly available encryption source code is not subject to the EAR (§734.3(b)(3) and §742.15(b)); the 2021 rule removed the e-mail notice for most such code but kept it for "non-standard cryptography". Most examples here only call standard libraries, which raises no question. Default: if the corpus ever holds an implementation of a non-standard scheme (for example a teaching cipher), the project sends the notice once and keeps its text in `docs/legal/`. This is unverified legal reading; a lawyer checks it at M1.

## 19.10 Takedown, corrections and contact

**How to report.** Each page has "Report a problem" (16), which opens a GitHub issue with the URL filled in and a choice: wrong, outdated, unclear, copied, rights, security, person, other. A rights holder or a person who does not want a public issue writes to the address on `/about/legal/` (Default: an alias that forwards to the maintainers). Security reports go to the channel in `SECURITY.md`.

| Report | Within | Action |
|---|---|---|
| a claim that text or code was copied | 2 working days | a person runs the copy check against the named source by hand. If the claim holds, the part is quarantined: it leaves the site at the next deploy, and a fresh version is queued with the source added to the copy-check corpus. The issue records the outcome |
| a DMCA notice or the like | 2 working days | the same, following the notice's procedure; the file is removed from the release archives in the next release; the notice is kept in `docs/legal/notices/` with personal details removed |
| a source asks the project to stop harvesting | 2 working days | the source's class becomes E, its harvester is turned off, its raw snapshots are deleted from R2 `programming-raw` unless needed as evidence of a dispute; its facts in cards stay (they are facts), with the citation changed to "source withdrawn" where the source asks |
| a trademark complaint | 5 working days | a person checks the use against 19.6 and changes the page if the use is not nominative |
| a request about a living person | 5 working days | the page is fixed or cut (19.8) |
| a security report about the site or the toolchain | per `SECURITY.md` | private handling, then a fix and a note |
| an error in content | no deadline | it enters the work queue with a high weight (14); the reader is credited on `/reviewers/` if they wish |

Removed content stays in git history, which is public. When a court order or a valid notice requires it, the maintainers rewrite history for that path only and publish a note that they did. A contributor whose work was removed may answer a notice in the same issue; the maintainers restore the work when the claim is withdrawn or when two maintainers find no copying.

## 19.11 Rules

Listed in 15 with the other rule families:

| Rule | Checks |
|---|---|
| C01 | the copy check: prose and code below the overlap thresholds of 19.4 against every cited source and the code corpora |
| C02 | `reuse lint` passes: every file has a licence from `REUSE.toml` |
| C03 | every source with `adapted: true` is class A, B or C, and the page footer carries the attribution that the licence asks |
| C04 | a person page's photo has a Wikimedia Commons licence of CC0, CC BY or CC BY-SA and a credit line; otherwise there is no photo |
| C05 | an exercise card or a task page has no run of more than 12 words that matches the captured statement of the exercise's site |
| C06 | no file in the MIT zone (code, probes, examples' code blocks) has its origin in a class B or D source in `prov/` (`op: adapt` with such a source) |
| C07 | no text field of a card or a register (`summary`, `title`, `note`) has its origin in a class B, C (attribution-required) or D source; card text is written fresh or comes from a class A source |
| C08 | every cited source has a row in `registers/sources.yaml` with a licence and a class |
| C09 | no page links to exploit code or a proof-of-concept repository; CVE cards link advisories only (19.7) |
| C10 | a person node about a living person has no field outside the allowed list of 19.8 |
