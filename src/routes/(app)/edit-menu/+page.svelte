<script lang="ts">
    import MenuRenderer from "$lib/components/menuRenderer.svelte";
    import MenuItemEditor from "$lib/components/menuItemEditor.svelte";
    import UiButton from "$lib/components/uiButton.svelte";
    import Navbar from "$lib/components/navbar.svelte";
    import { APIUrl, menuDefaultSections } from "../../../constants.js";

    import { page } from "$app/state";
    import { onMount } from "svelte";
    import UiButtonLink from "$lib/components/uiButtonLink.svelte";
    import { bearerToken } from "../stores.js";
    import { get } from "svelte/store";
    import { initKeycloak } from "$lib/auth.js";

    let menuId: string | null = null;

    let menu = $state({
        name: "Placeholder",
        date: new Date(),
    });
    let visibleItems = $state([] as Array<MenuItemType>);

    onMount(() => {
        menuId = page.url.searchParams.get("menu-id");
        initKeycloak();
        fetchMenu();
        fetchMenuItems();
    })

    const fetchMenu = async () => {
        menu = await fetch(`${APIUrl}/api/menu/${menuId}`)
        .then(resp => resp.json())
        .catch((error) => {
            console.log(error);
        });
    }

    const fetchMenuItems = async () => {
        await fetch(`${APIUrl}/api/menu/${menuId}/items`)
            .then((resp) => resp.json())
            .then(respJson => visibleItems = respJson)
            .catch((error) => {
                console.log(error);
            });
    }

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

    const createMenuItem = async () => {
        await fetch(`${APIUrl}/api/menu/${menuId}/create-item`, {
            method: "POST",
            headers: {
                Authorization: `Bearer: ${get(bearerToken)}`,
            },
            body: new URLSearchParams({name: "", description: ""}),
        }).catch((error) => {
            console.log(error);
        });

        fetchMenuItems();
    }

    let menuSections = $derived(itemsToSections(visibleItems));
</script>

{#snippet backButton()}
    <UiButtonLink text="← Back" href="/admin" />
{/snippet}

<Navbar mainItem={backButton} />

<div class="grid grid-cols-1 md:grid-cols-2 my-4">
    <div class="menu-editor md:px-6 py-6 px-4">
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
        <MenuRenderer menuName={menu.name} menuDate={new Date(menu.date)} sections={menuSections} />
    </section>
</div>