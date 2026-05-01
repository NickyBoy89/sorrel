<script lang="ts">
  import type { Snippet } from "svelte";

  type CardProps = {
    children: Snippet;
  };

  let { children }: CardProps = $props();
</script>

<div
  class="flex w-32 h-24 rounded-lg bg-neutral-700 items-center justify-center"
>
  {@render children()}
</div>
