# cmvp-tui

Terminal UI for searching NIST Cryptographic Module Validation Program (CMVP) validated modules.

> **Disclaimer:** This project is not affiliated with, endorsed by, or officially connected to NIST (National Institute of Standards and Technology). It is an independent tool that consumes publicly available CMVP data.



https://github.com/user-attachments/assets/382222d7-1a74-45c2-83b5-847b534f2c6a






## Installation

### Homebrew (macOS/Linux)

```bash
brew install ethanolivertroy/sectools/cmvp
```

### Scoop (Windows)

```powershell
scoop bucket add sectools https://github.com/ethanolivertroy/scoop-sectools
scoop install cmvp
```

### Download Binary

Download from [Releases](https://github.com/ethanolivertroy/cmvp-tui/releases):

| Platform | Binary |
|----------|--------|
| macOS (Apple Silicon) | `cmvp-darwin-arm64` |
| macOS (Intel) | `cmvp-darwin-amd64` |
| Linux (x64) | `cmvp-linux-amd64` |
| Linux (ARM64) | `cmvp-linux-arm64` |
| Windows (x64) | `cmvp-windows-amd64.exe` |

> **Note (macOS):** Remove the quarantine attribute before running: `xattr -d com.apple.quarantine cmvp`

### Go Install

```bash
go install github.com/ethanolivertroy/cmvp-tui@latest
```

### Build from Source

```bash
git clone https://github.com/ethanolivertroy/cmvp-tui.git
cd cmvp-tui
go build -o cmvp .
```

## Usage

```bash
cmvp
```

## Keys

| Key | Action |
|-----|--------|
| `/` | Filter/search |
| `j/k` or arrows | Navigate |
| `Enter` | View details |
| `d` | Toggle algorithm details (in detail view) |
| `Esc` | Back/clear filter |
| `q` | Quit |

## Data Source

Pulls from [NIST-CMVP-API](https://github.com/ethanolivertroy/NIST-CMVP-API) which mirrors NIST CMVP data.

## License

MIT
