# 05. Taxonomy

This document defines the two trees that hold the nodes, the language template, the levels, the core list, the crosswalks to outside classifications, and the border with algorithm and math2.

## 5.1 Two trees and one rule

Programming knowledge has two axes. One is the **idea** (closures, error handling, dependency injection, data races). The other is the **language** (how Python does it, how Rust does it). A single tree cannot hold both without duplicating one of them.

So there are two trees:

| Tree | Root | Holds | Example path |
|---|---|---|---|
| **Topic tree** | `nodes/<top>/` | concepts, definitions, tasks, patterns, techniques, practices, tools, libraries, formats, standards, weaknesses, works, people, organizations | `nodes/functions/closures-and-scope/closures/` |
| **Language trees** | `nodes/languages/<lang>/` | the language node and its features, one tree per language of tiers 1 to 3 | `nodes/languages/python/control-flow/match-statement/` |

The rule that connects them: **a fact about one language goes in that language's tree, and the topic page shows it through `realizes` edges.** A concept page never holds a long section about one language that a feature page should hold. The concept page's language section is short (the example, the one-line verdict, the differences), and it links the feature page for the rules.

Releases, diagnostics, api cards and proposals have owners (04.3), so they hang under their owner's language, tool or library, but they are cards and do not need a directory.

## 5.2 The topic tree

`topics/areas.yaml` holds the tree, with the same fields as algorithm (`id`, `name`, `summary`, `children`) and two new fields: `kinds` (the kinds this area holds most) and `crosswalk` (the CS2023, SWEBOK v4 and ACM CCS ids it maps, 5.7).

The top level has 36 areas, in 6 groups. The group is for the home page and the book; it is not part of the URL.

**A. Foundations of programs**

| id | Name | Holds (examples of sub-areas) |
|---|---|---|
| `foundations` | Foundations | what a program is, computation models in practice, bits and representation, numbers in machines (integers, overflow, floating point, decimal), characters and text, the machine model (CPU, memory hierarchy, instructions), the stack and the heap, compilation and interpretation at an overview, hello world in every language |
| `paradigms` | Paradigms | imperative, procedural, structured, object-oriented, functional, logic, declarative, dataflow, reactive, array, concatenative, event-driven, aspect-oriented, actor, multi-paradigm design |
| `syntax-and-semantics` | Syntax and semantics | lexical structure, grammars and notation (BNF, EBNF, PEG), expressions and statements, evaluation order, scope and binding, operational/denotational/axiomatic semantics at the practitioner level, undefined and implementation-defined behaviour |
| `types` | Types | type systems, static and dynamic checking, inference, generics and parametric polymorphism, subtyping and variance, algebraic data types, gradual typing, dependent and refinement types in practice, nullability, type classes and traits, type erasure and reification, units of measure |
| `data-and-values` | Data and values | primitive types, strings and string handling, collections in languages (arrays, lists, maps, sets, tuples, records), value and reference semantics, mutability and immutability, equality and identity, hashing in practice, serialization, copying (shallow and deep), dates and times, money and decimals |
| `functions-and-abstraction` | Functions and abstraction | functions and procedures, parameters and arguments (by value, by reference, by name, defaults, keyword, variadic), recursion and tail calls, first-class and higher-order functions, closures and scope, currying and partial application, lazy evaluation, continuations, generators and coroutines, modules and namespaces, information hiding |
| `objects-and-classes` | Objects and classes | objects, classes, prototypes, encapsulation, inheritance (single, multiple, mixins), interfaces and protocols, polymorphism and dispatch (single, double, multiple), composition over inheritance, constructors and initialization, object lifecycle, the expression problem |
| `memory-and-resources` | Memory and resources | allocation, manual memory management, garbage collection (tracing, reference counting, generational, concurrent), ownership and borrowing, RAII and deterministic destruction, smart pointers, memory safety, leaks, arenas and pools, resource handles (files, sockets, locks) |
| `errors-and-control` | Errors and control flow | control structures, exceptions, error values and result types, panics and aborts, assertions and contracts, non-local exits, cleanup (finally, defer, using, with), retries and timeouts |
| `concurrency-and-parallelism` | Concurrency and parallelism | processes and threads, shared memory and locks, atomics and memory models, message passing, async/await and event loops, futures and promises, actors, CSP and channels, structured concurrency, data parallelism, GPU programming in languages, races and deadlocks |
| `metaprogramming` | Metaprogramming | reflection and introspection, macros (textual, hygienic, procedural), templates and compile-time evaluation, code generation, annotations and attributes, decorators, eval, homoiconicity, DSLs embedded in a host language |

