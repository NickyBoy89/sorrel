<script lang="ts">
  let { mainItem = () => {}, items = () => {} } = $props()
</script>

<nav class="flex flex-row py-2">
  <div class="flex-grow">
      {@render mainItem()}
  </div>
  <div class="flex flex-row items-center space-x-4 underline">
    {@render items()}
  </div>
</nav>