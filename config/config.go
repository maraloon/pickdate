package config

import (
	"errors"
	"flag"
)

type Config struct {
	FirstWeekdayIsMo bool
	OutputFormat     string
}

func ValidateFlags() (Config, error) {
	var firstWeekday string
	flag.StringVar(&firstWeekday, "first-weekday", "mo", "Render calendar starting from selected weekday [mo/su]")
	var format string
	flag.StringVar(&format, "format", "yyyy/mm/dd", "Format of date output")
	flag.Parse()

	layout, err := transformDateLayout(format)
	if err != nil {
		return Config{}, err
	}

	return Config{
		FirstWeekdayIsMo: firstWeekday == "mo",
		OutputFormat:     layout,
	}, nil
}

func transformDateLayout(input string) (string, error) {
	dateformats := map[string]string{
		"yyyy/mm/dd": "2006/01/02",
		"Y/m/d":      "2006/01/02",
		"yyyy-mm-dd": "2006-01-02",
		"Y-m-d":      "2006-01-02",
		"F j, Y":     "January 2, 2006",
		"m/d/y":      "01/02/06",
		"M-d-y":      "Jan-02-06",
		"l":          "Monday",
		"D":          "Mon",
		"d":          "02",
		"j":          "2",
		"F":          "January",
		"M":          "Jan",
		"m":          "01",
		"n":          "1",
		"Y":          "2006",
		"y":          "06",
		// TODO: need some tweaks for realization
		// "N": "1 (for Monday) through 7 (for Sunday)",
		// "w": "0 (for Sunday) through 6 (for Saturday)",
		// "z": "0 through 365",
		// "W": "42 (the 42nd week in the year)",
		// "t": "Number of days in the given month	28 through 31",
		// "L": "Whether it's a leap year	1 if it is a leap year, 0 otherwise.",
	}

	if layout, exists := dateformats[input]; exists {
		return layout, nil
	} else if _, found := findKeyByValue(dateformats, input); found {
		return input, nil
	} else {
		return "", errors.New("wrong date layout")
	}
}

func findKeyByValue(m map[string]string, value string) (string, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return "", false
}
