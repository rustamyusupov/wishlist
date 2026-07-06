import { redirect, type Handle } from '@sveltejs/kit';
import { isSessionValid, SESSION_COOKIE } from '$lib/server/auth';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get(SESSION_COOKIE);
	event.locals.authenticated = token !== undefined && isSessionValid(token);

	const isProtected = event.route.id === '/new' || event.route.id === '/edit/[id]';
	if (!event.locals.authenticated && isProtected) redirect(303, '/login');
	if (event.locals.authenticated && event.url.pathname === '/login') redirect(303, '/');

	return resolve(event);
};
