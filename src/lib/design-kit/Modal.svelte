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
    "fixed left-0 top-0 w-full h-full bg-neutral-500/25 flex items-center justify-center z-50",
    { hidden: !visible },
  ]}
  onclick={handleClose}
  onkeydown={handleKeyDown}
  aria-label="Close modal"
></button>

<div class={["fixed center", { hidden: !visible }]}>
  <div
    class="flex flex-col mx-auto p-4 space-y-2 bg-neutral-900 rounded-md shadow-lg"
  >
    {@render children?.()}
  </div>
</div>

<style>
  .center {
    top: 50%;
    left: 50%;
  }
</style>
