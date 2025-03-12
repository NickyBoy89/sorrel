<script lang="ts">
	import '../../app.css';
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

<section class="grid grid-cols-8">
	<div class="xl:col-start-3 col-start-2 col-span-6 xl:col-span-4 h-screen">
		{@render children()}
	</div>
</section>