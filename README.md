# keyb

<p align="center">Global, custom hotkeys reference</p>

<p align="center">
	<img width="500" src="https://github.com/kencx/keyb/blob/master/assets/keyb.png?raw=true">
</p>

<p align="center">Create your own hotkeys reference sheet for quick reference from the terminal</p>

### Features
- Fuzzy finding
- Vim key bindings
- Export reference to stdout, yaml, json
- fzf, rofi support
- Fully customizable
- Single binary

### Non-Features
keyb does **not** support:

- Auto detection of hotkeys
- Setting of hotkeys for any applications
- Command selection

## Motivation

When I started out using tmux and bspwm, I couldn't remember all the hotkeys and
resorted to writing them down on paper, creating a physical cheatsheet. Then, I
thought: maybe there's a tool that does this. I didn't find one I liked so I
built `keyb`.

With `keyb`, I can list:

- Hotkeys that I occasionally forget
- Custom key combinations that I defined for my own workflow
- Hotkeys for tools that I'm new to

This idea is not novel - awesomeWM provides a helpful (but non configurable)
cheatsheet for its defined hotkeys.

## Install
keyb is still a work in progress. Use it at your own risk.

```bash
$ go get -u github.com/kencx/keyb
```

## Usage

```text
usage: keyb [options] [file]

Flags:
  -p, --print           Print to stdout
  -e, --export [file]   Export to file
  -k, --key [file]      Pass custom hotkeys file
  -c, --config [file]   Pass custom config file

  -h, --help            help for keyb
```

keyb requires a `yaml` file with a list of hot keys:

```yaml
- name: bspwm
  keybinds:
    - name: terminal
      key: Super + Return

- name: tmux
  prefix: ctrl + a
  keybinds:
    - name: split vertical
      key: "|"
    - name: split horizontal
      key: "-"
    - name: {next, prev} window
      key: shift + {>, <}
      ignore_prefix: true
```

A `prefix` key can be included for each category. This prefix will be appended
to the beginning of every keybind in that category. A key can choose to opt out
by including a `ignore_prefix=true` field.

Refer to the defaults provided in `examples` for more details.

`keyb` is best used when binding to a hotkey for quick referencing.
For example, with sxhkd and st:

```
super + slash
	st -c keys keyb
```

### Configuration
A config file is generated on the first run of `keyb`.

The defaults are:

```yml
defaults:
  keyb_path: "$HOME/.config/keyb/keyb.yml"
  debug: false
  reverse: false
  mouse: true
  title: ""
  prompt: "keys > "
  placeholder: "search..."
  prefix_sep: ";"
  sep_width: 4
```

Refer to `examples/config.yml` for more details.

| Hotkey                     | Description      |
|--------------------------- | ---------------- |
| <kbd>j, k / Up, Down</kbd> | Move cursor      |
| <kbd>Ctrl + u, d</kbd>     | Move half window |
| <kbd>Ctrl + b, f</kbd>     | Move full window |
| <kbd>H, M, L</kbd>         | Go to top, middle, bottom of screen |
| <kbd>g, G</kbd>            | Go to first, last line |
| <kbd>/</kbd>               | Enter search mode|
| <kbd>alt + d</kbd>         | Clear current search |
| <kbd>Esc</kbd>             | Exit search mode |
| <kbd>Ctrl + c, q</kbd>     | Quit		        |

## fzf, rofi

keyb also supports printing to stdout for use with other tools:

```bash
$ keyb -p | fzf
$ keyb -p | rofi -dmenu
```

## Acknowledgements
`keyb` is built with:

- [bubbletea](github.com/charmbracelet/bubbletea)
- [lipgloss](github.com/charmbracelet/lipgloss)
- [fuzzy](github.com/sahilm/fuzzy)
- [ansiterm](github.com/juju/ansiterm)

Alternatives that I found:

- [awesomeWM](https://github.com/awesomeWM/awesome)
- [showkeys](https://github.com/adamharmansky/showkeys)

## License
[MIT](LICENSE)
