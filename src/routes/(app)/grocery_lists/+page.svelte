<script lang="ts">
  import z from "zod";
  import { onMount } from "svelte";
  import { APIUrl } from "../../../constants";
  import Card from "$lib/components/ui/card.svelte";
  import UiButton from "$lib/components/uiButton.svelte";
  import { GroceryItem, GroceryList } from "$lib/grocery_list";

  let lists: Array<GroceryList> = $state([]);

  const fetchLists = async () => {
    await fetch(`${APIUrl}/api/v1/grocery_list`)
      .then((resp) => resp.json())
      .then((resp) => (lists = z.array(GroceryList).parse(resp)));
  };

  const fetchList = async (listId: number) => {
    return fetch(`${APIUrl}/api/v1/grocery_list/${listId}`)
      .then((resp) => resp.json())
      .then((resp) => z.array(GroceryItem).parse(resp));
  };

  onMount(() => {
    fetchLists();
  });
</script>

<div class="flex flex-col m-4 gap-y-4">
  {#each lists as list}
    <a href="/grocery_lists/list?id={list.id}">
      <Card>
        <div class="flex flex-row gap-x-2">
          <h2 class="font-mono text-xl font-semibold">{list.name}</h2>
          <div>
            &mdash; {#await fetchList(list.id)}-{:then items}{items.length}{:catch error}{error.message}{/await}
            items
          </div>
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
