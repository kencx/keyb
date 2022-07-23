# keyb

<p align="center">
	<img width="660" src="https://github.com/kencx/keyb/blob/master/assets/show.gif?raw=true">
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
usage: keyb [options] [file]

Flags:
  -p, --print           Print to stdout
  -e, --export [file]   Export to file
  -k, --key [file]      Pass custom hotkeys file
  -c, --config [file]   Pass custom config file
  -v, --version         Version info
  -h, --help            help for keyb
```

### Printing

keyb supports printing to stdout for use with other tools:

```bash
$ keyb -p | fzf
$ keyb -p | rofi -dmenu
```

### Search

- Enter search mode with `/` to perform fuzzy filtering on all rows
- Exit search mode again with `Esc`
- `Alt + d` clears the current filter

To perform filtering on section headings only, prefix the
search with `h:`. This will return all matching section headings with their
respective rows.

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

>Multiline fields are not supported at the moment.

### Configuration
keyb is customized with a `config.yml` file that is automatically generated in
the system's config directory.

By default, the following options are included.

| Option        | Default                  | Description |
| ------------- | ------------------------ | ----------- |
| `keyb_path`   | OS-dependent (see below) | keyb file path |
| `debug`       | `false`                  | Debug mode |
| `mouse`       | `true`                   | Mouse enabled |
| `reverse`     | `false`                  | Swap the name and key columns |
| `search_mode` | `false`                  | Start in search mode |
| `sort_keys`   | `false`                  | Sort keys alphabetically |
| `title`       | `""`                     | Title text |
| `prompt`      | `"keys > "`              | Search bar prompt text |
| `prompt_location` | `"top"`                | Location of search bar: `top, bottom` |
| `placeholder` | `"..."`                  | Search bar placeholder text |
| `prefix_sep`  | `";"`                    | Separator symbol between prefix and key |
| `sep_width`   | `4`                      | Separation width between columns |
| `margin`      | `0`                      | Space between window and border |
| `padding`     | `1`                      | Space between border and text |
| `border`      | `"hidden"`               | Border style: `normal, rounded, double, thick, hidden`|

>OS-specific keyb directory path:
>
>- Unix: `$XDG_CONFIG_HOME/keyb/keyb.yml`,
>- MacOS or Darwin: `$HOME/Library/Application Support/keyb/keyb.yml`,
>- Windows: `%Appdata%\keyb\keyb.yml`

#### Color
Both ANSI and hex color codes are supported.

| Color Option     | Default    | Description |
| ---------------- | ---------- | ----------- |
| `prompt`         | -          | Prompt text color |
| `cursor_fg`      | -          | Cursor foreground |
| `cursor_bg`      | -          | Cursor background |
| `filter_fg`      | `"#FFA066"`| Filter matching text foreground |
| `filter_bg`      | -          | Filter matching text background |
| `counter_fg`     | -          | Counter foreground |
| `counter_bg`     | -          | Counter background |
| `placeholder_fg` | -          | Placeholder foreground |
| `placeholder_bg` | -          | Placeholder background |
| `border_color`   | -          | Border color |

#### Hotkeys
Multiple keys may be set for a single binding, separated by commas.

| Hotkey                  | Default                    | Description      |
| ----------------------- | -------------------------- | ---------------- |
| `up`, `down`            | <kbd>j, k / Up, Down</kbd> | Move cursor      |
| `half_up, half_down`    | <kbd>Ctrl + u, d</kbd>     | Move half window |
| `full_up, full_down`    | <kbd>Ctrl + b, f</kbd>     | Move full window |
| `top, middle, bottom`   | <kbd>H, M, L</kbd>         | Go to top, middle, bottom of screen |
| `first_line, last_line` | <kbd>g, G</kbd>            | Go to first, last line |
| `search`                | <kbd>/</kbd>               | Enter search mode      |
| `clear_search`          | <kbd>Alt + d</kbd>         | Clear current search (remains in search mode) |
| `normal`                | <kbd>Esc</kbd>             | Exit search mode |
| `quit`                  | <kbd>Ctrl + c, q</kbd>     | Quit		      |

## Roadmap

- [x] Ability to customize keyb hotkeys
- [ ] Export to additional file formats (`json, toml, conf/ini` etc.)
- [ ] `-a, --add` flag to quickly add a single hotkey entry from the CLI
- [ ] Automatic parsing from online cheatsheet repos (eg. `cheat/cheatsheets`)

## Contributing

keyb requires Go 1.18.

Bug reports, feature requests and PRs are very welcome.

More examples for default keymaps in `examples/` are also welcome.

## Screenshots

<p align="center">
	<img width="660" src="https://github.com/kencx/keyb/blob/master/assets/keyb.jpg?raw=true">
</p>

## Similar Tools

- [showkeys](https://github.com/adamharmansky/showkeys) offers a keybinding popup similar to awesomewm
- [cheat](https://github.com/cheat/cheat) is a CLI alternative to view cheatsheets for
  commands and hotkeys for just about any topic
- Refer to [shortcut-pages](https://github.com/mt-empty/shortcut-pages), [cheat/cheatsheets](https://github.com/cheat/cheatsheets) for more cheatsheets

## License
[MIT](LICENSE)
