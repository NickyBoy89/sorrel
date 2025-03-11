<script lang="ts">
  import MenuItem from "./menuItem.svelte";
  let {menuName, menuDate = new Date(), sections = new Map<string, any>()} = $props();
  
</script>

<div class="menu inline-block text-center menu-font px-8 h-max justify-center w-full">
  <section class="menu-header mt-32 mb-8">
    <h1 class="menu-title">
      {menuName}
    </h1>
    <h3 class="menu-date">
      {menuDate.getDate()} {new Intl.DateTimeFormat("en-US", {month: "long"}).format(menuDate)} {menuDate.getFullYear()}
    </h3>
  </section>

  <section class="content flex-grow">
    {#each sections as [menuSectionHeader, menuItems]}
      <div class="menu-section my-4">
        <h2 class="menu-section-header">{menuSectionHeader}</h2>
        <ul>
          {#each menuItems as appetizer}
            <MenuItem name={appetizer.name} description={appetizer.description} />
          {/each}
        </ul>
      </div> 
    {/each}
  </section>
</div>

<style>
  .menu {
    background-color: #ffffff;
    color: #000000;
  }

  .menu-font {
    font-family: 'Times New Roman', Times, serif;
  }

  .menu-title {
    font-size: 4.25rem;
  }

  .menu-header {
    font-size: 1.25rem;
  }

  .menu-section-header {
    font-size: 3.25rem;
  }
</style>