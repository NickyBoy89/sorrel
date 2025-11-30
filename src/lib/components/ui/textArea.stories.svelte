<script module>
  import { defineMeta } from "@storybook/addon-svelte-csf";
  import { fn } from "storybook/test";

  import TextArea from "./textArea.svelte";

  const { Story } = defineMeta({
    title: "TextArea",
    component: TextArea,
    args: {
      onchange: fn(),
    },
    parameters: {
      backgrounds: {
        options: {},
      },
    },
  });
</script>

<Story
  name="Basic"
  args={{ placeholder: "This is a placeholder" }}
  globals={{ backgrounds: { value: "dark" } }}
/>
