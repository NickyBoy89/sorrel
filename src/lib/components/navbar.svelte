<script lang="ts">
  import UiButtonLink from "./uiButtonLink.svelte";

  let {backlinkHref = "/"} = $props();
</script>

<nav class="flex flex-row py-2">
  <div class="flex">
    <UiButtonLink text="← Back" href={backlinkHref} />
  </div>
  <!-- <div class="flex-grow">
      <a href="/">
          <img src="/static/img/logo.svg" width="72px" alt="">
      </a>
  </div>
  <div class="flex flex-row items-center space-x-4 underline">
      <a href="/static/Resume-Nicholas-Novak.pdf">Resume</a>
      <a href="/interactive">Interactive</a>
  </div> -->
</nav>