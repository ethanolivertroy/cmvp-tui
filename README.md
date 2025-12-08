# NIST CMVP CLI

Terminal UI for searching NIST Cryptographic Module Validation Program (CMVP) validated modules.

## Install

```bash
go install github.com/ethanolivertroy/nist-cmvp-cli@latest
```

Or build from source:

```bash
git clone https://github.com/ethanolivertroy/nist-cmvp-cli.git
cd nist-cmvp-cli
go build -o cmvp-cli .
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
| `Esc` | Back/clear filter |
| `q` | Quit |

## Data Source

Pulls from [NIST-CMVP-API](https://github.com/ethanolivertroy/NIST-CMVP-API) which mirrors NIST CMVP data.

- Active modules (~1,100)
- Historical modules (~4,000)
- Modules in process (~260)

## License

MIT
