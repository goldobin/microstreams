package rate

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Rate float64

func (r Rate) String() string {
	if r == 0.0 {
		return "0.000 msg/s"
	}

	if r >= 1 {
		return fmt.Sprintf("%.3f msg/s", r)
	}
	perMinute := r * 60
	if perMinute >= 1 {
		return fmt.Sprintf("%.3f msg/min", perMinute)
	}
	perHour := perMinute * 60
	if perHour >= 1 {
		return fmt.Sprintf("%.3f msg/h", perHour)
	}
	perDay := perHour * 24
	if perDay >= 1 {
		return fmt.Sprintf("%.3f msg/day", perDay)
	}

	return fmt.Sprintf("%.3f msg/day", perDay)
}

func (r Rate) Interval() time.Duration {
	return time.Duration(1e9 / r)
}

func (r Rate) Random(source rand.Source) time.Duration {
	var (
		rg    = rand.New(source)
		v     = r.Interval()
		randV = rg.Int64N(int64(v))
	)
	if randV <= 1 {
		return time.Duration(1)
	} else {
		return time.Duration(randV)
	}
}

const (
	PerSecond Rate = 1
	PerMinute Rate = 1 / 60.0
	PerHour   Rate = 1 / (60.0 * 60.0)
	PerDay    Rate = 1 / (60.0 * 60.0 * 24.0)
)
