# Implementation plan

Stack and requirements: see [decisions.md](decisions.md). Old Go code for
reference: `git show 6cd2630:<path>`.

## 1. Scaffold

- [x] SvelteKit (Svelte 5) + TypeScript via `npx sv create`, adapter-node
- [x] Prettier + ESLint, `.env` handling (`$env/dynamic/private`)
- [x] Dependencies pinned to exact versions (no `^`/`~`)
- [x] Base layout, favicon, global styles (new design)

## 2. Database

- [ ] Drizzle + better-sqlite3, drizzle-kit migrations
- [ ] Schema: `users`, `credentials`, `categories`, `currencies`, `wishes`
      (with `user_id`, `sort`), `prices` (history, `created_at`)
- [ ] Seed categories and currencies
- [ ] Migration script: import data from the existing `wishlist.db`,
      attach wishes to the owner user

## 3. Auth (WebAuthn)

- [ ] Registration: only while `users` is empty + one-time env token;
      store credential (public key, counter, transports)
- [ ] Login: challenge → verify → signed session cookie (1 year)
- [ ] Guard in `hooks.server.ts`: everything behind login except `/login`
- [ ] Logout; support several passkeys per user

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

- [ ] Dockerfile (node adapter), `DB_URL` volume as before
- [ ] Register passkeys on Mac and phone against the production domain

Each step ends with the app running (`npm run dev`) and verified by hand.
