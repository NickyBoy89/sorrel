<script lang="ts">
  import { page } from "$app/state";
  import { onMount } from "svelte";
  import { APIUrl } from "../../../../constants";
  import UiButton from "$lib/components/uiButton.svelte";
  import TextArea from "$lib/components/ui/textArea.svelte";
  import Select from "$lib/components/ui/select.svelte";

  let selectedIngredientId = $state(0);
  let quantity = $state("");

  let groceryListId: string | null = null;

  type GroceryItem = {
    id: number;
    name: string;
    category: string | null;
    quantity: string;
    checked: boolean;
  };

  type Ingredient = {
    id: number;
    name: string;
    category: string | null;
  };

  let items: Array<GroceryItem> = $state([]);

  let ingredientOptions: Array<Ingredient> = $state([]);

  const fetchItems = async () => {
    await fetch(`${APIUrl}/api/v1/grocery_list/${groceryListId}`)
      .then((resp) => resp.json())
      .then((respJson) => (items = respJson))
      .catch((error) => console.error(error));
  };

  const fetchIngredients = async () => {
    await fetch(`${APIUrl}/api/v1/ingredients`)
      .then((resp) => resp.json())
      .then((resp) => (ingredientOptions = resp));
  };

  onMount(() => {
    groceryListId = page.url.searchParams.get("id");

    fetchItems();

    fetchIngredients();
  });
</script>

<div class="flex flex-col">
  <ol>
    {#each items as ingredient}
      <div>
        Id: {ingredient.id}, Name: {ingredient.name}, Quantity: {ingredient.quantity}
      </div>
    {/each}
  </ol>
  <div class="flex flex-row grow gap-x-4 mx-4 items-center">
    <Select
      placeholder="Select or add an ingredient"
      options={ingredientOptions.map((ing) => ({
        value: ing.id.toString(),
        display: ing.name,
      }))}
      onchange={(selectedValue) => {
        selectedIngredientId = Number.parseInt(selectedValue);
      }}
    />
    <TextArea
      placeholder="Quantity..."
      onchange={(e) => (quantity = e.currentTarget.value)}
    />
    <UiButton
      text="Add"
      action={() =>
        fetch(`${APIUrl}/api/v1/grocery_list/${groceryListId}`, {
          method: "POST",
          body: JSON.stringify({
            id: selectedIngredientId,
            quantity,
            checked: false,
          }),
        })}
    />
  </div>
</div>
