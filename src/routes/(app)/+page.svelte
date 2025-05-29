<script lang="ts">;
    import MenuSelector from "$lib/components/menuSelector.svelte";
    import Navbar from "$lib/components/navbar.svelte";
    import UserStatus from "$lib/components/userStatus.svelte";
    import { toJsDate } from "$lib/tools.js";
    import { onMount } from "svelte";
    import { APIUrl } from "../../constants";
    import GroceryListRow from "$lib/components/grocery_list/groceryListRow.svelte";

    let username = $state("User...");
    let menus: Array<Menu> = $state([]);

    let groceryLists = $state([1]);

    onMount(() => {
        fetch(`${APIUrl}/api/user?${new URLSearchParams({
            userId: `${localStorage.getItem("userId")}`,
        })}`)
            .then((resp) => resp.json())
            .then((user) => username = user.display_name)
            .catch((error) => console.error(error));
        fetch(`${APIUrl}/api/menu/list`)
            .then((resp) => resp.json())
            .then((respJson) => menus = respJson)
            .catch((error) => console.error(error));
    })
</script>

<Navbar>
    <UserStatus userName={username} />
</Navbar>

<h1 class="text-4xl text-center mb-8 text-black dark:text-white">Shared With You</h1>
<div class="flex flex-col my-4 px-4 space-y-4">
    {#each menus as menu}
        <MenuSelector menuName={menu.name} menuDate={toJsDate(menu.date)} menuId={menu.id} canEdit={false}/>
    {/each}
</div>

<h1 class="text-4xl text-center mb-8 text-black dark:text-white">Shopping Lists</h1>

<div class="flex flex-col my-4 px-4 space-y-4">
    {#each groceryLists as groceryList}
        <GroceryListRow />
    {/each}
</div>

<!-- <h1 class="text-4xl text-center my-4">Shopping Lists</h1> -->