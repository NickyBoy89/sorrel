import type { PageLoad } from './$types';
import { backendRootURL } from '../../constants';

export const load: PageLoad = async ({ fetch }) => {
    const menus = await fetch(`${backendRootURL}/api/menu/list`)
    .then((resp) => resp.json()).catch((error) => {
        console.error(error);
        return {};
    });

    return {
        menus: menus,
    }
}