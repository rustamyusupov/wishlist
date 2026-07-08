import type { PriceChange } from '$lib/prices';

export type WishListItem = {
	id: number;
	name: string;
	link: string;
	categoryId: number;
	category: string;
	sort: number;
	amount: number | null;
	symbol: string | null;
	change: PriceChange | null;
};

export type WishGroup = {
	name: string;
	categoryId: number;
	items: WishListItem[];
};

export type WishDetail = {
	id: number;
	name: string;
	link: string;
	categoryId: number;
	sort: number;
	createdAt: string;
	amount: number | null;
	currencyId: number | null;
	history: { amount: number; code: string; symbol: string; createdAt: string }[];
};

export type Category = { id: number; name: string };
export type Currency = { id: number; code: string; symbol: string };
export type Option = { id: number; label: string };
export type Session = { authenticated: boolean; needsSetup: boolean };

export type WishInput = {
	name: string;
	link: string;
	categoryId: number;
	amount: number;
	currencyId: number;
};
