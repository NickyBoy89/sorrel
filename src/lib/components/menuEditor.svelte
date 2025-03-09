<script lang="ts">
    import UiButton from "./uiButton.svelte";

    let { menuName, menuDate = new Date(), id } = $props();
</script>

<div class="menu-item-editor rounded-md p-4 my-2 w-full">
  <div class="flex flex-row justify-between">
    <a href="/edit-menu/{id}" class="menu-title flex flex-row grow items-center">
      <p class="text-3xl mb-1 inline-block">{menuName}</p> - {new Intl.DateTimeFormat("en-US", {month: "short"}).format(menuDate)} {menuDate.getDate()} {menuDate.getFullYear()}
    </a>
    <div class="menu-options">
      <UiButton text="Edit" action={() => alert("Edit menu")}/>
    </div>
  </div>
</div>

<style>
  .menu-item-editor {
      background-color: #665c54;
  }
</style>