**B. Making languages**

| id | Name | Holds |
|---|---|---|
| `language-implementation` | Language implementation | lexing and parsing in practice, ASTs and IRs, interpreters (tree-walking, bytecode), compilers (front, middle, back end), JIT compilation, virtual machines, runtime systems, linking and loading, ABIs and calling conventions, garbage collector implementation, WebAssembly as a target, language servers |
| `languages` | Languages | the catalogue of all languages by family and lineage (the Algol family, the C family, Lisps, MLs, Smalltalk, the shell family, query languages, array languages, esoteric languages, ...), the genealogy, language design principles, the history of languages; the language trees hang here (5.3) |

**C. Building software**

| id | Name | Holds |
|---|---|---|
| `software-design` | Software design | design principles (SOLID, DRY, KISS, YAGNI, cohesion and coupling, separation of concerns), design patterns (creational, structural, behavioral), architecture styles and patterns, domain-driven design, API design, modularity, interfaces and contracts, design for change |
| `construction-practice` | Construction practice | naming, readable code, comments and documentation in code, code style and formatting, code smells, refactoring (the catalogue of refactorings), code review, technical debt, defensive programming, working with legacy code, pair and mob programming |
| `testing` | Testing | unit, integration, system and end-to-end testing, test doubles, test design techniques, property-based testing, fuzzing, mutation testing, snapshot and golden testing, coverage, TDD and BDD, test frameworks per language, flaky tests, testing concurrency |
| `debugging-and-observability` | Debugging and observability | debugging methods, debuggers, breakpoints and watchpoints, core dumps, logging, tracing, metrics, profiling as diagnosis, sanitizers, reproducible bug reports, post-mortems |
| `tools-and-workflow` | Tools and workflow | the shell and the command line, editors and IDEs, version control (Git in depth), build systems, package managers and dependency management, continuous integration and delivery, containers for developers, documentation tools, code search, dotfiles and environment management |
| `program-analysis-and-verification` | Program analysis and verification | static analysis, linters and type checkers as analyzers, abstract interpretation in practice, symbolic execution, model checking, deductive verification (Dafny, Why3, Frama-C, SPARK, Verus), proof assistants for programs, design by contract, runtime verification |
| `performance` | Performance | measuring (benchmarks, profilers), complexity in practice (linked to algorithm), caches and memory layout, allocation cost, vectorization, compiler optimizations, concurrency for speed, I/O performance, latency and throughput, performance of each runtime |
| `security` | Secure programming | the weakness catalogue (CWE), injection, memory-safety weaknesses, authentication and session handling in code, cryptography use (not the mathematics), secrets handling, input validation and output encoding, supply-chain security, secure defaults per language, sandboxing, threat modelling for programmers |
| `software-engineering-process` | Software engineering process | requirements in practice, estimation, methodologies (agile, Scrum, Kanban, XP, waterfall), releases and versioning, maintenance, metrics, documentation, team practice, laws of software (Brooks, Conway, Hyrum, ...), ethics and professional practice |
| `open-source-and-licensing` | Open source and licensing | free and open-source software, licences (every SPDX licence as a card), contribution practice, governance, foundations, copyright and patents for programmers, export rules at an overview |

**D. Data, systems and the network**

