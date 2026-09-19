# SeCure Rate Assuring Protocol with PrivacY

## Development environment

Scrappy runs on Linux with [Nix](https://nixos.org/download/). The demo needs
a TPM that supports the protocol and read/write access to `/dev/tpm0` for your
normal user. Configure that access through your host's TPM device permissions.
The browser also needs a graphical desktop session.

The shells share a revision and content hash pinned in `nix/pkgs.nix`; they do
not use your machine's `<nixpkgs>` channel. Nix downloads dependencies on first
use and reuses them across shells. Go dependencies are recorded in each module's
`go.mod` and `go.sum`. No container images are needed.

| Command (from the repository root) | Environment |
| --- | --- |
| `nix-shell` | Go, C compiler, pkg-config, SQLite, TPM tools |
| `nix-shell example/shell.nix` | Same environment, for the demo server |
| `nix-shell browser/shell.nix` | Base environment plus Chromium and jq |
| `nix-shell android/shell.nix` | Base environment plus JDK 11, Android SDK/NDK and gomobile |

Entering a shell only sets up the environment. It does not start applications
or initialize the TPM. The C compiler supplied by `mkShell` supports CGO
directly, so `steam-run` is unnecessary. Android tools are only downloaded when
you use the Android shell; that shell accepts the Android SDK license.

## Run the demo

### 1. Start the server

From the repository root:

```sh
nix-shell --run scrappy-server
```

This builds and runs the server from `example/`, where its template and
configuration paths resolve correctly. Startup initializes the demo's issuer
and signer, including TPM state. Keep it running during the browser demo.

If an existing TPM state prevents initialization, investigate it before clearing
the TPM: `tpm2_clear` removes keys and is not an automatic setup step.

### 2. Register the browser extension

Open another terminal at the repository root:

```sh
nix-shell browser/shell.nix
chromium --user-data-dir="$SCRAPPY_ROOT/.cache/chromium" chrome://extensions
```

Enable **Developer mode**, choose **Load unpacked**, and select this checkout's
`browser/extension/` directory. Copy the extension ID displayed on the page,
then close that Chromium window and run:

```sh
scrappy-browser YOUR_EXTENSION_ID
```

Replace `YOUR_EXTENSION_ID` with the actual 32-character ID. The command builds
the native messaging host, registers its absolute path and the extension ID in
the checkout's dedicated Chromium profile, and opens `http://localhost:8081/`.
Repeat registration if you move the checkout. The launcher and server use the
same working directory so they share the signer configuration.

Use `localhost` consistently: the verifier checks the exact origin, including
the hostname and port. No `/etc/hosts` entry, `xhost` change, root browser, or
disabled browser sandbox is needed.

## Android library

```sh
nix-shell android/shell.nix --run scrappy-android
```

This replaces the Compose `android_lib` service and writes
`core/android/scrappy_crypto.aar`. See [android/README.md](android/README.md).

## Benchmarks

Enter `nix-shell` at the repository root, then:

```sh
cd core
go test ./tests -benchmem -run='^$' -bench BenchmarkSignLog -benchtime 20x
go test ./tests -benchmem -run='^$' -bench BenchmarkVerifyLog -benchtime 20x
```

The lower-level ECDAA benchmarks belong to the separate
[`github.com/akakou/ecdaa`](https://github.com/akakou/ecdaa) repository.
