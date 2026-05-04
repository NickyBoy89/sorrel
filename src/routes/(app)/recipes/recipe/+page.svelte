<script lang="ts">
  import { page } from "$app/state";
  import { Recipe, RecipeId } from "$lib/recipe";
  import GroceryListItem from "$lib/components/groceryList/groceryListItem.svelte";
  import { onMount } from "svelte";
  import { Recipes } from "$lib/api/recipe";
  import { APIUrl } from "../../../../constants";
  import NavigationButton from "$lib/components/ui/navigationButton.svelte";
  import DotsThreeVertical from "phosphor-svelte/lib/DotsThreeVertical";

  const api = new Recipes(APIUrl);

  let recipeOptionsShown = $state(false);

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

{#if currentRecipe != null}
  <div class="flex flex-col relative">
    <div class="flex flex-row relative justify-between mx-4">
      <div class="font-semibold text-4xl my-2">
        {currentRecipe.name}
      </div>
      <button
        class="cursor-pointer hover:bg-neutral-800 flex justify-end items-center"
        onclick={() => (recipeOptionsShown = !recipeOptionsShown)}
      >
        <DotsThreeVertical class="text-xl" />
      </button>
    </div>
    <div
      class={[
        "mt-2 w-48 shadow-lg border border-neutral-500",
        { hidden: !recipeOptionsShown },
      ]}
    >
      This shows up sometimes
    </div>
    <div class="font-semibold text-xl mt-4">Ingredients</div>
    {#each currentRecipe.ingredients as ingredient}
      <GroceryListItem {ingredient} />
    {/each}
  </div>
{/if}
