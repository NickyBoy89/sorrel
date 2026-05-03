<script lang="ts">
  import type { Snippet } from "svelte";
  import type { ClassValue } from "svelte/elements";

  type CardProps = {
    children?: Snippet;
    fixed?: boolean;
    class?: ClassValue;
  };

  let { children, class: classValue, fixed = false }: CardProps = $props();
</script>

<div
  class={[
    "flex rounded-lg bg-neutral-700 items-center justify-center",
    { "w-32": fixed, "h-24": fixed },
    classValue,
  ]}
>
  {@render children?.()}
</div>
