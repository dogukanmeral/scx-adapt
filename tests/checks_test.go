package tests

import (
	"internal/checks"
	"slices"
	"testing"
)

type duplicateTest[T comparable] = struct {
	sli  []T
	b    bool
	dups []T
}

var intDuplicateTests = []duplicateTest[int]{
	{[]int{2, 3, 4, 4}, true, []int{4}},
	{[]int{}, false, []int{}},
	{[]int{1}, false, []int{}},
	{[]int{13, 13, 24, 24}, true, []int{13, 24}},
}

var stringDuplicateTests = []duplicateTest[string]{
	{[]string{"abc", "bcd", "abc"}, true, []string{"abc"}},
	{[]string{}, false, []string{}},
	{[]string{"abc"}, false, []string{}},
	{[]string{"abc", "bcd", "abc", "bcd"}, true, []string{"abc", "bcd"}},
}

func TestContainsDuplicate(t *testing.T) {
	for _, te := range intDuplicateTests {
		is, du := checks.ContainsDuplicate(te.sli)

		if is != te.b || !slices.Equal(du, te.dups) {
			t.Errorf("Slice: %v ||| Contains duplicate: %v (Expected: %v) ||| Duplicates: %v (Expected: %v)", te.sli, is, te.b, du, te.dups)
		}
	}

	for _, te := range stringDuplicateTests {
		is, du := checks.ContainsDuplicate(te.sli)

		if is != te.b || !slices.Equal(du, te.dups) {
			t.Errorf("Slice: %v ||| Contains duplicate: %v (Expected: %v) ||| Duplicates: %v (Expected: %v)", te.sli, is, te.b, du, te.dups)
		}
	}
}
