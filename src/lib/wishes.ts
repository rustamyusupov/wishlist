import type { WishGroup, WishInput, WishListItem } from '$lib/types';

export const groupByCategory = (wishes: WishListItem[]): WishGroup[] => {
	const groups: WishGroup[] = [];
	const index = new Map<number, number>();

	for (const wish of wishes) {
		const position = index.get(wish.categoryId);
		if (position === undefined) {
			index.set(wish.categoryId, groups.length);
			groups.push({ name: wish.category, categoryId: wish.categoryId, items: [wish] });
		} else {
			groups[position].items.push(wish);
		}
	}

	return groups;
};

export const parseWishInput = (data: FormData): WishInput | null => {
	const name = String(data.get('name') ?? '').trim();
	const link = String(data.get('link') ?? '').trim();
	const categoryId = Number(data.get('category'));
	const amount = Number(data.get('price'));
	const currencyId = Number(data.get('currency'));

	const valid =
		name !== '' &&
		link !== '' &&
		Number.isInteger(categoryId) &&
		categoryId > 0 &&
		Number.isFinite(amount) &&
		amount >= 0 &&
		Number.isInteger(currencyId) &&
		currencyId > 0;

	return valid ? { name, link, categoryId, amount, currencyId } : null;
};
