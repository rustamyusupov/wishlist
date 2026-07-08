<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';

	let { data, children } = $props();

	const logout = async () => {
		await fetch('/api/auth/logout', { method: 'POST' });
		await goto(resolve('/login'), { invalidateAll: true });
	};
</script>

<svelte:head>
	<title>Wishlist</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="app">
	<header class="header">
		<h1 class="brand"><a href={resolve('/')}>Wishlist</a></h1>
		{#if data.authenticated && page.route.id === '/'}
			<div class="controls">
				{#if page.url.searchParams.has('edit')}
					<a class="mode" href={resolve('/')} aria-label="Done editing">
						<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960" width="24" height="24">
							<path d="M382-240 154-468l57-57 171 171 367-367 57 57-424 424Z" />
						</svg>
					</a>
				{:else}
					<form method="GET" action={resolve('/')}>
						<input type="hidden" name="edit" value="" />
						<button class="mode" type="submit" aria-label="Edit list">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 -960 960 960"
								width="20"
								height="20"
							>
								<path
									d="M200-200h57l391-391-57-57-391 391v57Zm-80 80v-170l528-527q12-11 26.5-17t30.5-6q16 0 31 6t26 18l55 56q12 11 17.5 26t5.5 30q0 16-5.5 30.5T817-647L290-120H120Zm640-584-56-56 56 56Zm-141 85-28-29 57 57-29-28Z"
								/>
							</svg>
						</button>
					</form>
				{/if}
				<a class="add" href={resolve('/new')} aria-label="New wish">
					<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960" width="24" height="24">
						<path d="M440-440H200v-80h240v-240h80v240h240v80H520v240h-80v-240Z" />
					</svg>
				</a>
			</div>
		{/if}
	</header>
	<main class="main">
		{@render children()}
	</main>
	{#if data.authenticated}
		<footer class="footer">
			<button class="logout" type="button" onclick={logout}>Log out</button>
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

	.controls {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.mode,
	.add {
		display: grid;
		place-items: center;
		width: 44px;
		height: 44px;
		padding: 0;
		margin: -10px;
		cursor: pointer;
		background: none;
		border: none;
		fill: var(--color-text-muted);
		transition: fill 0.2s ease;
	}

	.mode:hover,
	.add:hover {
		fill: var(--color-accent);
	}

	.main {
		flex: 1;
	}

	.footer {
		display: flex;
		justify-content: center;
		font-size: 0.875rem;
	}

	.logout {
		padding: 15px 0;
		margin: -15px 0;
		font-size: 0.875rem;
		color: var(--color-accent);
		cursor: pointer;
		background: none;
		border: none;
		transition: color 0.2s ease;
	}

	.logout:hover {
		color: var(--color-accent-hover);
	}
</style>
