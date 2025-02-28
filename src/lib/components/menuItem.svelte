<script lang="ts">
  let {name, description} = $props();
</script>

<div class="menu-item py-1">
  <div class="item-name">{name}</div>
  <div class="item-description">{description}</div>
</div>

<style>
  .item-name {
    font-size: 2rem;
  }

  .item-description {
    font-size: 1.25rem;
    font-style: italic;
  }
</style>