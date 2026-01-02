<script lang="ts">
  import DotsThreeVertical from "phosphor-svelte/lib/DotsThreeVertical";
  import Backspace from "phosphor-svelte/lib/Backspace";
  import PencilSimple from "phosphor-svelte/lib/PencilSimple";

  import type { GroceryItem } from "$lib/grocery_list";

  let optionsOpen = $state(false);

  let {
    ingredient,
    checked = $bindable(false),
    ondelete,
  }: {
    ingredient: GroceryItem;
    checked?: boolean;
    ondelete?: () => void;
  } = $props();
</script>

<div class="flex flex-row items-center py-4">
  <label class="flex flex-row grow space-x-2">
    <input
      type="checkbox"
      id="checkbox"
      class={[
        "h-6 w-6 border border-neutral-700 rounded-full cursor-pointer",
        "circle-checkbox",
        { checked },
      ]}
      bind:checked
    />
    <div class="flex grocery-item-text">
      {ingredient.quantity}
      {ingredient.name}
    </div>
  </label>
  <div class="flex flex-col items-end">
    <button onclick={() => (optionsOpen = !optionsOpen)}>
      <DotsThreeVertical class="text-xl" />
    </button>
    {#if optionsOpen}
      <div
        class="flex flex-col bg-neutral-800 rounded-md divide-y-1 divide-neutral-700"
      >
        <div class="flex flex-row items-center gap-x-2 px-4 py-2">
          <div class="text-left grow">Edit Item</div>
          <PencilSimple />
        </div>
        <button
          class="flex flex-row text-red-500 cursor-pointer items-center gap-x-2 px-4 py-2"
          onclick={() => {
            ondelete?.();
            optionsOpen = false;
          }}
        >
          <div class="text-left grow">Delete Item</div>
          <Backspace />
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
  .circle-checkbox {
    appearance: none;
    position: relative;
    transition:
      background 0.1s,
      border 0.1s;
  }

  .circle-checkbox:checked {
    background: #d65d0e;
    border: 2px solid #d65d0e;
  }

  .circle-checkbox:checked::after {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0.35em;
    height: 0.65em;
    border: solid white;
    border-width: 0 0.15em 0.15em 0;
    transform: translate(-50%, -55%) rotate(45deg);
  }

  .grocery-item-text {
    transition: color 0.1s;
  }

  .circle-checkbox:checked ~ .grocery-item-text {
    color: var(--color-neutral-500);
  }
</style>
