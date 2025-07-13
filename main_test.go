package main

import (
	"testing"
	"time"
)

func mustParseHourMin(hm string) time.Time {
	t, err := time.Parse("15:04", hm)
	if err != nil {
		panic(err)
	}
	return t
}

// Fixed prayer times for a test day
var fixedPrayers = Prayers{
	{"Fajr", mustParseHourMin("04:15")},
	{"Sunrise", mustParseHourMin("05:45")},
	{"Dhuhr", mustParseHourMin("12:30")},
	{"Asr", mustParseHourMin("16:00")},
	{"Maghrib", mustParseHourMin("19:15")},
	{"Isha", mustParseHourMin("20:30")},
}

func TestCurrentAndNextPrayer(t *testing.T) {
	tests := []struct {
		now        string // time string in HH:MM
		wantCur    string
		wantNext   string
		nextOffset time.Duration
	}{
		// Before Fajr
		{"03:00", "Isha", "Fajr", 75 * time.Minute},
		// Just before Fajr
		{"04:10", "Isha", "Fajr", 5 * time.Minute},
		// Just after Fajr
		{"04:20", "Fajr", "Sunrise", 85 * time.Minute},
		// During Dhuhr
		{"13:00", "Dhuhr", "Asr", 180 * time.Minute},
		// Just before Isha
		{"20:25", "Maghrib", "Isha", 5 * time.Minute},
		// After Isha, should wrap to Fajr next day
		{"22:00", "Isha", "Fajr", 375 * time.Minute}, // 6h15m
		// Exactly at Maghrib
		{"19:15", "Maghrib", "Isha", 75 * time.Minute},
	}

	for _, tt := range tests {
		t.Run("Now "+tt.now, func(t *testing.T) {
			now := mustParseHourMin(tt.now)
			current := fixedPrayers.Current(now)
			next := fixedPrayers.Next(current)

			// Adjust next.Time if it's before now (wrap to next day)
			nextTime := next.Time
			if nextTime.Before(now) {
				nextTime = nextTime.Add(24 * time.Hour)
			}

			duration := nextTime.Sub(now)

			if current.Name != tt.wantCur {
				t.Errorf("Current() = %s; want %s", current.Name, tt.wantCur)
			}
			if next.Name != tt.wantNext {
				t.Errorf("Next() = %s; want %s", next.Name, tt.wantNext)
			}
			if abs(duration-tt.nextOffset) > time.Minute {
				t.Errorf("Duration to next = %v; want approx %v", duration, tt.nextOffset)
			}
		})
	}
}

// abs returns the absolute duration difference
func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
