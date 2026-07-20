# AGENTS.md — Working Instructions for AI Coding Agents

This file tells an AI agent (Claude Code or similar) exactly how to work
in this repository: how to branch, how to commit, how to open PRs, and
in what order to write code. It assumes the agent has already read
`docs/adr/` — if it hasn't, it should before doing anything else.

This document is operational, not architectural. For "why," see
`docs/adr/`. For "how to behave while working," read on.

## 0. Before starting any task

1. Read `docs/adr/0000-project-scope.md`. If the task isn't in the
   current milestone (`ROADMAP.md`), stop and flag this instead of
   implementing it. Do not silently expand scope.
2. Read `docs/adr/0002-high-level-architecture.md` and
   `docs/adr/0003-extensibility-model.md` if the task touches more than
   one layer (CLI, Runtime, Storage Driver, Deployment Adapter) or adds
   a new extension point.
3. Check `git status` and `git branch` — never start new work on top of
   uncommitted changes from an unrelated task, and never work directly
   on `main` or `staging`.

### Branch model overview

This repository uses two long-lived branches, not one:

```
feature branches  →  staging  →  main
   (PRs land here)   (integration)  (releases only)
```

- **`main`** always reflects the latest *released* state. It only ever
  receives merges from `staging`, via a release PR (see §4). Nobody
  branches directly off `main` for day-to-day work, and nobody opens a
  feature PR against it.
- **`staging`** is the integration branch. This is where all feature
  work, fixes, and docs changes land first. It is the default branch to
  branch from and the default PR target.
- **Feature branches** (`feat/...`, `fix/...`, `docs/...`, `chore/...`)
  branch off `staging` and PR back into `staging`.

If you are ever unsure which branch to target, the answer is
**`staging`**, unless you are specifically doing a release promotion
(§4.2).

## 1. Branching

- **Never commit directly to `main` or `staging`.** Always create a
  feature branch first.
- One branch per logical change. Do not combine an ADR-relevant
  architectural change with an unrelated bug fix in the same branch.
- Branch from an up-to-date `staging` — **not** `main`:

  ```bash
  git checkout staging
  git pull origin staging
  git checkout -b <type>/<short-description>
  ```

- Branch naming convention: `<type>/<short-description>`, where `<type>`
  matches the commit types below (see §2). Examples:

  ```bash
  git checkout -b feat/s3-storage-driver
  git checkout -b fix/webp-encoder-memory-leak
  git checkout -b docs/adr-0005-storage-driver-api
  git checkout -b chore/goreleaser-config
  ```

- If a task turns out to be larger than expected mid-way, stop and split
  it into a follow-up branch rather than letting one branch accumulate
  unrelated changes.
- Never branch off another in-progress feature branch unless the task is
  explicitly a follow-up that depends on unmerged work — and if you do,
  say so clearly in the PR description, since it changes the review
  order.

## 2. Commit discipline

**Commit every logical unit of work, no matter how small.** A commit
that adds a single struct field, fixes a typo in a doc, or adds one test
case is a valid, encouraged commit. Do not batch unrelated small changes
into one large commit "to save time" — small commits are what make
`git bisect`, code review, and `git blame` useful later. This applies
doubly to an AI agent: reviewers need to see the actual sequence of
reasoning/changes, not a single squashed diff they have to reverse-engineer.

### Commit message format

Use Conventional Commits:

```
<type>(<scope>): <short summary, imperative mood, no trailing period>

<optional body: what changed and why, not just what — the diff already
shows what changed>

<optional footer: Refs #123, or "See ADR-0005">
```

Types:

- `feat` — new functionality
- `fix` — bug fix
- `docs` — documentation only (including ADRs)
- `refactor` — code change that doesn't change behavior
- `test` — adding or fixing tests only
- `chore` — tooling, build config, dependency bumps
- `revert` — reverting a previous commit

Scope should name the layer or component: `runtime`, `cli`, `driver/s3`,
`adapter/cloudflare`, `adr`, etc.

Examples:

```
feat(driver/s3): implement GetObject with SigV4 signing

fix(runtime): return 413 instead of crashing on oversized source images

docs(adr): add ADR-0005 for Storage Driver API

test(runtime): add table-driven tests for resize fit modes

refactor(cli): extract config loading into its own package
```

### What NOT to do

- Do not write commit messages like `wip`, `fix stuff`, `updates`, or
  `changes per review`. Each commit message should be understandable in
  isolation, without needing the PR description for context.
- Do not amend or force-push commits that have already been reviewed by
  a human, unless explicitly asked to — this destroys review history.
- Do not combine a revert with new unrelated work in the same commit.

## 3. Keeping the branch in sync

Before opening a PR, and again if the PR sits open for more than a day:

```bash
git checkout staging
git pull origin staging
git checkout <your-branch>
git rebase staging
```

Prefer rebase over merge for keeping a feature branch current, to keep
history linear and make each commit's diff meaningful on its own. Resolve
conflicts commit-by-commit during rebase rather than squashing everything
into one conflict-resolution commit.

Never rebase onto `main` directly — `main` may be behind `staging` at
any given time by design (it only moves forward on release), so rebasing
a feature branch onto `main` can silently drop work that's already in
`staging`.

## 4. Opening a pull request

### 4.1 Feature PRs — target `staging`

Once the branch is ready:

