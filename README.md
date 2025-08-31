# Grafana Autodoc

[![Coverage Status](https://coveralls.io/repos/github/rastogiji/grafana-autodoc/badge.svg?branch=main)](https://coveralls.io/github/rastogiji/grafana-autodoc?branch=main) [![Release](https://img.shields.io/github/release/rastogiji/grafana-autodoc.svg)](https://github.com/rastogiji/grafana-autodoc/releases/latest) [![Go Report Card](https://goreportcard.com/badge/github.com/rastogiji/grafana-autodoc)](https://goreportcard.com/report/github.com/rastogiji/grafana-autodoc)

A command-line tool that automatically generates documentation from Grafana dashboard JSON files. It supports processing single files, directories, or glob patterns and outputs structured markdown documentation.

## Features

- 🚀 **Multiple input formats**: Single files, directories, or glob patterns
- 📝 **Markdown output**: Clean, structured documentation
- 🐳 **Docker support**: Containerized execution
- ⚡ **GitHub Action**: Automated documentation in CI/CD
- 🍺 **Homebrew**: Easy installation on macOS and Linux

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap rastogiji/tap
brew install grafana-autodoc
```

Or install directly:
```bash
brew install rastogiji/tap/grafana-autodoc
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/rastogiji/grafana-autodoc/releases):

```bash
# Linux x86_64
curl -L -o grafana-autodoc.tar.gz https://github.com/rastogiji/grafana-autodoc/releases/latest/download/grafana-autodoc_Linux_x86_64.tar.gz
tar -xzf grafana-autodoc.tar.gz
sudo mv grafana-autodoc /usr/local/bin/

# macOS Intel
curl -L -o grafana-autodoc.tar.gz https://github.com/rastogiji/grafana-autodoc/releases/latest/download/grafana-autodoc_Darwin_x86_64.tar.gz
tar -xzf grafana-autodoc.tar.gz
sudo mv grafana-autodoc /usr/local/bin/

# macOS Apple Silicon
curl -L -o grafana-autodoc.tar.gz https://github.com/rastogiji/grafana-autodoc/releases/latest/download/grafana-autodoc_Darwin_arm64.tar.gz
tar -xzf grafana-autodoc.tar.gz
sudo mv grafana-autodoc /usr/local/bin/
```

## Usage

### Command Line

```bash
# Process a single dashboard file
grafana-autodoc --input dashboard.json --output ./docs

# Process all JSON files in a directory
grafana-autodoc --input ./dashboards --output ./docs

# Process files matching a glob pattern
grafana-autodoc --input "./dashboards/*.json" --output ./docs

# Check version
grafana-autodoc --version

# Show help
grafana-autodoc --help
```

### GitHub Action

Use the GitHub Action to automatically generate documentation when dashboard files change:

```yaml
name: Generate Grafana Documentation

on:
  push:
    branches: [main]
    paths: ['dashboards/**/*.json']
  pull_request:
    paths: ['dashboards/**/*.json']

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        
      - name: Generate Documentation
        uses: rastogiji/grafana-autodoc@v1
        with:
          dashboard_files: './dashboards/*.json'
          output_dir: './docs'
      
      # Commit back to the repo or store it to a remote location
```

### Advanced GitHub Action Example

For processing only changed files:

```yaml
name: Generate Grafana Docs for Changed Files

on:
  push:
    branches: [main]
    paths: ['dashboards/**/*.json']

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      changed-files: ${{ steps.changed-files.outputs.all_changed_files }}
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
          
      - name: Get changed files
        id: changed-files
        uses: tj-actions/changed-files@v44
        with:
          files: dashboards/**/*.json
          
  generate-docs:
    runs-on: ubuntu-latest
    needs: detect-changes
    if: needs.detect-changes.outputs.changed-files != ''
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        
      - name: Generate Documentation
        uses: rastogiji/grafana-autodoc@v1
        with:
          dashboard_files: ${{ needs.detect-changes.outputs.changed-files }}
          output_dir: './docs'

      # Commit back to the repo or store it to a remote location
```

## Development

### Prerequisites

- Go 1.23 or later
- Make (optional, for using Makefile commands)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/rastogiji/grafana-autodoc.git
cd grafana-autodoc

# Install dependencies
go mod download

# Build the binary
make build

# Run tests
make test

```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 [Issues](https://github.com/rastogiji/grafana-autodoc/issues)
- 📧 [Email](mailto:animesh.rastogi54@gmail.com)
