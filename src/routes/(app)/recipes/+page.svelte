<script lang="ts">
  import Plus from "phosphor-svelte/lib/Plus";
  import { onMount } from "svelte";
  import { APIUrl } from "../../../constants";
  import { Recipe, RecipeId } from "$lib/recipe";
  import { page } from "$app/state";
  import Modal from "$lib/design-kit/Modal.svelte";
  import Card from "$lib/design-kit/Card.svelte";
  import { Recipes } from "$lib/api/recipe";
  import TextField from "$lib/design-kit/form-fields/TextField.svelte";
  import UiButton from "$lib/components/uiButton.svelte";

  type RecipeCard = {
    id: RecipeId;
    text: string;
    color: string;
  };

  let recipeId: RecipeId | null = $state(null);

  let recipes: Recipe[] = $state([]);

  let createRecipeModalVisible = $state(false);

  const api = new Recipes(APIUrl);

  const updateRecipes = async () => {
    const recipeIds = await api.findAllIds();

    recipes = await Promise.all(recipeIds.map(api.find.bind(api)));
  };

  onMount(async () => {
    const recipeIdParam = page.url.searchParams.get("id");
    if (recipeIdParam !== null) {
      recipeId = RecipeId.decode(Number.parseInt(recipeIdParam));
    }

    await updateRecipes();
  });

  let createdRecipeName: string = $state("");
</script>

{#snippet recipeCard(data: RecipeCard)}
  <a
    href="/recipes/recipe?id={data.id}"
    class="flex flex-col text-center m-4 cursor-pointer"
  >
    <Card class="w-32 h-24 {data.color}" />
    <div class="text-sm">{data.text}</div>
  </a>
{/snippet}

<Modal bind:visible={createRecipeModalVisible}>
  <Card>
    <div class="flex flex-col">
      <div class="text-xl font-semibold">Name</div>
      <TextField bind:initialValue={createdRecipeName} onchange={() => {}} />
      <UiButton
        action={async () => {
          await api.create({
            name: createdRecipeName,
            color: "#ffffff",
            ingredients: [],
          });

          createRecipeModalVisible = false;

          void updateRecipes();
        }}
        text="Create"
      />
    </div>
  </Card>
</Modal>

<div class="flex flex-wrap">
  <button
    class="flex flex-col text-center m-4 cursor-pointer"
    onclick={() => (createRecipeModalVisible = !createRecipeModalVisible)}
  >
    <Card class="w-32 h-24">
      <Plus />
    </Card>
    <div class="text-sm">Add new...</div>
  </button>
  {#each recipes as recipe}
    {@render recipeCard({
      id: recipe.id,
      text: recipe.name,
      color: recipe.color,
    })}
  {/each}
</div>
