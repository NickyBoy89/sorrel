<script lang="ts">
  import z from "zod";
  import { APIUrl } from "../../../constants";
  import { GroceryItem, GroceryList } from "$lib/grocery_list";

  import DotsThreeVertical from "phosphor-svelte/lib/DotsThreeVertical";
  import UiButton from "../uiButton.svelte";
  import ShareModal from "../shareModal.svelte";

  let beingEdited = $state(false);
  let optionsOpen = $state(false);
  let shareModalHidden = $state(true);

  let { list, ondelete }: { list: GroceryList; ondelete?: () => void } =
    $props();

  const fetchList = async (listId: number) => {
    return fetch(`${APIUrl}/api/v1/grocery_list/${listId}`)
      .then((resp) => resp.json())
      .then((resp) => z.array(GroceryItem).parse(resp));
  };
</script>

<ShareModal hidden={shareModalHidden} groceryListId={list.id} />
<div class="rounded-md bg-white dark:bg-zinc-800 border border-zinc-700">
  <div class="flex flex-col divide-y divide-neutral-700">
    <div class="flex flex-row items-center px-4 py-4 justify-between">
      {#if beingEdited}
        <input
          type="text"
          class="font-mono text-xl font-semibold bg-neutral-900 rounded-sm"
          value={list.name}
          onchange={(
            event: Event & { currentTarget: EventTarget & HTMLInputElement },
          ) => {
            const updateList: GroceryList = {
              id: list.id,
              name: event.currentTarget.value,
            };
            fetch(`${APIUrl}/api/v1/grocery_list/${list.id}`, {
              method: "PUT",
              body: JSON.stringify(updateList),
            }).then(() => {
              console.log("This is hit");
              list.name = updateList.name;
              beingEdited = false;
            });
          }}
        />
      {:else}
        <a
          href="/grocery_lists/list?id={list.id}"
          class="flex flex-row w-full gap-x-2"
        >
          <h2 class="font-mono text-xl font-semibold">{list.name}</h2>
          <div class="grow">
            &mdash; {#await fetchList(list.id)}-{:then items}{items.length}{:catch error}{error.message}{/await}
            items
          </div>
        </a>
      {/if}
      <button
        class="text-2xl cursor-pointer"
        onclick={() => (optionsOpen = !optionsOpen)}
      >
        <DotsThreeVertical />
      </button>
    </div>
    {#if optionsOpen}
      <div class="flex flex-row justify-between p-4">
        <UiButton text="Rename" action={() => (beingEdited = !beingEdited)} />
        <UiButton
          text="Share"
          action={() => (shareModalHidden = !shareModalHidden)}
        />
        <UiButton
          text="Delete"
          color="red"
          action={() =>
            fetch(`${APIUrl}/api/v1/grocery_list/${list.id}`, {
              method: "DELETE",
            }).then(ondelete)}
        />
      </div>
    {/if}
  </div>
</div>
