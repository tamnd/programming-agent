// SPDX-License-Identifier: Apache-2.0

package node

import "testing"

func TestIsBuildFile(t *testing.T) {
	for name, want := range map[string]bool{"go.mod": true, "x/Cargo.toml": true, "Impl.csproj": true, "merge_sort.go": false, "package.json": true, "README.md": false} {
		if got := IsBuildFile(name); got != want {
			t.Errorf("IsBuildFile(%q) = %v", name, got)
		}
	}
}

func TestIsTestFile(t *testing.T) {
	for name, want := range map[string]bool{
		"test_line_stats.py":      true,
		"line_stats_test.go":      true,
		"lineStats.test.js":       true,
		"LineStatsTest.java":      true,
		"LineStatsTests.cs":       true,
		"test_line_stats.bats":    true,
		"line_stats.bats":         true,
		"Get-LineStats.Tests.ps1": true,
		"line_stats.sh":           false,
		"Get-LineStats.ps1":       false,
		"line_stats.rs":           false,
	} {
		if got := IsTestFile(name); got != want {
			t.Errorf("IsTestFile(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestAPICollectionURLSegment(t *testing.T) {
	for _, c := range []struct{ module, want string }{
		{"java.util", "/java/api/java.util/"},
		{"reserved words", "/java/api/reserved-words/"},
		{"tsc --build", "/java/api/tsc--build/"},
		{"<stdio.h>", "/java/api/stdio.h/"},
		{"std::vec", "/java/api/std::vec/"},
		{"<>", "/java/api/builtins/"},
	} {
		k := &Card{Kind: "api", Owner: "java", Path: "x", Module: c.module}
		if got := k.Collection(); got != c.want {
			t.Errorf("module %q: got %q, want %q", c.module, got, c.want)
		}
	}
}
