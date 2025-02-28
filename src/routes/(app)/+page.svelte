<script lang="ts">
  import MenuRenderer from "$lib/components/menuRenderer.svelte";
</script>

<div class="grid grid-cols-2">
  <section class="menu-editor">
    This is some example text
  </section>
  
  <section class="menu-preview">
    <MenuRenderer menuName="Tomorrow Night's Dinner" />
  </section>
</div>