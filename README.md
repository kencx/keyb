# keyb

<p align="center">Quick search for application hotkeys at the tip of your fingers.</p>

<p align="center">Powered by koanf, bubbletea and lipgloss.</p>

[![keyb demo](https://asciinema.org/a/1fwoiNql5GBKF7lfmSpTpSQTJ.png)](https://asciinema.org/a/1fwoiNql5GBKF7lfmSpTpSQTJ)

### Features:
- Tabular reference for custom key bindings
- Export formatted output to stdout or file
- fzf, rofi support

### Planned:
- Prefix support
- Fuzzy search
- vim key bindings
- Support for json, toml
- More default cheatsheets

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

keyb requires a `yaml` file with your listed key bindings to work. List your keybindings in the file as so:
```yaml
# $XDG_CONFIG_HOME/keyb/keyb.yml
tmux:
  prefix: ctrl + a
  keybinds:
    - desc: split vertical
      key: "|"
    - desc: split horizontal
      key: "-"

bspwm:
  keybinds:
    - desc: terminal
      key: Super + Return
```
or refer to the defaults provided in `examples`.

>Note: Word wrapping is not well supported. Long text of > 88 characters is not
>recommended.

### Configuration
Configure your keyb instance with `$XDG_CONFIG_HOME/keyb/config`.

Refer to `examples/config` for the default config.

### Navigation

| Key Binding | Description  |
|------------ | ------------ |
| j,k/Up,Down | Move cursor  |
| G			  | Go to bottom |
| Ctrl + c, q | Quit		 |


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
