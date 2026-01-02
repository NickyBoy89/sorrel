<script lang="ts">
  import { page } from "$app/state";
  import { onMount } from "svelte";
  import { APIUrl } from "../../../../constants";
  import UiButton from "$lib/components/uiButton.svelte";
  import TextArea from "$lib/components/ui/textArea.svelte";
  import Select from "$lib/components/ui/select.svelte";
  import GroceryListItem from "$lib/components/groceryList/groceryListItem.svelte";
  import NavigationButton from "$lib/components/ui/navigationButton.svelte";

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

<NavigationButton text="All lists" href="/grocery_lists/" />
<div class="text-4xl ml-4 font-semibold">Groceries</div>
<div class="flex flex-col">
  <ol class="divide-y divide-neutral-800 my-2">
    {#each items as ingredient, index}
      <GroceryListItem
        {ingredient}
        ondelete={() =>
          fetch(
            `${APIUrl}/api/v1/grocery_list/${groceryListId}/item/${ingredient.id}`,
            { method: "DELETE" },
          ).then((_) => items.splice(index, 1))}
      />
    {/each}
  </ol>
  <div class="flex flex-col md:flex-row grow gap-4 mx-4 items-center">
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
        }).then((_) => fetchItems())}
    />
  </div>
</div>
