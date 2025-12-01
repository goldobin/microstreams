package ranges

import (
	"math/rand/v2"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Validation(t *testing.T) {
	tests := []struct {
		r       IntRange
		wantErr bool
	}{
		{r: IntRange{0, 1}},
		{r: IntRange{1, 2}},
		{r: IntRange{-1, 0}},
		{r: IntRange{-2, -1}},
		{r: IntRange{-100, -2}},
		{r: IntRange{2, 200}},
		{r: IntRange{0, 0}, wantErr: true},
		{r: IntRange{1, 1}, wantErr: true},
		{r: IntRange{-1, -1}, wantErr: true},
		{r: IntRange{-1, -2}, wantErr: true},
		{r: IntRange{2, 1}, wantErr: true},
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
		r    IntRange
		n    int
		want []IntRange
	}{
		{
			r: IntRange{0, 1},
			n: 0,
		},
		{
			r:    IntRange{0, 1},
			n:    1,
			want: []IntRange{{0, 1}},
		},
		{
			r:    IntRange{0, 1},
			n:    2,
			want: []IntRange{{0, 1}},
		},
		{
			r:    IntRange{0, 2},
			n:    5,
			want: []IntRange{{0, 1}, {1, 2}},
		},
		{
			r:    IntRange{0, 2},
			n:    1,
			want: []IntRange{{0, 2}},
		},
		{
			r: IntRange{0, 2},
			n: 2,
			want: []IntRange{
				{0, 1}, {1, 2},
			},
		},
		{
			r: IntRange{0, 3},
			n: 2,
			want: []IntRange{
				{0, 2}, {2, 3},
			},
		},
		{
			r: IntRange{0, 4},
			n: 2,
			want: []IntRange{
				{0, 2}, {2, 4},
			},
		},
		{
			r: IntRange{0, 5},
			n: 2,
			want: []IntRange{
				{0, 3}, {3, 5},
			},
		},
		{
			r: IntRange{0, 5},
			n: 3,
			want: []IntRange{
				{0, 2}, {2, 4}, {4, 5},
			},
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
		r = IntRange{0, 5}
	)

	randSrc := rand.NewPCG(1, 2)
	for i := 0; i < n; i++ {
		v := r.Random(randSrc)
		assert.GreaterOrEqual(t, v, r.Left)
		assert.Less(t, v, r.Right)
	}
}
