# Contributing

Thanks your help with building Firemodel!

## Code

Firemodel is written in go 1.11.

    go get -t -u github.com/mickeyreiss/firemodel/...

There are no other system dependencies required to develop or build firemodel.

## Testing

Firemodel uses standard go test for unit tests.

    go test ./...

These tests also run in CircleCI.

The majority of test coverage is based on generated fixtures.

Typical workflow:

1. Implement new feature.
2. Run tests: `go test`.
3. Review fixture diff.
4. If changes are acceptable, regenerate fixtures: `FIREMODEL_UPDATE_FIXTURES=true go test ./...`.
5. Commit code changes and fixture updates in a single commit.

## Releases

Firemodel is still in early development. There are not yet regular releases.

To cut a new binary release on GitHub, use `release.sh`.

Version numbers should use semver, and breaking changes—both to the CLI as well as to the firemodel output—should be avoided.
