<script lang="ts">
	import { enhance } from '$app/forms';
	import type { Snippet } from 'svelte';

	type Option = { id: number; label: string };
	type Values = {
		name?: string;
		link?: string;
		categoryId?: number;
		amount?: number;
		currencyId?: number;
	};

	let {
		action,
		categories,
		currencies,
		values = {},
		error,
		buttons
	}: {
		action?: string;
		categories: Option[];
		currencies: Option[];
		values?: Values;
		error?: string;
		buttons: Snippet;
	} = $props();
</script>

<form
	class="form"
	method="POST"
	{action}
	use:enhance={({ action: target, cancel }) => {
		if (target.search.includes('/delete') && !confirm('Delete this wish?')) cancel();
	}}
>
	<input name="name" type="text" placeholder="Description" value={values.name ?? ''} required />
	<input name="link" type="url" placeholder="Link" value={values.link ?? ''} required />
	<select name="category" required>
		<option disabled hidden selected={!values.categoryId} value="">Category</option>
		{#each categories as category (category.id)}
			<option value={category.id} selected={category.id === values.categoryId}>
				{category.label}
			</option>
		{/each}
	</select>
	<div class="money">
		<input
			name="price"
			type="number"
			min="0"
			step="any"
			placeholder="Price"
			value={values.amount ?? ''}
			required
		/>
		<select name="currency" required>
			<option disabled hidden selected={!values.currencyId} value="">Currency</option>
			{#each currencies as currency (currency.id)}
				<option value={currency.id} selected={currency.id === values.currencyId}>
					{currency.label}
				</option>
			{/each}
		</select>
	</div>

	{#if error}
		<p class="error" aria-live="polite">{error}</p>
	{/if}

	<div class="actions">
		{@render buttons()}
	</div>
</form>

<style>
	.form {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	input,
	select {
		width: 100%;
		height: 2.875rem;
		padding: 0 14px;
		font-size: 1rem;
		color: var(--color-text);
		background-color: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
	}

	select {
		padding-right: 44px;
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%238b84a3' stroke-width='2.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M6 9l6 6 6-6'/%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 14px center;
		background-size: 14px;
	}

	select:invalid {
		color: var(--color-text-muted);
	}

	.money {
		display: grid;
		grid-template-columns: 1fr 140px;
		gap: 12px;
	}

	.error {
		color: var(--color-danger);
	}

	.actions {
		display: flex;
		flex-direction: row-reverse;
		gap: 12px;
		margin-top: 4px;
	}

	.actions :global(button) {
		width: 7.5rem;
		min-height: 2.875rem;
		padding: 5px;
		font-size: 1rem;
		color: var(--color-accent);
		cursor: pointer;
		background-color: var(--color-accent-soft);
		border: none;
		border-radius: var(--radius);
		transition: filter 0.2s ease;
	}

	.actions :global(button:hover) {
		filter: brightness(0.96);
	}

	.actions :global(button.primary) {
		color: var(--color-bg);
		background-color: var(--color-accent);
	}

	.actions :global(button.danger) {
		width: auto;
		padding-inline: 0;
		margin-right: auto;
		color: var(--color-danger);
		background-color: transparent;
	}
</style>
