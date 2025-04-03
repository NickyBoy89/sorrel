<script lang="ts">
    import Fa from 'svelte-fa'
    import { faUser } from "@fortawesome/free-solid-svg-icons";

    let { userName } = $props();
</script>

<div class="rounded-md border border-zinc-200 dark:border-zinc-700 bg-white shadow-sm dark:bg-zinc-800 py-2 px-4">
    <div class="flex flex-row items-center gap-x-2">
        <Fa icon={faUser} />
        {userName}
    </div>
</div>