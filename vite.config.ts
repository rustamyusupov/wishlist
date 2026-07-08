import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');
	const target = env.API_PROXY || 'http://localhost:3100';
	const isLocalApi = target.startsWith('http://localhost');

	return {
		plugins: [
			sveltekit({
				compilerOptions: {
					runes: ({ filename }) =>
						filename.split(/[/\\]/).includes('node_modules') ? undefined : true
				},
				adapter: adapter()
			})
		],
		server: {
			port: 3000,
			host: true,
			proxy: {
				'/api': {
					target,
					changeOrigin: true,
					rewrite: isLocalApi ? (path) => path.replace(/^\/api/, '') : undefined
				}
			}
		},
		preview: {
			port: 3000,
			host: true
		}
	};
});
