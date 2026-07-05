<script lang="ts">
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>Wishlist</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="app">
	<header class="header">
		<h1 class="brand"><a href={resolve('/')}>Wishlist</a></h1>
	</header>
	<main class="main">
		{@render children()}
	</main>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		max-width: 640px;
		min-height: 100vh;
		padding: 16px;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.brand a {
		color: var(--color-text);
	}

	.brand a:hover {
		color: var(--color-accent);
		text-decoration: none;
	}

	.main {
		flex: 1;
	}
</style>
