# Runner images

`prog test` runs the builds, the native tests, the adapters and the custom checkers in these images (spec 08.7). Each command gets a new container with no network, a read-only root, a memory limit of 2 GiB, a CPU limit of 2 and a process limit of 512. The build tree is the only writable mount.

| Image | Languages | Toolchain |
|---|---|---|
| `prog-runner-python` | python | Python 3.14, pytest 9.2 |
| `prog-runner-go` | go | Go 1.27 |
| `prog-runner-rust` | rust | Rust 1.98 |
| `prog-runner-gcc` | c, cpp, bash | GCC 15, with ASan and UBSan; Bash 5.2 and bats-core |
| `prog-runner-node` | javascript, typescript | Node 24 |
| `prog-runner-jvm` | java, kotlin | JDK 25, Kotlin 2.4.21, JUnit console launcher 1.13.4 |
| `prog-runner-dotnet` | csharp, powershell | .NET SDK 10, xUnit 2.9.3, restored ahead of time; PowerShell 7 and Pester 5.7.1 |
| `prog-runner-lean` | lean | Lean 4.35.0 |

Each base image and each download is pinned by digest. The workflow `images.yml` builds the images, runs the sandbox tests on them, and pushes them to `ghcr.io/tamnd` with build provenance. The content repository pins each image by digest in `manifests/images.yaml`.

`prog run` builds more images on top of these for the blocks with `deps=` (spec 07.7). Each one is tagged `prog-deps-<lang>:<hash>` and has the label `prog.deps=1`. The Dockerfile text of each language is in `lang/deps.go`. The Rust image gets the crates of `Cargo.lock` with `cargo vendor` and compiles them once. The JVM image gets the jars of `pom.xml` with Maven, pinned by checksum, into `/opt/prog/lib`. The gcc image has no vcpkg, so a C++ block with `vcpkg.json` fails with a message that says so. To remove these images, run `docker image prune -a --filter label=prog.deps=1`.

The images are for `linux/amd64`. On another machine, give `prog test -unsafe-local` to run the tests with the tools of the machine.

The images expect these paths:

- `/opt/prog/junit.jar` is the JUnit console launcher.
- `/opt/kotlinc` is the Kotlin compiler.
- `/opt/nuget` holds the NuGet packages of the test project. It is read-only, so a restore cannot fetch a new package.
