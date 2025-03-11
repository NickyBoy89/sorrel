import { backendRootURL } from '../../../../constants';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
    const menuData = await fetch(`${backendRootURL}/api/menu/${params.menu_id}`);
    const menu = await menuData.json();

    const menuItemData = await fetch(`${backendRootURL}/api/menu/${params.menu_id}/items`)
    const items = await menuItemData.json()

    return {
        menu: menu,
        menuItems: items,
        menuId: params.menu_id
    }
}