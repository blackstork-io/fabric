package ast2md

import (
	"fmt"
	"testing"
)

func TestIntLen(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"positive number", 12345},
		{"negative number", -123},
		{"zero", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intLen(tt.n); got != len(fmt.Sprintf("%d", tt.n)) {
				t.Errorf("intLen(%d) = %v, want %v", tt.n, got, len(fmt.Sprintf("%d", tt.n)))
			}
		})
	}
}
