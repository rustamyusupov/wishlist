const formatter = new Intl.NumberFormat('ru-RU', {
	minimumFractionDigits: 2,
	maximumFractionDigits: 2
});

export const formatPrice = (amount: number, symbol: string) =>
	`${formatter.format(amount)}\u00A0${symbol}`;
