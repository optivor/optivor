# Contributing to Optivor

Thank you for considering a contribution. This document explains how the
project is organized so you can find the right place to start — and so
you don't need to read the entire codebase before you can contribute.

## Before anything else

Read these two documents. Most "why did you design it this way"
questions are already answered there:

- [`docs/adr/0000-project-scope.md`](./docs/adr/0000-project-scope.md) —
  what's in scope right now, what isn't.
- [`docs/adr/0001-project-philosophy.md`](./docs/adr/0001-project-philosophy.md) —
  the principles every design decision is checked against.

If you're proposing something that touches architecture or an extension
point, also read
[`docs/adr/0002-high-level-architecture.md`](./docs/adr/0002-high-level-architecture.md)
and
[`docs/adr/0003-extensibility-model.md`](./docs/adr/0003-extensibility-model.md).

## You can own a subsystem, not the whole codebase

Optivor is deliberately split into layers with narrow boundaries (CLI,
Runtime, Storage Drivers, Deployment Adapters — see ADR-0002). This means
you can contribute a Storage Driver, for example, by reading only the
driver interface and its tests — you don't need to understand the CLI or
any deployment adapter to do it well. If a PR review asks you to
understand a part of the codebase far outside what you're changing,
that's worth flagging as a process problem, not something to push
through.

## Where to start

| I want to... | Start here |
|---|---|
| Fix a bug in the runtime | Open an issue first using the Bug Report template, unless the fix is trivial and obvious |
| Add a new Storage Driver | Read `docs/adr/0003-extensibility-model.md`, then open a Driver Proposal issue before writing code |
| Add a new Deployment Adapter | Same as above — these are separate repositories once the pattern is established |
| Propose a new feature | Open a Feature Request issue; do not submit a PR before there's agreement it fits current scope (ADR-0000) |
| Improve docs | PRs welcome directly, no issue required |

## Scope discipline

Before proposing a feature, check
[`docs/adr/0000-project-scope.md`](./docs/adr/0000-project-scope.md) and
[`ROADMAP.md`](./ROADMAP.md). If what you want to build isn't in the
current milestone, it doesn't mean it's unwanted — it means it needs to
be discussed and scheduled first, so the project doesn't grow faster than
it can be reviewed and maintained. PRs that silently expand scope beyond
the current milestone will be asked to split or wait, not merged as-is.

## Pull request expectations

- Keep PRs scoped to one logical change. Large, multi-concern PRs are
  harder to review and more likely to stall.
- If your change affects an architectural boundary (e.g., adds a new
  kind of extension point, changes a driver interface), it likely needs
  a new ADR, not just a PR description. Ask in the issue first if you're
  unsure.
- Tests are expected for any change to the Runtime or a Storage Driver.
- Update relevant documentation in the same PR as the code change, not
  as a follow-up.

## Code review and ownership

See [`CODEOWNERS`](./CODEOWNERS) for who reviews which parts of the
repository. If you become a regular, trusted contributor to a specific
subsystem (a driver, an adapter), maintainership of that subsystem is
something we actively want to hand off — open an issue if you'd like to
discuss taking on that role.

## Code of conduct

Be respectful, assume good faith, and keep disagreements focused on the
technical merits. A formal Code of Conduct document will be added as the
contributor base grows; until then, this paragraph is the standard.
