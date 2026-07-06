<script lang="ts">
	import WishForm from '$lib/components/WishForm.svelte';

	let { data, form } = $props();
</script>

<svelte:head>
	<title>New wish — Wishlist</title>
</svelte:head>

<h2 class="title">New wish</h2>

<WishForm categories={data.categories} currencies={data.currencies} error={form?.error}>
	{#snippet buttons()}
		<button class="primary" type="submit">Add</button>
	{/snippet}
</WishForm>

<style>
	.title {
		margin-bottom: 16px;
	}
</style>
