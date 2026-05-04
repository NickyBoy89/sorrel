<script lang="ts">
  import { page } from "$app/state";
  import { Recipe, RecipeId } from "$lib/recipe";
  import GroceryListItem from "$lib/components/groceryList/groceryListItem.svelte";
  import { onMount } from "svelte";
  import { Recipes } from "$lib/api/recipe";
  import { APIUrl } from "../../../../constants";
  import NavigationButton from "$lib/components/ui/navigationButton.svelte";
  import PencilSimple from "phosphor-svelte/lib/PencilSimple";
  import Backspace from "phosphor-svelte/lib/Backspace";
  import DropdownOption from "$lib/design-kit/DropdownOption.svelte";
  import Dropdown from "$lib/design-kit/Dropdown.svelte";

  const api = new Recipes(APIUrl);

  let currentRecipe: Recipe | null = $state(null);

  onMount(async () => {
    const recipeIdParam = page.url.searchParams.get("id");

    if (recipeIdParam === null) {
      throw new Error("Null recipe id found");
    }

    const recipeId = RecipeId.decode(Number.parseInt(recipeIdParam));

    currentRecipe = await api.find(recipeId);
  });
</script>

<NavigationButton text="All recipes" href="/recipes/" />

{#snippet editIcon()}
  <PencilSimple />
{/snippet}

{#snippet deleteIcon()}
  <Backspace />
{/snippet}

{#if currentRecipe != null}
  <div class="flex flex-col">
    <div class="flex flex-row justify-between mx-4">
      <div class="font-semibold text-4xl my-2">
        {currentRecipe.name}
      </div>
      <Dropdown>
        <DropdownOption text="Edit recipe" icon={editIcon} />
        <DropdownOption text="Delete recipe" icon={deleteIcon} />
      </Dropdown>
    </div>
    <div class="font-semibold text-xl mt-4">Ingredients</div>
    {#each currentRecipe.ingredients as ingredient}
      <GroceryListItem {ingredient} />
    {/each}
  </div>
{/if}
