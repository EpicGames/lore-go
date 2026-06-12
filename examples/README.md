# Lore Go SDK Examples

This directory contains example applications demonstrating how to use the Lore Go SDK.

## Prerequisites

Each example requires the Lore native library. Run `go generate` to fetch the library before building:

```bash
# From the examples directory
cd examples
go generate ./...
```

This downloads the platform-specific native library into each example directory.

## Building Examples

From the examples directory:

```bash
go build -o fluent ./fluent
go build -o native ./native
go build -o poc ./poc
go build -o notifications ./notifications
```

## Running Examples

The `fluent`, `native`, and `notifications` examples each accept an optional
remote URL as the first command-line argument.

- **No argument** → fully offline run. The example creates a local repository
  and commits a file. Nothing is pushed; nothing is cloned.
- **With argument** (e.g. `lore://localhost`) → online run. The example also
  pushes the revision and, where applicable, clones the repository back or
  subscribes to notifications.

These examples do not perform authentication. If the remote requires it, run
`lore auth` from the CLI before invoking the example.

### Running a local Lore server

To exercise the online mode of these examples, you can run a Lore server
locally. The steps below build the server from source and configure it for
local development:

1. Clone the Lore repository and build the server in release mode:

   ```bash
   git clone https://github.com/EpicGames/lore.git
   cd lore
   cargo build --release
   ```

2. Create a local config file by copying the example:

   ```bash
   cp lore-server/config/local.toml.example lore-server/config/local.toml
   ```

3. Generate a random secret and set it as `presigned_url_hmac_key` in
   `lore-server/config/local.toml`:

   ```bash
   openssl rand -hex 32
   ```

4. Generate a self-signed TLS certificate (run from the directory where the
   server expects `cert.pem` and `key.pem`):

   ```bash
   openssl req \
     -subj '/CN=localhost:8443/O=Self signed/C=CH' \
     -new -newkey rsa:2048 -sha256 -days 365 -nodes -x509 \
     -keyout key.pem -out cert.pem
   ```

5. Start the server:

   ```bash
   RUST_LOG=info ./target/release/loreserver 2>&1 | tee /tmp/lore.log
   ```

The server is now reachable as `lore://localhost`, which you can pass to the
examples below.

### fluent - A full sequence of Lore operations using the fluent API

```bash
# Offline run
./examples/fluent/fluent

# Online run against a local server
./examples/fluent/fluent lore://localhost
```

### native - The same sequence using the low-level native (FFI) API

```bash
# Offline run
./examples/native/native

# Online run against a local server
./examples/native/native lore://localhost
```

### poc - Repository Status Example

```bash
# Run on current directory
./examples/poc/poc

# Run on specific repository
./examples/poc/poc /path/to/repository
```

### notifications - Notification Subscription Example

Demonstrates how to subscribe to and receive real-time notifications for
repository events (such as branch pushes, creates, and deletes). Notifications
require a remote server, so the offline run skips the subscription step.

```bash
# Offline run (creates local repo and commits, but no notifications)
./examples/notifications/notifications

# Online run against a local server
./examples/notifications/notifications lore://localhost
```
