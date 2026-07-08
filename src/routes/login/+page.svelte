<script lang="ts">
	import { startAuthentication, startRegistration } from '@simplewebauthn/browser';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let { data } = $props();
	let setupToken = $state('');
	let message = $state('');

	const post = async (path: string, body: unknown) => {
		const response = await fetch(path, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify(body)
		});
		if (!response.ok) {
			const details = await response.json().catch(() => null);
			throw new Error(details?.message ?? 'Request failed');
		}
		return response.json();
	};

	const report = (cause: unknown, fallback: string) => {
		if (cause instanceof Error && cause.name === 'NotAllowedError') return;
		message = cause instanceof Error ? cause.message : fallback;
	};

	const login = async () => {
		message = '';
		try {
			const options = await post('/api/auth/login', {});
			const assertion = await startAuthentication({ optionsJSON: options });
			await post('/api/auth/login', assertion);
			await goto(resolve('/'), { invalidateAll: true });
		} catch (cause) {
			report(cause, 'Sign-in failed');
		}
	};

	const register = async () => {
		message = '';
		try {
			const options = await post('/api/auth/register', { token: setupToken });
			const registration = await startRegistration({ optionsJSON: options });
			await post('/api/auth/register', { token: setupToken, response: registration });
			await goto(resolve('/'), { invalidateAll: true });
		} catch (cause) {
			report(cause, 'Registration failed');
		}
	};
</script>

<div class="login">
	<button class="primary" onclick={login}>Sign in with a passkey</button>

	{#if data.needsSetup}
		<div class="setup">
			<p class="hint">First run: register the owner passkey with the setup token.</p>
			<input type="password" placeholder="Setup token" autocomplete="off" bind:value={setupToken} />
			<button onclick={register} disabled={!setupToken}>Register this device</button>
		</div>
	{/if}

	{#if message}
		<p class="error" aria-live="polite">{message}</p>
	{/if}
</div>

<style>
	.login {
		display: flex;
		flex-direction: column;
		gap: 16px;
		max-width: 320px;
	}

	@media (max-width: 479px) {
		.login {
			max-width: none;
		}
	}

	button {
		min-height: 2.875rem;
		padding: 5px 15px;
		font-size: 1rem;
		color: var(--color-accent);
		cursor: pointer;
		background-color: var(--color-accent-soft);
		border: none;
		border-radius: var(--radius);
		transition: filter 0.2s ease;
	}

	button:hover:enabled {
		filter: brightness(0.96);
	}

	button:disabled {
		color: var(--color-text-muted);
		cursor: not-allowed;
	}

	button.primary {
		color: var(--color-bg);
		background-color: var(--color-accent);
	}

	.setup {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding-top: 16px;
		border-top: 1px solid var(--color-border);
	}

	.hint {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	input {
		height: 2.875rem;
		padding: 0 14px;
		font-size: 1rem;
		color: var(--color-text);
		background-color: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
	}

	.error {
		color: var(--color-danger);
	}
</style>
