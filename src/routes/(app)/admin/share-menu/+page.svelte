<script lang="ts">
    import UiButton from "$lib/components/uiButton.svelte";
    import { onMount } from "svelte";
    import Fa from 'svelte-fa'
    import { faSpinner, faCheck, faXmark, type IconDefinition } from "@fortawesome/free-solid-svg-icons";
    import { APIUrl } from "../../../../constants";
    import { bearerToken } from "../../stores";
    import { get } from "svelte/store";
    import { page } from "$app/state";

    type User = {
        id: number,
        display_name: string,
    };

    type Checkbox = Event & { currentTarget: EventTarget & HTMLInputElement};

    let menuId: number;

    let users = $state([] as Array<User>);
    let icons: Map<number, IconDefinition> = $state(new Map());

    let selected: Set<number> = new Set();

    onMount(() => {
        const rawId = page.url.searchParams.get("menu-id");
        if (rawId != null) {
            menuId = Number.parseInt(rawId);
        }
    });

    bearerToken.subscribe(token => {
        if (token == undefined) return;
        
        fetch(`${APIUrl}/api/users`, {
            headers: {
                Authorization: `Bearer: ${token}`
            },
        })
            .then((resp) => resp.json())
            .then((resp) => users = resp)
            .catch((error) => console.error(error));
    });

    const handleUserChecked = (event: Checkbox) => {
        const id = event?.currentTarget?.dataset.userid;
        if (id == undefined) return;

        const userId = Number.parseInt(id)
        
        if (event?.currentTarget?.checked) {
            selected.add(userId);
        } else {
            selected.delete(userId);
        }
    }

    const handleSendNotifications = () => {
        const selectedUserIds = Array.from(selected);

        selectedUserIds.forEach(userId => {
            icons.set(userId, faSpinner);
        });

        fetch(`${APIUrl}/api/menu/share`, {
            method: "POST",
            headers: {
                Authorization: `Bearer: ${get(bearerToken)}`,
            },
            body: JSON.stringify({
                menuId: menuId,
                users: selectedUserIds,
            }),
        }).then(resp => resp.json())
        .then((returnStatuses: Object) => {
            console.log(returnStatuses);

            for (let [userId, success] of Object.entries(returnStatuses)) {
                icons.set(Number.parseInt(userId), success ? faCheck : faXmark);
            }
        })
        .catch((error) => console.error(error))
    }
</script>

<div class="text-black dark:text-white">
    <ol>
    {#each users as user}
        <li>Id: {user.id}, Name: {user.display_name}<input type="checkbox" data-userid={user.id} onchange={handleUserChecked}>{#if icons.has(user.id)}<Fa icon={icons.get(user.id) as IconDefinition} />{/if}</li>
    {/each}
    </ol>
    
    <UiButton text="Share" action={handleSendNotifications} />
</div>