<script lang="ts">
    import UiButton from "./uiButton.svelte";
    import { onMount } from "svelte";
    import { backendRootURL } from "../../constants";

    let { menuName, menuDate = new Date(), id } = $props();

    let editorElement: HTMLElement | null = null;
    
    onMount(() => {
      editorElement = document.querySelector(".menu-editor");
    })

    function openEditor() {
      editorElement?.classList.remove("hidden");
      editorElement?.classList.add("flex");
    }

    function closeEditor() {
      editorElement?.classList.remove("flex");
      editorElement?.classList.add("hidden");
    }
</script>

<div class="menu-item-editor rounded-md p-4 my-2 w-full">
  <div class="flex flex-row justify-between">
    <a href="/edit-menu/{id}" class="menu-title flex flex-row grow items-center">
      <p class="text-3xl mb-1 inline-block">{menuName}</p> - {new Intl.DateTimeFormat("en-US", {month: "short"}).format(menuDate)} {menuDate.getDate()} {menuDate.getFullYear()}
    </a>
    <div class="menu-options">
      <UiButton text="Edit" action={openEditor}/>
    </div>
  </div>
  <div class="menu-editor hidden flex-col">
    <form action="{backendRootURL}/api/menu/{id}/edit" id="edit-menu-form" method="GET">
      <label for="edit-menu-name">Name:</label>
      <input type="text" id="edit-menu-name" name="name" class="block">
      <label for="edit-menu-date">Date:</label>
      <input type="date" id="edit-menu-date" name="date" class="block">
    </form>
    <UiButton text="Save" color="#458588" action={() => {
      document.getElementById("edit-menu-form")?.submit();
      closeEditor();
      }} />
  </div>
</div>

<style>
  .menu-item-editor {
      background-color: #665c54;
  }

  input {
      background-color: #3c3836;
      height: 2rem;
  }
</style>