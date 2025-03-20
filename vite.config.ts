import tailwindcss from "@tailwindcss/vite";
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

const config: import('vite').UserConfig = {
	plugins: [
		sveltekit(),
		tailwindcss(),
		SvelteKitPWA({
			strategies: "injectManifest"
		})
	]
};

export default config;