| id | Name | Holds |
|---|---|---|
| `data-and-persistence` | Data and persistence | SQL in depth, relational modelling, transactions and isolation, ORMs, NoSQL models, data formats (JSON, YAML, TOML, XML, CSV, Protocol Buffers, ...), regular expressions, text encodings and Unicode, files and file systems from a program, caching |
| `systems-programming` | Systems programming | the OS interface (system calls, POSIX, Win32), processes and signals, file descriptors and I/O models, memory mapping, IPC, networking sockets, low-level C and Rust, assembly language, linkers and object files, the kernel from a programmer's view |
| `web-development` | Web development | HTTP for programmers, HTML, CSS, the DOM and Web APIs, front-end frameworks, back-end frameworks, REST, GraphQL and RPC, authentication flows, web performance, accessibility in code, browsers as runtimes, WebAssembly on the web |
| `mobile-and-desktop` | Mobile and desktop | Android, iOS, cross-platform frameworks, desktop GUI toolkits, UI programming models (immediate, retained, declarative), app packaging and distribution |
| `distributed-and-cloud` | Distributed and cloud programming | distributed systems for programmers (RPC, consensus as used, idempotency, retries), microservices, message queues and streams, cloud platforms and services as programming targets, infrastructure as code, serverless, Kubernetes for developers, site reliability engineering practice |

**E. Domains**

| id | Name | Holds |
|---|---|---|
| `data-science-and-ml-programming` | Data science and ML programming | notebooks, dataframes, array programming (NumPy and its kin), plotting libraries, ML frameworks as programming tools (PyTorch, JAX, TensorFlow), data pipelines, reproducibility; the learning theory itself links out |
| `ai-assisted-programming` | AI-assisted programming | code completion, coding agents, prompting for code, tool use and MCP, evaluation of generated code, code LLM benchmarks, review of AI-written code, licensing and provenance of generated code, agent workflows (this project is one) |
| `graphics-and-games-programming` | Graphics and games programming | graphics APIs (OpenGL, Vulkan, Metal, Direct3D, WebGPU), shaders and shading languages, game loops, game engines as programming environments, audio programming |
| `scientific-and-hpc` | Scientific and high-performance computing | Fortran, Julia and numerical languages, MPI, OpenMP, CUDA and other accelerators, numerical pitfalls, reproducible scientific software |
| `embedded-and-hardware` | Embedded and hardware | microcontrollers, real-time programming, interrupts, bare-metal C and Rust, HDLs (Verilog, VHDL, Chisel), FPGA programming, PLC languages (IEC 61131-3), safety-critical coding standards (MISRA C, AUTOSAR C++14, JSF) |
| `domain-programming` | Domain programming | spreadsheets and formula languages, financial programming, bioinformatics programming, smart contracts, robotics (ROS), automation and scripting, education languages (Scratch, Logo), creative coding |

**F. People and history**

| id | Name | Holds |
|---|---|---|
| `history-and-culture` | History and culture | the history of programming by period, landmark works (the canon), people, organizations, events, folklore and jargon, programming culture and communities |
| `learning-programming` | Learning programming | how people learn to program, curricula (CS2023, the K-12 frameworks), common misconceptions (each as a node linked to the concept it is about), exercises and practice sites, notional machines |

Notes:

- The topic tree is a **Draft**. 21 (milestones) M1 fixes it. Each top area gets a `README` page and a sub-area list in `topics/areas.yaml` at M1, made by the catalog step (14.4).
- The sub-areas are 4 to 25 per top area. Estimate: 600 sub-areas over the tree.
- An item has one primary area and may have secondary areas (04.5). The rule for the primary area: the area whose sub-area name a reader would look under first. When two areas tie, the more specific wins (`sql-injection` goes under `security/injection`, not `data-and-persistence/sql`).

## 5.3 The language trees and the language template

Every language of tiers 1 to 3 (08.2) has a tree under `nodes/languages/<lang>/`. The sub-areas of a language tree are fixed by the **language template**, so that every language guide has the same shape, readers can move from one language to another at the same place, and the coverage report can show gaps per sub-area.

