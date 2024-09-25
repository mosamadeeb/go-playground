package main

import (
	"slices"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("array with fixed elements", func(t *testing.T) {
		numbers := [5]int{1, 2, 3, 4, 5}

		got := Sum(numbers[:])
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("array with random elements", func(t *testing.T) {
		numbers := [...]int{23, -9, 0, 3}

		got := Sum(numbers[:])
		want := 17

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("slice of numbers", func(t *testing.T) {
		numbers := []int{3, 2, 1}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})
}

func TestSumAll(t *testing.T) {
	t.Run("two slices", func(t *testing.T) {
		got := SumAll([]int{1, 2}, []int{0, 9})
		want := []int{3, 9}

		if !slices.Equal(got, want) {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("one slice", func(t *testing.T) {
		got := SumAll([]int{1, 1, 1})
		want := []int{3}

		if !slices.Equal(got, want) {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func TestSumAllTails(t *testing.T) {
	t.Run("multiple slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}

		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("empty slices", func(t *testing.T) {
		got := SumAllTails([]int{})
		want := []int{0}

		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
