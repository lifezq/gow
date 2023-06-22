package gow

import "testing"

func TestNew(t *testing.T) {

	if s := New(); s == nil || s.config == nil {
		t.Fatalf("New error %s\n", "...")
	}
}