```bash
git push -u origin <your-branch>
gh pr create \
  --base staging \
  --title "<type>(<scope>): <summary>" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

Then **fill in the PR template properly** — do not leave the checklist
boxes unchecked without addressing them. Specifically:

- The "Scope check" boxes from `PULL_REQUEST_TEMPLATE.md` must be
  answered honestly. If a change crosses a layer boundary it shouldn't
  (per ADR-0002), that's a signal to stop and reconsider the approach,
  not to check the box anyway.
- If the change introduces or modifies an extension point (per
  ADR-0003), the PR must either reference an existing ADR or include a
  new one in `docs/adr/` — a PR description alone is not sufficient for
  architecturally significant changes.
- Link the related issue with `Closes #<number>` if one exists.
- All required CI checks must be green before requesting review — see
  `docs/ci-cd.md`. Do not ask a human (or another agent) to review a PR
  with failing or pending checks.

#### Draft PRs for in-progress work

If a task is large enough to need visibility before it's finished, open
it as a draft:

```bash
gh pr create --draft --base staging --title "..." --body-file .github/PULL_REQUEST_TEMPLATE.md
```

Convert to ready-for-review only once tests pass locally, CI is green,
and the PR template checklist is genuinely satisfied:

```bash
gh pr ready <pr-number>
```

#### Merging a feature PR

- **Never open a feature PR against `main`.** All feature, fix, and docs
  work targets `staging`.
- Do not merge your own PR. Wait for the reviewer defined in
  `CODEOWNERS` for the paths you touched, even if you are confident the
  change is correct — this rule applies to AI agents the same as human
  contributors.
- Use squash merge only for PRs with noisy, non-meaningful intermediate
  commits (e.g., many `fix typo` commits from review feedback). Use a
  regular merge commit when the individual commits are each meaningful
  and worth preserving in history (this will usually be the case if §2
  was followed correctly).

### 4.2 Promoting `staging` to `main` (release PRs)

`main` only moves forward through a deliberate release PR, opened from
`staging`:

```bash
git checkout staging
git pull origin staging
gh pr create \
  --base main \
  --head staging \
  --title "release: promote staging to main" \
  --body "See CHANGELOG / commit list for what's included."
```

This is a separate, less frequent action from day-to-day feature work,
and an AI agent should only open a release PR when explicitly asked to,
not as part of finishing a regular feature task. Once merged, a version
tag is pushed on `main` to trigger the release pipeline — see
`docs/ci-cd.md` §"Continuous Deployment".

## 5. Order of operations when writing code

When implementing a feature that spans layers, **write and commit in
this order**, each step as its own commit or small set of commits:

1. **ADR first, if one is needed** (per ADR-0003's guidance on when a
   new extension point needs its own ADR). Commit this alone, as
   `docs(adr): ...`. Do not write implementation code in the same commit
   as a new ADR — the ADR should be reviewable independently.
2. **Interface / contract** — e.g., a new Storage Driver interface
   method, a new config field. Commit this before the implementation
   that uses it, so the diff shows the contract change in isolation.
3. **Core implementation** — the actual logic. If it's large, split by
   sub-concern (e.g., "add S3 GetObject" as one commit, "add S3 PutObject"
   as another) rather than one commit for the entire driver.
4. **Tests** — ideally alongside each implementation commit, not batched
   at the end. A commit that adds a function with no accompanying test
   coverage should be treated as incomplete unless the function is
   trivial (e.g., a pure config getter).
5. **Documentation** — update `docs/architecture/*.md` or `README.md` in
   the same PR, as a distinct `docs:` commit, not deferred to a
   follow-up PR.

### Respect layer boundaries while writing, not just at review time

Per `docs/adr/0002-high-level-architecture.md`:

- Code under `runtime/` must never import anything from `cli/` or any
  `deploy-*` adapter package.
- Code under a storage driver package must never import from `runtime/`
  — the dependency direction is Runtime → Storage Driver interface,
  never the reverse.
- If implementing a task seems to require violating one of these
  directions, that's a signal the task is mis-scoped or needs an
  interface change first (go back to step 2 above) — do not "just make
  it work" by adding a boundary-crossing import.

### Extension point discipline while writing

Per `docs/adr/0003-extensibility-model.md`:

- Storage Drivers: implement in-process, compiled into the runtime
  binary or a build variant — do not introduce dynamic loading or RPC
  for these.
- Deployment Adapters: implement as out-of-process, separately
  executable components — do not import provider SDKs into shared/core
  packages.
- Runtime Modules (transforms, cache backends, auth/policy): the
  mechanism is undecided as of ADR-0003. Until a follow-up ADR settles
  this, implement new transforms as plain internal Go functions/types,
  not behind a speculative plugin abstraction. Do not invent an
  extension mechanism unilaterally while implementing a feature — flag
  it and wait for the ADR.

## 6. CI/CD expectations

Every PR into `staging` runs the full CI suite automatically; every
merge from `staging` into `main` triggers the release pipeline. An AI
agent should never bypass, disable, or skip a CI check to get a PR
merged faster — a failing check on a small change is a signal to fix the
change, not to work around the check. Full details:
`docs/ci-cd.md`.

## 7. When something doesn't fit this document

If a task doesn't clearly map to the rules above — an ambiguous scope
question, an architectural judgment call, a case where following §5's
order seems wrong for the specific task — stop and ask, or open the PR
as a draft with an explicit note about the uncertainty, rather than
guessing and proceeding silently. Silent scope or architecture decisions
made by an AI agent are exactly what ADR-0000 through ADR-0003 exist to
prevent.