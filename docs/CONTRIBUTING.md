# Contributing

Contributions are welcome!
Just submit a pull request.

## Versioning

Starting from `v1.0.0`, [semantic versioning](https://semver.org) (including
metadata) is followed, including pre-release versions.

## Releasing

Releases are automated on GitHub Actions.
Once the release is published on GitHub Releases,
the GitHub actions workflow will automatically build the artifacts and update
the `go.dev` metadata.

For release changelogs, this repository uses Release Drafter to generate the
initial draft based on the PR titles and labels.
