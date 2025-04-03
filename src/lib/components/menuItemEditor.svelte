<script lang="ts">
    import { backendRootURL, menuDefaultSections } from "../../constants";
    import UiButton from "./uiButton.svelte";

    let { itemId, name: itemName = "", description: itemDesc = "", section: itemSection = "mains", onchange = () => {}} = $props();

    let name = itemName;
    let description = itemDesc;
    let section = $state(itemSection);

    const updateValues = async () => {
        await fetch(`${backendRootURL}/api/items/${itemId}/edit?${new URLSearchParams({
                name: name,
                description: description,
                section: section
            })}`, {
                method: "POST",
            })
            .catch((error) => {
            console.log(error);
        });

        onchange();
    };

    const deleteMenuItem = async () => {
        await fetch(`${backendRootURL}/api/items/${itemId}/delete`, { method: "POST" }).catch((error) => console.log(error));
        onchange();
    }
</script>

<div class="rounded-md p-4 bg-white dark:bg-zinc-800 border border-zinc-700">
    <form action="" class="flex flex-col">
        <div class="flex flex-row justify-between items-center">
            <select name="edit-item-section" id="edit-item-section" class="rounded-sm text-neutral-900 dark:text-white pl-2" value={section} onchange={(event) => {section = event?.target?.value; updateValues();}}>
                <option value="mains">Mains</option>
                <option value="desserts">Desserts</option>
                <option value="appetizers">Appetizers</option>
            </select>
            <div class="">
                <UiButton text="Delete" action={deleteMenuItem} />
            </div>
        </div>
        <input type="text" id="edit-item-name" name="edit-item-name" class="rounded-sm w-auto text-neutral-900 dark:text-white mt-2 pl-2" placeholder="Name of dish..." value={itemName} onchange={(event) => {name = event?.target?.value; updateValues();}}>
        <input type="text" id="edit-item-desc" class="rounded-sm text-neutral-900 dark:text-white mt-2 pl-2" placeholder="Description (Optional)" value="{itemDesc}" onchange={(event) => {description = event?.target?.value; updateValues();}}>
    </form>
</div>

<style>
    input,select {
        background-color: var(--color-white);
        border: 1px solid var(--color-zinc-700);
        height: 2rem;
    }

    @media (prefers-color-scheme: dark) {
        input,select {
            background-color: var(--color-zinc-900);
            border: 1px solid var(--color-zinc-700);
            height: 2rem;
        }
	}
</style>