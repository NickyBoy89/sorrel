<script lang="ts">
  let { text, action, color = "#d65d0e" } = $props();
</script>

<div class="button-container text-center text-neutral-900 dark:text-white">
  <button class="px-3 py-1.5 rounded-sm cursor-pointer" onclick={action} style="background-color: {color};">{text}</button>
</div>

<style>
  .button-container {
    font-size: 1rem;
  }
</style>