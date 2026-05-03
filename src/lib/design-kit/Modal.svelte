<script lang="ts">
  import type { Snippet } from "svelte";

  type ImportModalProps = {
    visible?: boolean;
    children?: Snippet;
  };

  let { visible = $bindable(false), children }: ImportModalProps = $props();

  const handleClose = () => {
    visible = false;
  };

  const handleKeyDown = (event: KeyboardEvent) => {
    if (event.key === "Escape") {
      handleClose();
    }
  };
</script>

<button
  class={[
    "fixed left-0 top-0 w-full h-full bg-neutral-500/25 flex items-center justify-center",
    { hidden: !visible },
  ]}
  onclick={handleClose}
  onkeydown={handleKeyDown}
  aria-label="Close modal"
>
</button>

<div class={["fixed center", { hidden: !visible }]}>
  {@render children?.()}
</div>

<style>
  .center {
    top: 50%;
    left: 50%;
  }
</style>
