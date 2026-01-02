<script lang="ts">
  import { onMount } from "svelte";
  import { APIUrl } from "../../../constants";
  import Card from "$lib/components/ui/card.svelte";
  import UiButton from "$lib/components/uiButton.svelte";

  let lists = $state([]);

  const fetchLists = async () => {
    await fetch(`${APIUrl}/api/v1/grocery_list`)
      .then((resp) => resp.json())
      .then((resp) => (lists = resp));
  };

  onMount(() => {
    fetchLists();
  });
</script>

<div class="flex flex-col mx-4">
  {#each lists as list}
    <a href="/grocery_lists/list?id={list}">
      <Card>
        <div class="flex flex-row gap-x-2">
          <h2 class="font-mono text-xl font-semibold">{list}</h2>
          <div>&mdash; 32 items</div>
        </div>
      </Card>
    </a>
  {/each}
  <UiButton
    text="Add new"
    action={() => {
      fetch(`${APIUrl}/api/v1/grocery_list`, { method: "POST" });
      fetchLists();
    }}
  />
</div>
