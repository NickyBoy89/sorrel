<script lang="ts">
  import UiButton from "$lib/components/uiButton.svelte";
  import { onMount, type Component } from "svelte";
  import Check from "phosphor-svelte/lib/Check";
  import Spinner from "phosphor-svelte/lib/Spinner";
  import X from "phosphor-svelte/lib/X";
  import { APIUrl } from "../../../../constants";
  import { page } from "$app/state";

  type User = {
    id: number;
    display_name: string;
  };

  type Checkbox = Event & { currentTarget: EventTarget & HTMLInputElement };

  let menuId: number;

  let users = $state([] as Array<User>);
  let icons: Map<number, Component> = $state(new Map());

  let selected: Set<number> = new Set();

  onMount(() => {
    const rawId = page.url.searchParams.get("menu-id");
    if (rawId != null) {
      menuId = Number.parseInt(rawId);
    }
  });

  fetch(`${APIUrl}/api/users`, {})
    .then((resp) => resp.json())
    .then((resp) => (users = resp))
    .catch((error) => console.error(error));

  const handleUserChecked = (event: Checkbox) => {
    const id = event?.currentTarget?.dataset.userid;
    if (id == undefined) return;

    const userId = Number.parseInt(id);

    if (event?.currentTarget?.checked) {
      selected.add(userId);
    } else {
      selected.delete(userId);
    }
  };

  const handleSendNotifications = () => {
    const selectedUserIds = Array.from(selected);

    selectedUserIds.forEach((userId) => {
      icons.set(userId, Spinner);
    });

    fetch(`${APIUrl}/api/menu/share`, {
      method: "POST",
      body: JSON.stringify({
        menuId: menuId,
        users: selectedUserIds,
      }),
    })
      .then((resp) => resp.json())
      .then((returnStatuses: Object) => {
        console.log(returnStatuses);

        for (let [userId, success] of Object.entries(returnStatuses)) {
          icons.set(Number.parseInt(userId), success ? Check : X);
        }
      })
      .catch((error) => console.error(error));
  };
</script>

<div class="text-black dark:text-white">
  <ol>
    {#each users as user}
      <li>
        Id: {user.id}, Name: {user.display_name}<input
          type="checkbox"
          data-userid={user.id}
          onchange={handleUserChecked}
        />{#if icons.has(user.id)}{icons.get(user.id)}{/if}
      </li>
    {/each}
  </ol>

  <UiButton text="Share" action={handleSendNotifications} />
</div>
