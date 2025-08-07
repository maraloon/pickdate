package main

import (
	"fmt"
	"os"
	"time"

	"github.com/maraloon/pickdate/config"
	"github.com/maraloon/pickdate/keymap"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	date   time.Time
	keys   keymap.KeyMap
	help   help.Model
	output string
	quit   bool
	config config.Config
}

type (
	Month []Week
	Week  [7]int
)

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

func firstDayOfMonth(year int, month time.Month, firstWeekDayIsMonday bool) int {
	firstDay := int(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday())
	if firstWeekDayIsMonday {
		firstDay = (firstDay + 6) % 7
	}
	return firstDay
}

func initialModel(config config.Config) *model {
	return &model{
		date:   config.StartAt,
		keys:   keymap.Keys,
		help:   help.New(),
		config: config,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quit = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Today):
			m.date = time.Now()
		case key.Matches(msg, m.keys.Left):
			m.date = m.date.AddDate(0, 0, -1)
		case key.Matches(msg, m.keys.Right):
			m.date = m.date.AddDate(0, 0, 1)
		case key.Matches(msg, m.keys.Down):
			m.date = m.date.AddDate(0, 0, 7)
		case key.Matches(msg, m.keys.Up):
			m.date = m.date.AddDate(0, 0, -7)
		case key.Matches(msg, m.keys.WeekStart):
			d := m.date.Day() - m.monthMap()[m.week()].firstDay()
			m.date = m.date.AddDate(0, 0, -d)
		case key.Matches(msg, m.keys.WeekEnd):
			d := m.monthMap()[m.week()].lastDay() - m.date.Day()
			m.date = m.date.AddDate(0, 0, d)
		case key.Matches(msg, m.keys.MonthStart):
			m.date = time.Date(m.date.Year(), m.date.Month(), 1, 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.keys.MonthEnd):
			m.date = time.Date(m.date.Year(), m.date.Month(), daysInMonth(m.date.Year(), m.date.Month()), 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.keys.MonthPrev):
			m.date = m.date.AddDate(0, -1, 0)
		case key.Matches(msg, m.keys.MonthNext):
			m.date = m.date.AddDate(0, 1, 0)
		case key.Matches(msg, m.keys.YearPrev):
			m.date = m.date.AddDate(-1, 0, 0)
		case key.Matches(msg, m.keys.YearNext):
			m.date = m.date.AddDate(1, 0, 0)
		case key.Matches(msg, m.keys.Select):
			m.output = m.date.Format(m.config.OutputFormat)
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) View() string {
	if m.output != "" || m.quit {
		return ""
	}

	var weekLegend string
	if m.config.FirstWeekdayIsMo {
		weekLegend = "Mo Tu We Th Fr Sa Su"
	} else {
		weekLegend = "Su Mo Tu We Th Fr Sa"
	}

	s := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5")).
		Render(
			fmt.Sprintf("    %s %d", m.date.Month(), m.date.Year())+"\n "+weekLegend,
		) + "\n"

	monthMap := m.monthMap()
	for _, week := range monthMap {
		for k, day := range week {
			if day == 0 {
				s += "   "
				continue
			}

			today := day == time.Now().Day() && m.date.Month() == time.Now().Month() && m.date.Year() == time.Now().Year()
			weekend := k >= 5
			if !m.config.FirstWeekdayIsMo {
				weekend = k == 0 || k == 6
			}
			focused := day == m.date.Day()
			style := lipgloss.NewStyle()

			if today {
				if focused {
					style = style.Background(lipgloss.Color("9")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Foreground(lipgloss.Color("9"))
				}
			} else if weekend {
				if focused {
					style = style.Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Foreground(lipgloss.Color("4"))
				}
			} else {
				if focused {
					style = style.Background(lipgloss.Color("3")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Foreground(lipgloss.Color("3"))
				}
			}
			s += style.Render(fmt.Sprintf(" %2d", day))
		}

		s += "\n"
	}

	if len(monthMap) == 4 {
		s += "\n\n"
	} else if len(monthMap) == 5 {
		s += "\n"
	}

	s += m.help.View(m.keys)

	w, h, _ := term.GetSize(int(os.Stderr.Fd()))
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, s)
}

func (m *model) monthMap() Month {
	daysInMonth := daysInMonth(m.date.Year(), m.date.Month())
	startDay := firstDayOfMonth(m.date.Year(), m.date.Month(), m.config.FirstWeekdayIsMo)

	monthMap := make(Month, 0)
	week := Week{}
	dayCounter := 1

	// Fill the first week with leading zeros
	for i := range startDay {
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
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	config, err := config.ValidateFlags()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	model := initialModel(config)
	p := tea.NewProgram(model, tea.WithOutput(tty))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	if model.quit {
		os.Exit(130)
	}

	fmt.Print(model.output)
}
