<script lang="ts">
    import MenuRenderer from "$lib/components/menuRenderer.svelte";
    import MenuItemEditor from "$lib/components/menuItemEditor.svelte";
    import UiButton from "$lib/components/uiButton.svelte";

    let { data } = $props();
</script>

<div class="grid grid-cols-2">
    <div class="menu-editor p-6">
        <div class="menu-editors grid gap-4">
            {#each data.menuItems as menuItem}
                <MenuItemEditor name={menuItem.name} description={menuItem.description} />
            {/each}
        </div>
        <div class="flex justify-between mt-4">
            <UiButton text="New" action={() => {alert("New menu item")}} color="#fb4934"/>
            <UiButton text="Save" action={() => {alert("Save menu items")}} color="#98971a"/>
        </div>
    </div>

    <section class="menu-preview">
        <MenuRenderer menuName={data.menu.name} menuDate={new Date(data.menu.date)}/>
    </section>
</div>