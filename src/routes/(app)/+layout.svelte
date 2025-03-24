<script lang="ts">
	import '../../app.css';
	import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { handleSubscribe, isSubscriptionValid } from '$lib/notificationManager';
	let { children } = $props();

	let userId;

	onMount(() => {
		userId = localStorage.getItem("userId");
		if (userId === null) {
			goto("/login");
		}

		if (!isSubscriptionValid()) {
			console.log("Resubscribing...");
			handleSubscribe(Number.parseInt(userId as string));
		}
	});
</script>

<link rel="stylesheet" href="/site.css">

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

	.main {
		background-color: var(--bg-color);
		color: var(--fg-color);
	}
</style>

<section class="main grid grid-cols-8">
	<div class="xl:col-start-3 col-start-2 col-span-6 xl:col-span-4 h-screen">
		{@render children()}
	</div>
</section>