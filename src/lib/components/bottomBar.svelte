<script lang="ts">
  import { type Snippet } from "svelte";

  type BottomBarProps = {
    children: Snippet;
  };

  let selectedTab = 0;

  let { children }: BottomBarProps = $props();
</script>

<div
  class="bottom-bar-container justify-evenly flex flex-row py-1 gap-x-8 bg-white border border-neutral-700"
>
  {@render children()}
</div>

<style>
  .bottom-bar-container {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 1;
  }
</style>
