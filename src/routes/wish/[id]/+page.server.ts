import { error } from '@sveltejs/kit';
import { ApiError, apiGet } from '$lib/server/api';
import type { WishDetail } from '$lib/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ request, params }) => {
	const cookie = request.headers.get('cookie') ?? '';
	const wish = await apiGet<WishDetail>(cookie, `/wishes/${params.id}`).catch((cause) => {
		if (cause instanceof ApiError && cause.status === 404) error(404, 'Wish not found');
		throw cause;
	});

	const history = wish.history.map((point) => ({ ...point, createdAt: new Date(point.createdAt) }));
	return { wish: { ...wish, history }, now: new Date() };
};
