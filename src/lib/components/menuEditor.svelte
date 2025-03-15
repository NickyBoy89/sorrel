<script lang="ts">
    import UiButton from "./uiButton.svelte";
    import { backendRootURL } from "../../constants";
    import { DateTime } from "luxon";

    let { menuName, menuDate, menuId, canEdit = true, relativeDate = true } = $props();

    let name = menuName;
    let date = $state(menuDate);

    let editorOpen = $state(false);

    const handleMenuChange = async () => {
      await fetch(`${backendRootURL}/api/menu/${menuId}/edit?${new URLSearchParams({
          name: name,
          date: date.toISOString().split("T")[0]
        })}`).catch((error) => {
        console.log(error);
      });
    }

    const displayedDate = (): string => {
      if (!relativeDate) {
        return date.toISOString().split("T")[0];
      }
      const relDate = DateTime.fromJSDate(date).toRelative({base: this});
      if (relDate === null) {
        return "unknown date"
      }
      return relDate
    }

    const toggleEditor = () => {
        editorOpen = !editorOpen;
    }
</script>

<div class="menu-item-editor rounded-md p-4 my-2 w-full">
  <div class="flex flex-row justify-between">
    <a href="/{canEdit ? "edit-menu" : "menu"}?menu-id={menuId}" class="menu-title flex flex-col sm:flex-row grow items-center">
      <div class="menu-name text-2xl">
        {menuName}
      </div>
      <div class="menu-date mx-4">
        {displayedDate()}
      </div>
    </a>
    <div class="menu-options flex items-center {canEdit ? "" : "hidden"}">
      <UiButton text="Edit" action={toggleEditor}/>
    </div>
  </div>
  <div class="menu-editor {editorOpen ? "" : "hidden"} flex-col">
    <div>
      <label for="edit-menu-name">Name:</label>
      <input type="text" id="edit-menu-name" name="name" class="block" value={menuName} onchange={handleMenuChange}>
      <label for="edit-menu-date">Date:</label>
      <input type="date" id="edit-menu-date" name="date" class="block date-editor" value={date.toISOString().split("T")[0]} onchange={handleMenuChange}>
    </div>
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