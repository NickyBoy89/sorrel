<script lang="ts">
  import * as arctic from "arctic";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { keycloak } from "$lib/auth";
  import Cookies from "js-cookie";

  onMount(async () => {
    const code = page.url.searchParams.get("code");
    const state = page.url.searchParams.get("state");

    const storedState = Cookies.get("state");
    const storedCodeVerifier = Cookies.get("code_verifier");

    if (
      code === null ||
      storedState == null ||
      state !== storedState ||
      storedCodeVerifier == null
    ) {
      // 400
      throw new Error("Invalid request");
    }

    try {
      const tokens = await keycloak.validateAuthorizationCode(
        code,
        storedCodeVerifier,
      );
      const accessToken = tokens.accessToken();
    } catch (e) {
      if (e instanceof arctic.OAuth2RequestError) {
        // Invalid authorization code, credentials, or redirect URI
        const code = e.code;
        // ...
      }
      if (e instanceof arctic.ArcticFetchError) {
        // Failed to call `fetch()`
        const cause = e.cause;
        // ...
      }
      // Parse error
    }
  });
</script>
