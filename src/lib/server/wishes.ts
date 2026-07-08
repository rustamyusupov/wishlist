import { asc, desc, eq, sql } from 'drizzle-orm';
import { priceChange } from '$lib/prices';
import { db } from '$lib/server/db';
import { categories, currencies, prices, wishes } from '$lib/server/db/schema';

export type WishInput = {
	name: string;
	link: string;
	categoryId: number;
	amount: number;
	currencyId: number;
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

export const listWishGroups = (includeEmpty = false) => {
	const rows = db
		.select({
			id: wishes.id,
			name: wishes.name,
			link: wishes.link,
			category: categories.name,
			amount: prices.amount,
			symbol: currencies.symbol
		})
		.from(wishes)
		.innerJoin(categories, eq(categories.id, wishes.categoryId))
		.leftJoin(
			prices,
			eq(
				prices.id,
				sql`(select id from prices where wish_id = ${wishes.id} order by created_at desc, id desc limit 1)`
			)
		)
		.leftJoin(currencies, eq(currencies.id, prices.currencyId))
		.orderBy(asc(categories.name), asc(wishes.sort), asc(wishes.name))
		.all();

	const priceRows = db
		.select({
			wishId: prices.wishId,
			amount: prices.amount,
			code: currencies.code,
			createdAt: prices.createdAt
		})
		.from(prices)
		.innerJoin(currencies, eq(currencies.id, prices.currencyId))
		.orderBy(asc(prices.createdAt), asc(prices.id))
		.all();

	const histories = new Map<number, typeof priceRows>();
	for (const price of priceRows) {
		const history = histories.get(price.wishId);
		if (history) history.push(price);
		else histories.set(price.wishId, [price]);
	}

	const now = new Date();
	const items = rows.map((row) => ({
		...row,
		change: priceChange(histories.get(row.id) ?? [], now)
	}));

	const grouped = new Map<string, typeof items>();
	if (includeEmpty) {
		const names = db
			.select({ name: categories.name })
			.from(categories)
			.orderBy(asc(categories.name))
			.all();
		for (const { name } of names) grouped.set(name, []);
	}
	for (const item of items) {
		const group = grouped.get(item.category);
		if (group) group.push(item);
		else grouped.set(item.category, [item]);
	}

	return [...grouped].map(([name, items]) => ({ name, items }));
};

export const listOptions = () => ({
	categories: db
		.select({ id: categories.id, label: categories.name })
		.from(categories)
		.orderBy(asc(categories.name))
		.all(),
	currencies: db
		.select({ id: currencies.id, label: currencies.code })
		.from(currencies)
		.orderBy(asc(currencies.id))
		.all()
});

export const getWish = (id: number) => {
	const wish = db.select().from(wishes).where(eq(wishes.id, id)).get();
	if (!wish) return null;

	const price = db
		.select({ amount: prices.amount, currencyId: prices.currencyId })
		.from(prices)
		.where(eq(prices.wishId, id))
		.orderBy(desc(prices.createdAt), desc(prices.id))
		.limit(1)
		.get();

	return {
		...wish,
		amount: price === undefined ? undefined : Math.round(price.amount * 100) / 100,
		currencyId: price?.currencyId
	};
};

export const getWishHistory = (id: number) => {
	const wish = db.select().from(wishes).where(eq(wishes.id, id)).get();
	if (!wish) return null;

	const history = db
		.select({
			amount: prices.amount,
			code: currencies.code,
			symbol: currencies.symbol,
			createdAt: prices.createdAt
		})
		.from(prices)
		.innerJoin(currencies, eq(currencies.id, prices.currencyId))
		.where(eq(prices.wishId, id))
		.orderBy(asc(prices.createdAt), asc(prices.id))
		.all();

	return { ...wish, history };
};

export const createWish = (input: WishInput) =>
	db.transaction((tx) => {
		const { sort } = tx
			.select({ sort: sql<number>`coalesce(max(${wishes.sort}), -1) + 1` })
			.from(wishes)
			.where(eq(wishes.categoryId, input.categoryId))
			.get()!;

		const { id } = tx
			.insert(wishes)
			.values({
				categoryId: input.categoryId,
				name: input.name,
				link: input.link,
				sort
			})
			.returning({ id: wishes.id })
			.get();

		tx.insert(prices)
			.values({ wishId: id, amount: input.amount, currencyId: input.currencyId })
			.run();
	});

export const updateWish = (id: number, input: WishInput) =>
	db.transaction((tx) => {
		const current = getWish(id);
		if (!current) return false;

		tx.update(wishes)
			.set({ name: input.name, link: input.link, categoryId: input.categoryId })
			.where(eq(wishes.id, id))
			.run();

		if (current.amount !== input.amount || current.currencyId !== input.currencyId) {
			tx.insert(prices)
				.values({ wishId: id, amount: input.amount, currencyId: input.currencyId })
				.run();
		}

		return true;
	});

export const reorderWishes = (categoryId: number, ids: number[]) =>
	db.transaction((tx) => {
		ids.forEach((id, index) => {
			tx.update(wishes).set({ categoryId, sort: index }).where(eq(wishes.id, id)).run();
		});
	});

export const deleteWish = (id: number) => db.delete(wishes).where(eq(wishes.id, id)).run();
