<script lang="ts">;
  import MenuEditor from "$lib/components/menuSelector.svelte";
  import UiButton from "$lib/components/uiButton.svelte";
  import { onMount } from "svelte";
  import { toJsDate } from "$lib/tools";
  import { bearerToken } from "../stores";
  import { get } from "svelte/store";
  import { APIUrl } from "../../../constants";

  let menus: Array<any> = $state([] as Array<MenuItemType>);

  const fetchMenus = async () => {
    await fetch(`${APIUrl}/api/menu/list`)
      .then((resp) => resp.json())
      .then((respJson) => menus = respJson)
      .catch((error) => {
        console.error(error);
        return [];
      });
  }

  onMount(() => {
    fetchMenus();
  })

  const createMenu = async () => {
    const now = new Date();

    await fetch(`${APIUrl}/api/menu/create`, {
        method: "POST",
        headers: {
          Authorization: `Bearer: ${get(bearerToken)}`
        },
        body: new URLSearchParams({
          name: "New Menu",
          date: `${now.getFullYear()}-${(now.getMonth() + 1).toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`,
        }),
    }).catch((error) => {
        console.log(error);
    });

    fetchMenus();
  }
</script>

<h1 class="text-4xl text-center my-4 text-black dark:text-white">Manage Menus</h1>
<div class="flex flex-col my-4 px-4">
  {#each menus as menu}
    <MenuEditor menuName={menu.name} menuDate={toJsDate(menu.date)} menuId={menu.id} canEdit={true} />
  {/each}
  <UiButton text="New" action={createMenu} />
</div>