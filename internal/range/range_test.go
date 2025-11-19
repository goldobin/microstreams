package _range

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Validation(t *testing.T) {
	tests := []struct {
		r       Range
		wantErr bool
	}{
		{r: Range{0, 1}},
		{r: Range{1, 2}},
		{r: Range{-1, 0}},
		{r: Range{-2, -1}},
		{r: Range{-100, -2}},
		{r: Range{2, 200}},
		{r: Range{0, 0}, wantErr: true},
		{r: Range{1, 1}, wantErr: true},
		{r: Range{-1, -1}, wantErr: true},
		{r: Range{-1, -2}, wantErr: true},
		{r: Range{2, 1}, wantErr: true},
	}

	for _, tt := range tests {
		err := tt.r.Validate()

		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestRange_Split(t *testing.T) {
	tests := []struct {
		r    Range
		n    int
		want []Range
	}{
		{
			r: Range{0, 1},
			n: 0,
		},
		{
			r:    Range{0, 1},
			n:    1,
			want: []Range{{0, 1}},
		},
		{
			r:    Range{0, 1},
			n:    2,
			want: []Range{{0, 1}},
		},
		{
			r:    Range{0, 2},
			n:    5,
			want: []Range{{0, 1}, {1, 2}},
		},
		{
			r:    Range{0, 2},
			n:    1,
			want: []Range{{0, 2}},
		},
		{
			r: Range{0, 2},
			n: 2,
			want: []Range{
				{0, 1}, {1, 2}},
		},
		{
			r: Range{0, 3},
			n: 2,
			want: []Range{
				{0, 2}, {2, 3}},
		},
		{
			r: Range{0, 4},
			n: 2,
			want: []Range{
				{0, 2}, {2, 4}},
		},
		{
			r: Range{0, 5},
			n: 2,
			want: []Range{
				{0, 3}, {3, 5}},
		},
		{
			r: Range{0, 5},
			n: 3,
			want: []Range{
				{0, 2}, {2, 4}, {4, 5}},
		},
	}

	for i, tt := range tests {
		result := tt.r.Split(tt.n)
		if len(tt.want) == 0 {
			if len(result) != 0 {
				t.Errorf("case %d: wantErr zero slice, got length %v", i, result)
			}
			continue
		}

		if !reflect.DeepEqual(result, tt.want) {
			t.Errorf("case %d: wantErr %v, got %v", i, tt.want, result)
		}
	}
}

func TestRange_PickRandom(t *testing.T) {
	var (
		n = 100
		r = Range{0, 5}
	)

	for i := 0; i < n; i++ {
		v := r.PickRandom()
		assert.GreaterOrEqual(t, v, r.Left)
		assert.Less(t, v, r.Right)
	}
}
