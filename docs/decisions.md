# Architecture decisions

## Context

The app is being rewritten from Go (net/http + html/template + SQLite) to a
frontend-centric stack. The old code is available in git history
(`git show 6cd2630:<path>`).

Goal — 50/50 pragmatism and learning: a modern stack without exotics.

## Principles

The codebase must read as a reference SvelteKit app — everything the official
docs recommend, nothing bespoke:

- Official tooling only: scaffolded and extended with the `sv` CLI,
  `npm` as the package manager
- Svelte 5 idioms: runes (`$state`, `$derived`, `$props`), snippets over slots
- Data flow the SvelteKit way: `load` functions, form actions with progressive
  enhancement (`use:enhance`) — pages work without JS where feasible
- Server code in `$lib/server`, validated env via `$env` modules
- Quality gates: `svelte-check`, ESLint (eslint-plugin-svelte) + Prettier;
  a11y warnings are errors. No automated tests, no code comments —
  verification is by hand

## Stack

| Area      | Decision                                 | Why                                                                                        |
| --------- | ---------------------------------------- | ------------------------------------------------------------------------------------------ |
| Framework | **SvelteKit (Svelte 5)**                 | Fullstack in one app: pages, form actions, API, sessions                                   |
| Language  | **TypeScript**                           | Typed models, safe refactoring                                                             |
| Styling   | **Plain CSS**, fresh redesign            | Svelte scopes styles per component; pastel lavender palette, dark theme via `light-dark()` |
| Database  | **SQLite** + **Drizzle ORM**             | Typed schema, `drizzle-kit push` workflow; prod data imported once into the new schema     |
| Auth      | **Passkey / WebAuthn** (@simplewebauthn) | Touch ID login, no password, no external provider. Own signed session cookie               |
| Deploy    | **Same VPS, Docker**, Node adapter       | Infrastructure unchanged                                                                   |

## Functional requirements

Reproduce as is:

- **Single user (the owner).** The list at `/` is public read-only;
  creating, editing, deleting and reordering require login
- One-year session
- Wish CRUD: name, link, category, price + currency, manual order (`sort`)
- Price history in `prices`
- Pages: `/` (list), `/new`, `/edit/[id]`, login/logout

New:

- **Drag-and-drop reordering** on desktop and mobile via **svelte-dnd-action**
  (HTML5 DnD doesn't work on touch). Reorder sends the new id order; server
  updates `sort` in one transaction

## WebAuthn notes

- Bootstrap: registration route works only while `credentials` is empty
  (plus a one-time env token) — otherwise the first visitor becomes the owner
- `credentials` table: public key, signature counter, transports;
  the owner can have several passkeys (Mac, phone)
- Passkey is domain-bound: `localhost` locally, production domain in prod;
  register from both devices or rely on iCloud passkey sync
