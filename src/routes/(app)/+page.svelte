<script lang="ts">;
  import MenuEditor from "$lib/components/menuEditor.svelte";
  import { backendRootURL } from "../../constants";
  let { data } = $props();
</script>

<h1 class="text-4xl">Manage Menus</h1>
<div class="flex flex-col my-4">
  {#each data.menus as menu}
    <MenuEditor menuName={menu.name} id={menu.id}/>
  {/each}
</div>