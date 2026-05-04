<script lang="ts">
  import type { Snippet } from "svelte";

  type DropdownOptionProps = {
    text: string;
    icon: Snippet;
  };

  let { icon, text }: DropdownOptionProps = $props();
</script>

<div class="flex flex-row items-center w-full gap-x-2 px-2 py-2">
  <div class="text-left grow">{text}</div>
  {@render icon()}
</div>
