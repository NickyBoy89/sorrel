import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    const menusData = await fetch("/api/menu/list");
    const menus = await menusData.json();

    return {
        menus: menus,
    }
}