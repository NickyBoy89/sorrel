<script lang="ts">
    let { children } = $props();
</script>

<div class="rounded-md p-4 bg-white dark:bg-zinc-800 border border-zinc-700">
    <div class="flex flex-col">
        {@render children?.()}
    </div>
</div>