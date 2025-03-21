import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    const menus = await fetch("/api/menu/list")
    .then((resp) => resp.json()).catch((error) => {
        console.error(error);
        return {};
    });

    return {
        menus: menus,
    }
}