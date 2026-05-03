<script lang="ts">
  let {
    placeholder = "",
    initialValue = $bindable(""),
    onchange,
  }: {
    placeholder?: string;
    initialValue?: string;
    onchange: (
      event: Event & { currentTarget: EventTarget & HTMLInputElement },
    ) => void;
  } = $props();
</script>

<input
  type="text"
  class="block rounded-sm pl-2 w-auto h-8 text-neutral-900 dark:text-white bg-white dark:bg-neutral-900 border border-neutral-700"
  {placeholder}
  bind:value={initialValue}
  {onchange}
/>
