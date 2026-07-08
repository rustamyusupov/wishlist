import { redirect, type Handle } from '@sveltejs/kit';
import { apiGet } from '$lib/server/api';
import type { Session } from '$lib/types';

export const handle: Handle = async ({ event, resolve }) => {
	const cookie = event.request.headers.get('cookie') ?? '';
	const session = await apiGet<Session>(cookie, '/auth/session').catch(() => null);
	event.locals.authenticated = session?.authenticated ?? false;

	const isProtected = event.route.id === '/new' || event.route.id === '/edit/[id]';
	if (!event.locals.authenticated && isProtected) redirect(303, '/login');
	if (event.locals.authenticated && event.url.pathname === '/login') redirect(303, '/');

	return resolve(event);
};
