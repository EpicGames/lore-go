<!-- Mirrors Lore CONTRIBUTING.md as of 2026-05 -->
# Contributing

The Lore Go SDK follows the Lore project contribution process. Issues and pull
requests are welcome in this repository.

## Prerequisites

- **Go 1.24+** — run `go version` to check
- **Python 3.10+** — required by the binding generator, with [uv](https://github.com/astral-sh/uv) for dependency management

## Build and test

```sh
uv pip install jinja2 pycparser
uv run python find_lorelib.py
uv run python generator/generate.py
go build -C lore_go
```

Run the tests:

```sh
go test -C lore_go ./... -v
```

## Formatting

Formatting is enforced by CI and must pass before any PR is merged:

```sh
go fmt ./...
```

## Before you code

For anything beyond a trivial fix, open a GitHub Issue and wait for a
maintainer to weigh in before investing significant effort. Changes to the
wire protocol or `lore-capi` belong in the
[Lore repository](https://github.com/EpicGames/lore), not here.

## Commit sign-off

Every commit must include a `Signed-off-by:` line. Add it with
`git commit -s`. The DCO, patent affirmation, copyright header rules, and
license compatibility policy are all defined in the canonical contributing
doc.

## Full contribution policy

The full PR process, review process, AI assistance policy, legal terms, and
community channels are published in the Lore repository:

→ [Lore CONTRIBUTING.md](https://github.com/EpicGames/lore/blob/main/CONTRIBUTING.md)
