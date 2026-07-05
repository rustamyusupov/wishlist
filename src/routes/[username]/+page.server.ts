import { error } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';
import { listWishGroups } from '$lib/server/wishes';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ params, locals }) => {
	const owner = db
		.select({ id: users.id, username: users.username })
		.from(users)
		.where(eq(users.username, params.username))
		.get();
	if (!owner) error(404, 'Not found');

	return {
		owner: owner.username,
		canEdit: locals.user?.id === owner.id,
		groups: listWishGroups(owner.id)
	};
};
