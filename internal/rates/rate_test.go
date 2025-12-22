package rates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		count    int64
		duration time.Duration
		wantRate Rate
	}{
		{
			name:     "case 1",
			count:    100,
			duration: time.Second,
			wantRate: Rate{Count: 100, Duration: time.Second},
		},
		{
			name:     "case 2",
			count:    0,
			duration: time.Second,
			wantRate: Rate{Count: 0, Duration: time.Second},
		},
		{
			name:     "case 3",
			count:    1,
			duration: 0,
			wantRate: Rate{Count: 1, Duration: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Rate{tt.count, tt.duration}
			assert.Equal(t, tt.wantRate, got)
		})
	}
}

//func TestNew_Panics(t *testing.T) {
//	tests := []struct {
//		name     string
//		count    int64
//		duration time.Duration
//	}{
//		{
//			name:     "case 1",
//			count:    -1,
//			duration: time.Second,
//		},
//		{
//			name:     "case 2",
//			count:    100,
//			duration: -time.Second,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			assert.Panics(t, func() {
//				Rate{tt.count, tt.duration}
//			})
//		})
//	}
//}

func TestRate_IsZero(t *testing.T) {
	tests := []struct {
		name string
		rate Rate
		want bool
	}{
		{
			name: "case 1",
			rate: Rate{0, time.Second},
			want: true,
		},
		{
			name: "case 2",
			rate: Rate{100, 0},
			want: true,
		},
		{
			name: "case 3",
			rate: Rate{0, 0},
			want: true,
		},
		{
			name: "case 4",
			rate: Rate{100, time.Second},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rate.IsZero()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRate_String(t *testing.T) {
	tests := []struct {
		name string
		rate Rate
		want string
	}{
		{
			name: "case 01",
			rate: Rate{0, time.Second},
			want: "0.00 msg/s",
		},
		{
			name: "case 02",
			rate: Rate{100, 0},
			want: "0.00 msg/s",
		},
		{
			name: "case 03",
			rate: Rate{100, time.Second},
			want: "100.00 msg/s",
		},
		{
			name: "case 04",
			rate: Rate{1, 2 * time.Second},
			want: "30.00 msg/min",
		},
		{
			name: "case 05",
			rate: Rate{30, time.Minute},
			want: "30.00 msg/min",
		},
		{
			name: "case 06",
			rate: Rate{1, 2 * time.Minute},
			want: "30.00 msg/h",
		},
		{
			name: "case 07",
			rate: Rate{1, 10 * time.Minute},
			want: "6.00 msg/h",
		},
		{
			name: "case 08",
			rate: Rate{1, 2 * time.Hour},
			want: "12.00 msg/day",
		},
		{
			name: "case 09",
			rate: Rate{1, 10 * time.Hour},
			want: "2.40 msg/day",
		},
		{
			name: "case 10",
			rate: Rate{1, 48 * time.Hour},
			want: "0.50 msg/day",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rate.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRate_Interval(t *testing.T) {
	tests := []struct {
		name string
		rate Rate
		want time.Duration
	}{
		{
			name: "case 1",
			rate: Rate{100, time.Second},
			want: 10 * time.Millisecond,
		},
		{
			name: "case 2",
			rate: Rate{1, time.Second},
			want: time.Second,
		},
		{
			name: "case 3",
			rate: Rate{2, time.Second},
			want: 500 * time.Millisecond,
		},
		{
			name: "case 4",
			rate: Rate{60, time.Minute},
			want: time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rate.Interval()
			assert.Equal(t, tt.want, got)
		})
	}
}
