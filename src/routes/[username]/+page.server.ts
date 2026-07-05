import { error, fail } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';
import { db } from '$lib/server/db';
import { categories, users } from '$lib/server/db/schema';
import { listWishGroups, reorderWishes } from '$lib/server/wishes';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = ({ params, locals, url }) => {
	const owner = db
		.select({ id: users.id, username: users.username })
		.from(users)
		.where(eq(users.username, params.username))
		.get();
	if (!owner) error(404, 'Not found');

	const canEdit = locals.user?.id === owner.id;
	return {
		owner: owner.username,
		canEdit,
		groups: listWishGroups(owner.id, canEdit && url.searchParams.has('edit'))
	};
};

export const actions: Actions = {
	reorder: async ({ request, params, locals }) => {
		if (locals.user?.username !== params.username) error(403, 'Forbidden');

		const data = await request.formData();
		const category = db
			.select({ id: categories.id })
			.from(categories)
			.where(eq(categories.name, String(data.get('category'))))
			.get();
		if (!category) return fail(400, { error: 'Unknown category' });

		const raw = String(data.get('ids') ?? '');
		const ids = raw === '' ? [] : raw.split(',').map(Number);
		if (ids.some((id) => !Number.isInteger(id) || id <= 0)) {
			return fail(400, { error: 'Invalid order' });
		}

		reorderWishes(locals.user.id, category.id, ids);
		return { ok: true };
	}
};
