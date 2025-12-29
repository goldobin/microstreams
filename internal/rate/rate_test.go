package rate_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goldobin/microstreams/internal/rate"
)

func TestRate_String(t *testing.T) {
	tests := []struct {
		name string
		rate rate.Rate
		want string
	}{
		{
			name: "case 1.1 - zero rate",
			rate: 0,
			want: "0.000 msg/s",
		},
		{
			name: "case 1.2 - one per second",
			rate: 1 * rate.PerSecond,
			want: "1.000 msg/s",
		},
		{
			name: "case 1.3 - high rate per second",
			rate: 100.0 * rate.PerSecond,
			want: "100.000 msg/s",
		},
		{
			name: "case 1.4 - fractional per second",
			rate: 10.5 * rate.PerSecond,
			want: "10.500 msg/s",
		},
		{
			name: "case 2.1 - half per second",
			rate: (1.0 / 2.0) * rate.PerSecond,
			want: "30.000 msg/min",
		},
		{
			name: "case 2.2 - thirty per minute",
			rate: 30 * rate.PerMinute,
			want: "30.000 msg/min",
		},
		{
			name: "case 2.3 - one per minute",
			rate: 1 * rate.PerMinute,
			want: "1.000 msg/min",
		},
		{
			name: "case 2.4 - fractional per minute",
			rate: 15.5 * rate.PerMinute,
			want: "15.500 msg/min",
		},
		{
			name: "case 3.1 - half per minute",
			rate: 0.5 * rate.PerMinute,
			want: "30.000 msg/h",
		},
		{
			name: "case 3.2 - one tenth per minute",
			rate: 0.1 * rate.PerMinute,
			want: "6.000 msg/h",
		},
		{
			name: "case 3.3 - one per hour",
			rate: 1 * rate.PerHour,
			want: "1.000 msg/h",
		},
		{
			name: "case 3.4 - twelve per hour",
			rate: 12 * rate.PerHour,
			want: "12.000 msg/h",
		},
		{
			name: "case 4.1 - half per hour",
			rate: 0.5 * rate.PerHour,
			want: "12.000 msg/day",
		},
		{
			name: "case 4.2 - 2.4 per day",
			rate: 2.4 * rate.PerDay,
			want: "2.400 msg/day",
		},
		{
			name: "case 4.3 - one per 48 hours",
			rate: (1.0 / 48.0) * rate.PerHour,
			want: "0.500 msg/day",
		},
		{
			name: "case 4.4 - one per day",
			rate: 1 * rate.PerDay,
			want: "1.000 msg/day",
		},
		{
			name: "case 4.5 - very low rate",
			rate: 0.1 * rate.PerDay,
			want: "0.100 msg/day",
		},
		{
			name: "case 4.6 - extremely low rate",
			rate: 0.001 * rate.PerDay,
			want: "0.001 msg/day",
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
		rate rate.Rate
		want time.Duration
	}{
		{
			name: "case 1.1 - high rate per second",
			rate: 100 * rate.PerSecond,
			want: 10 * time.Millisecond,
		},
		{
			name: "case 1.2 - one per second",
			rate: 1 * rate.PerSecond,
			want: time.Second,
		},
		{
			name: "case 1.3 - two per second",
			rate: 2 * rate.PerSecond,
			want: 500 * time.Millisecond,
		},
		{
			name: "case 1.4 - ten per second",
			rate: 10 * rate.PerSecond,
			want: 100 * time.Millisecond,
		},
		{
			name: "case 2.1 - sixty per minute",
			rate: 60 * rate.PerMinute,
			want: time.Second,
		},
		{
			name: "case 2.2 - one per minute",
			rate: 1 * rate.PerMinute,
			want: time.Minute,
		},
		{
			name: "case 2.3 - thirty per minute",
			rate: 30 * rate.PerMinute,
			want: 2 * time.Second,
		},
		{
			name: "case 2.4 - one per hour",
			rate: 1 * rate.PerHour,
			want: time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rate.Interval()
			assert.Equal(t, tt.want, got)
		})
	}
}
