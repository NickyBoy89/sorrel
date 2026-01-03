<script lang="ts">
  import { onMount } from "svelte";
  import { APIUrl } from "../../../../constants";

  onMount(() => {
    window.location.href = `${APIUrl}/api/auth/login`;
  });
</script>
