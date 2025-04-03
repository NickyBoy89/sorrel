<script lang="ts">
  let { mainItem = () => {}, children } = $props()
</script>

<nav class="flex flex-row py-2">
  <div class="flex-grow">
    {@render mainItem()}
  </div>
  <div class="flex flex-row items-center space-x-4 text-black dark:text-white">
    {@render children?.()}
  </div>
</nav>