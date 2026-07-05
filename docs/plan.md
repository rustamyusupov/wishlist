# Implementation plan

Stack and requirements: see [decisions.md](decisions.md). Old Go code for
reference: `git show 6cd2630:<path>`.

## 1. Scaffold

- [x] SvelteKit (Svelte 5) + TypeScript via `npx sv create`, adapter-node
- [x] Prettier + ESLint, `.env` handling (`$env/dynamic/private`)
- [x] Dependencies pinned to exact versions (no `^`/`~`)
- [x] Base layout, favicon, global styles (new design)

## 2. Database

- [x] Drizzle + better-sqlite3, `drizzle-kit push` workflow — `schema.ts`
      is the single source of truth, no migration files
- [x] Schema: `users`, `credentials`, `categories`, `currencies`, `wishes`
      (with `user_id`, `sort`), `prices` (history, `created_at`)
- [x] One-time import from the legacy prod snapshot — done; `local.db` now
      holds the production data (wishes, prices, seeded categories and
      currencies) in the new schema

## 3. Auth (WebAuthn)

- [x] Registration: only while no passkey exists + `SETUP_TOKEN` from env;
      store credential (public key, counter, transports)
- [x] Login: challenge → verify → signed session cookie (1 year)
- [x] Guard in `hooks.server.ts`: everything behind login except `/login`
- [x] Logout; several passkeys per user supported in schema and endpoints
      (Mac + phone covered by iCloud passkey sync)

## 4. Pages

- [ ] `/` — wish list grouped like the old home page, prices formatted
      with non-breaking space
- [ ] `/new`, `/edit/[id]` — form actions for create/update/delete
- [ ] `/login` — passkey flow
- [ ] Fresh design for list, forms and login: typography, colors, dark theme

## 5. Drag-and-drop

- [ ] svelte-dnd-action on the list (desktop + touch)
- [ ] Reorder action: new id order → update `sort` in one transaction
- [ ] Optimistic UI, keyboard accessibility

## 6. Deploy

- [ ] Dockerfile (node adapter), `DATABASE_URL` volume

Each step ends with the app running (`npm run dev`) and verified by hand.
