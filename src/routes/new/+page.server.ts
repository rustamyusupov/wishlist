import { fail, redirect } from '@sveltejs/kit';
import { apiGet, apiSend } from '$lib/server/api';
import type { Category, Currency } from '$lib/types';
import { parseWishInput } from '$lib/wishes';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ request }) => {
	const cookie = request.headers.get('cookie') ?? '';
	const [categories, currencies] = await Promise.all([
		apiGet<Category[]>(cookie, '/categories'),
		apiGet<Currency[]>(cookie, '/currencies')
	]);

	return {
		categories: categories.map((category) => ({ id: category.id, label: category.name })),
		currencies: currencies.map((currency) => ({ id: currency.id, label: currency.code }))
	};
};

export const actions: Actions = {
	default: async ({ request }) => {
		const cookie = request.headers.get('cookie') ?? '';
		const input = parseWishInput(await request.formData());
		if (!input) return fail(400, { error: 'All fields are required' });

		const result = await apiSend(cookie, 'POST', '/wishes', input);
		if (!result.ok) return fail(result.status, { error: result.message });

		redirect(303, '/');
	}
};
