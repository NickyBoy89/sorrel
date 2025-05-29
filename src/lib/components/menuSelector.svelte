<script lang="ts">
    import { get } from "svelte/store";
    import { APIUrl } from "../../constants";
    import UiButton from "./uiButton.svelte";
    import { DateTime } from "luxon";
    import { bearerToken } from "../../routes/(app)/stores";
    import TextArea from "./ui/textArea.svelte";
    import InteractiveBox from "./ui/interactiveBox.svelte";
    import UiButtonLink from "./uiButtonLink.svelte";

    let { menuName, menuDate, menuId, canEdit = false, relativeDate = true } = $props();

    let name = $state(menuName);
    let date: Date = $state(menuDate);

    let editorOpen = $state(false);

    const handleMenuChange = async () => {
      await fetch(`${APIUrl}/api/menu/${menuId}/edit?${new URLSearchParams({
          name: name,
          date: date.toISOString().split("T")[0]
        })}`, { 
          method: "POST",
          headers: {
            Authorization: `Bearer: ${get(bearerToken)}`,
          },
        }).catch((error) => {
        console.log(error);
      });
    };

    const handlePushMenu = async () => {
      await fetch(`${APIUrl}/api/menu/share`, { 
        method: "POST",
        headers: {
          Authorization: `Bearer: ${get(bearerToken)}`,
        },
        body: JSON.stringify({
          menuId: menuId,
          users: [-1],
        })
       })
      .catch((error) => {
        console.error(error);
      })
    }

    const displayedDate = (): string => {
      if (!relativeDate) {
        return date.toISOString().split("T")[0];
      }

      const relDate = DateTime.fromJSDate(date).toRelativeCalendar();
      if (relDate === null) {
        return "unknown date"
      }
      return relDate
    }

    const toggleEditor = () => {
        editorOpen = !editorOpen;
    }
</script>

<InteractiveBox>
  <div class="flex flex-row justify-between">
    <a href="/{canEdit ? "edit-menu" : "menu"}?menu-id={menuId}" class="menu-title flex flex-col sm:flex-row grow items-center">
      <div class="menu-name text-center lg:text-left text-2xl text-black dark:text-white">
        {name}
      </div>
      <div class="menu-date mx-4 text-black dark:text-white">
        {displayedDate()}
      </div>
    </a>
    <div class="menu-options flex items-center {canEdit ? "" : "hidden"}">
      <UiButton text="Edit" action={toggleEditor}/>
    </div>
  </div>
  <div class="menu-editor {editorOpen ? "" : "hidden"} flex-col">
    <div>
      <label for="edit-menu-name">Name:</label>
      <TextArea onchange={handleMenuChange} bind:initialValue={name} />
      <label for="edit-menu-date">Date:</label>
      <input type="date" id="edit-menu-date" name="date" class="block date-editor border-1 border-neutral-700 rounded-md" value={date.toISOString().split("T")[0]} onchange={(event: Event) => {
        const dateStr = (event?.target as HTMLInputElement).value;
        const [year, monthNumber, day] = dateStr.split("-");
        date = new Date(Number.parseInt(year), Number.parseInt(monthNumber) - 1, Number.parseInt(day));
        console.log(`Date changed to: ${date}`);
        
        handleMenuChange();
      }}>
    </div>
    <UiButton text="Close" color="#458588" action={toggleEditor} />
    <UiButtonLink text="Share with" color="#d79921" href="/admin/share-menu/?menu-id={menuId}" />
  </div>
</InteractiveBox>

<style>
  input {
      background-color: var(--color-white);
      height: 2rem;
  }

  @media (prefers-color-scheme: dark) {
		input {
      background-color: var(--color-neutral-600);
    }
	}
</style>
