import { error, json } from '@sveltejs/kit';
import { generateRegistrationOptions, verifyRegistrationResponse } from '@simplewebauthn/server';
import {
	CHALLENGE_COOKIE,
	CHALLENGE_MAX_AGE,
	createSessionToken,
	resolveRegistrant,
	SESSION_COOKIE,
	SESSION_MAX_AGE
} from '$lib/server/auth';
import { db } from '$lib/server/db';
import { credentials } from '$lib/server/db/schema';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request, url, cookies, locals }) => {
	const body = await request.json();
	const user = resolveRegistrant(locals.user, body.token);

	if (!body.response) {
		const registered = db.select({ id: credentials.id }).from(credentials).all();
		const options = await generateRegistrationOptions({
			rpName: 'Wishlist',
			rpID: url.hostname,
			userName: user.username,
			attestationType: 'none',
			excludeCredentials: registered,
			authenticatorSelection: { residentKey: 'preferred', userVerification: 'preferred' }
		});
		cookies.set(CHALLENGE_COOKIE, options.challenge, {
			path: '/login',
			maxAge: CHALLENGE_MAX_AGE
		});
		return json(options);
	}

	const expectedChallenge = cookies.get(CHALLENGE_COOKIE);
	if (!expectedChallenge) error(400, 'Challenge expired, try again');

	const { verified, registrationInfo } = await verifyRegistrationResponse({
		response: body.response,
		expectedChallenge,
		expectedOrigin: url.origin,
		expectedRPID: url.hostname,
		requireUserVerification: false
	});
	if (!verified || !registrationInfo) error(400, 'Passkey verification failed');

	const { credential } = registrationInfo;
	db.insert(credentials)
		.values({
			id: credential.id,
			userId: user.id,
			publicKey: Buffer.from(credential.publicKey),
			counter: credential.counter,
			transports: credential.transports ? JSON.stringify(credential.transports) : null
		})
		.run();

	cookies.delete(CHALLENGE_COOKIE, { path: '/login' });
	cookies.set(SESSION_COOKIE, createSessionToken(user.id), {
		path: '/',
		maxAge: SESSION_MAX_AGE
	});

	return json({ ok: true });
};
