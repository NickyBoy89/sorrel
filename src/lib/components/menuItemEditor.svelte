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

<div class="rounded-md p-4 bg-white dark:bg-neutral-700 border border-gray-300 dark:border-none">
    <form action="" class="flex flex-col">
        <div class="flex flex-row justify-between items-center">
            <select name="edit-item-section" id="edit-item-section" class="rounded-sm text-neutral-900 dark:text-white bg-neutral-300 dark:bg-neutral-800 pl-2" value={section} onchange={(event) => {section = event?.target?.value; updateValues();}}>
                <option value="mains">Mains</option>
                <option value="desserts">Desserts</option>
                <option value="appetizers">Appetizers</option>
            </select>
            <div class="">
                <UiButton text="Delete" action={deleteMenuItem} />
            </div>
        </div>
        <input type="text" id="edit-item-name" name="edit-item-name" class="rounded-sm w-auto bg-neutral-300 dark:bg-neutral-800 text-neutral-900 dark:text-white mt-2 pl-2" placeholder="Name of dish..." value={itemName} onchange={(event) => {name = event?.target?.value; updateValues();}}>
        <input type="text" id="edit-item-desc" class="rounded-sm bg-neutral-300 dark:bg-neutral-800 text-neutral-900 dark:text-white mt-2 pl-2" placeholder="Description (Optional)" value="{itemDesc}" onchange={(event) => {description = event?.target?.value; updateValues();}}>
    </form>
</div>

<style>
    #edit-item-section {
        height: 1.5rem;
    }

    input {
        height: 2rem;
    }
</style>