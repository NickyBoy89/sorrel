<script lang="ts">
    import { APIUrl } from "../../../../constants";

</script>

<div class="text-black dark:text-white">
    {#await fetch(`${APIUrl}/api/users`).then((resp) => resp.json())}
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
</div>