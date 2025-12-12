# NIST CMVP CLI

Terminal UI for searching NIST Cryptographic Module Validation Program (CMVP) validated modules.



https://github.com/user-attachments/assets/382222d7-1a74-45c2-83b5-847b534f2c6a






## Install

### Download Binary (Recommended)

Download the latest release for your platform:

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/ethanolivertroy/NIST-CMVP-CLI/releases/latest/download/cmvp-darwin-arm64 -o cmvp
chmod +x cmvp
sudo mv cmvp /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -L https://github.com/ethanolivertroy/NIST-CMVP-CLI/releases/latest/download/cmvp-darwin-amd64 -o cmvp
chmod +x cmvp
sudo mv cmvp /usr/local/bin/
```

**Linux (AMD64)**
```bash
curl -L https://github.com/ethanolivertroy/NIST-CMVP-CLI/releases/latest/download/cmvp-linux-amd64 -o cmvp
chmod +x cmvp
sudo mv cmvp /usr/local/bin/
```

**Windows**

Download `cmvp-windows-amd64.exe` from the [releases page](https://github.com/ethanolivertroy/NIST-CMVP-CLI/releases).

### Go Install

```bash
go install github.com/ethanolivertroy/nist-cmvp-cli@latest
```

### Build from Source

```bash
git clone https://github.com/ethanolivertroy/nist-cmvp-cli.git
cd nist-cmvp-cli
go build -o cmvp .
```

## Usage

```bash
./cmvp-cli
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
