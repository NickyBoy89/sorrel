<script lang="ts">
  import CaretLeft from "phosphor-svelte/lib/CaretLeft";

  let {
    text,
    href,
  }: {
    text: string;
    href: string;
  } = $props();
</script>

<a class="flex flex-row items-center text-lg px-2 py-4 color-active" {href}>
  <CaretLeft />
  {text}
</a>

<style>
  .color-active {
    color: #d65d0e;
  }
</style>
