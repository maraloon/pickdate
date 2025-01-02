package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

type keyMap struct {
	Quit       key.Binding
	Help       key.Binding
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Today      key.Binding
	WeekStart  key.Binding
	WeekEnd    key.Binding
	MonthStart key.Binding
	MonthEnd   key.Binding
	Select     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.MonthStart, k.MonthEnd, k.WeekStart, k.WeekEnd},
		{k.Select, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc/ctrl-c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("k/↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("j/↓", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("h/←", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("l/→", "right"),
	),
	Today: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "today"),
	),
	WeekStart: key.NewBinding(
		key.WithKeys("H", "^"),
		key.WithHelp("H/^", "week start"),
	),
	WeekEnd: key.NewBinding(
		key.WithKeys("L", "$"),
		key.WithHelp("L/$", "week end"),
	),
	MonthStart: key.NewBinding(
		key.WithKeys("K", "g"),
		key.WithHelp("K/g", "month start"),
	),
	MonthEnd: key.NewBinding(
		key.WithKeys("J", "G"),
		key.WithHelp("J/G", "month end"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select (copy to buffer)"),
	),
}

type model struct {
	date     time.Time
	selected bool

	keys keyMap
	help help.Model
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

		keys: keys,
		help: help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
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
		case key.Matches(msg, m.keys.Select):
			m.selected = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.selected {
		output := fmt.Sprintf("%d/%02d/%02d\n", m.date.Year(), int(m.date.Month()), m.date.Day())
		termenv.Copy(output)
		return output
	}

	s := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5")).
		Render(
			fmt.Sprintf("   %s %d", m.date.Month(), m.date.Year())+"\nMo Tu We Th Fr Sa Su",
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
			focused := day == m.date.Day()
			var style = lipgloss.NewStyle()

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
			s += style.Render(fmt.Sprintf("%2d ", day))
		}

		s += "\n"
	}

	if len(monthMap) == 4 {
		s += "\n\n"
	} else if len(monthMap) == 5 {
		s += "\n"
	}

	// currentWeekMap := m.monthMap()[m.week()]
	// left := currentWeekMap[0]
	// right := currentWeekMap[6]
	// s += "\n"
	// s += lipgloss.NewStyle().Render(fmt.Sprintf("day: %d\n", m.date.Day()))
	// s += lipgloss.NewStyle().Render(fmt.Sprintf("left: %d\n", left))
	// s += lipgloss.NewStyle().Render(fmt.Sprintf("right: %d\n", right))
	// s += lipgloss.NewStyle().Render(fmt.Sprintf("week: %d\n", m.week()))

	helpView := m.help.View(m.keys)
	// return helpView
	s += helpView

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, s)
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
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
