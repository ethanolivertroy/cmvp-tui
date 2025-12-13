# cmvp-tui

Terminal UI for searching NIST Cryptographic Module Validation Program (CMVP) validated modules.

> **Disclaimer:** This project is not affiliated with, endorsed by, or officially connected to NIST (National Institute of Standards and Technology). It is an independent tool that consumes publicly available CMVP data.



https://github.com/user-attachments/assets/382222d7-1a74-45c2-83b5-847b534f2c6a






## Install

### Homebrew (macOS/Linux)

```bash
brew install ethanolivertroy/sectools/cmvp
```

### Download Binary

Download the latest release for your platform:

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/ethanolivertroy/cmvp-tui/releases/latest/download/cmvp-darwin-arm64 -o cmvp
chmod +x cmvp
xattr -d com.apple.quarantine cmvp  # Remove macOS quarantine flag
sudo mv cmvp /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -L https://github.com/ethanolivertroy/cmvp-tui/releases/latest/download/cmvp-darwin-amd64 -o cmvp
chmod +x cmvp
xattr -d com.apple.quarantine cmvp  # Remove macOS quarantine flag
sudo mv cmvp /usr/local/bin/
```

> **Note:** The `xattr` command removes the quarantine attribute that macOS adds to downloaded files. Without it, Gatekeeper will block the unsigned binary.

**Linux (AMD64)**
```bash
curl -L https://github.com/ethanolivertroy/cmvp-tui/releases/latest/download/cmvp-linux-amd64 -o cmvp
chmod +x cmvp
sudo mv cmvp /usr/local/bin/
```

**Windows**

Download `cmvp-windows-amd64.exe` from the [releases page](https://github.com/ethanolivertroy/cmvp-tui/releases).

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
