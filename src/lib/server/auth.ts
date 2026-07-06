import { createHmac, timingSafeEqual } from 'node:crypto';
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import { db } from '$lib/server/db';
import { credentials } from '$lib/server/db/schema';

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

export const createSessionToken = () => {
	const expires = Math.floor(Date.now() / 1000) + SESSION_MAX_AGE;
	return `${expires}.${sign(String(expires))}`;
};

export const isSessionValid = (token: string) => {
	const [expires, signature] = token.split('.');
	if (!expires || !signature) return false;
	const expected = Buffer.from(sign(expires));
	const actual = Buffer.from(signature);
	if (expected.length !== actual.length || !timingSafeEqual(expected, actual)) return false;
	return Number(expires) >= Date.now() / 1000;
};

export const hasCredentials = () =>
	db.select({ id: credentials.id }).from(credentials).limit(1).get() !== undefined;

export const assertRegistrationAllowed = (authenticated: boolean, setupToken: unknown) => {
	if (authenticated) return;
	if (hasCredentials() || !env.SETUP_TOKEN || setupToken !== env.SETUP_TOKEN) {
		error(403, 'Registration is closed');
	}
};
