<script lang="ts">
	import { formatDate, formatPrice } from '$lib/format';

	type Point = { amount: number; createdAt: Date };

	let { points, symbol, now }: { points: Point[]; symbol: string; now: Date } = $props();

	const HEIGHT = 240;
	const MARGIN = { top: 12, right: 12, bottom: 28 };
	const DAY_MS = 24 * 60 * 60 * 1000;

	const monthFormatter = new Intl.DateTimeFormat('en-GB', { month: 'short' });

	let svg = $state<SVGSVGElement>();
	let width = $state(0);
	let active = $state<number | null>(null);

	const w = $derived(width || 640);

	const niceStep = (rough: number) => {
		const power = 10 ** Math.floor(Math.log10(rough));
		const unit = rough / power;
		return (unit >= 7.07 ? 10 : unit >= 3.16 ? 5 : unit >= 1.41 ? 2 : 1) * power;
	};

	const xDomain = $derived.by(() => {
		const start = points[0].createdAt.getTime();
		const end = now.getTime();
		return { start, end: end > start ? end : start + DAY_MS };
	});

	const yScale = $derived.by(() => {
		const amounts = points.map((point) => point.amount);
		const low = Math.min(...amounts);
		const high = Math.max(...amounts);
		const step = niceStep((high - low || high * 0.2 || 1) / 4);
		const min = Math.floor(low / step) * step;
		const count = Math.round((Math.ceil(high / step) * step - min) / step);
		const ticks = Array.from({ length: count + 1 }, (_, index) => min + index * step);
		const decimals = step >= 1 ? 0 : Math.ceil(-Math.log10(step));
		const formatter = new Intl.NumberFormat('ru-RU', {
			minimumFractionDigits: decimals,
			maximumFractionDigits: decimals
		});
		return {
			min,
			max: min + count * step,
			ticks,
			labels: ticks.map((tick) => formatter.format(tick))
		};
	});

	const marginLeft = $derived(
		16 + Math.max(...yScale.labels.map((label) => label.length), 2) * 7.5
	);
	const plotWidth = $derived(w - marginLeft - MARGIN.right);
	const plotHeight = HEIGHT - MARGIN.top - MARGIN.bottom;

	const x = (time: number) =>
		marginLeft + ((time - xDomain.start) / (xDomain.end - xDomain.start)) * plotWidth;
	const y = (value: number) =>
		MARGIN.top + (1 - (value - yScale.min) / (yScale.max - yScale.min)) * plotHeight;

	const monthTicks = $derived.by(() => {
		const start = new Date(xDomain.start);
		const months = [];
		for (
			let month = new Date(start.getFullYear(), start.getMonth() + 1, 1);
			month.getTime() <= xDomain.end;
			month = new Date(month.getFullYear(), month.getMonth() + 1, 1)
		) {
			months.push(month);
		}
		const step = Math.ceil(months.length / 6) || 1;
		return months.filter((_, index) => (months.length - 1 - index) % step === 0);
	});

	const monthLabel = (month: Date, index: number) => {
		const name = monthFormatter.format(month);
		if (index > 0 && month.getMonth() !== 0) return name;
		return `${name} ’${String(month.getFullYear() % 100).padStart(2, '0')}`;
	};

	const path = $derived.by(() => {
		let d = `M ${x(points[0].createdAt.getTime())} ${y(points[0].amount)}`;
		for (const point of points.slice(1)) {
			d += ` H ${x(point.createdAt.getTime())} V ${y(point.amount)}`;
		}
		return `${d} H ${x(xDomain.end)}`;
	});

	const locate = (event: PointerEvent) => {
		if (!svg) return;
		const pointer = event.clientX - svg.getBoundingClientRect().left;
		let nearest = 0;
		for (let index = 1; index < points.length; index += 1) {
			const current = Math.abs(x(points[index].createdAt.getTime()) - pointer);
			const best = Math.abs(x(points[nearest].createdAt.getTime()) - pointer);
			if (current < best) nearest = index;
		}
		active = nearest;
	};

	const tooltipShift = $derived.by(() => {
		if (active === null) return '-50%';
		const position = x(points[active].createdAt.getTime());
		if (position < 90) return '0%';
		if (position > w - 90) return '-100%';
		return '-50%';
	});
</script>

<div class="chart" bind:clientWidth={width}>
	<svg
		bind:this={svg}
		role="group"
		aria-label="Price history chart"
		width={w}
		height={HEIGHT}
		onpointermove={locate}
		onpointerdown={locate}
		onpointerleave={() => (active = null)}
	>
		{#each yScale.ticks as tick, index (tick)}
			<line class="grid" x1={marginLeft} x2={w - MARGIN.right} y1={y(tick)} y2={y(tick)} />
			<text class="tick" x={marginLeft - 8} y={y(tick)} text-anchor="end" dy="0.35em"
				>{yScale.labels[index]}</text
			>
		{/each}
		{#each monthTicks as month, index (month.getTime())}
			<text class="tick" x={x(month.getTime())} y={HEIGHT - 8} text-anchor="middle"
				>{monthLabel(month, index)}</text
			>
		{/each}
		{#if active !== null}
			<line
				class="crosshair"
				x1={x(points[active].createdAt.getTime())}
				x2={x(points[active].createdAt.getTime())}
				y1={MARGIN.top}
				y2={MARGIN.top + plotHeight}
			/>
		{/if}
		<path class="line" d={path} />
		{#each points as point, index (index)}
			<!-- focus shows the tooltip; the table below carries the same values -->
			<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
			<circle
				class="dot"
				role="img"
				tabindex="0"
				aria-label="{formatPrice(point.amount, symbol)} — {formatDate(point.createdAt)}"
				cx={x(point.createdAt.getTime())}
				cy={y(point.amount)}
				r={active === index ? 5 : 4}
				onfocus={() => (active = index)}
				onblur={() => (active = null)}
			/>
		{/each}
	</svg>
	{#if active !== null}
		<div
			class="tooltip"
			style:left="{x(points[active].createdAt.getTime())}px"
			style:top="{y(points[active].amount) - 12}px"
			style:transform="translate({tooltipShift}, -100%)"
		>
			<strong>{formatPrice(points[active].amount, symbol)}</strong>
			<span>{formatDate(points[active].createdAt)}</span>
		</div>
	{/if}
</div>

<style>
	.chart {
		position: relative;
	}

	svg {
		touch-action: pan-y;
	}

	.grid {
		stroke: var(--color-border);
		stroke-width: 1;
	}

	.crosshair {
		stroke: var(--color-text-muted);
		stroke-width: 1;
	}

	.tick {
		font-size: 0.75rem;
		font-variant-numeric: tabular-nums;
		fill: var(--color-text-muted);
	}

	.line {
		fill: none;
		stroke: var(--color-accent);
		stroke-width: 2;
		stroke-linecap: round;
		stroke-linejoin: round;
	}

	.dot {
		fill: var(--color-accent);
		stroke: var(--color-bg);
		stroke-width: 2;
	}

	.tooltip {
		position: absolute;
		display: flex;
		flex-direction: column;
		padding: 6px 10px;
		pointer-events: none;
		white-space: nowrap;
		background-color: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
	}

	.tooltip strong {
		font-size: 0.9375rem;
		font-variant-numeric: tabular-nums;
	}

	.tooltip span {
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}
</style>
