<div align="center">
  <img src="assets/gq-logo.png" alt="gq logo" height="30%" width="500" />
<br />

[![Go Version](https://img.shields.io/badge/go-1.25+-blue)]()
[![License](https://img.shields.io/badge/license-MIT-green)]()
![Coverage](https://raw.githubusercontent.com/jmpargana/gq/refs/heads/badges/.badges/main/coverage.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/jmpargana/gq)
![Release](https://img.shields.io/github/v/release/jmpargana/gq)
![Homebrew](https://img.shields.io/badge/homebrew-available-brightgreen)
![Security](https://img.shields.io/github/actions/workflow/status/jmpargana/gq/govulncheck.yml?label=security)


</div>

## gq CLI

`gq` is a jq-like command-line tool written in Go, inspired by one of John Crickett’s coding challenges. It aims to provide a fast, simple, and expressive way to query and transform JSON data from the command line.

This project is a work in progress and is intentionally being developed as a learning and exploration exercise.

In particular, `gq` has served as a personal case study to practice:

* CLI design and common command-line patterns
* Building parsers
* Performance analysis and optimization in Go

An article detailing these learnings and trade-offs is planned.


## Installation

### Homebrew

`gq` is available via Homebrew using the `jmpargana/tools` tap:

```sh
brew tap jmpargana/tools
brew install gq
```

### Docker

A Docker image is published at:

```
jmpargana/gq
```

You can pull it with:

```sh
docker pull jmpargana/gq
```


## Usage

`gq` is designed to work in a similar way to `jq`, reading JSON from standard input and applying a query expression.

### Basic example

```sh
cat input.json | gq '.[0]'
```

### Using with Docker

When using Docker, make sure to keep stdin open:

```sh
cat input.json | docker run -i jmpargana/gq '.[0]'
```

### Example JSON

Given an input file `input.json`:

```json
[
  { "name": "Alice", "age": 30 },
  { "name": "Bob", "age": 25 }
]
```

You can run:

```sh
cat input.json | gq '[{newName: .[] | .name}]' # or '[.[] | {newName: .name}]'
```

Output:

```text
[
  {
    "newName": "Alice"
  },
  {
    "newName": "Bob"
  }
]
```


## Development

### Local development

Local development follows standard Go CLI conventions.

Requirements:

* Go (recent version recommended)
* Make

Common tasks are available via the `Makefile`:

```sh
make build
make test
```

This will build the `gq` binary and run the test suite.

### Docker development

A `Dockerfile` is included and can be used to build and run the tool locally:

```sh
docker build -t gq .
cat input.json | docker run -i gq '.[0]'
```

The Docker image is built using a minimal base and is ready to run without additional dependencies.


## Release & Distribution

Releases are automated using GoReleaser and include:

* Homebrew formula updates
* Multi-platform binaries
* Docker images
