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

		(async () => {
			console.log("Testing subscription...");

			const isValid = await isSubscriptionValid();

			console.log(`Valid: ${isValid}`);

			if (!isValid) {
				console.log("Resubscribing...");
				await handleSubscribe(Number.parseInt(userId as string));
			}
			
		})()
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
			--bg-color: #212121;
		} 
	}
</style>

{@render children()}