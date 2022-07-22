package rng

import (
	"reflect"
	"testing"
)

func TestRange_Validation(t *testing.T) {
	cases := []struct {
		r        Range
		expected func(error) bool
	}{
		{
			r:        Range{0, 1},
			expected: noError,
		},
		{
			r:        Range{1, 2},
			expected: noError,
		},
		{
			r:        Range{-1, 0},
			expected: noError,
		},
		{
			r:        Range{-2, -1},
			expected: noError,
		},
		{
			r:        Range{-100, -2},
			expected: noError,
		},
		{
			r:        Range{2, 200},
			expected: noError,
		},
		{
			r:        Range{0, 0},
			expected: validationError,
		},
		{
			r:        Range{1, 1},
			expected: validationError,
		},
		{
			r:        Range{-1, -1},
			expected: validationError,
		},
		{
			r:        Range{-1, -2},
			expected: validationError,
		},
		{
			r:        Range{2, 1},
			expected: validationError,
		},
	}

	for i, c := range cases {
		err := c.r.Validate()
		errExpected := c.expected(err)
		if !errExpected {
			t.Errorf("case %d: error expected %v, got %v", i, errExpected, err)
		}
	}

}

func TestRange_Split(t *testing.T) {
	cases := []struct {
		r        Range
		n        int
		expected []Range
	}{
		{
			r:        Range{0, 1},
			n:        0,
			expected: []Range{},
		},
		{
			r:        Range{0, 1},
			n:        1,
			expected: []Range{{0, 1}},
		},
		{
			r: Range{0, 1},
			n: 2,
			expected: []Range{
				{0, 1},
			},
		},
		{
			r: Range{0, 2},
			n: 5,
			expected: []Range{
				{0, 1},
				{1, 2},
			},
		},
		{
			r:        Range{0, 2},
			n:        1,
			expected: []Range{{0, 2}},
		},
		{
			r: Range{0, 2},
			n: 2,
			expected: []Range{
				{0, 1},
				{1, 2},
			},
		},
		{
			r: Range{0, 3},
			n: 2,
			expected: []Range{
				{0, 2},
				{2, 3},
			},
		},
		{
			r: Range{0, 4},
			n: 2,
			expected: []Range{
				{0, 2},
				{2, 4},
			},
		},
		{
			r: Range{0, 5},
			n: 2,
			expected: []Range{
				{0, 3},
				{3, 5},
			},
		},
		{
			r: Range{0, 5},
			n: 3,
			expected: []Range{
				{0, 2},
				{2, 4},
				{4, 5},
			},
		},
	}

	for i, c := range cases {
		result := c.r.Split(c.n)
		if len(c.expected) == 0 {
			if len(result) != 0 {
				t.Errorf("case %d: expected zero slice, got length %v", i, result)
			}
			continue
		}

		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("case %d: expected %v, got %v", i, c.expected, result)
		}
	}
}

func TestRange_PickRandom(t *testing.T) {
	n := 100
	r := Range{0, 5}
	for i := 0; i < n; i++ {
		v := r.PickRandom()
		if v < r.Left || v >= r.Right {
			t.Errorf("expected value between 0 and 100, got %v", v)
		}
	}
}

func noError(err error) bool {
	return err == nil
}

func validationError(err error) bool {
	return err != nil
}
