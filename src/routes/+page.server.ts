import { error, fail } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';
import { db } from '$lib/server/db';
import { categories } from '$lib/server/db/schema';
import { listWishGroups, reorderWishes } from '$lib/server/wishes';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals, url }) => ({
	groups: listWishGroups(locals.authenticated && url.searchParams.has('edit'))
});

export const actions: Actions = {
	reorder: async ({ request, locals }) => {
		if (!locals.authenticated) error(403, 'Forbidden');

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

		reorderWishes(category.id, ids);
		return { ok: true };
	}
};
