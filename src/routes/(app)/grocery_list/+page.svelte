<script lang="ts">
    import { page } from "$app/state";
    import type { GroceryItem } from "$lib/grocery_list";
    import { onMount } from "svelte";
    import { APIUrl } from "../../../constants";
    import Navbar from "$lib/components/navbar.svelte";
    import UiButtonLink from "$lib/components/uiButtonLink.svelte";
    import GroceryListItem from "$lib/components/groceryList/groceryListItem.svelte";

    let groceryListId: string | null;

    let groceryListItems: Array<GroceryItem> = $state([]);
    
    onMount(() => {
        groceryListId = page.url.searchParams.get("id");

        fetch(`${APIUrl}/api/v1/grocery_list/${groceryListId}/items`)
            .then(resp => resp.json())
            .then(respJson => groceryListItems = respJson)
            .catch(error => console.error(error));
    })

    const combineCategories = (items: GroceryItem[]): Map<string, GroceryItem[]> => {
        let combined: Map<string, GroceryItem[]> = new Map();

        items.forEach((item) => {
            const category = item.category == null ? "Uncategorized" : item.category;

            if (!combined.has(category)) {
                combined.set(category, []);
            }

            combined.get(category)?.push(item);
        })

        return combined;
    }

    let categortizedItems = $derived(combineCategories(groceryListItems));
</script>

{#snippet categoryHeader(name: string)}
    <div class="font-bold text-lg w-full header-underlined text-neutral-500 mt-8">{name}</div>
{/snippet}

<Navbar>
    {#snippet mainItem()}
        <UiButtonLink text="Back" href="/" />
    {/snippet}
</Navbar>

<div class="flex flex-col ml-4">
    {#each categortizedItems as [categoryName, categoryItems]}
        <div class="flex flex-col w-full">
            {@render categoryHeader(categoryName)}
            <div class="flex flex-col">
                {#each categoryItems as item}
                    <GroceryListItem text={item.name}/>
                {/each}
            </div>
        </div>
    {/each}
</div>

<style>
    .header-underlined {
        border-bottom: 1px solid var(--color-neutral-800);
    }
</style>