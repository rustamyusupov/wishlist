const formatter = new Intl.NumberFormat('ru-RU', {
	minimumFractionDigits: 2,
	maximumFractionDigits: 2
});

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
	day: 'numeric',
	month: 'short',
	year: 'numeric'
});

export const formatPrice = (amount: number, symbol: string) =>
	`${formatter.format(amount)}\u00A0${symbol}`;

export const formatDate = (date: Date) => dateFormatter.format(date);

export const formatPercent = (value: number) => (value < 1 ? '<1%' : `${Math.round(value)}%`);
