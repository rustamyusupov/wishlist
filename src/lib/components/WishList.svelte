<script lang="ts">
	import { dndzone, type DndEvent } from 'svelte-dnd-action';
	import { flip } from 'svelte/animate';
	import { deserialize } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { formatPrice } from '$lib/format';

	type Item = {
		id: number;
		name: string;
		link: string;
		amount: number | null;
		symbol: string | null;
	};
	type Group = { name: string; items: Item[] };

	let { groups, editable = false }: { groups: Group[]; editable?: boolean } = $props();

	const FLIP_DURATION = 150;

	let lists = $derived(groups.map((group) => ({ ...group, items: [...group.items] })));

	const reorder = (name: string, items: Item[]) => {
		lists = lists.map((group) => (group.name === name ? { ...group, items } : group));
	};

	const consider = (name: string, event: CustomEvent<DndEvent<Item>>) => {
		reorder(name, event.detail.items);
	};

	const finalize = async (name: string, event: CustomEvent<DndEvent<Item>>) => {
		reorder(name, event.detail.items);

		const before = groups.find((group) => group.name === name)?.items ?? [];
		const after = event.detail.items;
		if (after.map((item) => item.id).join(',') === before.map((item) => item.id).join(',')) {
			return;
		}

		const body = new FormData();
		body.set('category', name);
		body.set('ids', after.map((item) => item.id).join(','));
		const response = await fetch('?/reorder', {
			method: 'POST',
			headers: { 'x-sveltekit-action': 'true' },
			body
		});
		deserialize(await response.text());
		await invalidateAll();
	};
</script>

{#snippet row(wish: Item)}
	{#if editable}
		<a class="edit" href={resolve('/edit/[id]', { id: String(wish.id) })} aria-label="Edit">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960" width="18" height="18">
				<path
					d="M200-200h57l391-391-57-57-391 391v57Zm-80 80v-170l528-527q12-11 26.5-17t30.5-6q16 0 31 6t26 18l55 56q12 11 17.5 26t5.5 30q0 16-5.5 30.5T817-647L290-120H120Zm640-584-56-56 56 56Zm-141 85-28-29 57 57-29-28Z"
				/>
			</svg>
		</a>
	{/if}
	<p class="text">
		{#if editable}{wish.name}{:else}<a href={wish.link}>{wish.name}</a
			>{/if}{#if wish.amount !== null && wish.symbol !== null}<span class="price"
				>&nbsp;– {formatPrice(wish.amount, wish.symbol)}</span
			>{/if}
	</p>
{/snippet}

{#if lists.length === 0}
	<p class="empty">No wishes yet.</p>
{/if}

{#each lists as group (group.name)}
	<section class="category">
		<h2 class="heading">{group.name}</h2>
		{#if editable}
			<ul
				class="list"
				use:dndzone={{
					items: group.items,
					flipDurationMs: FLIP_DURATION,
					dropTargetStyle: {}
				}}
				onconsider={(event) => consider(group.name, event)}
				onfinalize={(event) => finalize(group.name, event)}
			>
				{#each group.items as wish (wish.id)}
					<li class="wish draggable" animate:flip={{ duration: FLIP_DURATION }}>
						{@render row(wish)}
					</li>
				{/each}
			</ul>
		{:else}
			<ul class="list">
				{#each group.items as wish (wish.id)}
					<li class="wish">
						{@render row(wish)}
					</li>
				{/each}
			</ul>
		{/if}
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

	ul.list:empty {
		min-height: 36px;
		border: 1px dashed var(--color-border);
		border-radius: var(--radius);
	}

	.wish {
		display: flex;
		align-items: flex-start;
		gap: 8px;
		padding-block: 6px;
	}

	.draggable {
		cursor: grab;
		user-select: none;
		background-color: var(--color-bg);
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
