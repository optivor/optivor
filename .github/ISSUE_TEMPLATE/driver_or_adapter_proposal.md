---
name: Storage Driver / Deployment Adapter proposal
about: Propose a new Storage Driver or Deployment Adapter
title: "[Driver/Adapter] "
labels: extension-proposal
---

## Type

- [ ] Storage Driver
- [ ] Deployment Adapter
- [ ] Runtime Module (see `docs/adr/0003-extensibility-model.md` — the
      implementation mechanism for these is still being decided; please
      discuss before starting implementation)

## What are you proposing to support

E.g., "Storage Driver for Azure Blob Storage" or "Deployment Adapter for
Fly.io".

## Have you read ADR-0003 (Extensibility Model)?

- [ ] Yes — I understand which category this falls into and the
      constraints that apply to it (in-process vs. out-of-process,
      hot-path vs. deploy-time-only).

## Interface / integration approach

Briefly describe how this will plug into the existing layer boundary
(see `docs/adr/0002-high-level-architecture.md`). You don't need a full
design doc — just enough for a maintainer to sanity-check the approach
before you invest time in it.

## Ownership

Are you intending to maintain this driver/adapter on an ongoing basis?
Optivor is designed so contributors can own an entire subsystem (see
`CONTRIBUTING.md`) — if you'd like to take on long-term ownership, say so
here.
