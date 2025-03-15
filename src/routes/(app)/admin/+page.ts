import type { PageLoad } from './$types';
import { backendRootURL } from '../../constants';

export const load: PageLoad = async ({ fetch }) => {
    const menusData = await fetch(`${backendRootURL}/api/menu/list`);
    const menus = await menusData.json();

    return {
        menus: menus,
    }
}