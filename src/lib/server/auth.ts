import { createHmac, timingSafeEqual } from 'node:crypto';
import { error } from '@sveltejs/kit';
import { asc } from 'drizzle-orm';
import { env } from '$env/dynamic/private';
import { db } from '$lib/server/db';
import { credentials, users } from '$lib/server/db/schema';

export const SESSION_COOKIE = 'session';
export const SESSION_MAX_AGE = 365 * 24 * 60 * 60;
export const CHALLENGE_COOKIE = 'challenge';
export const CHALLENGE_MAX_AGE = 5 * 60;

const secret = () => {
	if (!env.SESSION_SECRET) throw new Error('SESSION_SECRET is not set');
	return env.SESSION_SECRET;
};

const sign = (payload: string) =>
	createHmac('sha256', secret()).update(payload).digest('base64url');

export const createSessionToken = (userId: number) => {
	const expires = Math.floor(Date.now() / 1000) + SESSION_MAX_AGE;
	const payload = `${userId}.${expires}`;
	return `${payload}.${sign(payload)}`;
};

export const getSessionUserId = (token: string) => {
	const [userId, expires, signature] = token.split('.');
	if (!userId || !expires || !signature) return null;
	const expected = Buffer.from(sign(`${userId}.${expires}`));
	const actual = Buffer.from(signature);
	if (expected.length !== actual.length || !timingSafeEqual(expected, actual)) return null;
	if (Number(expires) < Date.now() / 1000) return null;
	return Number(userId);
};

export const hasCredentials = () =>
	db.select({ id: credentials.id }).from(credentials).limit(1).get() !== undefined;

export const resolveRegistrant = (
	user: App.Locals['user'],
	setupToken: unknown
): NonNullable<App.Locals['user']> => {
	if (user) return user;
	if (hasCredentials() || !env.SETUP_TOKEN || setupToken !== env.SETUP_TOKEN) {
		error(403, 'Registration is closed');
	}
	return (
		db
			.select({ id: users.id, username: users.username })
			.from(users)
			.orderBy(asc(users.id))
			.limit(1)
			.get() ??
		db
			.insert(users)
			.values({ username: 'owner' })
			.returning({
				id: users.id,
				username: users.username
			})
			.get()
	);
};
