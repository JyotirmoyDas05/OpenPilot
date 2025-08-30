# OpenPilot

<p align="center">
    <a href="https://github.com/surya/openpilot/releases"><img src="https://img.shields.io/github/release/surya/openpilot" alt="Latest Release"></a>
    <a href="https://github.com/surya/openpilot/actions"><img src="https://github.com/surya/openpilot/workflows/build/badge.svg" alt="Build Status"></a>
</p>

<p align="center">Your new coding IDE, now available in your favourite terminal.<br />Your tools, your code, and your workflows, wired into your LLM of choice.</p>

## Features

- **Multi-Model:** choose from a wide range of LLMs or add your own via OpenAI- or Anthropic-compatible APIs
- **Flexible:** switch LLMs mid-session while preserving context
- **Session-Based:** maintain multiple work sessions and contexts per project
- **LSP-Enhanced:** OpenPilot uses LSPs for additional context, just like you do
- **Extensible:** add capabilities via MCPs (`http`, `stdio`, and `sse`)
- **Works Everywhere:** first-class support in every terminal on macOS, Linux, Windows (PowerShell and WSL), FreeBSD, OpenBSD, and NetBSD

## Installation

Use a package manager:

### Homebrew

```bash
brew install surya/tap/openpilot
```

### NPM

Install OpenPilot globally:

```bash
npm install -g @suryastra/openpilot
```

Run OpenPilot one-off (no install):

```bash
npx @suryastra/openpilot --help
```

Pin a specific released binary:

```bash
OPENPILOT_VERSION=0.1.0 npx @suryastra/openpilot --version
```

## Arch Linux (btw)

```bash
yay -S openpilot-bin
```

## Nix

```bash
nix run github:numtide/nix-ai-tools#openpilot
```

Windows users:

```bash
# Winget
winget install surya.openpilot

# Scoop
scoop bucket add surya https://github.com/surya/scoop-bucket.git
scoop install openpilot
```

<details>
<summary><strong>Nix (NUR)</strong></summary>

