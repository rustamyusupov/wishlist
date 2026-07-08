import { fail } from '@sveltejs/kit';
import { apiGet, apiSend } from '$lib/server/api';
import type { Category, WishListItem } from '$lib/types';
import { groupByCategory } from '$lib/wishes';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ request, locals, url }) => {
	const cookie = request.headers.get('cookie') ?? '';
	const wishes = await apiGet<WishListItem[]>(cookie, '/wishes');
	const groups = groupByCategory(wishes);

	if (locals.authenticated && url.searchParams.has('edit')) {
		const categories = await apiGet<Category[]>(cookie, '/categories');
		const present = new Set(groups.map((group) => group.categoryId));
		for (const category of categories) {
			if (!present.has(category.id)) {
				groups.push({ name: category.name, categoryId: category.id, items: [] });
			}
		}
	}

	groups.sort((a, b) => a.name.localeCompare(b.name));
	return { groups };
};

export const actions: Actions = {
	reorder: async ({ request }) => {
		const cookie = request.headers.get('cookie') ?? '';
		const data = await request.formData();
		const categoryId = Number(data.get('categoryId'));
		const raw = String(data.get('ids') ?? '');
		const ids = raw === '' ? [] : raw.split(',').map(Number);

		const invalid =
			!Number.isInteger(categoryId) ||
			categoryId <= 0 ||
			ids.some((id) => !Number.isInteger(id) || id <= 0);
		if (invalid) return fail(400, { error: 'Invalid order' });
		if (ids.length === 0) return { ok: true };

		const result = await apiSend(cookie, 'PUT', '/wishes/order', { categoryId, ids });
		if (!result.ok) return fail(result.status, { error: result.message });

		return { ok: true };
	}
};
