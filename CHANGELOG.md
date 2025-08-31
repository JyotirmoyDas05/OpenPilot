## <small>1.1.8 (2025-08-31)</small>

* fix(npm): handle HTTP redirects in postinstall script ([984016e](https://github.com/JyotirmoyDas05/OpenPilot/commit/984016e))

## <small>1.1.7 (2025-08-31)</small>

* fix: prevent dirty working directory by using --no-save for npm installs ([77cc848](https://github.com/JyotirmoyDas05/OpenPilot/commit/77cc848))

## <small>1.1.6 (2025-08-31)</small>

* Merge branch 'main' of https://github.com/JyotirmoyDas05/OpenPilot ([ab3334e](https://github.com/JyotirmoyDas05/OpenPilot/commit/ab3334e))
* fix: commit lockfile changes before GoReleaser to prevent dirty working directory ([ff31801](https://github.com/JyotirmoyDas05/OpenPilot/commit/ff31801))

## <small>1.1.5 (2025-08-31)</small>

* Merge branch 'main' of https://github.com/JyotirmoyDas05/OpenPilot ([437f3bf](https://github.com/JyotirmoyDas05/OpenPilot/commit/437f3bf))
* fix: add build artifacts to gitignore and clean working directory in workflow ([fc1171f](https://github.com/JyotirmoyDas05/OpenPilot/commit/fc1171f))

## <small>1.1.4 (2025-08-31)</small>

* fix: update GoReleaser comment for better asset description ([ef466d3](https://github.com/JyotirmoyDas05/OpenPilot/commit/ef466d3))
* chore: remove lint workflow configuration ([6b98226](https://github.com/JyotirmoyDas05/OpenPilot/commit/6b98226))

## <small>1.1.3 (2025-08-31)</small>

* fix: remove invalid --skip-validate flag from goreleaser command ([69e46cd](https://github.com/JyotirmoyDas05/OpenPilot/commit/69e46cd))

## <small>1.1.2 (2025-08-31)</small>

* fix: integrate GoReleaser with semantic-release using exec plugin for proper release automation ([9968cb9](https://github.com/JyotirmoyDas05/OpenPilot/commit/9968cb9))

## <small>1.1.1 (2025-08-31)</small>

* fix: prevent GitHub release asset conflicts by cleaning existing assets and specifying upload patter ([f666af0](https://github.com/JyotirmoyDas05/OpenPilot/commit/f666af0))

## 1.1.0 (2025-08-31)

* feat: enhance npm README with improved documentation, platform support details, and better troublesh ([35588e9](https://github.com/JyotirmoyDas05/OpenPilot/commit/35588e9))

# Changelog

All notable changes will be documented automatically by semantic-release.

## <small>1.0.5 (2025-08-31)</small>

- fix: clean up duplicate changelog entries and improve structure ([c28cd5a](https://github.com/JyotirmoyDas05/OpenPilot/commit/c28cd5a))
- fix: improve changelog formatting with better categorization and structure ([ffcd097](https://github.com/JyotirmoyDas05/OpenPilot/commit/ffcd097))

## <small>1.0.4 (2025-08-31)</small>

### Bug Fixes

- Add node_modules to .gitignore and remove from tracking ([e788498](https://github.com/JyotirmoyDas05/OpenPilot/commit/e788498))
- Correct indentation in GoReleaser build ldflags ([a201223](https://github.com/JyotirmoyDas05/OpenPilot/commit/a201223))
- Update default repository in postinstall script to JyotirmoyDas05/openpilot ([de56df3](https://github.com/JyotirmoyDas05/OpenPilot/commit/de56df3))
- Update goreleaser config for new project structure ([318c532](https://github.com/JyotirmoyDas05/OpenPilot/commit/318c532))
- Update GoReleaser config for proper asset naming and file structure ([c7c72fe](https://github.com/JyotirmoyDas05/OpenPilot/commit/c7c72fe))
- Update GoReleaser to use x86_64 instead of amd64 for asset naming ([262e71b](https://github.com/JyotirmoyDas05/OpenPilot/commit/262e71b))
- Update semantic-release workflow to install GoReleaser and build binaries ([a73f08f](https://github.com/JyotirmoyDas05/OpenPilot/commit/a73f08f))
- Update workflow to handle both files and directories in asset duplication ([a106b7f](https://github.com/JyotirmoyDas05/OpenPilot/commit/a106b7f))
- Use GoReleaser snapshot mode to avoid tag conflicts with semantic-release ([25920de](https://github.com/JyotirmoyDas05/OpenPilot/commit/25920de))

### Documentation

- Add comment for asset naming fix ([b1e4992](https://github.com/JyotirmoyDas05/OpenPilot/commit/b1e4992))
- Update comment for release readiness ([bed165b](https://github.com/JyotirmoyDas05/OpenPilot/commit/bed165b))

### Code Refactoring

- Simplify GoReleaser configuration by removing unused fields and hooks ([9d26c4f](https://github.com/JyotirmoyDas05/OpenPilot/commit/9d26c4f))

## <small>1.0.3 (2025-08-31)</small>

- fix: enhance asset name patterns in postinstall script for better compatibility and delete redundant ([c195273](https://github.com/JyotirmoyDas05/OpenPilot/commit/c195273))

## <small>1.0.2 (2025-08-31)</small>

- fix: enhance asset name patterns and improve archive extraction handling in postinstall script ([d2642ca](https://github.com/JyotirmoyDas05/OpenPilot/commit/d2642ca))

## <small>1.0.1 (2025-08-31)</small>

- fix: enhance asset handling in release workflow and improve error fallback in postinstall script ([f22c69c](https://github.com/JyotirmoyDas05/OpenPilot/commit/f22c69c))

## 1.0.0 (2025-08-31)

- fix: 2nd fix for release ([72ee7e4](https://github.com/JyotirmoyDas05/OpenPilot/commit/72ee7e4))
- fix: 3rd fix for semantic release ([2df8c32](https://github.com/JyotirmoyDas05/OpenPilot/commit/2df8c32))
- fix: 4th fix for semantic release involving Building GO Binaries ([ec294b3](https://github.com/JyotirmoyDas05/OpenPilot/commit/ec294b3))
- fix: adding package-lock.json ([b4e62c6](https://github.com/JyotirmoyDas05/OpenPilot/commit/b4e62c6))
- feat: Add initial configuration schema and provider settings ([2fa02ed](https://github.com/JyotirmoyDas05/OpenPilot/commit/2fa02ed))
- feat: initial semantic-release setup ([f231d5d](https://github.com/JyotirmoyDas05/OpenPilot/commit/f231d5d))
- docs: Update model listing reference to Awesome-LLM and change badge image ([fade335](https://github.com/JyotirmoyDas05/OpenPilot/commit/fade335))
