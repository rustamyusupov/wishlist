import { redirect, type Handle } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';
import { getSessionUserId, SESSION_COOKIE } from '$lib/server/auth';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get(SESSION_COOKIE);
	const userId = token ? getSessionUserId(token) : null;

	event.locals.user = userId
		? (db
				.select({ id: users.id, username: users.username })
				.from(users)
				.where(eq(users.id, userId))
				.get() ?? null)
		: null;

	const isPublicRoute =
		event.route.id === '/[username]' ||
		event.url.pathname === '/login' ||
		event.url.pathname.startsWith('/login/');

	if (!event.locals.user && !isPublicRoute) redirect(303, '/login');
	if (event.locals.user && event.url.pathname === '/login') redirect(303, '/');

	return resolve(event);
};
