<script lang="ts">
  let { text, action } = $props();
</script>

<button class="px-4 rounded-sm cursor-pointer text-2xl" onclick={action}>{text}</button>

<style>
  button {
    display: inline-block;
    background-color: #d65d0e;
  }
</style>