OpenPilot is available via [NUR](https://github.com/nix-community/NUR) in `nur.repos.surya.openpilot`.

You can also try out OpenPilot via `nix-shell`:

```bash
# Add the NUR channel.
nix-channel --add https://github.com/nix-community/NUR/archive/main.tar.gz nur
nix-channel --update

# Get OpenPilot in a Nix shell.
nix-shell -p '(import <nur> { pkgs = import <nixpkgs> {}; }).repos.surya.openpilot'
```

</details>

<details>
<summary><strong>Debian/Ubuntu</strong></summary>

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.surya.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/surya.gpg
echo "deb [signed-by=/etc/apt/keyrings/surya.gpg] https://repo.surya.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/surya.list
sudo apt update && sudo apt install openpilot
```

</details>

<details>
<summary><strong>Fedora/RHEL</strong></summary>

```bash
echo '[surya]
name=Surya
baseurl=https://repo.surya.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.surya.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/surya.repo
sudo yum install openpilot
```

</details>

Or, download it:

- [Packages][releases] are available in Debian and RPM formats
- [Binaries][releases] are available for Linux, macOS, Windows, FreeBSD, OpenBSD, and NetBSD

[releases]: https://github.com/surya/openpilot/releases

Or just install it with Go:

```
go install github.com/JyotirmoyDas05/openpilot@latest
```

> [!WARNING]
> Productivity may increase when using OpenPilot and you may find yourself nerd
> sniped when first using the application. If the symptoms persist, join the
> [Discord][discord] and nerd snipe the rest of us.

## Getting Started

The quickest way to get started is to grab an API key for your preferred
provider such as Anthropic, OpenAI, Groq, or OpenRouter and just start
OpenPilot. You'll be prompted to enter your API key.

That said, you can also set environment variables for preferred providers.

| Environment Variable       | Provider                                           |
| -------------------------- | -------------------------------------------------- |
| `ANTHROPIC_API_KEY`        | Anthropic                                          |
| `OPENAI_API_KEY`           | OpenAI                                             |
| `OPENROUTER_API_KEY`       | OpenRouter                                         |
| `GEMINI_API_KEY`           | Google Gemini                                      |
| `VERTEXAI_PROJECT`         | Google Cloud VertexAI (Gemini)                     |
| `VERTEXAI_LOCATION`        | Google Cloud VertexAI (Gemini)                     |
| `GROQ_API_KEY`             | Groq                                               |
| `AWS_ACCESS_KEY_ID`        | AWS Bedrock (Claude)                               |
| `AWS_SECRET_ACCESS_KEY`    | AWS Bedrock (Claude)                               |
| `AWS_REGION`               | AWS Bedrock (Claude)                               |
| `AZURE_OPENAI_ENDPOINT`    | Azure OpenAI models                                |
| `AZURE_OPENAI_API_KEY`     | Azure OpenAI models (optional when using Entra ID) |
| `AZURE_OPENAI_API_VERSION` | Azure OpenAI models                                |

### By the Way

Is there a provider you’d like to see in OpenPilot? Is there an existing model that needs an update?

OpenPilot’s default model listing is managed in [Catwalk](https://github.com/surya/catwalk), an community-supported, open source repository of OpenPilot-compatible models, and you’re welcome to contribute.

<a href="https://github.com/surya/catwalk"><img width="174" height="174" alt="Catwalk Badge" src="https://github.com/user-attachments/assets/95b49515-fe82-4409-b10d-5beb0873787d" /></a>

## Configuration

OpenPilot runs great with no configuration. That said, if you do need or want to
customize OpenPilot, configuration can be added either local to the project itself,
or globally, with the following priority:

1. `.openpilot.json`
2. `openpilot.json`
3. `$HOME/.config/openpilot/openpilot.json` (Windows: `%USERPROFILE%\AppData\Local\openpilot\openpilot.json`)

Configuration itself is stored as a JSON object:

```json
{
  "this-setting": { "this": "that" },
  "that-setting": ["ceci", "cela"]
}
```

As an additional note, OpenPilot also stores ephemeral data, such as application state, in one additional location:

```bash
# Unix
$HOME/.local/share/openpilot/openpilot.json

# Windows
%LOCALAPPDATA%\openpilot\openpilot.json
```

### LSPs

OpenPilot can use LSPs for additional context to help inform its decisions, just
like you would. LSPs can be added manually like so:

```json
{
  "$schema": "https://surya.land/openpilot.json",
  "lsp": {
    "go": {
      "command": "gopls",
      "env": {
        "GOTOOLCHAIN": "go1.24.5"
      }
    },
    "typescript": {
      "command": "typescript-language-server",
      "args": ["--stdio"]
    },
    "nix": {
      "command": "nil"
    }
  }
}
```

### MCPs

OpenPilot also supports Model Context Protocol (MCP) servers through three
transport types: `stdio` for command-line servers, `http` for HTTP endpoints,
and `sse` for Server-Sent Events. Environment variable expansion is supported
using `$(echo $VAR)` syntax.

```json
{
  "$schema": "https://surya.land/openpilot.json",
  "mcp": {
    "filesystem": {
      "type": "stdio",
      "command": "node",
      "args": ["/path/to/mcp-server.js"],
      "env": {
        "NODE_ENV": "production"
      }
    },
    "github": {
      "type": "http",
      "url": "https://example.com/mcp/",
      "headers": {
        "Authorization": "$(echo Bearer $EXAMPLE_MCP_TOKEN)"
      }
    },
    "streaming-service": {
      "type": "sse",
      "url": "https://example.com/mcp/sse",
      "headers": {
        "API-Key": "$(echo $API_KEY)"
      }
    }
  }
}
```

### Ignoring Files

OpenPilot respects `.gitignore` files by default, but you can also create a
`.openpilotignore` file to specify additional files and directories that OpenPilot
should ignore. This is useful for excluding files that you want in version
control but don't want OpenPilot to consider when providing context.

The `.openpilotignore` file uses the same syntax as `.gitignore` and can be placed
in the root of your project or in subdirectories.

### Allowing Tools

By default, OpenPilot will ask you for permission before running tool calls. If
you'd like, you can allow tools to be executed without prompting you for
permissions. Use this with care.

```json
{
  "$schema": "https://surya.land/openpilot.json",
  "permissions": {
    "allowed_tools": [
      "view",
      "ls",
      "grep",
      "edit",
      "mcp_context7_get-library-doc"
    ]
  }
}
```

You can also skip all permission prompts entirely by running OpenPilot with the
`--yolo` flag. Be very, very careful with this feature.

### Local Models

Local models can also be configured via OpenAI-compatible API. Here are two common examples:

#### Ollama

```json
{
  "providers": {
    "ollama": {
      "name": "Ollama",
      "base_url": "http://localhost:11434/v1/",
      "type": "openai",
      "models": [
        {
          "name": "Qwen 3 30B",
          "id": "qwen3:30b",
          "context_window": 256000,
          "default_max_tokens": 20000
        }
      ]
    }
  }
}
```

#### LM Studio

```json
{
  "providers": {
    "lmstudio": {
      "name": "LM Studio",
      "base_url": "http://localhost:1234/v1/",
      "type": "openai",
      "models": [
        {
          "name": "Qwen 3 30B",
          "id": "qwen/qwen3-30b-a3b-2507",
          "context_window": 256000,
          "default_max_tokens": 20000
        }
      ]
    }
  }
}
```

### Custom Providers

OpenPilot supports custom provider configurations for both OpenAI-compatible and
Anthropic-compatible APIs.

#### OpenAI-Compatible APIs

Here’s an example configuration for Deepseek, which uses an OpenAI-compatible
API. Don't forget to set `DEEPSEEK_API_KEY` in your environment.

```json
{
  "$schema": "https://surya.land/openpilot.json",
  "providers": {
    "deepseek": {
      "type": "openai",
      "base_url": "https://api.deepseek.com/v1",
      "api_key": "$DEEPSEEK_API_KEY",
      "models": [
        {
          "id": "deepseek-chat",
          "name": "Deepseek V3",
          "cost_per_1m_in": 0.27,
          "cost_per_1m_out": 1.1,
          "cost_per_1m_in_cached": 0.07,
          "cost_per_1m_out_cached": 1.1,
          "context_window": 64000,
          "default_max_tokens": 5000
        }
      ]
    }
  }
}
```

#### Anthropic-Compatible APIs

Custom Anthropic-compatible providers follow this format:

```json
{
  "$schema": "https://surya.land/openpilot.json",
  "providers": {
    "custom-anthropic": {
      "type": "anthropic",
      "base_url": "https://api.anthropic.com/v1",
      "api_key": "$ANTHROPIC_API_KEY",
      "extra_headers": {
        "anthropic-version": "2023-06-01"
      },
      "models": [
        {
          "id": "claude-sonnet-4-20250514",
          "name": "Claude Sonnet 4",
          "cost_per_1m_in": 3,
          "cost_per_1m_out": 15,
          "cost_per_1m_in_cached": 3.75,
          "cost_per_1m_out_cached": 0.3,
          "context_window": 200000,
          "default_max_tokens": 50000,
          "can_reason": true,
          "supports_attachments": true
        }
      ]
    }
  }
}
```

### Amazon Bedrock

OpenPilot currently supports running Anthropic models through Bedrock, with caching disabled.

- A Bedrock provider will appear once you have AWS configured, i.e. `aws configure`
- OpenPilot also expects the `AWS_REGION` or `AWS_DEFAULT_REGION` to be set
- To use a specific AWS profile set `AWS_PROFILE` in your environment, i.e. `AWS_PROFILE=myprofile openpilot`

### Vertex AI Platform

Vertex AI will appear in the list of available providers when `VERTEXAI_PROJECT` and `VERTEXAI_LOCATION` are set. You will also need to be authenticated:

```bash
gcloud auth application-default login
```

To add specific models to the configuration, configure as such:

```json
{
  "$schema": "https://surya.land/openpilot.json",
  "providers": {
    "vertexai": {
      "models": [
        {
          "id": "claude-sonnet-4@20250514",
          "name": "VertexAI Sonnet 4",
          "cost_per_1m_in": 3,
          "cost_per_1m_out": 15,
          "cost_per_1m_in_cached": 3.75,
          "cost_per_1m_out_cached": 0.3,
          "context_window": 200000,
          "default_max_tokens": 50000,
          "can_reason": true,
          "supports_attachments": true
        }
      ]
    }
  }
}
```

## A Note on Claude Max and GitHub Copilot

OpenPilot only supports model providers through official, compliant APIs. We do not
support or endorse any methods that rely on personal Claude Max and GitHub Copilot
accounts or OAuth workarounds, which may violate Anthropic and Microsoft’s
Terms of Service.

We’re committed to building sustainable, trusted integrations with model
providers. If you’re a provider interested in working with us,
[reach out](mailto:vt100@surya.sh).

## Logging

Sometimes you need to look at logs. Luckily, OpenPilot logs all sorts of
stuff. Logs are stored in `./.openpilot/logs/openpilot.log` relative to the project.

The CLI also contains some helper commands to make perusing recent logs easier:

```bash
# Print the last 1000 lines
openpilot logs

# Print the last 500 lines
openpilot logs --tail 500

# Follow logs in real time
openpilot logs --follow
```

Want more logging? Run `openpilot` with the `--debug` flag, or enable it in the
config:

```json
{
  "$schema": "https://surya.land/openpilot.json",
  "options": {
    "debug": true,
    "debug_lsp": true
  }
}
```

## Whatcha think?

We’d love to hear your thoughts on this project. Need help? We gotchu. You can find us on:

- [Twitter](https://twitter.com/suryacli)
- [Discord][discord]
- [Slack](https://surya.land/slack)
- [The Fediverse](https://mastodon.social/@suryacli)

[discord]: https://surya.land/discord

## License

[FSL-1.1-MIT](https://github.com/surya/openpilot/raw/main/LICENSE)

---

Part of [Surya](https://surya.land).

<a href="https://surya.land/"><img alt="The Surya logo" width="400" src="https://stuff.surya.sh/surya-banner-next.jpg" /></a>

<!--prettier-ignore-->
Surya热爱开源 • Surya loves open source
