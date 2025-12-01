package rates

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Rate struct {
	count    int64
	duration time.Duration
}

func (r Rate) IsZero() bool {
	return r.count == 0 || r.duration == 0
}

func (r Rate) String() string {
	if r.IsZero() {
		return "0.00 msg/s"
	}

	perSec := float64(r.count) / r.duration.Seconds()
	if perSec >= 1 {
		return fmt.Sprintf("%.2f msg/s", perSec)
	}
	perMinute := perSec * 60
	if perMinute >= 1 {
		return fmt.Sprintf("%.2f msg/min", perMinute)
	}
	perHour := perMinute * 60
	if perHour >= 1 {
		return fmt.Sprintf("%.2f msg/h", perHour)
	}
	perDay := perHour * 24
	if perDay >= 1 {
		return fmt.Sprintf("%.2f msg/day", perDay)
	}

	return fmt.Sprintf("%.2f msg/day", perDay)
}

func New(count int64, duration time.Duration) Rate {
	if count < 0 {
		panic("count must be >= 0")
	}
	if duration < 0 {
		panic("duration must be >= 0")
	}
	return Rate{count: count, duration: duration}
}

func (r Rate) Interval() time.Duration {
	return time.Duration(float64(r.duration) / float64(r.count))
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
