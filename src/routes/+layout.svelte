<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';

	let { data, children } = $props();
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
	{#if data.user}
		<footer class="footer">
			<span>{data.user.username}</span>
			<span>·</span>
			<form method="POST" action={resolve('/logout')} use:enhance>
				<button class="logout" type="submit">Log out</button>
			</form>
		</footer>
	{/if}
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		gap: 24px;
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

	.footer {
		display: flex;
		align-items: baseline;
		justify-content: center;
		gap: 8px;
		font-size: 14px;
		color: var(--color-text-muted);
	}

	.logout {
		padding: 15px 0;
		margin: -15px 0;
		font-size: 14px;
		color: var(--color-accent);
		cursor: pointer;
		background: none;
		border: none;
		transition: color 0.2s ease;
	}

	.logout:hover {
		color: var(--color-text);
	}
</style>
