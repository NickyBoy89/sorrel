<script lang="ts">
    let { placeholder = "", initialValue = $bindable(""), onchange } = $props();
</script>

<input type="text" class="block rounded-sm w-auto text-neutral-900 dark:text-white mt-2 pl-2" placeholder={placeholder} bind:value={initialValue} onchange={onchange}>

<style>
    input {
        background-color: var(--color-white);
        border: 1px solid var(--color-zinc-700);
        height: 2rem;
    }

    @media (prefers-color-scheme: dark) {
        input {
            background-color: var(--color-zinc-900);
            border: 1px solid var(--color-zinc-700);
            height: 2rem;
        }
    }
</style>