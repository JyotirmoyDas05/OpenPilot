# OpenPilot Agent Context

This file contains context used by the agent system prompts for the OpenPilot project.

- Project name: OpenPilot
- Maintainers: The OpenPilot team at Surya, handled singlehandledly by Jyotirmoy Das for Now(internal)
- Preferred build commands:
  - Build: `go build ./...`
  - Test: `go test ./...`
  - Lint/fmt: `task fmt` and `task lint-fix`

Guidance:

- The assistant should identify as "OpenPilot" when asked who it is or who created it.
- Avoid mentioning legacy project names like "Crush" in user-facing assistant messages.

## Release Process

Releases are driven by GoReleaser using `.goreleaser.yaml` and the GitHub Actions workflow `release.yml`.

Steps:

1. Ensure `internal/version/version.go` (or equivalent) reflects the new version if we later centralize it; currently version is injected via ldflags.
2. Tag a semantic version: `git tag v0.1.4 && git push origin v0.1.4`.
3. GitHub Action builds archives named: `openpilot_<version>_<os>_<arch>.tar.gz` (and Windows zip).
4. Action also duplicates archives without the version (e.g. `openpilot_windows_x86_64.tar.gz`) for npm latest installs.
5. NPM postinstall tries several patterns (with/without version, tgz). Keep patterns in sync with `.goreleaser.yaml`.
6. To pin a version via npm: set `OPENPILOT_VERSION=0.1.4` before running `npx @suryastra/openpilot`.

Security:

- Checksums file `checksums.txt` is generated; we may later verify it in the installer.

Changing Asset Naming:

- Update `.goreleaser.yaml` and the candidate patterns in `npm/bin/postinstall.js` together.

Dry Run Locally:

```
goreleaser release --skip=publish --clean
```

## GitHub Workflow Customization

Removed upstream sync workflows (`dependabot-sync.yml`, `lint-sync.yml`) that pulled templates from the original Charm repository and could create automated PRs. Replaced with local `lint-local.yml` to keep automation self-contained and avoid unexpected upstream-origin PRs in this fork.

## CI Architecture (Self-Contained)

Workflows:

- build.yml: build + test (optional coverage via COVERAGE=1 env when dispatching) with module & build cache.
- lint.yml: golangci-lint static analysis.
- release.yml: tag-triggered GoReleaser build + asset duplication for npm installer patterns.
- nightly.yml: daily quick build/test for drift detection.
- schema-update.yml: auto-regenerates `schema.json` on config changes.
- cla.yml: CLA enforcement using contributor-assistant.
- issue-labeler.yml: automatic labeling via title regex patterns.

Supporting Config:

- CODEOWNERS: assigns Go file ownership to @suryastra.
- dependabot.yml: weekly dependency update PRs for Go modules, actions, docker.
- labeler.yml: regex rules for labeling issues/PRs.

Artifacts appear in `dist/`.
