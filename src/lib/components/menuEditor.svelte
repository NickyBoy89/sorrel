<script lang="ts">
    import UiButton from "./uiButton.svelte";
    import { backendRootURL } from "../../constants";

    let { menuName, menuDate = new Date(), id } = $props();

    let name = menuName;
    let date = $state(menuDate);

    let editorOpen = $state(false);

    const handleMenuChange = async () => {
      await fetch(`${backendRootURL}/api/menu/${id}/edit?${new URLSearchParams({
          name: name,
          date: date.toISOString().split("T")[0]
        })}`).catch((error) => {
        console.log(error);
      });
    }

    const toggleEditor = () => {
        editorOpen = !editorOpen;
    }
</script>

<div class="menu-item-editor rounded-md p-4 my-2 w-full">
  <div class="flex flex-row justify-between">
    <a href="/edit-menu?menu-id={id}" class="menu-title flex flex-col sm:flex-row grow items-center">
      <div class="menu-name text-2xl">
        {menuName}
      </div>
      <div class="menu-date mx-4">
        {new Intl.DateTimeFormat("en-US", {month: "short"}).format(date)} {date.getDate()} {date.getFullYear()}
      </div>
    </a>
    <div class="menu-options flex items-center">
      <UiButton text="Edit" action={toggleEditor}/>
    </div>
  </div>
  <div class="menu-editor {editorOpen ? "" : "hidden"} flex-col">
    <form action="{backendRootURL}/api/menu/{id}/edit" id="edit-menu-form" method="GET">
      <label for="edit-menu-name">Name:</label>
      <input type="text" id="edit-menu-name" name="name" class="block" value={menuName} onchange={handleMenuChange}>
      <label for="edit-menu-date">Date:</label>
      <input type="date" id="edit-menu-date" name="date" class="block date-editor" value={date.toISOString().split("T")[0]} onchange={handleMenuChange}>
    </form>
    <UiButton text="Close" color="#458588" action={toggleEditor} />
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