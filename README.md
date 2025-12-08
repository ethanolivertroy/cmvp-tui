# NIST CMVP CLI

Terminal UI for searching NIST Cryptographic Module Validation Program (CMVP) validated modules.



https://github.com/user-attachments/assets/382222d7-1a74-45c2-83b5-847b534f2c6a






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

## License

MIT
