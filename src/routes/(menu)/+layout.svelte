<script lang="ts">
  import '../../app.css';
	let { children } = $props();
</script>

{@render children()}

<style>
  /* :global(html) {
		background-color: var(--fg-color);
		color: var(--bg-color);
	} */
</style>