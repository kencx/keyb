# Configuration

keyb will accept the following config in decreasing priority:

- `-c FILE` flag
- The default config path `$XDG_CONFIG_HOME/keyb/config.yml` (see note)
- The default configuration (see [default.yml](default.yml))

Note: If `$XDG_CONFIG_HOME` is set, it will be prioritized and used in Unix and Darwin
systems. Otherwise, keyb will fall back to the default OS config directory
defined as such:

- Unix: `$XDG_CONFIG_HOME/keyb/`,
- MacOS/Darwin: `$HOME/Library/Application Support/keyb/`,
- Windows: `%Appdata%\keyb\`

## Options

| Option        | Default                  | Description |
| ------------- | ------------------------ | ----------- |
| `keyb_path`   | OS-dependent (see above) | keyb file path |
| `debug`       | `false`                  | Debug mode |
| `reverse`     | `false`                  | Swap the name and key columns |
| `mouse`       | `true`                   | Mouse enabled |
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

### Color
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

### Hotkeys
Multiple keys may be set for a single binding, separated by commas.

| Hotkey                  | Default                    | Description      |
| ----------------------- | -------------------------- | ---------------- |
| `up`, `down`            | <kbd>j, k / Up, Down</kbd> | Move cursor      |
| `up_focus`, `down_focus`| <kbd>Ctrl + j, ctrl + k </kbd> | Move cursor in search mode |
| `half_up, half_down`    | <kbd>Ctrl + u, d</kbd>     | Move half window (also works in search mode) |
| `full_up, full_down`    | <kbd>Ctrl + b, f</kbd>     | Move full window (also works in search mode) |
| `top, middle, bottom`   | <kbd>H, M, L</kbd>         | Go to top, middle, bottom of screen |
| `first_line, last_line` | <kbd>g, G</kbd>            | Go to first, last line |
| `search`                | <kbd>/</kbd>               | Enter search mode      |
| `clear_search`          | <kbd>Alt + d</kbd>         | Clear current search (remains in search mode) |
| `normal`                | <kbd>Esc</kbd>             | Exit search mode |
| `quit`                  | <kbd>Ctrl + c, q</kbd>     | Quit		      |

These hotkeys configure the cursor behaviour in the search bar only:

| Hotkey                  | Default                     | Description      |
| ----------------------- | --------------------------- | ---------------- |
| `cursor_word_forward`     | <kbd>alt+right, alt+f</kbd> | Move forward by word |
| `cursor_word_backward`    | <kbd>alt+left, alt+b</kbd>  | Move backward by word |
| `cursor_delete_word_backward` | <kbd>alt+backspace</kbd> | Delete word backward |
| `cursor_delete_word_forward`  | <kbd>alt+delete</kbd>   | Delete word forward |
| `cursor_delete_after_cursor`  | <kbd>alt+k</kbd>        | Delete after cursor |
| `cursor_delete_before_cursor` | <kbd>alt+u</kbd>        | Delete before cursor |
| `cursor_line_start`       | <kbd>home, ctrl+a</kbd>     | Move cursor to start |
| `cursor_line_end`         | <kbd>end, ctrl+e</kbd>      | Move cursor to end |
| `cursor_paste`            | <kbd>ctrl+v</kbd>           | Paste into search bar|
