# Contributing

Contributions are welcome. Please open an issue before submitting a pull request
so we can discuss the change first.

## Development setup
```bash
git clone https://github.com/vukan322/docforge
cd docforge
go mod tidy
go test ./...
```

## Guidelines

- Keep the public API clean and minimal
- All logic belongs in `internal/ooxml`, not in `docforge.go`
- All new features must include tests
- New test fixtures go in `testdata/templates/`
- Run `golangci-lint run && go test ./...` before submitting
- Follow existing code style — no comments in code
