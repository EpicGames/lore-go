# Lore Go SDK 

## About
This repository contains tools to exend Lore with Go. 

Lore is an open source version control system that is designed for unprecedented scalability of both data and teams. It is optimized for projects that combine code with large binary assets, including games and entertainment, and caters for the needs of developers and artists alike. 

For full Lore documentation, architecture details, and contribution guidelines, visit the [main Lore repository](https://github.com/EpicGames/lore).


## Install

### Stable Release

```bash
go get github.com/EpicGames/lore-go@latest
```

### Nightly Build

Nightly builds are published as tagged versions of the form `v<X.Y.Z>-nightly-<REV>`. Pin a specific nightly with:

```bash
go get github.com/EpicGames/lore-go@v0.8.2-nightly-3496
```

### Install the Lore native library

The Go SDK binds against the Lore C library (`liblore.so` / `liblore.dylib` / `lore.dll`). Run `fetch-lore-lib` once per build to install the matching version next to your application binary:

```bash
go run github.com/EpicGames/lore-go/cmd/fetch-lore-lib
```

- `-o <output_dir>` — destination directory (defaults to the current directory)
- `-os <target_os>` and `-arch <target_arch>` — fetch the library for a different platform (e.g. `-os linux -arch amd64` for cross-compilation)

To automate this, add a directive to your `main.go` so `go generate ./...` installs the library:

```go
//go:generate go run github.com/EpicGames/lore-go/cmd/fetch-lore-lib
```

Source priority for `fetch-lore-lib`:

1. `LORE_LIB_PATH` — if set, the file at that path is copied. This env var is also honored by the SDK at runtime (point it at a `.so` / `.dylib` / `.dll` file and both fetch-time and runtime use it).
2. Otherwise the library is downloaded from `LORE_RELEASE_BASE_URL`, falling back to the URL baked into the SDK at build time (see [Generate the Go bindings](#generate-the-go-bindings)). The URL is constructed as `<base>/<versionTag>/<artifactName>`.

At runtime the SDK also searches next to the compiled executable, so a `go generate`-installed library is found automatically.

## Minimal example

The default package (`github.com/EpicGames/lore-go`) exposes the high-level fluent API. A low-level, C-like wrapper around the underlying FFI is also available under `github.com/EpicGames/lore-go/native` for advanced use cases.

```go
import (
    "fmt"

    "github.com/EpicGames/lore-go"
    "github.com/EpicGames/lore-go/types"
)

lore.LogConfigure(&types.LoreLogConfigFFI{
    File:     true,
    FilePath: "/path/to/log/directory",
    Level:    types.LoreLogLevel_DEBUG,
})

globals := types.LoreGlobalArgsFFI{
    RepositoryPath: "/path/to/local/repository",
}
args := types.LoreRepositoryStatusArgsFFI{
    Staged: true,
    Scan:   true,
}
_, err := lore.RepositoryStatus(&globals, &args).
    Callback(func(event types.LoreEvent) {
        if event.Tag == types.LoreEventTag_REPOSITORY_STATUS_FILE {
            fmt.Println(event.Data)
        }
    }).
    Wait()
```

For comprehensive examples, see [examples/fluent/fluent.go](examples/fluent/fluent.go) (fluent) and [examples/native/native.go](examples/native/native.go) (low-level).

## Contributing

### Set up your dev environment

1. Clone the Lore Go SDK repository:

```bash
git clone https://github.com/EpicGames/lore-go
```

2. (Optional) Create a Python virtual environment for the binding generator:

```bash
uv venv .venv
source .venv/bin/activate
```

3. Install the Python modules used by the binding generator:

```bash
uv pip install jinja2 pycparser
```

### Get the Lore library

The SDK binds against the Lore C library. Pick one of the two options below depending on whether you're also modifying the Lore core.

#### Option A — build the library from Lore source

Use this when you're changing the Lore C/Rust core alongside the Go SDK.

1. Clone [Lore's repository](https://github.com/EpicGames/lore) and build it:

```bash
cargo build --release
```

#### Option B — fetch a pre-built Lore library

Use this when you only need to develop the Go SDK against an existing Lore version.

1. Download the header and binaries from [Lore's repository](https://github.com/EpicGames/lore) release page.

### Generate the Go bindings

1. Point `LORE_BUILD_PATH` at the library directory from the previous section:

```bash
export LORE_BUILD_PATH="<path-to>/lore/"
```

2. (Optional) Set `LORE_RELEASE_BASE_URL` and `LORE_VERSION` to bake a default download base URL into `cmd/fetch-lore-lib/version.go`. End users of the published SDK will use this URL when running `fetch-lore-lib` without overriding it themselves. If unset, generation falls back to the public Lore release URL.

```bash
export LORE_RELEASE_BASE_URL="https://github.com/EpicGames/lore/releases/download"
export LORE_VERSION="0.8.2"
```

3. Generate the bindings and build the SDK:

```bash
uv run python find_lorelib.py
uv run python generator/generate.py
go build -C lore_go
```

4. Any edits you now make under `lore_go/` are picked up by re-running `go build`. If you change anything under `generator/templates/` or pull a new Lore pre-built binary, re-run step 3 to regenerate the bindings.

### Run the examples

With the dev environment set up, a Lore library available, and the Go bindings generated, run an example from the repository root:

```bash
export LORE_LIB_PATH="lore/lib/lorelib-arm64-apple-darwin.dylib"
go generate -C examples ./...
go build -C examples -o fluent ./fluent
examples/fluent/fluent
```

To run the low-level native example instead, swap `fluent` for `native`.

### Run the test suite

```bash
go test -C lore_go ./... -v
```

## Releasing

The project is released using the `Release Lore Go SDK` GitHub Action. The workflow runs `validate-cr.yml` to generate the SDK against the requested Lore version (with `LORE_RELEASE_BASE_URL` baked into `cmd/fetch-lore-lib/version.go`), then commits the result and tags it.
