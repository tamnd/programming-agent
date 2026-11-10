// SPDX-License-Identifier: Apache-2.0

package lang

import "testing"

func TestWithTool(t *testing.T) {
	p := Programs["java"].WithTool("java", "/opt/jdk/bin/java")
	if p.Steps[1][0] != "/opt/jdk/bin/java" || p.Tools[1] != "/opt/jdk/bin/java" || p.Steps[0][0] != "javac" {
		t.Errorf("got %+v", p)
	}
	if Programs["java"].Steps[1][0] != "java" {
		t.Error("WithTool changed the original program")
	}
}

func TestClean(t *testing.T) {
	for _, tc := range []struct{ dir, in, want string }{
		{"/var/folders/x/T/prog-run-1", "Test run for /private/var/folders/x/T/prog-run-1/bin/a.dll\r\nat /var/folders/x/T/prog-run-1/a.cs\n", "Test run for $WORK/bin/a.dll\nat $WORK/a.cs\n"},
		{"/tmp/prog-run-2", "/private/tmp/prog-run-2/main.py and /tmp/prog-run-2/main.py", "$WORK/main.py and $WORK/main.py"},
		// Only the temporary directories of macOS have a real path under /private.
		{"/home/ci/prog-run-3", "/private/home/ci/prog-run-3/main.py", "/private$WORK/main.py"},
	} {
		if got := clean(tc.in, tc.dir); got != tc.want {
			t.Errorf("%s: %q, want %q", tc.dir, got, tc.want)
		}
	}
}
