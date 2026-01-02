<script lang="ts">
  import CaretDown from "phosphor-svelte/lib/CaretDown";
  import { APIUrl } from "../../../constants";

  type SelectOption = {
    value: string;
    display: string;
  };

  let expanded = $state(false);

  let content = $state("");

  const handleOptionSelected = (option: SelectOption) => {
    content = option.display;
    onchange(option.value);
  };

  let {
    placeholder,
    options = [],
    onchange,
  }: {
    placeholder?: string;
    options: Array<SelectOption>;
    onchange: (updatedValue: string) => void;
  } = $props();
</script>

<div
  class="flex flex-row grow w-auto min-h-8 pl-2 rounded-sm bg-white dark:bg-neutral-900 border border-neutral-700 items-center"
>
  <div class="flex flex-col grow">
    <input
      type="text"
      class="flex grow"
      {placeholder}
      bind:value={content}
      onfocus={() => (expanded = true)}
      onblur={() => {
        // This needs to be delayed a bit, because if the user takes an action somewhere else, we need to give them a bit until we unexpand and remove the options
        setTimeout(() => (expanded = false), 200);
      }}
    />
    {#if expanded}
      {#if content.length > 0 && options.find((option) => option.display === content) === undefined}
        <button
          class="flex"
          onclick={async () => {
            await fetch(`${APIUrl}/api/v1/ingredients`, {
              method: "POST",
              body: JSON.stringify({ name: content }),
            })
              .then((resp) => resp.text())
              .then((createdId) => {
                const newOption: SelectOption = {
                  value: createdId,
                  display: content,
                };
                options.push(newOption);
                handleOptionSelected(newOption);
              });
          }}>Add "{content}" as new ingredient...</button
        >
      {/if}
      {#each options as option}
        <button class="flex" onclick={() => handleOptionSelected(option)}
          >{option.display}</button
        >
      {/each}
    {/if}
  </div>
  <button
    class={["select-caret", { expanded }]}
    onclick={() => (expanded = !expanded)}
  >
    <CaretDown />
  </button>
</div>

<style>
  .select-caret {
    transition: rotate 0.2s ease-in-out;
  }

  .select-caret.expanded {
    rotate: 180deg;
  }
</style>
