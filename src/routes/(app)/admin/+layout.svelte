<script lang="ts">
	import '../../../app.css';
	import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { handleSubscribe, isSubscriptionValid } from '$lib/notificationManager';
    import { initKeycloak } from '$lib/auth';
	let { children } = $props();
    
	let userId;

	onMount(() => {
		initKeycloak();

		userId = localStorage.getItem("userId");
		if (userId === null) {
			goto("/login");
			return;
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

{@render children()}