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
  -s, --strip           Strip ANSI chars (only for print/export)
  -k, --key [file]      Pass custom hotkeys file
  -c, --config [file]   Pass custom config file

  -h, --help            help for keyb
```

keyb requires a `yaml` file with a list of hotkeys:

```yaml
# $XDG_CONFIG_HOME/keyb/keyb.yml
tmux:
  prefix: ctrl + a
  keybinds:
    - desc: split vertical
      key: "|"
    - desc: split horizontal
      key: "-"
    - desc: {next, prev} window
      key: shift + {>, <}
      ignore_prefix: true

bspwm:
  keybinds:
    - desc: terminal
      key: Super + Return
```

A `prefix` key can be included for each category. This prefix will be appended
to the beginning of every keybind in that category. A key can choose to opt out
by including a `ignore_prefix=true` field.

Refer to the defaults provided in `examples` for more details.

Finally, bind `keyb` to a hotkey for quick reference. For
example, with sxhkd and st:

```
super + slash
	st -c keys keyb
```

### Configuration
Configure your keyb instance in `$XDG_CONFIG_HOME/keyb/config`.

Refer to `examples/config` for more details.

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

keyb also supports the export of a formatted output for use with other programs
like fzf:
```bash
$ keyb -p | fzf
```

For rofi, the output must be stripped of ansi formatting first:
```bash
$ keyb -e output.txt -s
$ cat output.txt | rofi -dmenu
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
