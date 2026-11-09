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
