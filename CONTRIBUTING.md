# Contributing

## Table of Contents

- [Guidelines](#guidelines)
- [Getting started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Developing](#developing)
- [Testing](#testing)

## Guidelines

See [GUIDELINES.md](GUIDELINES.md) for more information.

## Getting started

### Prerequisites

1. [Install go](https://go.dev/doc/install)

### Installation

1. Clone the repository:

```shell script
git clone https://github.com/MrSquaare/boba.git
cd boba
```

2. Install dependencies:

```shell script
go mod download
```

## Developing

Use one of the [examples](examples) to test your changes or create a new one.

```shell script
go run examples/basic/main.go
```

## Testing

Lint the code:

> Prerequisites: [Install golangci-lint](https://golangci-lint.run/welcome/install/)

```shell script
golangci-lint run
```

Auto-fix the code:

```shell script
golangci-lint run --fix
```
