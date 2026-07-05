<script lang="ts">
	import { page } from '$app/state';
	import WishList from '$lib/components/WishList.svelte';

	let { data } = $props();

	const editable = $derived(data.canEdit && page.url.searchParams.has('edit'));
</script>

<svelte:head>
	<title>{data.owner} — Wishlist</title>
</svelte:head>

<WishList groups={data.groups} {editable} />
