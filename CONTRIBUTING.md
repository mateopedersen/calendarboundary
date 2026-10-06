# Contributing

Open an issue before changing the public API. Include the relevant civil dates, expected ISO week metadata, and whether the behavior should affect display rows or date arithmetic.

Before submitting a change, run:

```sh
gofmt -w *.go
go test ./...
go vet ./...
go test -race ./...
```

Add table-driven regression cases for leap rules, New Year crossovers, configurable week starts, and row-count boundaries. Keep the module dependency-free unless a concrete need justifies a dependency.
