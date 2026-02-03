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
  import List from "phosphor-svelte/lib/List";
  import Oven from "phosphor-svelte/lib/Oven";
  let { children } = $props();

  let userId;

  onMount(() => {
    userId = localStorage.getItem("userId");
    if (userId === null) {
      goto("/login");
      return;
    }

    (async () => {
      const isValid = await isSubscriptionValid();

      if (!isValid) {
        await handleSubscribe(Number.parseInt(userId));
      }
    })();
  });
</script>

<link rel="stylesheet" href="/site.css" />

{@render children()}
<BottomBar
  options={[
    { icon: House, optionName: "/" as const },
    { icon: Oven, optionName: "/recipes" as const },
    { icon: List, optionName: "/grocery_lists" as const },
    { icon: User, optionName: "/profile" as const },
  ] as const}
  selected={"home"}
/>
