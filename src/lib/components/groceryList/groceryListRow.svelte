<script lang="ts">
    import { onMount } from "svelte";
    import InteractiveBox from "../ui/interactiveBox.svelte";
    import { APIUrl } from "../../../constants";

    let { groceryListId } = $props();

    let itemCount = $state(-1);

    onMount(() => {
        fetch(`${APIUrl}/api/v1/grocery_list/${groceryListId}`)
            .then(resp => resp.json())
            .then(respJson => itemCount = respJson.size)
            .catch(error => console.error(error));
    })
</script>

<InteractiveBox>
    <a href="/grocery_list/?id={groceryListId}">
        <div class="flex flex-row justify-between items-center">
            <div class="text-xl font-semibold">Groceries</div>
            <div class="text-md">{itemCount} remaining</div>
        </div>
    </a>
</InteractiveBox>