export type PricePoint = {
	amount: number;
	code: string;
	symbol: string;
	createdAt: Date;
};

export type PriceChange = {
	direction: 'up' | 'down';
	percent: number;
	low: boolean;
};

export const currencyTail = <T extends { code: string }>(history: T[]): T[] => {
	if (history.length === 0) return [];

	let start = history.length - 1;
	while (start > 0 && history[start - 1].code === history[start].code) start -= 1;

	return history.slice(start);
};
