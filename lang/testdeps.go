// SPDX-License-Identifier: Apache-2.0

package lang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// The test runners of the test blocks with deps= (spec 07.11). They use the packages that WithDeps gets for the block.

// cargoTestBuild builds the test binary of the block and its crates. The crates are compiled already, so only the block compiles.
var cargoTestBuild = []string{"cargo", "test", "--offline", "--locked", "--no-run", "--quiet"}

// cargoTestRun runs the tests one at a time, so that the order of the report is the order of the tests.
var cargoTestRun = []string{"cargo", "test", "--offline", "--locked", "--", "--test-threads=1"}

// rustTestDeps makes the test block the src/main.rs of a Cargo project with the Cargo.toml and the Cargo.lock of the node, and runs cargo test. A test build needs no main function.
func rustTestDeps(p Program, files map[string][]byte, dir string) (Program, error) {
	if _, err := cargoBin(files["Cargo.toml"]); err != nil {
		return p, err
	}
	p = cloneSteps(p)
	p.Main = "src/main.rs"
	p.Tools = []string{"cargo"}
	build := cargoTestBuild
	if dir == "" {
		build = []string{"sh", "-c", "cp -pR " + imageDepsDir + "/. . && exec " + strings.Join(cargoTestBuild, " ")}
	} else {
		p.Seed = dir
	}
	p.Steps = [][]string{build, slices.Clone(cargoTestRun)}
	p.Setup = 0
	return p, nil
}

// cppTestDeps compiles the test file with the test framework of the vcpkg.json and runs the test binary. The framework library has the main function.
func cppTestDeps(p Program, files map[string][]byte, _ string) (Program, error) {
	libs, run, err := cppTestFramework(files["vcpkg.json"])
	if err != nil {
		return p, err
	}
	p = cloneSteps(p)
	build := []string{"c++", "-std=c++23", "-Wall", "-I" + imageCppInclude, "-L" + imageCppLib, "-o", "main_test", ArgSources}
	p.Steps = [][]string{append(build, libs...), run}
	return p, nil
}

// cppTestFramework finds GoogleTest or Catch2 in the dependencies of a vcpkg.json. It returns the libraries to link and the command that runs the tests.
func cppTestFramework(manifest []byte) (libs, run []string, err error) {
	var m struct {
		Dependencies []json.RawMessage `json:"dependencies"`
	}
	if err := json.Unmarshal(manifest, &m); err != nil {
		return nil, nil, fmt.Errorf("deps: vcpkg.json: %w", err)
	}
	for _, raw := range m.Dependencies {
		var name string
		if json.Unmarshal(raw, &name) != nil {
			var dep struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(raw, &dep); err != nil {
				return nil, nil, fmt.Errorf("deps: vcpkg.json: %w", err)
			}
			name = dep.Name
		}
		switch name {
		case "gtest":
			return []string{"-lgtest_main", "-lgtest", "-pthread"}, []string{"./main_test"}, nil
		case "catch2":
			// Catch2 3 can run the tests in a random order. The order of the file gives a stable report.
			return []string{"-lCatch2Main", "-lCatch2"}, []string{"./main_test", "--order", "decl"}, nil
		}
	}
	return nil, nil, errors.New("deps: the vcpkg.json of a C++ test block names no test framework; add gtest or catch2 to its dependencies")
}

// imageNuGet is the folder of the NuGet packages in the dotnet image (images/dotnet).
const imageNuGet = "/opt/nuget"

// nugetMissing matches the errors of a NuGet restore that cannot find a package, or cannot read a source.
var nugetMissing = regexp.MustCompile(`\berror NU1(101|102|103|301)\b`)

// nugetFolder is the folder of the NuGet packages on this machine: NUGET_PACKAGES, or .nuget/packages in the home directory.
func nugetFolder() (string, error) {
	if d := os.Getenv("NUGET_PACKAGES"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".nuget", "packages"), nil
}

// csharpTestDeps puts the .csproj of the node into the work directory and runs dotnet test. The restore reads the packages from the folder of the packages only, so it needs no network. The folder is /opt/nuget in the image, and the NuGet folder of this machine in local mode.
func csharpTestDeps(p Program, files map[string][]byte, dir string) (Program, error) {
	source := imageNuGet
	if dir != "" {
		var err error
		if source, err = nugetFolder(); err != nil {
			return p, err
		}
	}
	p = cloneSteps(p)
	p.Files = map[string][]byte{}
	for name, b := range files {
		if strings.HasSuffix(name, ".csproj") {
			p.Files[name] = b
		}
	}
	p.Steps = [][]string{
		{"dotnet", "build", "-nologo", "-v", "q", "-clp:NoSummary;ErrorsOnly", "-p:NuGetAudit=false", "-p:RestoreSources=" + source},
		{"dotnet", "test", "--no-build", "--nologo", "-v", "q"},
	}
	return p, nil
}

// csharpPrepare checks that the packages of the .csproj are in the NuGet folder of this machine. The restore reads only that folder, so it has no network.
func csharpPrepare(ctx context.Context, dir string) error {
	if _, err := exec.LookPath("dotnet"); err != nil {
		return DepsSkip{"csharp", "dotnet is not installed"}
	}
	source, err := nugetFolder()
	if err != nil {
		return err
	}
	out, err := runIn(ctx, dir, "dotnet", "restore", "-nologo", "-v", "q", "-p:NuGetAudit=false", "-p:RestoreSources="+source)
	if err != nil {
		if nugetMissing.Match(out) {
			return DepsSkip{"csharp", "the packages of the .csproj are not all in the NuGet folder of this machine, and the run has no network; " + depsNoLocalSupportHint}
		}
		return fmt.Errorf("deps: dotnet restore of the .csproj failed:\n%s", tail(string(out), 2000))
	}
	for _, d := range []string{"obj", "bin"} {
		if err := os.RemoveAll(filepath.Join(dir, d)); err != nil {
			return err
		}
	}
	return nil
}
