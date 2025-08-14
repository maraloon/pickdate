package main

import (
	"fmt"
	"os"

	"github.com/maraloon/datepicker"
	"github.com/maraloon/pickdate/color"
	"github.com/maraloon/pickdate/config"
	"github.com/maraloon/pickdate/keymap"
	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	exitCancel = iota + 1
	exitSelect
)

type mainModel struct {
	cal  *datepicker.Model
	help help.Model
	quit *int8
}

func InitModel(cal *datepicker.Model) *mainModel {
	q := int8(0)
	return &mainModel{
		cal:  cal,
		help: cal.Help,
		quit: &q,
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.cal.Init())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "ctrl+c", "q", "esc":
			*m.quit = exitCancel
			return m, tea.Quit
		case "enter":
			*m.quit = exitSelect
			return m, tea.Quit
		default:
			m.cal.Update(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	if *m.quit > 0 {
		return ""
	}

	var s string
	s += lipgloss.JoinVertical(lipgloss.Center, lipgloss.NewStyle().Render(m.cal.View()))
	s += "\n" + m.help.View(keymap.Keys) + "\n"

	w, h, _ := term.GetSize(int(os.Stderr.Fd()))
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, s)
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
	config.HideHelp = true

	colors, err := color.ValidateStdin()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	cal := datepicker.InitModel(config, colors)
	model := InitModel(cal)
	p := tea.NewProgram(model, tea.WithOutput(tty))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	switch *model.quit {
	case exitCancel:
		os.Exit(130)
	case exitSelect:
		fmt.Print(model.cal.CurrentValue())
	}
}
