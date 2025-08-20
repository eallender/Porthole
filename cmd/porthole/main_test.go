package main

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("version constant should not be empty")
	}
	
	if !strings.Contains(version, "0.1.0") {
		t.Errorf("expected version to contain '0.1.0', got: %s", version)
	}
}