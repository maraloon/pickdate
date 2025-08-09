# Pickdate

fastly select date via vim-motions and print it to `stdout`

![showcase](readme/preview.png) 

## Install

### arch linux:

``` bash
yay -S install pickdate 
```

### manualy:

```bash
git clone git@github.com:maraloon/pickdate.git
cd pickdate
go install
```

## Usage idea

It's for what i develop this app. Terminal-based notes. Open (or create) file for selected date

```bash
#!/usr/bin/env sh
prev_selected_date=$(date +"%Y/%m/%d")
while true
do 
    selected_date=$(pickdate -m --start-at $prev_selected_date) || exit 1
    nvim "$HOME/diary/$selected_date.md" # opens diary/2025/01/15.md
    prev_selected_date=$selected_date # alows to stay in selected date after quit editor
done
```

![usage](readme/usage.gif) 

## Flags

```
Usage: pickdate [OPTIONS]

Options:
  -f, --format string     Format of date output (default "yyyy/mm/dd")
  -h, --help              Help
  -m, --monday            Monday as first day of week
      --start-at string   Pointed date on enter (default today)
  -s, --sunday            Sunday as first day of week (default true)
```

### `--format` values

You can use both left and right format types


|   Format     | Go Layout         |
|--------------|-------------------|
| `yyyy/mm/dd` | `2006/01/02`      |
| `Y/m/d`      | `2006/01/02`      |
| `yyyy-mm-dd` | `2006-01-02`      |
| `Y-m-d`      | `2006-01-02`      |
| `F j, Y`     | `January 2, 2006` |
| `m/d/y`      | `01/02/06`        |
| `M-d-y`      | `Jan-02-06`       |
| `l`          | `Monday`          |
| `D`          | `Mon`             |
| `d`          | `02`              |
| `j`          | `2`               |
| `F`          | `January`         |
| `M`          | `Jan`             |
| `m`          | `01`              |
| `n`          | `1`               |
| `Y`          | `2006`            |
| `y`          | `06`              |


## Set custom colors for days

You can send string to stdin to set custom colors for each day.
String format: `color1:day1,day2;color2:day3,day31`
`color` - `[0-15]` or hex (`#b16286`)
`day` - `2006/01/02`

example: `echo "#b16286:2025/08/10,2025/08/11;#d79920:2025/08/12" | pickdate`

## TODO

- [x] Show today, style
- [x] Help menu
- [ ] Jumps
    - [x] Jump to today
    - [ ] Month jump
        - [x] p, n
        - [ ] m[1-12]<cr>
    - [ ] Year jump
        - [x] P, N
        - [ ] y[1-12]<cr>
    - [ ] Jump in line: 3l - 3 days later
    - [ ] Jump n month up/down: 3ml/3m<down> - 3 month down 
    - [ ] Jump lines: 2j - 2 weeks later
    - [ ] Jump to selected day: `d[1-31]`/`31g`/`31<cr>` will jump on 31th day of current month
- [ ] Lists
    - [ ] Month list (M)
    - [ ] Year list (Y)
- [ ] View
    - [ ] Show 3 month view
    - [ ] Show full year view
    - [ ] Change colors via config
- [ ] Toggle fullsceen (WithAltScreen)
- [x] Center align
- [ ] Toggle week start, monday or sunday
- [ ] CLI opts
    - [x] Week first day
    - [x] Output date format 
    - [x] `--start-at date`
    - [ ] Fullscreen
- aur
    - [x] stable
    - [ ] git version

## Made with

<p><a href="https://stuff.charm.sh/bubbletea/bubbletea-4k.png"><img src="https://github.com/charmbracelet/bubbletea/assets/25087/108d4fdb-d554-4910-abed-2a5f5586a60e" width="313" alt="Bubble Tea Title Treatment"></a></p>
