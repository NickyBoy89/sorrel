import tailwindcss from "@tailwindcss/vite";
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const config: import('vite').UserConfig = {
	plugins: [
		sveltekit(),
		tailwindcss(),
	]
};

export default config;
