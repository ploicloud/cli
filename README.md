# Ploi Cloud CLI

`ploicloud` (short alias: `pcctl`) is the official command-line interface for [Ploi Cloud](https://ploi.cloud). Every REST endpoint is exposed as a subcommand, auto-generated from the platform's OpenAPI spec.

```sh
ploicloud login
ploicloud applications list
ploicloud applications deploy 42
pcctl whoami
```

---

## Install

**Homebrew:**

```sh
brew install ploicloud/tap/ploicloud
```

**Direct binary:** grab the right archive for your OS/arch from [Releases](https://github.com/ploicloud/cli/releases/latest) and drop `ploicloud` (and `pcctl`) somewhere on your `$PATH`.

## Update

```sh
brew upgrade ploicloud
```

Or re-download the latest binary from the Releases page.

---

## Run and test locally

You need Go 1.25+ and `make`.

```sh
git clone git@github.com:ploicloud/cli.git
cd cli
make build
./bin/ploicloud --help
```

`make build` does three things: regenerates the cobra command tree from the committed OpenAPI spec, compiles `cmd/ploicloud`, and copies the result to a second binary named `pcctl`.

**Useful targets:**

```sh
make build         # regenerate + compile into ./bin/{ploicloud,pcctl}
make generate      # only regenerate ./internal/commands/zz_generated.go
make sync-spec     # curl https://ploi.cloud/docs/api.json -> gen/spec.json
make test          # go test ./...
make clean         # rm -rf bin/ dist/ zz_generated.go
```

**Point the CLI at a non-production API** (e.g. local dev or a staging deploy):

```sh
./bin/ploicloud --api-url=https://ploi.test applications list
```

The `--api-url` flag is honored on every subcommand. To persist it:

```sh
mkdir -p ~/.ploicloud
cat > ~/.ploicloud/config.toml <<EOF
api_url = "https://ploi.test"
EOF
```

**Custom OAuth client_id** at build time (the production UUID is hardcoded as the default):

```sh
make build CLIENT_ID=<some-other-passport-client-uuid>
```

**Running a single test package:**

```sh
go test ./internal/auth/ -v
```

---

## Release

One command, no prompts:

```sh
make release
```

What it does:

1. Aborts if the working tree has uncommitted changes.
2. Pulls latest from `origin` (if upstream is set).
3. Runs `make sync-spec` to pull the freshest OpenAPI spec from `https://ploi.cloud/docs/api.json`.
4. Regenerates commands, builds, runs tests. Fails fast if anything breaks.
5. Auto-bumps the patch version from the latest git tag (`v0.0.4 → v0.0.5`).
6. Commits `gen/spec.json` if it changed.
7. Tags the commit and pushes with `--follow-tags`.

The tag push triggers `.github/workflows/release.yml`, which runs goreleaser:

- Builds darwin/linux/windows × amd64/arm64 binaries.
- Uploads archives containing both `ploicloud` and `pcctl` to a new GitHub Release.
- Updates the `ploicloud/homebrew-tap` formula, so `brew upgrade ploicloud` immediately sees the new version.

**Required GitHub secrets** (configured in repo settings):

- `HOMEBREW_TAP_GITHUB_TOKEN` — PAT with `repo` scope on `ploicloud/homebrew-tap` so goreleaser can push the formula update.
- `OAUTH_CLIENT_ID` — the Passport public client UUID (`019e5393-…`); baked into release binaries via ldflags.

---

## How the generated commands work

The OpenAPI spec at `gen/spec.json` is the source of truth. `gen/cmdgen/main.go` parses it at build time and emits `internal/commands/zz_generated.go` — one cobra `Command` per operation, with path params as positional args and query/body params as typed flags. The CLI sends bearer-authenticated HTTPS requests via `internal/client/transport.go` and prints JSON. No runtime spec parsing, no reflection.

When the platform API changes, `make release` re-fetches the spec and regenerates everything in lockstep with the version bump.
