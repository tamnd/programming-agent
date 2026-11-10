# programming-agent

[![CI](https://github.com/tamnd/programming-agent/actions/workflows/ci.yml/badge.svg)](https://github.com/tamnd/programming-agent/actions/workflows/ci.yml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/tamnd/programming-agent/badge)](https://scorecard.dev/viewer/?uri=github.com/tamnd/programming-agent)
[![Go Reference](https://pkg.go.dev/badge/github.com/tamnd/programming-agent.svg)](https://pkg.go.dev/github.com/tamnd/programming-agent)

This repository holds `prog`, the toolchain of [programming.tamnd.com](https://programming.tamnd.com). The content is in [tamnd/programming](https://github.com/tamnd/programming).

`prog` reads the content tree and checks it against the specification. It runs each code example and compares the output with the recorded output. It renders the site and deploys it.

The toolchain started as a copy of [algorithm-agent](https://github.com/tamnd/algorithm-agent), the toolchain of algorithm.tamnd.com. The packages are copies. They do not import algorithm-agent, so the two projects can change at their own speed.

## Install

You need Go 1.27 or later.

```sh
go install github.com/tamnd/programming-agent/cmd/prog@latest
```

You can also download a binary from the [releases](https://github.com/tamnd/programming-agent/releases). Each binary has a build provenance attestation. To verify a binary, run `gh attestation verify prog_linux_amd64.tar.gz -R tamnd/programming-agent`.

## Use

Run these commands in a clone of tamnd/programming.

| Command | What it does |
| --- | --- |
| `prog new concept closures -area functions-and-abstraction/closures-and-scope` | Makes a node directory with the sections of its kind. |
| `prog promote rust/errors/E0382` | Makes a node from a card. The node keeps the tag and the facts, and the row leaves its register. |
| `prog fmt` | Puts node.md, the registers and the topics in canonical form. Add `-check` to only report. |
| `prog mint` | Gives each provisional id `NEW-<slug>` a permanent tag. Add `-dry-run` to see the tags first. |
| `prog mv <slug> <new-slug>` | Renames or moves a node. The old URL redirects. The tag does not change. |
| `prog check` | Runs the audit rules of spec 15. Rule X08 is a level finding when more than 10% of the code blocks of a node are `norun`. The limit is 40% in embedded-and-hardware, graphics-and-games-programming, mobile-and-desktop and distributed-and-cloud. A `norun` fragment of a group that runs does not count. |
| `prog run closures` | Runs the code blocks of a page and compares the output with the output blocks. Add `-update` to write the output blocks from the runs. The blocks run in the sandbox, in CI, or on your computer with `-unsafe-local`. |
| `prog compat probe match-statement` | Runs the probes of the compat entries of a node in the pinned image of each release, and writes the results to `manifests/probes.jsonl`. A probe must pass on the release that the entry gives and fail on the release before. It needs Docker or Podman. Add `-all` to run the probes that are already verified again. |
| `prog run json -sandbox docker` | A block with `deps=deps/rust/Cargo.lock` uses the packages of a lock file of the node (spec 07.7). Rust reads `Cargo.lock` and `Cargo.toml`, Java reads `pom.xml`, and C++ reads `vcpkg.json`. C# reads a `.csproj`, for a test block only. In the sandbox, prog builds an image with the packages on top of the image of the language. The tag is a hash of the language image and the lock file, so the next run uses the same image. The image build is the only step with network. The block runs in the new image with no network. With `-unsafe-local`, Rust uses cargo with `--offline` and Java uses Maven with `-o`, so the packages must be in the local cache. C++ has no local support. A block that cannot get its packages, or a block of another language with `deps=`, is skipped with the reason `deps`. It does not fail. |
| `prog run unit-testing` | A `test` block is a test file (spec 07.11). prog runs it with the test runner of the language and compares the report with the next `output` block. Python uses `python3 -m unittest -v`, and its report goes to stderr, so `-update` writes `output stream=stderr`. JavaScript and TypeScript use `node --test`. Go uses `go test -v` in a module with the name `example`. Rust uses `rustc --test`, or `cargo test` with `deps=`. Java uses the JUnit console launcher, from `/opt/prog/junit.jar` in the image or from the cache on your computer. With a `pom.xml`, prog leaves out the JUnit jars of the pom, because the launcher has JUnit already. C# needs `deps=` with a `.csproj` that names the test framework, such as xunit, and runs `dotnet test`. The packages come from `/opt/nuget` in the image or from the NuGet folder on your computer. C++ needs `deps=` with a `vcpkg.json` that names `gtest` or `catch2`, and runs only in the sandbox. The normalizers replace the timings with `$TIME` and the paths of the run with `$WORK`. A block passes when each test passes and the report matches. A `test fails expect=test-failure` block passes when at least one test fails. A language with no test runner, such as C, is skipped with the reason `runner`. A test framework that prog cannot get is skipped with the reason `deps`. |
| `prog status closures` | Shows the ladders of a node and what the next level needs. |
| `prog site serve` | Renders the site and serves it on your computer. It renders again when a file changes. |
| `prog site check` | Renders the site and checks each internal link. |
| `prog site deploy` | Renders the site, checks the limits of the host and uploads the tree as the static assets of the Cloudflare Worker `programming` with a pinned wrangler. With `-preview`, it uploads a version with a preview alias and does not change production. It reads `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` from the environment and does not print them. |

`prog help` shows all commands. A command that is planned but not built yet exits with code 3 and names its milestone.

The exit codes are: 0 for success, 1 when a check finds a problem, 2 for a usage error and 3 for a planned command or a missing tool.

## Layout

| Directory | Contents |
| --- | --- |
| `cmd/prog` | The command. |
| `node` | Reads the content tree and the card registers. |
| `check` | The audit rules of spec 15. |
| `status` | The ladders of a node (spec 04.9). |
| `lang` | Build trees and native tests for each language. |
| `contract` | The contracts and the shared cases of tasks (spec 10). |
| `formal` | Builds the Lean package and checks the axiom policy. |
| `viz` | Bundles the scenes with esbuild and checks them in a JavaScript runtime. |
| `compat` | Reads the compat entries, the implementation register and the probe results (spec 09), and runs the probes. |
| `examples` | Reads the code blocks of the pages, runs them and checks the output (spec 07). |
| `markup` | Markdown, code blocks and math. |
| `katex` | Renders TeX to HTML on the server. |
| `site` | The templates, assets and pages of the site. |
| `spec` | The specification. The site publishes it at [/spec/](https://programming.tamnd.com/spec/). |
| `testdata/corpus` | A small content tree for the tests. |

## Develop

```sh
make test     # go vet and go test
make lint     # golangci-lint
make site     # render testdata/corpus into build/site
```

Read [CONTRIBUTING.md](CONTRIBUTING.md) before you open a pull request. To report a security problem, read [SECURITY.md](SECURITY.md).

To release, push a tag `vX.Y.Z` and run `scripts/release.sh vX.Y.Z`. A patch release comes every few merged pull requests. A minor release comes at the end of each milestone.

## Roadmap

The work is in six milestones, M0 to M5, from November 2026 to December 2028. [Spec 21](https://programming.tamnd.com/spec/21-milestones/) gives the goal and the exit criteria of each one. The [milestones](https://github.com/tamnd/programming-agent/milestones) list the open issues.

## Licence

The code is under the [Apache License 2.0](LICENSE). The scene package `viz/packages/core` is under the MIT licence, so scenes in the content repository can use it. The parts from other projects are listed in [THIRD_PARTY.md](THIRD_PARTY.md).
