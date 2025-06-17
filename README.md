# TUI datepicker

fastly select date via vim-motions and print it to `stdout`
> WIP but work enough. See `todo.md` for planning functional

![showcase](readme/preview.png) 

## Install

```bash
git clone git@github.com:maraloon/tui-datepicker.git
go install
```

## Usage idea

It's for what i develop this app. Terminal-based notes. Open (or create) file for selected date

```bash
#!/usr/bin/env sh
tui-datepicker
nvim "$HOME/diary/$(wl-paste).md" # opens smth like ~/diary/2025/01/15.md
```

![usage](readme/usage.gif) 

## Made with

<p><a href="https://stuff.charm.sh/bubbletea/bubbletea-4k.png"><img src="https://github.com/charmbracelet/bubbletea/assets/25087/108d4fdb-d554-4910-abed-2a5f5586a60e" width="313" alt="Bubble Tea Title Treatment"></a></p>
