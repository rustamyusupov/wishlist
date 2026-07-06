import { fail, redirect } from '@sveltejs/kit';
import { createWish, listOptions, parseWishInput } from '$lib/server/wishes';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = () => listOptions();

export const actions: Actions = {
	default: async ({ request }) => {
		const input = parseWishInput(await request.formData());
		if (!input) return fail(400, { error: 'All fields are required' });

		createWish(input);
		redirect(303, '/');
	}
};
