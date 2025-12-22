package ranges

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Int struct {
	Left  int
	Right int // exclusive
}

func (r Int) String() string {
	return fmt.Sprintf("[%d, %d)", r.Left, r.Right)
}

func (r Int) Validate() error {
	if r.Left >= r.Right {
		return fmt.Errorf("left must be less than right")
	}
	return nil
}

func (r Int) Random(src rand.Source) int {
	rg := rand.New(src)
	return rg.IntN(r.Right-r.Left) + r.Left
}

func (r Int) Length() int {
	return r.Right - r.Left
}

func (r Int) Split(n int) []Int {
	if n <= 0 {
		return nil
	}

	if n == 1 {
		return []Int{r}
	}

	if n > r.Length() {
		n = r.Length()
	}

	result := make([]Int, n)
	l := int(math.Ceil(float64(r.Length()) / float64(n)))

	for i := 0; i < n; i++ {
		left := i * l
		right := (i + 1) * l

		if right > r.Right {
			right = r.Right
		}

		result[i] = Int{
			Left:  left,
			Right: right,
		}
	}

	return result
}
