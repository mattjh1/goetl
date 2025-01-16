# Goetl

Goetl is a simple ETL tool to extract, transform, and load data for RAG (Retrieval-Augmented Generation) solutions.

## Installation

### Prerequisites

- **Go** installed on your system.
- **Make** for managing builds and installation.

### Build and Install

```bash
make install
```

Builds the binary and installs it to `~/.local/bin`.

### Uninstall

```bash
make uninstall
```

Removes the binary from `~/.local/bin`.

## Configuration

The default configuration file is located at:  
`$HOME/.goetl/goetl.yaml`

You can specify a different configuration file using the `--config` flag:

```bash
goetl etl --config /path/to/your/config.yaml
```

## Usage

Run the `goetl` CLI with:

```bash
goetl [command] [flags]
```

Example:

```bash
goetl etl --source_path ./data --tika_server_url http://localhost:9998
```

For detailed help:

```bash
goetl --help
```
