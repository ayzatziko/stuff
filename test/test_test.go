package test

import "testing"

func TestAppendA(t *testing.T) {
	helpered(t, "A", "B")
	notHelpered(t, "A", "B")
}

func helpered(t *testing.T, want, got string) {
	t.Helper()

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func notHelpered(t *testing.T, want, got string) {
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
