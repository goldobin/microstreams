package ranges

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type IntRange struct {
	Left  int
	Right int // exclusive
}

func (r IntRange) String() string {
	return fmt.Sprintf("[%d, %d)", r.Left, r.Right)
}

func (r IntRange) Validate() error {
	if r.Left >= r.Right {
		return fmt.Errorf("left must be less than right")
	}
	return nil
}

func (r IntRange) Random(src rand.Source) int {
	rg := rand.New(src)
	return rg.IntN(r.Right-r.Left) + r.Left
}

func (r IntRange) Length() int {
	return r.Right - r.Left
}

func (r IntRange) Split(n int) []IntRange {
	if n <= 0 {
		return nil
	}

	if n == 1 {
		return []IntRange{r}
	}

	if n > r.Length() {
		n = r.Length()
	}

	result := make([]IntRange, n)
	l := int(math.Ceil(float64(r.Length()) / float64(n)))

	for i := 0; i < n; i++ {
		left := i * l
		right := (i + 1) * l

		if right > r.Right {
			right = r.Right
		}

		result[i] = IntRange{
			Left:  left,
			Right: right,
		}
	}

	return result
}
