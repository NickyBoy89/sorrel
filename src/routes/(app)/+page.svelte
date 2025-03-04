<script lang="ts">
  import MenuRenderer from "$lib/components/menuRenderer.svelte";
  import MenuItemEditor from "$lib/components/menuItemEditor.svelte";
  import UiButton from "$lib/components/uiButton.svelte";
</script>

<div class="grid grid-cols-2">
  <div class="menu-editor p-6">
    <div class="menu-editors grid gap-4">
      <MenuItemEditor />
      <MenuItemEditor />
    </div>
    <div class="flex justify-between mt-4">
      <UiButton text="New" action={() => {alert("New menu item")}} color="#fb4934"/>
      <UiButton text="Save" action={() => {alert("Save menu items")}} color="#98971a"/>
    </div>
  </div>
  
  <section class="menu-preview">
    <MenuRenderer menuName="Tomorrow Night's Dinner" />
  </section>
</div>