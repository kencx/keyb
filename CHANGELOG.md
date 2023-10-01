# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

[v0.2.0]: https://github.com/kencx/keyb/compare/v0.1.0...v0.2.0
[v0.3.0]: https://github.com/kencx/keyb/compare/v0.2.0...v0.3.0
[v0.4.0]: https://github.com/kencx/keyb/compare/v0.3.0...v0.4.0
[v0.4.1]: https://github.com/kencx/keyb/compare/v0.4.0...v0.4.1
[v0.5.0]: https://github.com/kencx/keyb/compare/v0.4.1...v0.5.0
[v0.6.0]: https://github.com/kencx/keyb/compare/v0.5.0...v0.6.0
[v0.7.0]: https://github.com/kencx/keyb/compare/v0.6.0...v0.7.0

## [v0.7.0]

### Fixed
- Exiting now clears the output

### Chore
- Replace deprecated `--rm-dist` with `--clean` flag in Goreleaser action

## [v0.6.0]

### Changed
- Bump goreleaser/goreleaser-action from 4 to 5 [#25](https://github.com/kencx/keyb/pull/25)
- Simplify config package
- Bump actions/checkout from 3 to 4 [#24](https://github.com/kencx/keyb/pull/24)
- Update Go to v1.21
- Update bubbletea to v0.24.2

### Fixed
- Fix environment variables not being expanded in the config's `settings.keyb_path`.

## [v0.5.0]
### Added
- Add ability to customize search bar cursor bindings.

### Changed
- Default configuration file is no longer automatically generated. `keyb` will
  default to reading any config file at `$XDG_CONFIG_HOME/keyb/config.yml`. If
  not present, the default configuration is used.

## [v0.4.1]
### Fixed
- Fix inability to delete characters in search bar [1385049](https://github.com/kencx/keyb/commit/138504964bad8f8827c5f5e9c1572298d4d5e102)

## [v0.4.0]
### BREAKING
- Use `XDG_CONFIG_HOME` environment variable in macOS if set [#18](https://github.com/kencx/keyb/pull/18)

### Added
- Add ability to move cursor in search mode
  [55dd7ad](https://github.com/kencx/keyb/commit/55dd7adead29316d3952e7c19bb5b15546394668)
- Add ability to customize cursor movement keys in search mode
  [55dd7ad](https://github.com/kencx/keyb/commit/55dd7adead29316d3952e7c19bb5b15546394668)

### Changed
- Update dependencies [#19](https://github.com/kencx/keyb/pull/19) and
  [3cba5b8](https://github.com/kencx/keyb/commit/3cba5b801acd617e9d1c37734582f3f15d2ec41b)

### Fixed
- Fix search bar cursor not blinking when focused [#20](https://github.com/kencx/keyb/pull/20)

## [v0.3.0]
### Added
- Add "add" subcommand to quick add keybinds [50b66d7](https://github.com/kencx/keyb/commit/50b66d7a78c4a08a9cb5ad5bd02d909b7b27ae53), [b9167474](https://github.com/kencx/keyb/commit/b9167474c9c5d12ed8ea0ca9630489fa7266bebe)

## [v0.2.0]
### Added
- Add counter styling [#2](https://github.com/kencx/keyb/pull/2)
- Add placeholder styling [#4](https://github.com/kencx/keyb/pull/4)
- Add ability to customize keyb key bindings [43ae9b8](https://github.com/kencx/keyb/commit/43ae9b83fbf5cae367ab74614fa42fce79817165)
- Add `sort_keys` option to sort keys alphabetically [#7](https://github.com/kencx/keyb/pull/7)
- Add ability to customize prompt location [019f6ca](https://github.com/kencx/keyb/commit/019f6cad03ada6507e6585e4f4403826dcd23212)
- Add `search_mode` option to start keyb in search mode [d6c53e1](https://github.com/kencx/keyb/commit/d6c53e1b908f05f6c0f7836068b4b6bbe1e8a451)

<!-- ### Changed -->
<!---->
<!-- ### Removed -->
<!---->
<!-- ### Fixed -->
