import { apiGet } from '$lib/server/api';
import type { Session } from '$lib/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ request }) => {
	const cookie = request.headers.get('cookie') ?? '';
	const session = await apiGet<Session>(cookie, '/auth/session');
	return { needsSetup: session.needsSetup };
};
