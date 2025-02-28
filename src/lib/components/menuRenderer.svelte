<script lang="ts">
  import MenuItem from "./menuItem.svelte";
  let {items = new Map()} = $props();
</script>

<div class="grid grid-cols-8">
  <div class="col-start-3 col-span-4 border-debu flex flex-col h-screen text-center menu-font">
    <div class="menu-header mt-32 mb-16">
      <h1 class="menu-title">
        Victoria Night
      </h1>
      <h3 class="menu-date">
        1 February 2025
      </h3>
    </div>

    <section class="content flex-grow">
      {#each items as [menuSectionHeader, menuItems]}
        <div class="menu-section my-4">
          <h2 class="menu-header">{menuSectionHeader}</h2>
          <ul>
            {#each menuItems as appetizer}
              <MenuItem name={appetizer.name} description={appetizer.description} />
            {/each}
          </ul>
        </div> 
      {/each}
    </section>
  </div>
</div>

<style>
  .menu-font {
    font-family: 'Times New Roman', Times, serif;
  }

  .menu-title {
    font-size: 6.25rem;
  }

  .menu-header {
    font-size: 3.25rem;
  }
</style>