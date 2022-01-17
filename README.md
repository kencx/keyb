# keyb

`keyb` is a global hotkey cheatsheet inspired by awesomeWM. 

It offers a helpful reference to your custom key bindings for any program by parsing all listed key bindings and displaying them in a neat, tabular interface.

<p align="center">
<img src="https://raw.githubusercontent.com/kencx/keyb/master/assets/keyb.png" alt="showcase">
</p>

### Features:
- Categorize your key bindings for easy reference
- Export formatted output to stdout or file (compatible with fzf or rofi)

### Planned:
- Prefix support
- Fuzzy search
- vim key bindings
- Support for json, toml
- More default cheatsheets

## Install
keyb is still a work in progress and not ready.

```bash
$ go get -u github.com/kencx/keyb
```

## Usage

```bash
usage: keyb [options] [file]

Flags:
  -p, --print		Print to stdout
  -e, --export [file]	Export to file
  -k, --key [file]	Pass custom key bindings file
  -c, --config [file]	Pass custom config file

  -h, --help		help for keyb
```


keyb requires a `yaml` file with your listed key bindings to work. List your keybindings in the file as so:
```yaml
tmux:
  prefix: ctrl + a
  keybinds:
    - desc: split vertical
      key: "|"
    - desc: split horizontal
      key: "-"
bspwm:
  keybinds:
    - desc: open configs
      key: super + z
```
or refer to the defaults provided in `examples`.

### Navigation

| Key Binding | Description  |
|------------ | ------------ |
| h,j,k,l     | Move cursor  |
| Arrow keys  | Move cursor	 |
| G			  | Go to bottom |
| Ctrl + c, q | Quit		 |


## fzf, rofi

keyb also supports exporting of a formatted output for use with other programs
like fzf:
```bash
$ keyb -e output
$ cat output | fzf
```

For rofi, the output must be stripped of ansi formatting first:
```bash
$ keyb -e output -s
$ cat output | rofi -dmenu
```

## Configuration
The default config path is `$HOME/.config/keyb`.


## Inspiration
- [awesomeWM](https://github.com/awesomeWM/awesome)
- [showkeys](https://github.com/adamharmansky/showkeys)

## License
[MIT](LICENSE)
