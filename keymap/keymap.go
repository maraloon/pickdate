package keymap

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
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
	MonthPrev  key.Binding
	MonthNext  key.Binding
	YearPrev   key.Binding
	YearNext   key.Binding
	Select     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.MonthStart, k.MonthEnd, k.WeekStart, k.WeekEnd},
		{k.MonthPrev, k.MonthNext, k.YearPrev, k.YearNext},
		{k.Today, k.Select, k.Help, k.Quit},
	}
}

var Keys = KeyMap{
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
	MonthPrev: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev month"),
	),
	MonthNext: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next month"),
	),
	YearPrev: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "prev year"),
	),
	YearNext: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "next year"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}
