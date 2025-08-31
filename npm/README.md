# @suryastra/openpilot

OpenPilot CLI wrapper for Node.js that automatically downloads the native binary for your platform from GitHub Releases.

## Quick Start

### Global Installation

Install globally to use the `openpilot` command anywhere:

```bash
npm install -g @suryastra/openpilot
openpilot --help
```

### One-time Usage

Run without installing globally:

```bash
npx @suryastra/openpilot --help
```

### Project Dependency

Add to your Node.js project:

```bash
npm install @suryastra/openpilot
```

## Configuration

Environment variables for customization:

| Variable            | Description                | Default                    |
| ------------------- | -------------------------- | -------------------------- |
| `OPENPILOT_VERSION` | Pin a specific version     | `latest`                   |
| `OPENPILOT_REPO`    | Override GitHub repository | `JyotirmoyDas05/openpilot` |
| `OPENPILOT_DEBUG`   | Enable verbose logging     | `0`                        |

### Version Pinning Example

```bash
OPENPILOT_VERSION=1.0.6 npx @suryastra/openpilot --version
```

## How It Works

1. **Postinstall**: Downloads the appropriate binary archive from GitHub Releases
2. **Extraction**: Unpacks and places the `openpilot` binary in the package bin directory
3. **Execution**: JavaScript shim spawns the native binary with your arguments

## Platform Support

- **Windows**: x86_64, ARM64 (requires Windows 10+ with tar support)
- **macOS**: x86_64 (Intel), ARM64 (Apple Silicon)
- **Linux**: x86_64, ARM64

## Troubleshooting

### Common Issues

- **Node.js Version**: Requires Node.js 18 or higher
- **Corporate Proxies**: Set `HTTPS_PROXY` environment variable if downloads are blocked
- **Windows Extraction**: Requires `tar` command (included in Git for Windows and Windows 10+)

### Debug Mode

Enable detailed logging to troubleshoot download issues:

```bash
OPENPILOT_DEBUG=1 npm install @suryastra/openpilot
```

### Manual Installation

If automatic download fails, manually download the binary:

1. Visit [GitHub Releases](https://github.com/JyotirmoyDas05/openpilot/releases)
2. Download the appropriate archive for your platform
3. Extract and place the binary as:
   ```
   node_modules/@suryastra/openpilot/bin/openpilot[.exe]
   ```

## License

See [LICENSE.md](../../LICENSE.md) in the project root.
