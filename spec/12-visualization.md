# 12. Visualization

A visualization shows what a program does inside the machine: the frames on the call stack, the objects on the heap and the references between them, the moment a value is moved or borrowed, the order in which threads and tasks run, the bytes of a struct. These are the things that a programmer must imagine and that prose explains badly. This document extends the visualization system of algorithm (2150/10: the trace protocol, the primitives, the player, the checks) with the primitives that programming needs, with traces recorded from real runs, with diagrams generated from code and data, and with screenshots. Where this document is silent, 2150/10 applies.

## 12.1 Principles

1. **The picture is the real run.** A scene of an example shows the run that CI executed (07.4): the same code, the same image, the same output. A picture that the build cannot tie to a run is a figure (V1), not a scene (V2).
2. **Record first, author second.** Where a tracer exists for the language (12.4), the scene is a recorded trace of the example, drawn by the player. An author writes a scene by hand only for models that no tracer can see (the event loop's queues inside a runtime, a borrow, a cache line), and the check of 12.11 ties it to a real run.
3. **The algorithm principles stay.** The reader drives, one step is one line, trace then draw, the same primitives everywhere, readable without colour, motion or script (2150/10.1).
4. **Show the model the page teaches.** The heap view of a Python page shows names bound to objects. The heap view of a C page shows bytes at addresses. The primitives have views for each model (12.3), and a scene uses the view of its language's model.
5. **Generated, never drawn by hand.** Class diagrams, sequence diagrams, railroad diagrams, query plans, git graphs and the genealogy come from code, grammars, tool output or the graph. A diagram file drawn by hand in `figures/` is allowed only for ideas that have no data source, and it must have a text alternative (V09).

What we learn from others, without copying their code or assets (19): Python Tutor (Philip Guo) shows that a view of frames and objects teaches references and aliasing better than prose; the Rust community's ownership diagrams show that a timeline of moves and borrows makes the borrow checker readable; explain.depesz.com and similar tools show that a query plan is easier to read as a tree with costs than as text.

## 12.2 Scene sources

| Source | What it is | Used for | Tied to the run by |
|---|---|---|---|
| **recorded** | a trace that a tracer writes while the example runs in its image (12.4) | stack, heap, references, closures, recursion, exceptions, generators, threads in languages with a tracer | the same CI run: the trace and the output block come from one execution (V03) |
| **authored** | `viz/scene.ts`, as in algorithm, using `@prog/viz` | the event loop, ownership and borrowing, memory models, GC phases, isolation anomalies in SQL, anything a tracer cannot see | the line-sync check against the real run (V06) and the declared failures (V07) |
| **generated** | a diagram that the build makes from a source: code, a grammar, tool output, the graph | class and sequence diagrams, ASTs, railroad diagrams, query plans, the git DAG, memory layout, the genealogy (12.5, 12.6) | the generator runs on the node's own code or data in the image |
| **captured** | a screenshot of a GUI or a rendered page (12.7) | GUI examples, HTML and CSS rendering | the scheduled job that runs the example |

`@prog/viz` is `@prog/viz` (the same package source, copied as in 03) plus the primitives of 12.3. It lives in `tamnd/programming-agent` under MIT. The trace protocol (steps, events, `line`, `note`, `vars`, `phase`, counters) and the state vocabulary (`active`, `compare`, `visited`, `path`, `rejected`, `invariant`, ...) are the same as in algorithm, so the player is one player.

## 12.3 New primitives

| Primitive | Draws | Views | Events (beyond the algorithm set) | Used by (examples) |
|---|---|---|---|---|
| `frames` | the call stack: one box per frame with the function name, the line, the local names and their values | `column` (newest on top), `tree` (for recursion trees) | `call(fn, args)`, `ret(value)`, `unwind(exc)`, `bind(frame, name, ref)`, `rebind`, `del` | functions, recursion, scope, closures, exceptions |
| `heap` (objects) | objects as boxes with a type and fields; references as arrows from names and fields | `names` (Python, JS, Java, C#, Go: names point to objects), `boxes` (values inline, pointers as arrows: C, C++, Rust, Go structs), `bytes` (a memory range with addresses) | `alloc(id, type, fields)`, `field(id, k, ref)`, `free(id)`, `collect(ids)`, `refcount(id, n)`, `addr(id, address)` | references, aliasing, copies, mutability, GC, pointers, dangling pointers |
| `ownership` | one lane per value or binding, with spans for "owned", "borrowed shared", "borrowed mutable", "moved out" | `timeline` | `own(var, value)`, `move(from, to)`, `borrow(var, kind, span)`, `drop(var)`, `conflict(line, code)` | ownership, moves, borrowing, lifetimes, RAII, C++ move semantics |
| `lanes` | one lane per thread, task or process; steps as blocks on the lanes; shared variables as rows; happens-before arrows | `lanes`, `interleaving` (one linear order of the steps, chosen by the reader) | `spawn(t)`, `join(t)`, `run(t, line)`, `read(t, var)`, `write(t, var, v)`, `lock(t, m)`, `unlock`, `wait`, `notify`, `send(ch, v)`, `recv`, `race(a, b)`, `deadlock(cycle)` | threads, data races, locks, deadlock, channels, actors, memory models |
| `eventloop` | the call stack, the task queue, the microtask queue, the timers and the I/O sources of an event-driven runtime | `js` (the HTML event loop model), `node` (the libuv phases), `asyncio`, `generic` | `enqueue(q, task)`, `dequeue(q)`, `tick`, `await(promise)`, `resolve(promise, v)`, `phase(name)` | async/await, promises, callbacks, `setTimeout` ordering, Python asyncio |
| `layout` | a type's bytes: fields as ranges at offsets, padding as hatched ranges, alignment marks, the total size | `row` (bytes in a line), `grid` (by word) | `field(name, offset, size)`, `pad(offset, size)`, `align(n)` | struct padding, alignment, `repr(C)`, unions, endianness, enum tags |
| `dag` (git) | commits as nodes, parents as edges, branches and tags as labels, `HEAD`, the index and the working tree as side panels | `lanes` (one column per branch), `graph` | `commit(id, parents, msg)`, `ref(name, id)`, `head(ref)`, `stage(path)`, `checkout`, `merge`, `rebase(from, onto)`, `conflict(path)` | Git branching, merge, rebase, reset, reflog, cherry-pick |
| `ast` | a parse tree or an abstract syntax tree, linked to the source text ranges | `tidy`, `indented` | `node(id, kind, range)`, `edge(parent, child, field)`, `focus(range)` | parsing, precedence, macros, the language-implementation area |
| `derivation` | an inference-rule tree: judgments with the rule name on each bar | `rules` | `judge(id, ctx, term, type)`, `rule(name, premises)`, `unify(a, b, subst)`, `fail(reason)` | type checking, type inference, the formal sections (11) |
| `automaton` | an NFA or a DFA with states, transitions and the input with a cursor; for backtracking engines, the stack of choice points | `graph`, `backtrack` | `state(id, accept)`, `trans(a, b, label)`, `active(ids)`, `consume(i)`, `choice(push or pop)`, `match(range)` | regex matching, backtracking, catastrophic backtracking and ReDoS, lexers |
| `plan` | a query plan tree with the operator, the estimated and the actual rows and costs | `tree`, `flame` | `op(id, kind, est, actual)`, `edge(parent, child)` | SQL query plans, indexes, joins |
| `sequence` | a sequence diagram: participants as columns, messages as arrows in time order | `uml` | `participant(id, name)`, `msg(from, to, label, kind)`, `activate`, `deactivate` | HTTP request and response, RPC, method calls between objects, TLS handshake as programmers see it |
| `classes` | a class diagram: types with members, and the relations (inherits, implements, composes, uses) | `uml` | none (static) | patterns (06.3), object-oriented features |
| `railroad` | a syntax diagram for a grammar production | `railroad` | none (static) | feature synopses (06.3), grammars (11.5), regex syntax |
| `bits` (extended) | a machine word with the sign, exponent and mantissa fields for IEEE 754, two's complement for integers, UTF-8 bytes of a code point | `int`, `float`, `utf8` | `set(i)`, `field(name, range)`, `overflow` | integer overflow, floating point, text encoding |

Rules for primitives are as in algorithm: a scene uses only the primitives of `@prog/viz`, and a new primitive is a change to the package with its own review (V13).

## 12.4 Recorded traces

A tracer runs an example inside its image and writes a trace file. The tracer runs in the same CI job as the example's run (07.4), so the trace, the output and the run status come from one execution.

**Tracers by language** (Default; the list grows with tier 2):

| Language | Tracer | How | From |
|---|---|---|---|
| Python | `prog-trace-py` | `sys.monitoring` (PEP 669) events for line, call, return, raise and yield; `sys.settrace` for versions before 3.12 | M2 |
| JavaScript, TypeScript | `prog-trace-js` | source instrumentation: the compiled JS is rewritten to call a recorder at each statement and call (Default); the V8 inspector protocol with stepping as the fallback for code that the rewrite cannot handle | M2 |
| Java, Kotlin | `prog-trace-jvm` | the Java Debug Interface (JDI): step events, frames and object fields through a debugger that attaches to the run | M3 |
| C, C++ | `prog-trace-native` | gdb with its Python API: step by line, read frames, locals and the memory of pointed-to objects; `rr` to record once and replay, so stepping does not change timing | M3 |
| Rust, Go, others | the native tracer with gdb or lldb (Rust), Delve (Go) | as for C | M4, Default |

**The trace file** (`build/traces/<node>/<example-id>.json`, gzip) holds the algorithm trace format (2150/10.2.3) with these additions:

| Field | Meaning |
|---|---|
| `run` | the run key of 07.4 (image digest, wrapper, flags, block hash) and the run status |
| `source` | the file and the line of each step, in the code that the reader sees (the wrapper's lines are removed and the line numbers mapped, 07.3) |
| `frames` | for each step: the frames with the function, the line and the local bindings; a binding is a primitive value or a reference to a heap id |
| `heap` | for each step that changes it: the objects with id, type and fields, as a diff from the step before |
| `threads` | the thread or task of each step, for `lanes` |
| `out` | the stdout and stderr chunks written at each step |
| `truncated` | `true` when the run passed the step budget (12.10); the trace keeps the first steps and the last steps |

**Filters.** A tracer records only the example's own code. Steps inside the standard library are folded into one step ("calls `sorted`, returns ..."). Objects of the runtime that the reader did not make (modules, the builtins) are hidden. The per-language filters are in `languages/<lang>/trace.yaml`.

**Values.** Values are shown as the language prints them (`repr` in Python, the debugger's format in C), with the per-language normalizers of 07.6 for addresses and identity hashes. In the `bytes` view, real addresses are kept, because the address is the point; the trace stores them, and the build marks them as "addresses from one run".

A trace is published only for a block at X2 or better. When the block's output changes (a new image), the trace is recorded again in the same run.

## 12.5 Generated diagrams

| Diagram | Source | Generator | Where it appears |
|---|---|---|---|
| class diagram | the node's example code | symbol extraction with tree-sitter (the `treesitter` field of 08.4): types, members, inheritance and implementation; composition from field types | pattern pages (06.3, Solution), object-oriented features |
| sequence diagram | a recorded trace (12.4): calls between objects of the example's types | the trace filtered to calls across types | pattern pages, RPC and HTTP pages |
| AST and parse tree | the example's code | the language's own parser where it has one in the image (Python `ast`, the TypeScript compiler API, `go/ast`, `syn` for Rust), tree-sitter otherwise | parsing, precedence, metaprogramming, language implementation |
| railroad diagram | the grammar productions of the feature (11.5) | our renderer, from EBNF or PEG | feature synopses (06.3), grammar pages, regex syntax |
| type derivation | the formal development (11) or the inference engine of `Prog.Inference` | the derivation printed by the Lean code | type checking and inference pages |
| query plan | the example's SQL run in the image | `EXPLAIN (ANALYZE, FORMAT JSON)` in PostgreSQL, `EXPLAIN QUERY PLAN` in SQLite, and the plan output of the other engines that the page shows | SQL pages, indexes, joins |
| git DAG | a shell example of Git commands run in the image | the build reads the repository after each command (`git log --all` with a fixed format, `git status`) and emits `dag` events, one step per command | Git pages (06.3, tool) |
| memory layout | a probe program in the image that prints the size, the alignment and the offset of each field | per target: `x86_64` and `aarch64` images | struct layout, padding, `repr(C)`, unions |
| automaton | the regex of the example | our Thompson construction and subset construction for the NFA and DFA; a backtracking simulator for the `backtrack` view | regex pages, ReDoS |
| request and response sequence | the log of the local fixture server (07.5, `fixture=http`) | the build reads the requests and responses that the fixture saw | HTTP, REST, client library pages |

Rules:

- A generated diagram is rebuilt when its source changes. It is never edited by hand.
- A diagram from our own simulator (the automaton, the backtracking view) must agree with the real engine: the match result of the simulator equals the match result of the engine in the image, on the example and on the cases of the node (V08).
- A diagram from tool output records the tool's version and shows it ("plan from PostgreSQL 18").

## 12.6 The genealogy and the release timeline

The genealogy is the graph of `influenced_by` and `variant_of` edges between languages (04.6), drawn by year.

- **Data.** Each edge has a source (04.6). The year of a language is `appeared`. Releases come from `releases.yaml` (09.1) for languages of tiers 1 to 3 and from release cards for the others.
- **Views.** The family view (one family of 08.7, 10 to 200 languages) is a layered graph with time on the x axis. The language view shows one language with its ancestors and descendants to depth 2, and its release timeline as a row of marks with the support window (09.3) shaded. The full view (all languages with an edge, about 5,000, Estimate) is a canvas with zoom and a filter by family and by tier, loaded only on `/languages/genealogy/`.
- **Static first.** Each family view and each language view is an SVG made at build time, with the text alternative of 12.9. The interactive view adds filters, hover and links on top of the same layout.
- **Release timeline.** The language page (06.3, "Versions and support") shows the releases with their dates, end-of-life dates and the computed status (09.5). The release timeline of several languages on one axis is part of the comparison pages (08.8).
- Layout is deterministic: the same graph gives the same picture, so a diff of the SVG shows real changes only.

## 12.7 Screenshots

GUI examples (07.5, `norun reason=gui`) and HTML and CSS pages that show rendering get captured images instead of scenes.

- A scheduled job runs the example with a virtual display (Xvfb for Linux GUI toolkits) or in headless Chromium and Firefox (for HTML and CSS), at fixed sizes (a desktop and a phone width), and saves PNG and AVIF.
- The image records the toolkit or browser version, the OS image and the date, and shows "captured on <date>".
- A capture whose pixels change between two runs with the same inputs is flagged as unstable, and the page uses the older capture until a person looks at it.
- A screenshot is a V1 figure. It never counts as V2.
- Screenshots of other people's products are not used, except under nominative use for a tool's own interface, captured from our own run (19).

## 12.8 The player and code sync

The player is the algorithm player (2150/10.3) with two changes:

- **Code sync uses the real code.** The code panel shows the example as the reader copies it (07.3), and the highlighted line is the line of the trace's `source`. There is no pseudocode panel on programming pages. When a concept page shows the same scene in several languages (a closure in Python, JavaScript and Rust), each language has its own trace, and the language tabs switch the trace.
- **Side panels.** A scene may show up to three primitives at once, in a fixed layout: the main view (for example `heap`), the `frames` column on the left, and the output console at the bottom (from the trace's `out`). The console shows the output up to the current step, so the reader sees which line printed what.

Input editing works only for scenes with an input schema and a browser runtime (07.10): the reader changes the input, the browser runs the tracer of that runtime (Default: Pyodide for Python at M3; others later), and the player draws the new trace. For other languages the presets are the inputs.

## 12.9 Accessibility

As in algorithm (2150/10.8.1), and:

- **Text alternative for every scene and figure.** The text view (`t`) of each new primitive is defined: `frames` as a nested list of frames and bindings; `heap` as a table of objects with their fields, and a list of "name refers to object" lines; `ownership` as a table of variables with the line ranges of each state; `lanes` as a table of steps with the thread, the line and the shared reads and writes; `eventloop` as the contents of each queue; `layout` as a table of fields with offset, size and padding; `dag` as the output of `git log --graph --oneline`; `ast` and `derivation` as indented outlines; `automaton` as a transition table; `plan` as the indented plan text; `sequence` as a numbered list of messages; `classes` as a list of types with members and relations; `railroad` as the EBNF text. Each generated SVG has the same text in a `<desc>` and in a fold under the figure.
- **Keyboard step control.** `→` and `←` step, `Home` and `End` go to the start and the end, `Space` plays and pauses, `[` and `]` jump to the previous and next phase, `t` toggles the text view, `1`, `2`, `3` move focus between the panels. Every control is also a button with a label.
- **Narration.** Each step of a recorded trace gets a generated note ("`counter` is bound to a new int object 1"), from templates per event and language, so that a screen reader can follow a recorded scene that no author wrote notes for. Authored scenes have notes written by the author, as in algorithm.
- **Colour.** The state marks of 2150/10.2.5 apply. Ownership states have patterns (solid for owned, dotted for shared borrow, double for mutable borrow, crossed for moved out).

## 12.10 Size budget

| Budget | Value |
|---|---|
| `@prog/viz` runtime, gzip | at most 90 KB (the 60 KB of `@prog/viz` plus at most 30 KB for the new primitives), loaded once, only when a player scrolls into view |
| one authored scene bundle, gzip | at most 15 KB |
| one trace file, gzip | at most 200 KB; the default preset's trace at most 50 KB |
| steps in a recorded trace | at most 2,000; the tracer stops at the limit and sets `truncated` |
| heap objects drawn at once | at most 150; more objects are grouped ("+ 40 more `int` objects") |
| one generated SVG | at most 150 KB before gzip; the genealogy full view is a canvas with tiled data, at most 2 MB gzip, loaded only on its own page |
| screenshots | AVIF at most 120 KB each, PNG fallback at most 400 KB |
| playback and first draw | as in algorithm: 60 frames per second, first draw within 100 ms of load, no layout shift |

All numbers are Default and are checked by V10.

## 12.11 Checks

`prog viz check <node>` runs these checks in CI for each changed scene, trace or diagram, and the nightly build runs them all. They are rules of the audit (15).

| Rule | Check |
|---|---|
| V01 | an authored scene builds, imports only `@prog/viz`, and uses no DOM, network, timer or unseeded random |
| V02 | determinism: the same input gives the same trace hash twice; a recorded trace gives the same trace hash on a re-run in the same image, after the address and identity normalizers |
| V03 | a recorded trace comes from the same run as the example's output block: the `run` key matches the run key of the block in `build/runs/`, and the trace's `out` equals the recorded output after normalization (07.6) |
| V04 | every `source` line of a trace is a line of the example as the reader sees it (no wrapper lines) |
| V05 | the presets of an authored scene reach every line of the example that the scene claims to show |
| V06 | **line sync for authored scenes**: the sequence of lines of the authored scene equals the sequence of lines of a recorded or instrumented run of the same example, with repeats merged; for languages with no tracer, the build instruments the example with print markers at each line, runs it, and compares |
| V07 | **declared failures agree**: an `ownership` `conflict(line, code)` event, a `lanes` `race` or `deadlock` event, or an `unwind` event matches a `fails` block of the node with that signature (07.2): the same diagnostic code at the same line, the same sanitizer kind (ThreadSanitizer for a race), the same exception type |
| V08 | a diagram from our own simulator agrees with the real engine or tool on the example and the node's cases (12.5) |
| V09 | every scene, trace and figure has a text view or a text alternative; every step has a note (written or generated); every state used is in the vocabulary |
| V10 | the budgets of 12.10 hold |
| V11 | the filmstrip and the social image render, as in algorithm |
| V12 | each control of an explorable has a label and a range, and each value of the range draws without error |
| V13 | the scene uses only the primitives of `@prog/viz` |
| V14 | a generated diagram's source hash matches the current source (no stale diagram), and the tool version it records is the version in the image |
| V15 | a screenshot is younger than the newest image of its toolkit or browser by at most 90 days, or the page shows it as "captured with <old version>" |

V03 and V07 are the checks that matter most. They make a scene a picture of the real run, and they make an error marker in a scene the same error that the compiler or the sanitizer gave.

## 12.12 The V ladder and what needs a scene

The levels are the ones of 04.9:

| Level | It needs |
|---|---|
| V0 | no scene, no trace, no figure |
| V1 | static figures: generated diagrams, screenshots, or a filmstrip from a scene that fails a check needed for V2 |
| V2 | an interactive scene (recorded or authored) or an explorable that passes V01 to V14 as they apply, with step back, code sync on the real code, the text view and a note on each step |
| V3 | V2, and a person who is not the author reviewed it for clarity and accessibility on this exact hash (15.2) |

**By kind**, for a node of the core:

| Kind | Needs |
|---|---|
| `concept` with `viz: expected` | V2 (and V3 for "complete", 04.9) |
| `pattern` | V1: the generated class diagram and sequence diagram (06.3) |
| `feature` | optional; V2 for the features that `realizes` a `viz: expected` concept in a tier-1 language, through the concept's scene in that language |
| `task` | optional |
| `tool` | V2 for the core concepts of Git (the object model, branches, merge, rebase, reset); V1 for other core tools (a screenshot of first use when it has a UI) |
| `language` | V1: the genealogy view and the release timeline, generated |
| `weakness` | V1; V2 for the memory-safety weaknesses (a `heap` scene in the `bytes` view of the overflow or the use after free, from the `fails` example under ASan) |
| `format` | V1 for encodings (the `bits` view of UTF-8), otherwise none |
| others | none |

**Concepts with `viz: expected`** (the first list, about 60, Estimate; the curator sets the field and the editors approve it at M1):

| Area | Concepts |
|---|---|
| foundations | two's complement and integer overflow, floating point representation, text encoding and UTF-8, stack and heap |
| functions-and-abstraction | the call stack, recursion, tail calls, scope and binding, closures, generators and coroutines, lazy evaluation, continuations |
| data-and-values | value and reference semantics, aliasing, mutability, shallow and deep copy, equality and identity |
| objects-and-classes | dynamic dispatch and vtables, inheritance and method resolution order, prototypes |
| memory-and-resources | pointers, dangling pointers, manual allocation, reference counting and cycles, tracing GC (mark and sweep, copying, generational), ownership and moves, borrowing, RAII, memory layout and alignment |
| errors-and-control | exceptions and stack unwinding, `finally` and `defer` order |
| concurrency-and-parallelism | threads and interleaving, data races, locks and mutual exclusion, deadlock, channels, actors, the event loop, async/await, promises and microtasks, memory models and atomics |
| syntax-and-semantics, types | operator precedence and the AST, evaluation order, type checking, type inference |
| data-and-persistence | regex matching, backtracking and ReDoS, SQL joins, indexes and query plans, transaction isolation anomalies |
| tools-and-workflow | Git commits and branches, merge, rebase, the index |
| web-development | the HTTP request lifecycle, DOM event propagation, the CSS box model, flexbox, grid |

## 12.13 Authoring

1. `prog new viz <node> [--from-trace <example-id>]`. With `--from-trace`, the command records the trace of the example and writes `viz/scene.yaml`: the example id, the primitives and views, the presets, and the note templates to override. No code is needed for a recorded scene.
2. For an authored scene, the command writes `viz/scene.ts` with one `v.line` call per line of the example, as in algorithm.
3. `prog viz serve <node>` opens the player with reload on save and the trace inspector.
4. `prog viz check <node>` runs the checks of 12.11 and prints the V level before and after.
5. The `visualizer` agent (14) follows the same steps. It prefers a recorded scene and writes an authored scene only when the concept is in the list of models that no tracer can see (the event loop, ownership, memory models, GC phases, isolation anomalies).
