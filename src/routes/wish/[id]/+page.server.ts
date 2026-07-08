import { error } from '@sveltejs/kit';
import { getWishHistory } from '$lib/server/wishes';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ params }) => {
	const wish = getWishHistory(Number(params.id));
	if (!wish) error(404, 'Wish not found');

	return { wish, now: new Date() };
};
