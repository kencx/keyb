# keyb

<p align="center">
	<img width="660" src="https://github.com/kencx/keyb/blob/master/assets/compressed.gif?raw=true">
</p>

<p align="center">Create and view your own custom hotkey cheatsheet in the terminal</p>

<p align="center">
	<img src="https://goreportcard.com/badge/github.com/kencx/keyb">
	<img src="https://github.com/kencx/keyb/actions/workflows/test.yml/badge.svg?branch=master">
</p>

### Features
- Lightweight and quick
- Fully customizable
- Fuzzy filtering
- Vim key bindings
- Export to stdout for fzf, rofi support

### Non-Features
keyb does **not** support:

- Auto detection of hotkeys
- Command selection

## Motivation

I had trouble remembering the various hotkeys that I sometimes use. It got
annoying to look them up so I resorted to writing them down on a paper
cheatsheet. Then, I thought: maybe there's a tool that does this better. I
didn't find one I liked so I built keyb.

With keyb, I can list:

- Hotkeys that I occasionally forget or am new to
- Custom key combinations that I defined for my own workflow
- Short commands that I sometimes use

It is best used as a popup cheatsheet.

## Install
keyb supports Linux, MacOS and Windows.

### Compiled Binary
Download a compiled binary from the [releases](https://github.com/kencx/keyb/releases) page.

### Install with Go

```bash
$ go install github.com/kencx/keyb@latest
```

### Build from source

```bash
$ git clone https://github.com/kencx/keyb
$ cd keyb
$ make build
```

## Usage

```text
usage: keyb [options] <command>

Options:
  -p, --print     Print to stdout
  -e, --export    Export to file
  -k, --key       Key bindings at custom path
  -c, --config    Config file at custom path
  -v, --version   Version info
  -h, --help      help for keyb

Commands:
  a, add          Add keybind to keyb file
```

### Search

- Enter search mode with `/` to perform fuzzy filtering on all rows
- Exit search mode again with `Esc`
- `Alt + d` clears the current filter

To perform filtering on section headings only, prefix the
search with `h:`. This will return all matching section headings with their
respective rows.

### Printing

keyb supports printing to stdout for use with other tools:

```bash
$ keyb -p | fzf
$ keyb -p | rofi -dmenu
```

### keyb File

keyb requires a `yaml` file with a list of hotkeys to work. A default file is
generated in your system's config directory if no other file is specified.

Hotkeys are classified into sections with a name and (optional) prefix field.
When displayed, sections are sorted by alphabetical order while the keys
themselves are arranged in their defined order.

```yaml
- name: bspwm
  keybinds:
    - name: terminal
      key: Super + Return
```

The prefix is a key combination that will be prepended to every hotkey in the
section. A key can choose to opt out by including a `ignore_prefix: true` field.
Prefixes are useful for applications with a common leading hotkey like tmux.

```yaml
- name: tmux
  prefix: ctrl + b
  keybinds:
    - name: Create new window
      key: c
    - name: Prev, next window
      key: Shift + {←, →}
      ignore_prefix: true
```

Refer to the `examples` for more examples.

>Multiline fields are not supported!

### Quick Add

```text
usage: keyb [-k file] add [app; name; key]

Options:
  -b, --binding  Key binding
  -p, --prefix   Ignore prefix
```

You can quick add bindings from the command line to a specified file. If `-k
file` is given and exists, the new keybind will be appended to the file.
Otherwise, `keyb_path` defined in `config.yml` will be used.

```bash
$ keyb add -b "kitty; open terminal; super + enter"
```

When adding a new keybind, the app name, keybind name and keybind must be
specified. It is separated by `;` and wrapped in quotes (to prevent parsing errors).

## Configuration

keyb can be customized with a config file at the default OS config
directory (i.e. `$XDG_CONFIG_HOME/keyb/config.yml`). If no such file exists, the
default configuration will be used.

See [config](examples/config/README.md) for all configuration options.

## Roadmap

- [x] Ability to customize keyb hotkeys
- [x] `a, add` subcommand to quickly add a single hotkey entry from the CLI
- [ ] Export to additional file formats (`json, toml, conf/ini` etc.)
- [ ] Support multiple keyb files or directories

## Contributing

keyb requires Go 1.21. Bug reports, feature requests and PRs are very welcome.

## Similar Tools

- [showkeys](https://github.com/adamharmansky/showkeys) offers a keybinding popup similar to awesomewm
- [cheat](https://github.com/cheat/cheat) is a CLI alternative to view cheatsheets for
  commands and hotkeys for just about any topic
- Refer to [shortcut-pages](https://github.com/mt-empty/shortcut-pages), [cheat/cheatsheets](https://github.com/cheat/cheatsheets) for more cheatsheets

## License
[MIT](LICENSE)
