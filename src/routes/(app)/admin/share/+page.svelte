<script lang="ts">
    import UiButton from "$lib/components/uiButton.svelte";
    import { backendRootURL } from "../../../../constants";

</script>

{#await fetch(`${backendRootURL}/api/users`).then((resp) => resp.json())}
Loading...
{:then users}
<ol>
{#each users as user}
<li>Id: {user.id}, Name: {user.display_name}</li>
{/each}
</ol>
{:catch error}
There was an error fetching the users: {error}
{/await}

<UiButton text="Share" action={() => {}}/>