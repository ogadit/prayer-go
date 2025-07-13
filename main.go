package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// --------------------------- Structs ---------------------------

type Config struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Method  int    `json:"method"`
	School  int    `json:"school"`
}

type Response struct {
	Data struct {
		Timings map[string]string `json:"timings"`
		Date    struct {
			Readable string `json:"readable"`
			Hijri    struct {
				Day   string `json:"day"`
				Month struct {
					English string `json:"en"`
				} `json:"month"`
				Year string `json:"year"`
			} `json:"hijri"`
		} `json:"date"`
	} `json:"data"`
}

type Prayer struct {
	Name string
	Time time.Time
}

type Prayers []Prayer

// --------------------------- Main ---------------------------

func main() {
	args := os.Args[1:]
	const cachePath = "prayer_cache.json"

	// Load config
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Attempt to read from cache
	var respData *Response
	respData, err = readCachedPrayerData(cachePath)
	if err != nil {
		// Cache missing or stale; fetch from API
		respData, err = fetchPrayerData(config)
		if err != nil {
			log.Fatalf("Failed to fetch data: %v", err)
		}
		_ = writePrayerCache(cachePath, respData)
	}

	// Parse timings
	prayers, err := parsePrayers(respData.Data.Timings)
	if err != nil {
		log.Fatalf("Failed to parse prayers: %v", err)
	}

	now := time.Now()
	current := prayers.Current(now)
	next := prayers.Next(current)

	// Hijri date
	hijri := fmt.Sprintf(
		"%s %s %s A.H.",
		addOrdinal(respData.Data.Date.Hijri.Day),
		respData.Data.Date.Hijri.Month.English,
		respData.Data.Date.Hijri.Year,
	)

	// Duration & format
	nextTime := next.Time
	nowTime := mustParseTodayTime(now)
	if nextTime.Before(nowTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}
	duration := formatDuration(nextTime.Sub(nowTime))
	nextLine := fmt.Sprintf(
		"%s in %s at %s",
		next.Name,
		duration,
		next.Time.Format("3:04 PM"),
	)
	if len(args) != 0 {
		if args[0] == "-a" {
			fmt.Printf("%s in %s\n", next.Name, duration)
		}
	} else {
		fmt.Printf("%s Time\n", current.Name)
		fmt.Println(nextLine)
		fmt.Println(hijri)
		fmt.Println()
	}
}

func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func fetchPrayerData(cfg *Config) (*Response, error) {
	url := fmt.Sprintf(
		"http://api.aladhan.com/v1/timingsByCity?city=%s&country=%s&method=%d&school=%d",
		cfg.City, cfg.Country, cfg.Method, cfg.School)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &result, nil
}

// --------------------------- Parsing & Helpers ---------------------------

func parsePrayers(timings map[string]string) (Prayers, error) {
	important := map[string]bool{
		"Fajr": true, "Sunrise": true, "Dhuhr": true,
		"Asr": true, "Maghrib": true, "Isha": true,
	}

	prayers := Prayers{}
	location := time.Now().Location()

	for name, raw := range timings {
		if !important[name] {
			continue
		}

		clean := strings.Split(raw, " ")[0]
		t, err := time.ParseInLocation("15:04", clean, location)
		if err != nil {
			log.Printf("Failed to parse %s: %v", name, err)
			continue
		}
		prayers = append(prayers, Prayer{Name: name, Time: t})
	}

	sort.Slice(prayers, func(i, j int) bool {
		return prayers[i].Time.Before(prayers[j].Time)
	})

	return prayers, nil
}

func (p Prayers) Current(now time.Time) Prayer {
	nowClean := mustParseTodayTime(now)
	if nowClean.Before(p[0].Time) {
		return p[len(p)-1]
	}
	current := p[0]
	for _, v := range p {
		if !v.Time.After(nowClean) {
			current = v
		}
	}
	return current
}

func (p Prayers) Next(current Prayer) Prayer {
	for i, v := range p {
		if v == current {
			if i+1 < len(p) {
				return p[i+1]
			}
			return p[0]
		}
	}
	return p[0]
}
