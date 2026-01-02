<script lang="ts">
  import z from "zod";
  import { onMount } from "svelte";
  import { APIUrl } from "../../../constants";
  import UiButton from "$lib/components/uiButton.svelte";
  import { GroceryList } from "$lib/grocery_list";
  import GroceryListRow from "$lib/components/groceryList/groceryListRow.svelte";

  let lists: Array<GroceryList> = $state([]);

  const fetchLists = async () => {
    await fetch(`${APIUrl}/api/v1/grocery_list`)
      .then((resp) => resp.json())
      .then((resp) => (lists = z.array(GroceryList).parse(resp)));
  };

  onMount(() => {
    fetchLists();
  });
</script>

<div class="flex flex-col m-4 gap-y-4">
  {#each lists as list}
    <GroceryListRow {list} ondelete={() => fetchLists()} />
  {/each}
  <UiButton
    text="Add new"
    action={() => {
      fetch(`${APIUrl}/api/v1/grocery_list`, { method: "POST" });
      fetchLists();
    }}
  />
</div>
