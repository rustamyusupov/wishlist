<script lang="ts">
	import { page } from '$app/state';
	import WishList from '$lib/components/WishList.svelte';

	let { data } = $props();

	const editable = $derived(data.authenticated && page.url.searchParams.has('edit'));
</script>

<WishList groups={data.groups} {editable} />
