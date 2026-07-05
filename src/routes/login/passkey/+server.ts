import { error, json } from '@sveltejs/kit';
import {
	generateAuthenticationOptions,
	verifyAuthenticationResponse
} from '@simplewebauthn/server';
import { eq } from 'drizzle-orm';
import {
	CHALLENGE_COOKIE,
	CHALLENGE_MAX_AGE,
	createSessionToken,
	SESSION_COOKIE,
	SESSION_MAX_AGE
} from '$lib/server/auth';
import { db } from '$lib/server/db';
import { credentials } from '$lib/server/db/schema';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request, url, cookies }) => {
	const body = await request.json();

	if (!body.response) {
		const options = await generateAuthenticationOptions({
			rpID: url.hostname,
			userVerification: 'preferred'
		});
		cookies.set(CHALLENGE_COOKIE, options.challenge, {
			path: '/login',
			maxAge: CHALLENGE_MAX_AGE
		});
		return json(options);
	}

	const expectedChallenge = cookies.get(CHALLENGE_COOKIE);
	if (!expectedChallenge) error(400, 'Challenge expired, try again');

	const credential = db.select().from(credentials).where(eq(credentials.id, body.id)).get();
	if (!credential) error(400, 'Unknown passkey');

	const { verified, authenticationInfo } = await verifyAuthenticationResponse({
		response: body,
		expectedChallenge,
		expectedOrigin: url.origin,
		expectedRPID: url.hostname,
		requireUserVerification: false,
		credential: {
			id: credential.id,
			publicKey: new Uint8Array(credential.publicKey),
			counter: credential.counter,
			transports: credential.transports ? JSON.parse(credential.transports) : undefined
		}
	});
	if (!verified) error(400, 'Passkey verification failed');

	db.update(credentials)
		.set({ counter: authenticationInfo.newCounter })
		.where(eq(credentials.id, credential.id))
		.run();

	cookies.delete(CHALLENGE_COOKIE, { path: '/login' });
	cookies.set(SESSION_COOKIE, createSessionToken(credential.userId), {
		path: '/',
		maxAge: SESSION_MAX_AGE
	});

	return json({ ok: true });
};
