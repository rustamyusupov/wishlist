import { error, fail, redirect } from '@sveltejs/kit';
import { ApiError, apiGet, apiSend } from '$lib/server/api';
import type { Category, Currency, WishDetail } from '$lib/types';
import { parseWishInput } from '$lib/wishes';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ request, params }) => {
	const cookie = request.headers.get('cookie') ?? '';
	const [wish, categories, currencies] = await Promise.all([
		apiGet<WishDetail>(cookie, `/wishes/${params.id}`).catch((cause) => {
			if (cause instanceof ApiError && cause.status === 404) error(404, 'Wish not found');
			throw cause;
		}),
		apiGet<Category[]>(cookie, '/categories'),
		apiGet<Currency[]>(cookie, '/currencies')
	]);

	return {
		wish: {
			name: wish.name,
			link: wish.link,
			categoryId: wish.categoryId,
			amount: wish.amount ?? undefined,
			currencyId: wish.currencyId ?? undefined
		},
		categories: categories.map((category) => ({ id: category.id, label: category.name })),
		currencies: currencies.map((currency) => ({ id: currency.id, label: currency.code }))
	};
};

export const actions: Actions = {
	update: async ({ request, params }) => {
		const cookie = request.headers.get('cookie') ?? '';
		const input = parseWishInput(await request.formData());
		if (!input) return fail(400, { error: 'All fields are required' });

		const result = await apiSend(cookie, 'PUT', `/wishes/${params.id}`, input);
		if (!result.ok) {
			if (result.status === 404) error(404, 'Wish not found');
			return fail(result.status, { error: result.message });
		}

		redirect(303, '/');
	},
	delete: async ({ request, params }) => {
		const cookie = request.headers.get('cookie') ?? '';
		await apiSend(cookie, 'DELETE', `/wishes/${params.id}`);
		redirect(303, '/');
	}
};
