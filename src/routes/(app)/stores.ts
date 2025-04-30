import { writable, type Writable } from "svelte/store";

export const bearerToken: Writable<string | undefined> = writable(undefined);