import { error, fail, redirect } from '@sveltejs/kit';
import { deleteWish, getWish, listOptions, parseWishInput, updateWish } from '$lib/server/wishes';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = ({ params, locals }) => {
	const wish = getWish(Number(params.id), locals.user!.id);
	if (!wish) error(404, 'Wish not found');

	return { wish, ...listOptions() };
};

export const actions: Actions = {
	update: async ({ request, params, locals }) => {
		const input = parseWishInput(await request.formData());
		if (!input) return fail(400, { error: 'All fields are required' });

		const updated = updateWish(Number(params.id), locals.user!.id, input);
		if (!updated) error(404, 'Wish not found');

		redirect(303, '/');
	},
	delete: ({ params, locals }) => {
		deleteWish(Number(params.id), locals.user!.id);
		redirect(303, '/');
	}
};
