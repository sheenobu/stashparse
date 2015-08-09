package stashparse

import (
	"os"
	"testing"
)

// TestParseComplex parses the complex configuration
func TestParseComplex(t *testing.T) {
	var config Config
	r, _ := os.Open("complex.conf")
	decoder := NewDecoder(r)
	err := decoder.Decode(&config)

	if err != nil {
		t.Errorf("%s", err)
	}
}
