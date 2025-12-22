package rates

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Rate struct {
	Count    int64
	Duration time.Duration
}

func (r Rate) IsZero() bool {
	return r.Count == 0 || r.Duration == 0
}

func (r Rate) String() string {
	if r.IsZero() {
		return "0.00 msg/s"
	}

	perSec := float64(r.Count) / r.Duration.Seconds()
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

func (r Rate) Interval() time.Duration {
	return time.Duration(float64(r.Duration) / float64(r.Count))
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
