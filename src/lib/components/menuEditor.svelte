<script lang="ts">
    import UiButton from "./uiButton.svelte";

    let { menuName, menuDate = new Date() } = $props();
</script>

<div class="menu-item-editor rounded-md p-4 my-2 w-full">
  <form action="" class="flex flex-row justify-between">
      <div class="menu-name text-4xl">{menuName}</div>
      <div class="menu-options">
        <UiButton text="Edit" action={() => alert("Edit pressed")}/>
      </div>
  </form>
</div>

<style>
  .menu-item-editor {
      background-color: #665c54;
  }

  .menu-section-selector {
      background-color: #3c3836;
  }

  .item-edit-input {
      background-color: #3c3836;
      height: 2rem;
  }
</style>