| # | Sub-area id | Name | Holds (features, api groups) | Tier 1 | Tier 2 | Tier 3 |
|---|---|---|---|---|---|---|
| 1 | `overview` | Overview and history | the language node, its design goals, history, editions and releases, implementations | ✓ | ✓ | ✓ |
| 2 | `getting-started` | Getting started | install, toolchain, hello world, the REPL, running a program, project layout | ✓ | ✓ | ✓ |
| 3 | `syntax` | Lexical structure and syntax | tokens, comments, identifiers, literals, operators and precedence, statements and expressions, the grammar | ✓ | ✓ | ✓ |
| 4 | `values-and-types` | Values and types | the built-in types, numbers, text, booleans, null and its friends, type conversions, the type system of the language | ✓ | ✓ | ✓ |
| 5 | `variables-and-scope` | Variables, binding and scope | declarations, constants, scope rules, shadowing, lifetime | ✓ | ✓ | ✓ |
| 6 | `control-flow` | Control flow | conditionals, loops, pattern matching, jumps, labelled blocks | ✓ | ✓ | ✓ |
| 7 | `functions` | Functions | definitions, parameters, returns, closures, recursion, lambdas, overloading | ✓ | ✓ | ✓ |
| 8 | `collections` | Built-in collections | arrays, lists, maps, sets, tuples, slices, ranges, iteration protocols, comprehensions | ✓ | ✓ | ✓ |
| 9 | `strings-and-text` | Strings and text | string types, encodings, formatting and interpolation, regex support, Unicode handling | ✓ | ✓ | |
| 10 | `user-types` | User-defined types | structs, records, classes, enums, unions, interfaces, traits, protocols, generics | ✓ | ✓ | ✓ |
| 11 | `errors` | Errors and exceptions | the error model, the exception hierarchy, error values, panics, cleanup | ✓ | ✓ | ✓ |
| 12 | `memory` | Memory and resources | allocation, the GC or ownership model, destructors, resource handling, unsafe code | ✓ | ✓ | |
| 13 | `modules-and-packages` | Modules, packages and builds | the module system, imports, visibility, packages, the package manager, the build tool, publishing | ✓ | ✓ | ✓ |
| 14 | `concurrency` | Concurrency and parallelism | threads, async, channels, the memory model of the language, synchronization in its stdlib | ✓ | ✓ | |
| 15 | `metaprogramming` | Metaprogramming | macros, reflection, annotations, decorators, code generation, compile-time evaluation | ✓ | ✓ | |
| 16 | `stdlib` | The standard library | a map of the standard library; the api cards hang here (one api group per module or package) | ✓ | ✓ | ✓ |
| 17 | `io-and-os` | I/O, files and the OS | files, paths, processes, environment, time, networking from the stdlib | ✓ | ✓ | |
| 18 | `testing` | Testing | the test frameworks, assertions, mocking, test runners and conventions | ✓ | ✓ | |
| 19 | `tooling` | Tooling | compilers and interpreters, debuggers, profilers, linters, formatters, LSP, notebooks | ✓ | ✓ | |
| 20 | `ecosystem` | Ecosystem | the main libraries and frameworks by need, the community, the governance and proposal process | ✓ | ✓ | |
| 21 | `performance` | Performance | the cost model, optimization guidance, the profiling workflow, the compilation and runtime knobs | ✓ | | |
| 22 | `interop` | Interoperability | FFI, calling C, embedding, the ABI, interop with the platform languages | ✓ | ✓ | |
| 23 | `idioms-and-style` | Idioms and style | the style guide, the idioms, what code in this language "should" look like, naming conventions | ✓ | ✓ | ✓ |
| 24 | `pitfalls` | Pitfalls | the known traps, each as a feature or a short node, linked to diagnostics and weaknesses | ✓ | ✓ | |
| 25 | `versions` | Versions and migration | the release nodes, the migration guides between major versions, deprecations and removals | ✓ | ✓ | ✓ |

Rules:

