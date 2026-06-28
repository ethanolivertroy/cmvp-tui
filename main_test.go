package main

import "testing"

func TestDefaultVersion(t *testing.T) {
	if version == "" {
		t.Fatal("version must not be empty")
	}
}
