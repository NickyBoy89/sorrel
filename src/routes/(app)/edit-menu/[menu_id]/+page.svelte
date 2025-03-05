<script lang="ts">
    import MenuRenderer from "$lib/components/menuRenderer.svelte";
    import MenuItemEditor from "$lib/components/menuItemEditor.svelte";
    import UiButton from "$lib/components/uiButton.svelte";

    let { data } = $props();

    async function createMenuItem() {
        await fetch(`http://localhost:9031/api/menu/${data.menuId}/create-item`, {
            method: "POST",
            body: new URLSearchParams({name: "", description: ""}),
        }).catch((error) => {
            console.log(error);
        });
    }

    async function saveMenuItems() {
        await fetch(`http://localhost:9031/api/menu/${data.menuId}/items/edit/1`, {
            method: "POST",
            body: new URLSearchParams({name: "Something", description: "Else"}),
        }).catch((error) => {
            console.log(error);
        })
    }
</script>

<div class="grid grid-cols-2">
    <div class="menu-editor p-6">
        <div class="menu-editors grid gap-4">
            {#each data.menuItems as menuItem}
                <MenuItemEditor name={menuItem.name} description={menuItem.description} />
            {/each}
        </div>
        <div class="flex justify-between mt-4">
            <UiButton text="New" action={createMenuItem} color="#fb4934"/>
            <UiButton text="Save" action={saveMenuItems} color="#98971a"/>
        </div>
    </div>

    <section class="menu-preview">
        <MenuRenderer menuName={data.menu.name} menuDate={new Date(data.menu.date)}/>
    </section>
</div>