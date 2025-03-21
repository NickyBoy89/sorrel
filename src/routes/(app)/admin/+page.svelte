<script lang="ts">;
  import MenuEditor from "$lib/components/menuEditor.svelte";
  import UiButton from "$lib/components/uiButton.svelte";
  import { backendRootURL } from "../../../constants";
  let { data } = $props();

  const fetchMenus = fetch(`${backendRootURL}/api/menu/list`).then((resp) => resp.json());

  const createMenu = async () => {
    const now = new Date();

    await fetch(`${backendRootURL}/api/menu/create`, {
        method: "POST",
        body: new URLSearchParams({
          name: "New Menu",
          date: `${now.getFullYear()}-${now.getMonth().toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`,
        }),
    }).catch((error) => {
        console.log(error);
    });
    }
</script>

<h1 class="text-4xl text-center my-4">Manage Menus</h1>
<div class="flex flex-col my-4">
  {#await fetchMenus}
  Loading...
  {:then menus}
  {#each menus as menu}
    <MenuEditor menuName={menu.name} menuDate={new Date(menu.date)} menuId={menu.id}/>
  {/each}
  {:catch error}
  Error loading menus: {error}
  {/await}
  <UiButton text="New" action={createMenu} />
</div>