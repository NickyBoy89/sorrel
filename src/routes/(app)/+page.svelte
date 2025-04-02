<script lang="ts">;
    import MenuEditor from "$lib/components/menuEditor.svelte";
    import { toJsDate } from "$lib/tools.js";
    import { backendRootURL } from "../../constants.js";
</script>

<h1 class="text-4xl text-center my-4 text-black dark:text-white">Shared With You</h1>
<div class="flex flex-col my-4">
    {#await fetch(`${backendRootURL}/api/menu/list`).then((resp) => resp.json())}
        Loading...
    {:then menus}
    {#each menus as menu}
        <MenuEditor menuName={menu.name} menuDate={toJsDate(menu.date)} menuId={menu.id} canEdit={false}/>
    {/each}
    {:catch error}
        There was an error fetching the menus: {error}
    {/await}
</div>

<!-- <h1 class="text-4xl text-center my-4">Shopping Lists</h1> -->