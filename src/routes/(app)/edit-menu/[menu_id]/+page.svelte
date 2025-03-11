<script lang="ts">
    import MenuRenderer from "$lib/components/menuRenderer.svelte";
    import MenuItemEditor from "$lib/components/menuItemEditor.svelte";
    import UiButton from "$lib/components/uiButton.svelte";
    import Navbar from "$lib/components/navbar.svelte";
    import { backendRootURL, menuDefaultSections } from "../../../../constants.js";

    let { data } = $props();

    let visibleItems = $state(data.menuItems as Array<MenuItemType>);

    const toUppercase = (s: string) => {
        return s.charAt(0).toUpperCase() + s.slice(1);
    }

    const itemsToSections = (items: Array<MenuItemType>): Map<string, Array<MenuItemType>> => {
        let sections = new Map();

        items.forEach(menuItem => {
            const sectionName = toUppercase(menuItem.section == null ? menuDefaultSections : menuItem.section);
            if (!sections.has(sectionName)) {
                sections.set(sectionName, []);
            }
            sections.get(sectionName).push(menuItem);
        });

        return sections
    }

    let menuSections = $derived(itemsToSections(visibleItems));

    const fetchMenuItems = async () => {
        await fetch(`${backendRootURL}/api/menu/${data.menuId}/items`)
            .then((resp) => resp.json())
            .then((respJson) => {
                visibleItems = respJson;
            })
            .catch((error) => {
                console.log(error);
            });
    }

    const createMenuItem = async () => {
        await fetch(`${backendRootURL}/api/menu/${data.menuId}/create-item`, {
            method: "POST",
            body: new URLSearchParams({name: "", description: ""}),
        }).catch((error) => {
            console.log(error);
        });

        fetchMenuItems();
    }
</script>

<Navbar backlinkHref="/" />

<div class="grid grid-cols-2">
    <div class="menu-editor p-6">
        <div class="menu-editors grid gap-4">
            {#each visibleItems as menuItem (menuItem.id)}
                <MenuItemEditor itemId={menuItem.id} name={menuItem.name} description={menuItem.description} section={menuItem.section} onchange={fetchMenuItems} />
            {/each}
        </div>
        <div class="flex justify-between mt-4">
            <UiButton text="New" action={createMenuItem} color="#fb4934"/>
        </div>
    </div>

    <section class="menu-preview">
        <MenuRenderer menuName={data.menu.name} menuDate={new Date(data.menu.date)} sections={menuSections} />
    </section>
</div>