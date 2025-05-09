<script lang="ts">
    import UiButton from "$lib/components/uiButton.svelte";
    import { onMount } from "svelte";
    import Fa from 'svelte-fa'
    import { faSpinner, type IconDefinition } from "@fortawesome/free-solid-svg-icons";
    import { APIUrl } from "../../../../constants";

    type User = {
        id: number,
        display_name: string,
    };

    let users = $state([] as Array<User>);

    let icons = $state([] as Array<null | IconDefinition>);
    let selected = [];

    onMount(() => {
        fetch(`${APIUrl}/api/users`)
            .then((resp) => resp.json())
            .then((resp) => users = resp)
            .catch((error) => console.error(error));

        selected = Array(users.length).fill(false);
        icons = Array(users.length).fill(null);
    });

    const handleUserChecked = (event: Event & { currentTarget: EventTarget & HTMLInputElement}) => {
        const ind = event?.currentTarget?.dataset.index;
        if (ind == null) return;
        selected[Number.parseInt(ind)] = event?.currentTarget?.checked;
    }

    const handleSendNotifications = () => {
        users.forEach((userId, ind) => {
            icons[ind] = faSpinner;
        });
    };
</script>

<div class="text-black dark:text-white">
    <ol>
    {#each users as user, ind}
    <li>Id: {user.id}, Name: {user.display_name}<input type="checkbox" data-index={ind} onchange={handleUserChecked}>{#if icons[ind] != null}<Fa icon={icons[ind]} />{/if}</li>
    {/each}
    </ol>
    
    <UiButton text="Share" action={() => {}}/>
</div>