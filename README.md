# keyb

<p align="center">Quick search for application hotkeys at the tip of your fingers.</p>

<p align="center">Powered by Go, bubbletea and lipgloss.</p>

[![keyb demo](https://asciinema.org/a/1fwoiNql5GBKF7lfmSpTpSQTJ.png)](https://asciinema.org/a/1fwoiNql5GBKF7lfmSpTpSQTJ)

### Features:
- Tabular reference for custom key bindings
- Prefix support
- Export formatted output to stdout or file
- fzf, rofi support

### Planned:
- Fuzzy search
- vim key bindings
- Support for json, toml
- More default cheatsheets

>`keyb` does not scrape your applications for hotkeys. All hotkeys must be
>listed manually.

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
  -k, --key [file]      Pass custom key bindings file
  -c, --config [file]   Pass custom config file

  -h, --help            help for keyb
```

keyb requires a `yaml` file with your listed key bindings to work. List your keybindings in the file:

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
      ignore_prefix = true

bspwm:
  keybinds:
    - desc: terminal
      key: Super + Return
```

A `prefix` key can be included for each category. This prefix will be appended
to the beginning of every keybind in that category. A key can choose to opt out
by including a `ignore_prefix=true` field.

Refer to the defaults provided in `examples` for more details.

>Note: Word wrapping is not well supported. Long text of > 88 characters is not
>recommended.

### Configuration
Configure your keyb instance in `$XDG_CONFIG_HOME/keyb/config`.

Refer to `examples/config` for more details.

| Key Binding | Description      |
|------------ | ------------     |
| j,k/Up,Down | Move cursor      |
| Ctrl + u, d | Move half window |
| G			  | Go to bottom     |
| Ctrl + c, q | Quit		     |

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

## Inspiration
- [awesomeWM](https://github.com/awesomeWM/awesome)
- [showkeys](https://github.com/adamharmansky/showkeys)

## License
[MIT](LICENSE)
