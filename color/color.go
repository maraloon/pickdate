package color

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/maraloon/datepicker"
)

func ValidateStdin() (datepicker.Colors, error) {
	colors := make(datepicker.Colors)

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return colors, nil
	}

	reader := bufio.NewReader(os.Stdin)
	stdin, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	groups := strings.Split(strings.TrimSpace(stdin), ";")
	for _, group := range groups {
		parts := strings.SplitN(group, ":", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid format: " + group)
		}

		color := parts[0]
		days := strings.Split(parts[1], ",")

		for _, dayStr := range days {
			day, err := parseDate(dayStr)
			if err != nil {
				return nil, err
			}
			colors[day.Format("2006/01/02")] = color
		}
	}

	return colors, nil
}

func parseDate(dateStr string) (time.Time, error) {
	t, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
