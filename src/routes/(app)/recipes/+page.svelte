<script lang="ts">
  import Plus from "phosphor-svelte/lib/Plus";
  import { onMount } from "svelte";
  import { APIUrl } from "../../../constants";
  import z from "zod";
  import { Recipe, RecipeId } from "$lib/recipe";
  import { page } from "$app/state";
  import GroceryListItem from "$lib/components/groceryList/groceryListItem.svelte";

  type RecipeCard = {
    id: RecipeId;
    text: string;
    color: string;
  };

  let recipeId: RecipeId | null = $state(null);
  let currentRecipe: Recipe | null = $state(null);

  let recipes: Recipe[] = $state([]);

  const allRecipeIds = z.array(RecipeId);

  const fetchRecipe = async (recipeId: RecipeId) =>
    fetch(`${APIUrl}/api/v1/recipes/${recipeId}`)
      .then((resp) => resp.json())
      .then((resp) => Recipe.decode(resp));

  onMount(async () => {
    const recipeIdParam = page.url.searchParams.get("id");
    if (recipeIdParam !== null) {
      recipeId = RecipeId.decode(Number.parseInt(recipeIdParam));
    }

    const recipeIds = await fetch(`${APIUrl}/api/v1/recipes`)
      .then((resp) => resp.json())
      .then((resp) => allRecipeIds.decode(resp));

    recipes = await Promise.all(recipeIds.map(fetchRecipe));
  });
</script>

{#snippet recipeCard(data: RecipeCard)}
  <a
    href="/recipes?id={data.id}"
    onclick={async () => {
      (recipeId = data.id), (currentRecipe = await fetchRecipe(data.id));
    }}
    class="flex flex-col text-center m-4 cursor-pointer"
  >
    <div class="w-32 h-24 rounded-lg {data.color}"></div>
    <div class="text-sm">{data.text}</div>
  </a>
{/snippet}

<div class="flex flex-wrap">
  {#each recipes as recipe}
    {@render recipeCard({
      id: recipe.id,
      text: recipe.name,
      color: recipe.color,
    })}
  {/each}
  <button class="flex flex-col text-center m-4 cursor-pointer">
    <div
      class="flex w-32 h-24 rounded-lg bg-neutral-700 items-center justify-center"
    >
      <Plus />
    </div>
    <div class="text-sm">Add new...</div>
  </button>
</div>

{#if currentRecipe !== null}
  <div class="flex flex-col">
    <div class="font-semibold text-4xl mb-4">{currentRecipe.name}</div>
    <div class="font-semibold text-xl mt-4">Ingredients</div>
    {#each currentRecipe.ingredients as ingredient}
      <GroceryListItem {ingredient} />
    {/each}
  </div>
{/if}
