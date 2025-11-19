package _range

import (
	"fmt"
	"math"
	"math/rand"
)

type Range struct {
	Left  int
	Right int // exclusive
}

func (r Range) String() string {
	return fmt.Sprintf("[%d, %d)", r.Left, r.Right)
}

func (r Range) Validate() error {
	if r.Left >= r.Right {
		return fmt.Errorf("left must be less than right")
	}
	return nil
}

func (r Range) PickRandom() int {
	return rand.Intn(r.Right-r.Left) + r.Left
}

func (r Range) Length() int {
	return r.Right - r.Left
}

func (r Range) Split(n int) []Range {
	if n <= 0 {
		return nil
	}

	if n == 1 {
		return []Range{r}
	}

	if n > r.Length() {
		n = r.Length()
	}

	result := make([]Range, n)
	l := int(math.Ceil(float64(r.Length()) / float64(n)))

	for i := 0; i < n; i++ {
		left := i * l
		right := (i + 1) * l

		if right > r.Right {
			right = r.Right
		}

		result[i] = Range{
			Left:  left,
			Right: right,
		}
	}

	return result
}
