<script lang="ts">
  import z from "zod";
  import { User } from "$lib/system";
  import { APIUrl } from "../../constants";
  import { onMount } from "svelte";
  import UiButton from "./uiButton.svelte";

  let selectedUser: number | null = $state(null);

  let users: Array<User> = $state([]);

  let {
    groceryListId,
    hidden = true,
  }: { groceryListId: number; hidden?: boolean } = $props();

  const fetchUsers = async () => {
    await fetch(`${APIUrl}/api/users`)
      .then((resp) => resp.json())
      .then((resp) => (users = z.array(User).parse(resp)));
  };

  onMount(() => {
    fetchUsers();
  });
</script>

<button
  class={[
    "fixed left-0 top-0 w-full h-full bg-neutral-500/25 items-center justify-center",
    { hidden },
  ]}
  onclick={() => (hidden = true)}
  aria-label="Close modal"
></button>
<div class={["fixed center", { hidden }]}>
  <div class="flex flex-col mx-auto p-4 space-y-2 bg-neutral-900 rounded-md">
    <div class="text-xl font-semibold">Share with</div>
    <select
      class="rounded-sm border px-2 h-8 border-neutral-700"
      onchange={(event) =>
        (selectedUser = Number.parseInt(event.currentTarget.value))}
    >
      {#each users as user}
        <option value={user.id}>{user.display_name}</option>
      {/each}
    </select>
    <UiButton
      action={() => {
        fetch(`${APIUrl}/api/v1/grocery_list/${groceryListId}/share`, {
          method: "POST",
          body: JSON.stringify([selectedUser]),
        });
      }}
      text="Share"
    />
  </div>
</div>

<style>
  .center {
    top: 50%;
    left: 50%;
  }
</style>
