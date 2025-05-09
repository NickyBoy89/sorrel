<script lang="ts">
	import '../../../app.css';
	import { goto } from '$app/navigation';
    import { onMount } from 'svelte';
    import { handleSubscribe, isSubscriptionValid } from '$lib/notificationManager';
	let { children } = $props();
    import Keycloak, { type KeycloakInitOptions } from 'keycloak-js';
    import { bearerToken } from '../stores';
    import { keycloakClientId, keycloakRealm, keycloakUrl } from '../../../constants';

	let userId;

	let instance = {
		url: keycloakUrl,
		realm: keycloakRealm,
		clientId: keycloakClientId
	};
	const initOptions: KeycloakInitOptions = {
		onLoad: "check-sso",
	};

	let keycloakLoaded = $state(false);

	onMount(() => {
		const keycloak = new Keycloak(instance);
		keycloak.init(initOptions)
			.then((auth) => {
				if (auth) {
					console.log("Auth successful!");
				}
				bearerToken.set(keycloak.token);		
				keycloakLoaded = true;
			})
			.catch((error) => { 
				console.error("Auth not successful")
				console.error(error)
			 });
			
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