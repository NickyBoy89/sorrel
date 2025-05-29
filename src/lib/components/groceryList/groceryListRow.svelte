<script lang="ts">
    import InteractiveBox from "../ui/interactiveBox.svelte";
    import UiButton from "../uiButton.svelte";

    let expanded = $state(true);

    let numShared = 3;

    const toggleRowExpand = () => {
        expanded = !expanded;
    };
</script>

<InteractiveBox>
    <div class="flex flex-row justify-between">
        
        <div class="flex flex-col space-y-4">
            <div class="flex flex-row items-baseline">
                <div class="flex flex-row text-2xl pl-2 pr-2 font-bold items-center">
                    Groceries:
                </div>
                <div class="text-lg font-medium mr-4">
                    Apples, Oranges, Pears...
                </div>
                <div class="text-sm">
                    &mdash; 2 items
                </div>
            </div>
            {#if expanded}
            <div class="flex flex-col text-center">
                <div class="text-lg font-bold">Contains:</div>
                <ul class="list-inside list-disc">
                    <li>Apples</li>
                    <li>Oranges</li>
                    <li>Pears</li>
                </ul>
            </div>
            {/if}
        </div>
        <div class="flex flex-col space-y-4">
            <div class="flex flex-row items-center space-x-4">
                <div>Shared with {numShared} people</div>
                <UiButton text={expanded ? "Less" : "More"} action={toggleRowExpand}/>
            </div>
            {#if expanded}
            <div class="flex flex-col text-center">
                <div class="text-lg font-bold">Shared with:</div>
                <ul class="list-inside list-disc ">
                    <li>Nicholas Novak</li>
                    <li>Your Mother</li>
                    <li>Allen Price</li>
                </ul>
            </div>
            {/if}
        </div>
    </div>
</InteractiveBox>