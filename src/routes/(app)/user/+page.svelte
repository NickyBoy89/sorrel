<script lang="ts">;
    import MenuEditor from "$lib/components/menuEditor.svelte";
    let { data } = $props();
</script>

<h1 class="text-4xl text-center my-4">Shared With You</h1>
<div class="flex flex-col my-4">
    {#each data.menus as menu}
        <MenuEditor menuName={menu.name} menuDate={new Date(menu.date)} id={menu.id} canEdit={false}/>
    {/each}
</div>

<h1 class="text-4xl text-center my-4">Shopping Lists</h1>