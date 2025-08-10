# pkgen

## Install
```shell
go install github.com/ifnotnil/pkgen/cmd/pkgen@latest
```

### Running Modes

#### In `//go:generate`
Inside a `//go:generate` comment line this:

```golang
//go:generate pkgen --template pkgpath
```

This will result running only for the current package



#### Project level
Running the `pkgen` in the project level will result running for all the packages recursively.

```shell
pkgen --template pkgpath
```

Or with custom template

```shell
pkgen --template-custom /path/to/template.tmpl
```

## Custom Templates 

Each template runs with the struct returned from [`golang.org/x/tools/go/packages`](https://github.com/golang/tools/blob/8866876b956fadd4905eb7f49d5d5301d0bc7644/go/packages/packages.go#L419)

