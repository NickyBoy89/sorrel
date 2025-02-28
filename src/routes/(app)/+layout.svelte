<script lang="ts">
	import '../app.css';
	import Navbar from '$lib/components/navbar.svelte';
	let { children } = $props();
</script>

<style>
	:root {
		--bg-color: #ffffff;
		--fg-color: #000000;
	}

	@media (prefers-color-scheme: dark) {
		:root {
			--fg-color: #ffffff;
			--bg-color: #1d2021;
		} 
	}

	:global(html) {
		background-color: var(--bg-color);
		color: var(--fg-color);
	}
</style>

<Navbar />



{@render children()}