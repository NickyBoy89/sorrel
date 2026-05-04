<script lang="ts">
  import DotsThreeVertical from "phosphor-svelte/lib/DotsThreeVertical";
  import type { Snippet } from "svelte";

  type DropdownProps = {
    expanded?: boolean;
    children?: Snippet;
  };

  let { expanded = $bindable(false), children }: DropdownProps = $props();
</script>

<div class="relative flex flex-col justify-center items-start">
  <button
    class="cursor-pointer hover:bg-neutral-800 flex justify-end items-center px-2 h-full rounded-lg"
    onclick={() => (expanded = !expanded)}
  >
    <DotsThreeVertical class="text-xl" />
  </button>
  <div
    class={[
      "absolute flex flex-col right-0 top-16 w-48 shadow-lg rounded-lg bg-neutral-700 divide-y-1 divide-neutral-800",
      { hidden: !expanded },
    ]}
  >
    {@render children?.()}
  </div>
</div>