- A ✓ means the sub-area must exist for the language at its tier, with at least its README and the features that the catalog step finds (14.4). A missing ✓ sub-area fails rule E12 at C3 for the language.
- A language may add sub-areas after the 25 (for example `ownership-and-borrowing` for Rust, `templates` for C++, `pointers` for C and Go, `the-borrow-checker` stays a feature, `jvm` for Java and Kotlin, `ffi-and-unsafe` for Rust). An added sub-area is in `languages/<lang>/areas.yaml`, with a reason.
- A sub-area that does not apply to a language (for example `memory` for SQL) says so in its README, with one line of reason, and has no features. It counts as done.
- Domain languages of tier 1 (SQL, shell, regex, HTML and CSS, PowerShell, 08.2) use a **domain template** that replaces some rows. SQL: `queries`, `joins`, `aggregation`, `window-functions`, `ddl`, `transactions`, `indexes`, `procedural-extensions`, `dialects`. Regex: `syntax`, `character-classes`, `quantifiers`, `groups-and-backreferences`, `lookaround`, `flags`, `flavours`, `performance-and-redos`. Shell: `commands-and-pipelines`, `quoting-and-expansion`, `variables`, `control-flow`, `functions`, `redirection`, `job-control`, `portability`, `the-core-utilities`. 08.5 lists all of them.

## 5.4 The core list

The core is the set of items that must be complete before the site is called "complete for the core". It answers "what does every programmer need?" It is the measure for M2 and M3 (21).

| Core part | What | Size (Estimate) | Source of the list |
|---|---|---|---|
| Core concepts | concepts and definitions that appear in at least 2 of: the CS2023 Software Development Fundamentals and Foundations of Programming Languages knowledge units, SWEBOK v4 Software Construction, the index of 3 of the canon textbooks (13.3), the concept lists of 3 tier-1 language references | 400 | the canon match step (14.5), as in algorithm |
| Core features | for each tier-1 language, the features that its own tutorial or "tour" covers | 120 per language, 1,800 for 15 languages | the official tutorials (Python tutorial, The Rust Book, A Tour of Go, Java tutorial, ...) |
| Core tasks | the tasks of the core task list: the intersection of Programming-Idioms, Rosetta task names, Go by Example and Rust by Example titles, plus the "how do I" queries of 17.9 with high volume | 300 | 10.1 |
| Core patterns | the GoF 23, the POSA and PoEAA core, the main concurrency and resilience patterns, the main smells and anti-patterns | 150 | the canon (13.3) |
| Core refactorings and smells | the refactoring catalogue and the code smells | 66 refactorings and 24 smells | Fowler's catalogue (refactoring.com) and the smells of *Refactoring* 2nd ed., names only |
| Core weaknesses | CWE Top 25 (2025) and the OWASP Top 10 (2025) | 35 | CWE and OWASP |
| Core tools | the tools in every tier-1 language's toolchain row, Git, the main editors, the main CI services, the main debuggers and profilers | 120 | `registers/languages.yaml` and the developer surveys |
| Core formats and standards | JSON, YAML, TOML, XML, CSV, Protocol Buffers, UTF-8 and Unicode, HTTP, URI, semver, the main language standards | 40 | by hand, with reasons |
| Core practices | principles, laws and methods that the canon names | 80 | the canon |

Total core: about 3,000 items (Estimate). The core is computed and stored in `registers/core.yaml` with the reason for each item, as in algorithm 05.5. A person approves it at the end of M1.

## 5.5 Order of sub-areas and items

As in algorithm 05.4: the sub-areas of an area and the items of a sub-area are listed in a **reading order**, not alphabetically. The order is from the prereqs DAG first, then by level, then by a hand-set `order` field. Language trees follow the template order (5.3).

## 5.6 Levels

| Level | Who it is for | Example |
|---|---|---|
| `intro` | a person in the first year of programming | variables, loops, functions, strings, hello world |
| `core` | a working programmer | closures, exceptions, generics, Git branching, unit tests, SQL joins |
| `advanced` | a senior or specialist programmer | memory models, variance, macros, lock-free code, JIT tiers |
| `research` | a PL or SE researcher, or the edge of practice | effect handlers, session types, gradual typing soundness |

