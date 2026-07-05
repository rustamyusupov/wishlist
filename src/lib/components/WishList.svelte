<script lang="ts">
	import { resolve } from '$app/paths';
	import { formatPrice } from '$lib/format';

	type Group = {
		name: string;
		items: {
			id: number;
			name: string;
			link: string;
			amount: number | null;
			symbol: string | null;
		}[];
	};

	let { groups, editable = false }: { groups: Group[]; editable?: boolean } = $props();
</script>

{#if groups.length === 0}
	<p class="empty">No wishes yet.</p>
{/if}

{#each groups as group (group.name)}
	<section class="category">
		<h2 class="heading">{group.name}</h2>
		<ul class="list">
			{#each group.items as wish (wish.id)}
				<li class="wish">
					{#if editable}
						<a class="edit" href={resolve('/edit/[id]', { id: String(wish.id) })} aria-label="Edit">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 -960 960 960"
								width="18"
								height="18"
							>
								<path
									d="M200-200h57l391-391-57-57-391 391v57Zm-80 80v-170l528-527q12-11 26.5-17t30.5-6q16 0 31 6t26 18l55 56q12 11 17.5 26t5.5 30q0 16-5.5 30.5T817-647L290-120H120Zm640-584-56-56 56 56Zm-141 85-28-29 57 57-29-28Z"
								/>
							</svg>
						</a>
					{/if}
					<p class="text">
						<a href={wish.link}>{wish.name}</a
						>{#if wish.amount !== null && wish.symbol !== null}<span class="price"
								>&nbsp;– {formatPrice(wish.amount, wish.symbol)}</span
							>{/if}
					</p>
				</li>
			{/each}
		</ul>
	</section>
{/each}

<style>
	.empty {
		color: var(--color-text-muted);
	}

	.category {
		margin-bottom: 24px;
	}

	.heading {
		margin-bottom: 8px;
	}

	.list {
		display: flex;
		flex-direction: column;
	}

	.wish {
		display: flex;
		align-items: flex-start;
		gap: 8px;
		padding-block: 6px;
	}

	.edit {
		display: grid;
		place-items: center;
		flex-shrink: 0;
		width: 28px;
		height: 28px;
		margin-block: -2px;
		fill: var(--color-text-muted);
		transition: fill 0.2s ease;
	}

	.edit:hover {
		fill: var(--color-accent);
	}

	.text {
		overflow-wrap: anywhere;
	}

	.price {
		color: var(--color-text);
		font-variant-numeric: tabular-nums;
	}
</style>
