# pkgen
[![CI Status](https://github.com/ifnotnil/pkgen/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/ifnotnil/pkgen/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/ifnotnil/pkgen/graph/badge.svg?token=c0O5dL2fpQ)](https://codecov.io/gh/ifnotnil/pkgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/ifnotnil/pkgen)](https://goreportcard.com/report/github.com/ifnotnil/pkgen)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ifnotnil/pkgen)](https://pkg.go.dev/github.com/ifnotnil/pkgen)

Generate a file inside each package. Using pre-made templates or custom ones.

## Install
```shell
go install github.com/ifnotnil/pkgen/cmd/pkgen@latest
```

### Running Modes

#### In `//go:generate`
Inside a `//go:generate` comment line this:

```golang
//go:generate pkgen --template '<template>'
```

This will result running only for the current package



#### Project level
Running the `pkgen` in the project level will result running for all the packages recursively.

```shell
pkgen --template '<template>'
```

Or with a custom template

```shell
pkgen --template-custom /path/to/template.tmpl
```

## Templates

### Built-in Templates

| Template    | File                                       | Description |
|-------------|--------------------------------------------|-------------|
| `pkgpath`   | [pkgpath.tmpl](templates/pkgpath.tmpl)     | Simple template that generates the full package path as a string constant |
| `oteltrace` | [oteltrace.tmpl](templates/oteltrace.tmpl) | Basic OpenTelemetry tracing setup with tracer only. It creates a package level tracer, using the full package path as name. |
| `otel`      | [otel.tmpl](templates/otel.tmpl)           | Full OpenTelemetry setup with tracer, meter, and logger for observability. It creates a package level tracer, meter and logger, using the full package path as name. |

### Custom Templates

Each template is rendered provided the struct returned from [`golang.org/x/tools/go/packages`](https://github.com/golang/tools/blob/8866876b956fadd4905eb7f49d5d5301d0bc7644/go/packages/packages.go#L419)