The level is a field of the front matter (04.5), set by the curator and checked by the critic. It controls the reading order and the default view of search (17.4).

## 5.7 Crosswalks

The build keeps crosswalks to the classifications that readers and other sites use. Each crosswalk is a register; each row maps one outside id to one or more item tags, with the kind of match (exact, broader, narrower, related).

| Crosswalk | Outside system | Licence of the outside ids | Used for |
|---|---|---|---|
| `crosswalk/cs2023.yaml` | ACM/IEEE/AAAI CS2023 knowledge areas and units (SDF, FPL, SE, SEP, PDC, SEC, SF, ...) | names of a public curriculum, cited | the coverage report and the learning paths (16.12) |
| `crosswalk/swebok-v4.yaml` | SWEBOK v4 knowledge areas and topics (18 areas) | names cited | the coverage of the software-engineering areas |
| `crosswalk/acm-ccs-2012.yaml` | ACM CCS 2012, "Software and its engineering" and "Theory of computation › Semantics" subtrees | free to use with credit | the link from papers to items |
| `crosswalk/cwe.yaml` | CWE ids, all versions since 4.0 | CWE terms of use (free with credit) | the weakness kind |
| `crosswalk/wikidata.yaml` | Wikidata QIDs for each item that has one | CC0 | `same_as`, people and languages |
| `crosswalk/pldb.yaml` | PLDB ids | public domain | language cards (08.3) |
| `crosswalk/linguist.yaml` | GitHub Linguist language names | MIT | file extensions and syntax colour |
| `crosswalk/so-tags.yaml` | Stack Overflow tags with counts | CC BY-SA (counts are facts) | search terms and the "missing" signal (18.6) |
| `crosswalk/rosetta.yaml` | Rosetta Code task names | names only (13.5) | the task universe (10.1) |
| `crosswalk/spdx.yaml` | SPDX licence and exception ids | CC0 | the licence cards |
| `crosswalk/algorithm.yaml` | algorithm.tamnd.com tags | CC0 (our own) | the border (5.8) |

Rules:

- The crosswalks drive the "no missing" checks of 18.4. A row of an outside list that maps to no item is a gap. The universe register (18.3) takes it in with an outcome.
- An outside system's text is never copied into a node. Only ids and names are kept, with the source.

## 5.8 The border with algorithm and math2

The three sites are one library with three wings. The border rules:

| Question | Goes to | Example |
|---|---|---|
| Is it a procedure that solves a problem, with correctness and cost as the subject? | algorithm | merge sort, Dijkstra, a hash table as a data structure, consistent hashing |
| Is it how a language, a library or a tool does something, or how programmers work? | programming | Python's `sorted` and `key=`, `java.util.HashMap`'s API and its resize behaviour, Go maps and iteration order |
| Is it a mathematical object or theorem? | math2 | lambda calculus as mathematics, the Curry-Howard correspondence as a theorem, Hoare logic soundness |
| Is it the programmer's use of a mathematical idea? | programming, with a link | type inference in practice, Hoare triples in a code review, floating point rounding in code |

- Items that sit on both sides get a node on each site, each about its own side, and the two link each other with `same_as.algorithm` or `same_as.math2`. Example: "hash table" on algorithm (the data structure, its analysis) and "maps and dictionaries" on programming (the language facilities, iteration order, hashing of keys, the pitfalls).
- The universe register (18.3) records the outcome `elsewhere:algorithm` or `elsewhere:math2` for an item that belongs only to the other site, with the target tag. Such items are linked from search and from the area pages, so a reader of programming who looks for "quicksort" finds the way to it.
- The `algorithm` register of languages (`registers/languages.yaml`) and the programming language nodes stay in sync: the programming language node is the source of truth for toolchain versions, and the algorithm register copies from it at each release of programming (20.6).
