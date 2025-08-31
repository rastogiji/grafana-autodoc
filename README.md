# Grafana Autodoc
[![Coverage Status](https://coveralls.io/repos/github/rastogiji/grafana-autodoc/badge.svg?branch=main)](https://coveralls.io/github/rastogiji/grafana-autodoc?branch=main)
This tool generates documentation for Grafana Dashboards from their json representation. Even though you can use this tool as a standalone binary or as a kubectl plugin, it is intended to be used a Github Action to generate documentation for your Grafana Dashboards as code.

## Example Usage
```yaml
name: Generate Grafana Docs
on:
  push:
    branches:
      - master
    paths:
      - '**/*.json'
jobs:
  get-changes:
    runs-on: ubuntu-latest
    outputs:
      matrix: ${{ steps.set-output.outputs.matrix }}
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      - name: Changed files
        id: changed-files
        uses: tj-actions/changed-files@v44
        with:
          matrix: true
          quotepath: false
      - name: Set output in the matrix format
        id: set-output
        run: |
          JSON=$(echo '${{ steps.changed-files.outputs.all_changed_files }}' | jq -c '{files: .}')
          echo "matrix=$JSON" >> $GITHUB_OUTPUT
          echo "$GITHUB_OUTPUT"

  generate-docs:
    runs-on: ubuntu-latest
    needs: get-changes
    strategy:
      matrix:
        file: ${{ fromJson(needs.get-changes.outputs.matrix).files }}
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      - name: Print Outputs
        run: |
          echo "$GITHUB_OUTPUT"
      - name: Generate Documentation
        uses: rastogiji/grafana-autodoc@v0
        with:
          dashboard: ${{ matrix.file }}
      - name: Get MD file
        run: |
          ls -la
```

### Key Features:

1. Triggers on pushes to the `master` branch that include changes to JSON files.
2. Uses `tj-actions/changed-files` to detect modified JSON files.
3. Creates a matrix of changed files for parallel processing.
4. Generates documentation for each changed Grafana dashboard JSON file using `rastogiji/grafana-autodoc@main`.

### How it works:

1. The `get-changes` job identifies modified JSON files and prepares a matrix.
2. The `generate-docs` job runs for each changed file, generating documentation using the Grafana Autodoc tool.

This will generate <dashboard-name>.md files in the Github Action workspace. You can then use further steps to either commit these files to a new repo, store them in an S3 bucket, etc
