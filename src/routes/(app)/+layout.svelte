<script lang="ts">
  import "../../app.css";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import {
    handleSubscribe,
    isSubscriptionValid,
  } from "$lib/notificationManager";
  import BottomBar from "$lib/components/bottomBar.svelte";
  import House from "phosphor-svelte/lib/House";
  import User from "phosphor-svelte/lib/User";
  let { children } = $props();

  let userId;

  let selectedTab = $state(0);

  onMount(() => {
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
    })();
  });
</script>

<link rel="stylesheet" href="/site.css" />

{@render children()}
<BottomBar
  options={[
    { icon: House, optionName: "home" },
    { icon: User, optionName: "profile" },
  ]}
  selected={"home"}
/>
