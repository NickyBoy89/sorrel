<script lang="ts">
    import { get } from "svelte/store";
    import { APIUrl } from "../../constants";
    import { bearerToken } from "../../routes/(app)/stores";
    import InteractiveBox from "./ui/interactiveBox.svelte";
    import TextArea from "./ui/textArea.svelte";
    import UiButton from "./uiButton.svelte";

    let { itemId, name: itemName = "", description: itemDesc = "", section: itemSection = "mains", onchange = () => {}} = $props();

    let name = itemName;
    let description = itemDesc;
    let section = $state(itemSection);

    const updateValues = async () => {
        await fetch(`${APIUrl}/api/items/${itemId}/edit?${new URLSearchParams({
                name: name,
                description: description,
                section: section
            })}`, {
                method: "POST",
                headers: {
                    Authorization: `Bearer: ${get(bearerToken)}`
                },
            })
            .catch((error) => {
            console.log(error);
        });

        onchange();
    };

    const deleteMenuItem = async () => {
        await fetch(`${APIUrl}/api/items/${itemId}/delete`, { method: "POST", headers: { Authorization: `Bearer: ${get(bearerToken)}` } }).catch((error) => console.log(error));
        onchange();
    }
</script>

<InteractiveBox>
    <div class="flex flex-row justify-between items-center">
        <select name="edit-item-section" id="edit-item-section" class="rounded-sm text-neutral-900 dark:text-white pl-2" value={section} onchange={(event) => {section = event?.currentTarget?.value; updateValues();}}>
            <option value="mains">Mains</option>
            <option value="desserts">Desserts</option>
            <option value="appetizers">Appetizers</option>
        </select>
        <div class="">
            <UiButton text="Delete" action={deleteMenuItem} />
        </div>
    </div>
    <TextArea placeholder="Name of dish..." initialValue={itemName} onchange={(event: Event & { currentTarget: EventTarget & HTMLInputElement }) => {name = event?.currentTarget?.value; updateValues();}} />
    <TextArea placeholder="Description (Optional)" initialValue={itemDesc} onchange={(event: Event & { currentTarget: EventTarget & HTMLInputElement }) => {description = event?.currentTarget?.value; updateValues();}} />
</InteractiveBox>

<style>
    select {
        background-color: var(--color-white);
        border: 1px solid var(--color-zinc-700);
        height: 2rem;
    }

    @media (prefers-color-scheme: dark) {
        select {
            background-color: var(--color-zinc-900);
            border: 1px solid var(--color-zinc-700);
            height: 2rem;
        }
	}
</style>