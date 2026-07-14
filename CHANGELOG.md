# Changelog

All notable changes to `chronary-go` will be documented in this file starting with the soft-launch release.

## 0.7.0 — 2026-07-14

- Add `BookingPageService` (CRUD) for agent-created public scheduling links, plus a `BookingPages` usage counter (#1036).

## 0.6.0 — 2026-07-12

- Add `ConnectionLinksService` for create/get/cancel human calendar setup requests.
- Add `AllowStale` query support and provider-neutral availability health fields; correct free-slot decoding to `start`/`end`.

## 0.1.2 — 2026-05-18

- Add `CONTRIBUTING.md` to the public mirror documenting that this repo is generated from a private monorepo; PRs are welcome as proof-of-concept but can't be merged directly. No behavioral change.
