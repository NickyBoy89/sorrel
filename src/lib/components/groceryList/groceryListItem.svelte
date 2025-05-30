<script lang="ts">
  let { text, checked = false } = $props();
</script>

<form class="">
    <label class="flex flex-row space-x-2 cursor-pointer select-none align-middle font-semibolds">
    <input type="checkbox" name="checked" id="checked" class="circle-checkbox cursor-pointer my-auto" checked={checked}>
    <span class="align-middle text-lg w-full grocery-item-text header-underlined py-2">{text}</span>
    </label>
</form>

<style>
    .circle-checkbox {
      appearance: none;
      width: 1.25rem;
      height: 1.25rem;
      border: 0.1em solid var(--color-neutral-700);
      border-radius: 50%;
      cursor: pointer;
      position: relative;
      transition: background 0.1s, border 0.1s;
    }

    .circle-checkbox:checked {
      background: #d65d0e;
      border: 2px solid #d65d0e;
    }

    .circle-checkbox:checked::after {
      content: '';
      position: absolute;
      top: 50%;
      left: 50%;
      width: 0.35em;
      height: 0.65em;
      border: solid white;
      border-width: 0 0.15em 0.15em 0;
      transform: translate(-50%, -55%) rotate(45deg);
    }

    .grocery-item-text {
        transition: color 0.1s;
    }

    .circle-checkbox:checked ~ .grocery-item-text {
        color: var(--color-neutral-500);
    }

    .header-underlined {
        border-bottom: 1px solid var(--color-neutral-800);
    }
</style>