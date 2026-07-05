<script lang="ts">
	import WishForm from '$lib/components/WishForm.svelte';

	let { data, form } = $props();
</script>

<svelte:head>
	<title>Edit wish — Wishlist</title>
</svelte:head>

<h2 class="title">Edit wish</h2>

<WishForm
	action="?/update"
	categories={data.categories}
	currencies={data.currencies}
	values={data.wish}
	error={form?.error}
>
	{#snippet buttons()}
		<button class="primary" type="submit">Save</button>
		<button class="danger" type="submit" formaction="?/delete" formnovalidate>Delete</button>
	{/snippet}
</WishForm>

<style>
	.title {
		margin-bottom: 16px;
	}
</style>
