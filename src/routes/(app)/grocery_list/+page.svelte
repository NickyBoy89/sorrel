<script lang="ts">
    import { page } from "$app/state";
    import type { GroceryItem } from "$lib/grocery_list";
    import { onMount } from "svelte";
    import { APIUrl } from "../../../constants";
    import Navbar from "$lib/components/navbar.svelte";
    import UiButtonLink from "$lib/components/uiButtonLink.svelte";
    import GroceryListItem from "$lib/components/groceryList/groceryListItem.svelte";

    let groceryListId: string | null;

    let groceryListItems = $state([] as Array<GroceryItem>);

    onMount(() => {
        groceryListId = page.url.searchParams.get("id");

        fetch(`${APIUrl}/api/v1/grocery_list/${groceryListId}`)
            .then(resp => resp.json())
            .then(respJson => groceryListItems = respJson)
            .catch(error => console.error(error));
    })
</script>

{#snippet groceryListCategoryHeader(headerName: string)}
    <h1 class="text-4xl text-black dark:text-white">{headerName}</h1>
{/snippet}

<Navbar>
    {#snippet mainItem()}
        <UiButtonLink text="Back" href="/" />
    {/snippet}
</Navbar>


{@render groceryListCategoryHeader("Produce")}

<ul>
    <GroceryListItem />
    <GroceryListItem />
    <GroceryListItem />
    <GroceryListItem />
    <GroceryListItem />
    <!-- {#each groceryListItems as groceryItem}
        <li>{groceryItem.name}</li>
    {/each} -->
</ul>

{@render groceryListCategoryHeader("Meats")}

<ul>
    <GroceryListItem />
    <GroceryListItem />
    <GroceryListItem />
</ul>