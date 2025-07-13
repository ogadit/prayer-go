package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// mustParseTodayTime takes a full time and returns the same hour:minute on today's date.
// This helps you strip date info when comparing only on prayer times.
func mustParseTodayTime(t time.Time) time.Time {
	str := t.Format("15:04")
	result, err := time.ParseInLocation("15:04", str, t.Location())
	if err != nil {
		log.Fatalf("failed to parse current time: %v", err)
	}
	return result
}

// formatDuration returns a clean human-readable version of a time.Duration
// e.g., "1 Hour 30 Minutes" or "3 Minutes"
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	if h == 0 && m == 0 {
		return "less than a minute"
	}

	if h > 0 {
		return fmt.Sprintf("%d Hour%s %d Minute%s", h, plural(h), m, plural(m))
	}
	return fmt.Sprintf("%d Minute%s", m, plural(m))
}

// plural returns "s" for plural, "" for singular
func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// addOrdinal adds the English ordinal suffix to a numeric day
// e.g., "1" => "1st", "2" => "2nd"
func addOrdinal(day string) string {
	switch day {
	case "1", "21", "31":
		return day + "st"
	case "2", "22":
		return day + "nd"
	case "3", "23":
		return day + "rd"
	default:
		return day + "th"
	}
}

func readCachedPrayerData(cachePath string) (*Response, error) {
	file, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err := json.Unmarshal(file, &resp); err != nil {
		return nil, err
	}

	// Check if the cache is for today
	today := time.Now().Format("02 Jan 2006")
	if resp.Data.Date.Readable != today {
		return nil, fmt.Errorf("cache is stale")
	}

	return &resp, nil
}

func writePrayerCache(cachePath string, data *Response) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath, bytes, 0644)
}
