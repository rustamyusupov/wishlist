<script lang="ts">
	import PriceChart from '$lib/components/PriceChart.svelte';
	import { formatDate, formatPrice } from '$lib/format';
	import { currencyTail } from '$lib/prices';

	let { data } = $props();

	const tail = $derived(currencyTail(data.wish.history));
	const records = $derived([...data.wish.history].reverse());
</script>

<svelte:head>
	<title>{data.wish.name} — Wishlist</title>
</svelte:head>

<h2 class="title"><a href={data.wish.link}>{data.wish.name}</a></h2>

{#if tail.length > 0}
	<PriceChart points={tail} symbol={tail[0].symbol} now={data.now} />
{:else}
	<p class="empty">No price records yet.</p>
{/if}

{#if data.wish.history.length > tail.length}
	<p class="note">The chart shows only {tail[0].symbol} prices.</p>
{/if}

{#if records.length > 0}
	<table class="records">
		<thead>
			<tr><th>Date</th><th class="amount">Price</th></tr>
		</thead>
		<tbody>
			{#each records as record, index (index)}
				<tr>
					<td>{formatDate(record.createdAt)}</td>
					<td class="amount">{formatPrice(record.amount, record.symbol)}</td>
				</tr>
			{/each}
		</tbody>
	</table>
{/if}

<style>
	.title {
		margin-bottom: 16px;
	}

	.empty,
	.note {
		color: var(--color-text-muted);
	}

	.note {
		margin-top: 12px;
		font-size: 0.75rem;
	}

	.records {
		width: 100%;
		margin-top: 24px;
		border-collapse: collapse;
	}

	.records th {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--color-text-muted);
		text-align: left;
	}

	.records th,
	.records td {
		padding: 6px 0;
		border-bottom: 1px solid var(--color-border);
	}

	.records .amount {
		font-variant-numeric: tabular-nums;
		text-align: right;
	}
</style>
