# pkgen

### Running Modes

#### In `//go:generate`
Inside a `//go:generate` comment line this:

```golang
//go:generate pkgen --template pkgpath
```

This will result running only for the current package



#### Project level
Running the `pkgen` in the project level will result running for all the packages recursively.

