<script lang="ts">
  import MenuItem from "./menuItem.svelte";
  let {menuName, menuDate = new Date(), sections = new Map<string, any>()} = $props();
  
</script>

<div class="menu inline-block text-center menu-font px-8 h-max justify-center w-full h-full">
  <section class="menu-date mt-16 lg:mt-32 mb-8">
    <h1 class="text-6xl sm:text-7xl">
      {menuName}
    </h1>
    <h3 class="text-xl">
      {menuDate.getDate()} {new Intl.DateTimeFormat("en-US", {month: "long"}).format(menuDate)} {menuDate.getFullYear()}
    </h3>
  </section>

  <section class="content flex-grow">
    {#each sections as [menuSectionHeader, menuItems]}
      <div class="menu-section my-4">
        <h2 class="text-5xl">{menuSectionHeader}</h2>
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
</style>