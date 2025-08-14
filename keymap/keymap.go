package keymap

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/maraloon/datepicker"
)

type KeyMap struct {
	datepicker.KeyMap
	Quit   key.Binding
	Select key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	bindings := k.KeyMap.FullHelp()
	bindings[3] = []key.Binding{k.Today, k.Select, k.Help, k.Quit}
	return bindings
}

var Keys = KeyMap{
	KeyMap: datepicker.Keys,
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc/ctrl-c", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}
