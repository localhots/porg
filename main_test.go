package main

import (
	"testing"
)

func TestToTimeFormat(t *testing.T) {
	assertEqualStr(t, "dist/photos/06/01/02-150405", toTimeFormat("dist/photos/%Y/%m/%d-%H%M%S"))
	assertEqualStr(t, "dist/2006-01-02/150405", toTimeFormat("dist/%y-%m-%d/%H%M%S"))
}

func assertEqualStr(t *testing.T, exp, out string) {
	t.Helper()
	if exp != out {
		t.Errorf("Expected %q, got %q", exp, out)
	}
}