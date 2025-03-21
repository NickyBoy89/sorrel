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
            })}`)
            .catch((error) => {
            console.log(error);
        });

        onchange();
    };

    const deleteMenuItem = async () => {
        await fetch(`${backendRootURL}/api/items/${itemId}/delete`).catch((error) => console.log(error));
        onchange();
    }
</script>

<div class="menu-item-editor rounded-md p-4">
    <form action="" class="flex flex-col">
        <div class="flex flex-row justify-between items-center">
            <select name="edit-item-section" id="edit-item-section" class="rounded-sm menu-section-selector pl-2" value={section} onchange={(event) => {section = event?.target?.value; updateValues();}}>
                <option value="mains">Mains</option>
                <option value="desserts">Desserts</option>
                <option value="appetizers">Appetizers</option>
            </select>
            <div class="">
                <UiButton text="Delete" action={deleteMenuItem} />
            </div>
        </div>
        <input type="text" id="edit-item-name" name="edit-item-name" class="ounded-sm w-auto mt-2 pl-2" placeholder="Name of dish..." value={itemName} onchange={(event) => {name = event?.target?.value; updateValues();}}>
        <input type="text" id="edit-item-desc" class="rounded-sm mt-2 pl-2" placeholder="Description (Optional)" value="{itemDesc}" onchange={(event) => {description = event?.target?.value; updateValues();}}>
    </form>
</div>

<style>
    .menu-item-editor {
        background-color: #665c54;
    }

    .menu-section-selector {
        background-color: #3c3836;
    }

    #edit-item-section {
        height: 1.5rem;
    }

    input {
        background-color: #3c3836;
        height: 2rem;
    }
</style>