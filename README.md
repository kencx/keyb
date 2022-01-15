# keyb

`keyb` is a custom configured, global hotkey cheatsheet, inspired by
awesomeWM's very own cheatsheet.

With `keyb`, you can list any hotkeys from any program, all in one place!

## Usage
List your custom keybindings in a config file. `keyb` supports yaml & json
config files.

```yaml
tmux:
  prefix: ctrl + a
  keybinds:
    - command: split vertical
      key: "|"
    - command: split horizontal
      key: "-"
bspwm:
  keybinds:
    - command: open configs
	  key: super + z
```

## Inspiration
- [awesomeWM](https://github.com/awesomeWM/awesome)
- [showkeys](https://github.com/adamharmansky/showkeys)
