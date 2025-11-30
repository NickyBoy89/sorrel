<script lang="ts">
  import { type Component } from "svelte";

  type BottomBarOption = {
    icon: Component;
    optionName: string;
  };

  type BottomBarProps = {
    options: readonly BottomBarOption[];
    selected: (typeof options)[number]["optionName"];
  };

  let { options, selected = $bindable() }: BottomBarProps = $props();
</script>

<div
  class="fixed-bar justify-evenly flex flex-row py-1 gap-x-8 dark:bg-neutral-900 bg-white border border-neutral-700"
>
  {#each options as option}
    {@const OptionIcon = option.icon}
    <OptionIcon
      size={32}
      weight={selected === option.optionName ? "fill" : "regular"}
      class="cursor-pointer"
      onclick={() => (selected = option.optionName)}
    />
  {/each}
</div>

<style>
  .fixed-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 1;
  }
</style>
