<script lang="ts">
  let { text, href, color = "#d65d0e" } = $props();
</script>

<a href={href} class="button-container text-center">
  <button class="px-3 py-1.5 rounded-sm cursor-pointer" style="background-color: {color};">{text}</button>
</a>

<style>
  .button-container {
    font-size: 1rem;
  }
</style>