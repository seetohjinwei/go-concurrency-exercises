package main

import "testing"

func TestSame(t *testing.T) {
	t.Run("same tree", func(t *testing.T) {
		got := Same(NewTree(1), NewTree(1))
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("diff tree", func(t *testing.T) {
		got := Same(NewTree(1), NewTree(2))
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
