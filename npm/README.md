# @suryastra/openpilot

Node.js shim for the OpenPilot CLI. Installs the native binary for your OS/arch from GitHub Releases.

## Install

For Global Install

```bash
npm install -g @suryastra/openpilot
```

### or

For One time use

```bash
npx @suryastra/openpilot
```

Environment variables:

- `OPENPILOT_VERSION` pin a version (default latest)
- `OPENPILOT_REPO` override GitHub repo (default `suryastra/openpilot`)
- `OPENPILOT_DEBUG` set to `1` for verbose postinstall logging

Set `OPENPILOT_VERSION` to pin a version:

```bash
OPENPILOT_VERSION=0.1.0 npx @suryastra/openpilot --version
```

## How it works

- Postinstall script downloads the archive from GitHub Releases
- Extracts and places the `openpilot` binary into the package bin dir
- `openpilot` JavaScript shim spawns the binary

## Troubleshooting

- Ensure Node 18+
- Corporate proxies may block the download: set `HTTPS_PROXY`
- Windows extraction requires `tar` in PATH (included in recent Git and Windows 10+). If missing, manually download release.
- If download fails, enable debug: `OPENPILOT_DEBUG=1 npm i @suryastra/openpilot` to see attempted asset names.
- The installer tries several archive name patterns (with/without version, .tar.gz/.tgz). If none match, manually fetch from Releases and place binary as `node_modules/@suryastra/openpilot/bin/openpilot`.

## License

See project root LICENSE.
