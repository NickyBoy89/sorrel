<script lang="ts">
  import { goto } from "$app/navigation";
  import UiButton from "$lib/components/uiButton.svelte";
</script>

<UiButton text="Log in" action={() => goto("/auth/login")} />
