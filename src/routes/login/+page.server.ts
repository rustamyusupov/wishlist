import { hasCredentials } from '$lib/server/auth';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = () => ({ needsSetup: !hasCredentials() });
