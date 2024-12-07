package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	date     time.Time
	selected bool
}

type Month []Week
type Week [7]int

func (week Week) firstDay() int {
	for _, day := range week {
		if day != 0 {
			return day
		}
	}
	return 0
}

func (week Week) lastDay() int {
	for i := 6; i >= 0; i-- {
		if week[i] != 0 {
			return week[i]
		}
	}
	return 0
}

func (m model) week() int {
	firstDay := time.Date(m.date.Year(), m.date.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstWeekday := firstDay.Weekday()
	if firstWeekday == 0 {
		firstWeekday = 7
	}
	return (m.date.Day() + int(firstWeekday-2)) / 7
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func firstDayOfMonth(year int, month time.Month) int {
	return (int(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday()) + 6) % 7
	// TODO: this return is work when week start from Sunday, so we can easy implement it
	// return int(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday())
}

func initialModel() model {
	return model{
		date:     time.Now(),
		selected: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			m.date = m.date.AddDate(0, 0, -1)
		case "right", "l":
			m.date = m.date.AddDate(0, 0, 1)
		case "down", "j":
			m.date = m.date.AddDate(0, 0, 7)
		case "up", "k":
			m.date = m.date.AddDate(0, 0, -7)
		case "^", "H":
			d := m.date.Day() - m.monthMap()[m.week()].firstDay()
			m.date = m.date.AddDate(0, 0, -d)
		case "$", "L":
			d := m.monthMap()[m.week()].lastDay() - m.date.Day()
			m.date = m.date.AddDate(0, 0, d)
		case "g", "J":
			m.date = time.Date(m.date.Year(), m.date.Month(), 1, 0, 0, 0, 0, time.UTC)
		case "G":
			m.date = time.Date(m.date.Year(), m.date.Month(), daysInMonth(m.date.Year(), m.date.Month()), 0, 0, 0, 0, time.UTC)
		case "enter":
			m.selected = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.selected {
		return fmt.Sprintf("%d/%02d/%02d\n", m.date.Year(), int(m.date.Month()), m.date.Day())
	}

	var legend = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5"))
	s := legend.Render(fmt.Sprintf("   %s %d", m.date.Month(), m.date.Year()))
	s += "\n"
	s += legend.Render("Mo Tu We Th Fr Sa Su")
	s += "\n"

	var style = lipgloss.NewStyle()

	for _, week := range m.monthMap() {
		for k, day := range week {
			focused := day == m.date.Day()
			today := day == time.Now().Day() && m.date.Month() == time.Now().Month() && m.date.Year() == time.Now().Year()
			if day == 0 {
				s += "   "
			} else {
				if focused {
					if today {
						s += style.Background(lipgloss.Color("9")).
							Foreground(lipgloss.Color("0")).
							Render(fmt.Sprintf("%2d ", day))
					} else {
						s += style.Background(lipgloss.Color("15")).
							Foreground(lipgloss.Color("0")).
							Render(fmt.Sprintf("%2d ", day))
					}
				} else if today {
					s += style.Foreground(lipgloss.Color("9")).
						Render(fmt.Sprintf("%2d ", day))
				} else if k >= 5 {
					s += style.Foreground(lipgloss.Color("4")).
						Render(fmt.Sprintf("%2d ", day))
				} else {
					s += style.Render(fmt.Sprintf("%2d ", day))
				}
			}
		}

		s += "\n"
	}

	currentWeekMap := m.monthMap()[m.week()]
	left := currentWeekMap[0]
	right := currentWeekMap[6]
	s += "\n"
	s += style.Render(fmt.Sprintf("day: %d\n", m.date.Day()))
	s += style.Render(fmt.Sprintf("left: %d\n", left))
	s += style.Render(fmt.Sprintf("right: %d\n", right))
	s += style.Render(fmt.Sprintf("week: %d\n", m.week()))

	return s
}

func (m model) monthMap() Month {
	daysInMonth := daysInMonth(m.date.Year(), m.date.Month())
	startDay := firstDayOfMonth(m.date.Year(), m.date.Month())

	monthMap := make(Month, 0)
	week := Week{}
	dayCounter := 1

	// Fill the first week with leading zeros
	for i := 0; i < startDay; i++ {
		week[i] = 0
	}

	// Fill the days of the month
	for dayCounter <= daysInMonth {
		week[startDay] = dayCounter
		dayCounter++
		startDay++

		// If the week is full, add it to the weeks slice and reset
		if startDay == 7 {
			monthMap = append(monthMap, week)
			week = Week{}
			startDay = 0
		}
	}

	// Add the last week if it has any days
	if startDay > 0 {
		monthMap = append(monthMap, week)
	}

	return monthMap
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
