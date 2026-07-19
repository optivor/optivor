# CI/CD

This document describes the automated checks and release pipeline for
Optivor. It complements `AGENTS.md` (which tells an agent *when* to open
PRs and against which branch) by describing what actually runs
automatically once a PR exists, and how a release gets built and
published.

Branch model recap (see `AGENTS.md` §0 for full detail):

```
feature branches  →  staging  →  main
```

CI runs on **every PR into `staging`**. CD runs on **every merge/tag on
`main`**. `staging` itself does not deploy or publish anything — it is
purely an integration and test surface.

## Continuous Integration (on PRs into `staging`)

All of the following are required status checks — a PR cannot be merged
into `staging` until every one of them is green. None of these should be
made optional as a convenience; if a check is too slow or too flaky to
be a hard requirement, that's a problem to fix in the check, not a
reason to make it advisory.

### 1. Build

```bash
go build ./...
```

Note: CGO builds require both `libvips-dev` and `libheif-dev` installed on the host or builder container for full WebP and AVIF format support.

Run on at least the primary target platform in early V0; expand to a
build matrix (linux/darwin/windows × amd64/arm64) once cross-compilation
matters for the milestone in progress (see `ROADMAP.md`).

### 2. Static analysis

```bash
go vet ./...
golangci-lint run
```

`golangci-lint` config should enable at minimum: `staticcheck`, `errcheck`,
`unused`, `gosimple`, `govet`. Lint failures block merge — no
"fix in a follow-up" exceptions.

### 3. Tests

```bash
go test ./... -race -cover
```

- `-race` is non-negotiable given the runtime's concurrency model (see
  `docs/adr/0004-technology-choices.md` on why Go was chosen partly for
  its concurrency model — untested races undermine that choice).
- Coverage is reported but not hard-gated by a percentage threshold in
  V0; revisit once the codebase is large enough for a threshold to be
  meaningful rather than arbitrary.

### 4. Layer boundary check (encodes ADR-0002)

A custom script (not a third-party tool) that walks `go list -deps` (or
equivalent import graph) and fails CI if:

- anything under `runtime/` imports `cli/` or any `deploy-*` package,
- any storage driver package imports `runtime/`.

This turns the dependency-direction rule in
`docs/adr/0002-high-level-architecture.md` from a documentation promise
into something CI actually enforces, rather than relying on every
reviewer to catch a violation by eye.

### 5. Dependency vulnerability scan

```bash
govulncheck ./...
```

Run on every PR, not just on a schedule — catching a vulnerable
dependency at merge time is cheaper than catching it after release.

### 6. Commit message lint

Enforce the Conventional Commits format described in `AGENTS.md` §2 with
a commit-lint action on the PR's commit range. This exists specifically
because an AI agent (or a human in a hurry) can otherwise slip into
`wip` / `fix stuff` commit messages under time pressure — CI should catch
this rather than relying on review discipline alone.

### 7. ADR-awareness check (soft check, not a hard gate)

A bot comment (not a blocking check) that flags when a PR touches a path
matching a new driver, adapter, or extension-point package
(`drivers/*`, `deploy-*`, anything under a `modules/` or `plugins/`-style
path once one exists) without any corresponding change under
`docs/adr/`. This is a nudge for the human reviewer, not an automatic
block — some such PRs are legitimately covered by an existing ADR
already and shouldn't be forced to add a new one.

## Integration tests (on `staging`, post-merge)

Unlike the checks above, these run *after* merge into `staging`, on a
schedule or on every push to `staging` — they're slower and involve real
external services, so they don't gate individual feature PRs one by one,
but they do gate the release promotion in §4.2 of `AGENTS.md`.

- Spin up a real S3-compatible backend (MinIO) as a CI service container
  and run the Storage Driver test suite against it, not just against
  mocks. Mocked unit tests catch logic bugs; this catches SigV4 and
  protocol-level bugs that mocks can't.
- Run a basic end-to-end smoke test: start the standalone runtime binary
  against the MinIO container, request an image, assert it comes back
  resized and correctly encoded.

If integration tests on `staging` go red, that's a signal to fix
`staging` before opening a release PR to `main` — it should not be
treated as a background issue to get to eventually.

## Continuous Deployment (on `main`)

`main` only changes via the release PR flow described in `AGENTS.md`
§4.2. Once that PR is merged:

1. **Tag the release** on `main` using semantic versioning
   (`vMAJOR.MINOR.PATCH`), consistent with the milestone in
   `ROADMAP.md` (e.g., V0 releases as `v0.x.y`, pre-1.0).
2. **GoReleaser runs on tag push**, and handles:
   - cross-compiled binaries for supported platforms,
   - checksums for every artifact,
   - a GitHub Release with generated release notes,
   - Docker image build and push (multi-arch) to the project's
     container registry,
   - Homebrew tap formula update,
   - APT/RPM package generation once those distribution channels are
     actually set up (not required for early V0 releases).
3. **Changelog** is generated from Conventional Commit history between
   the previous tag and this one — this is a direct payoff of enforcing
   commit message format in CI (§6 above); without it, changelog
   generation would require manual work every release.

No automatic deployment of a *running service* happens as part of this
pipeline in V0 — Optivor is distributed as a binary/image for users to
run themselves (see `docs/adr/0001-project-philosophy.md` on BYOS), not
deployed to infrastructure Optivor operates. This section will need a
new subsection once a Deployment Adapter's own CD (e.g., publishing a
Cloudflare Worker template) is added — that is out of scope until
`ROADMAP.md`'s V0.2 milestone.

## Branch protection rules

### `staging`

- Require all CI checks from the "Continuous Integration" section above
  to pass.
- Require at least one review from the relevant `CODEOWNERS` entry.
- Require the branch to be up to date with `staging` before merging
  (i.e., re-run checks after a rebase, not just at PR-open time).
- Disallow force-push.
- Disallow direct pushes, including from maintainers — everything goes
  through a PR, no exceptions, so history and review trail stay
  consistent.

### `main`

Everything `staging` requires, plus:

- Only accept merges from `staging` via the release PR — no feature
  branch may target `main` directly (this should be enforced by
  repository settings, not just by convention in `AGENTS.md`).
- Require integration tests (not just unit tests) to be green on
  `staging` at the point of promotion.
- Require a passing GoReleaser dry-run (`goreleaser release --snapshot
  --clean`) as part of the release PR's checks, so a broken release
  config is caught before tagging, not after.

## Secrets and environments

- Test credentials (e.g., MinIO access keys used only in CI) are stored
  as repository or environment secrets, scoped to the `staging`
  environment, and are never real production credentials — there's
  nothing "production" about them since Optivor itself doesn't run a
  hosted service.
- Release-time secrets (container registry push credentials, Homebrew
  tap deploy key, any code-signing key) are scoped to a `release`
  environment with required reviewers, so tagging a release is not
  sufficient on its own to publish artifacts without a human
  acknowledging it during early releases. This can be relaxed once the
  release process has enough of a track record to trust unattended.

## What CI does *not* do (yet)

To keep this proportional to `ROADMAP.md`'s current milestone:

- No performance/benchmark regression gating yet — revisit once the
  Runtime Module mechanism (ADR-0003) is decided and there's a stable
  hot path worth benchmarking.
- No automated deployment-adapter end-to-end tests against real cloud
  accounts yet — these are expensive and adapter-specific, and belong
  with each adapter's own repository once `ROADMAP.md`'s V0.2+ milestones
  are reached.