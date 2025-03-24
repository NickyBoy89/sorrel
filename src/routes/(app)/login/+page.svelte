<script lang="ts">
    import Fa from 'svelte-fa'
    import { faBell } from '@fortawesome/free-solid-svg-icons'
    import { handleSubscribe } from '$lib/notificationManager';
    import { goto } from '$app/navigation';

    let statusMessage = $state("");
    let statusBorder = $state("#665c54");

    let inviteCode: number;

    let notificationsPanelOpen = $state(false);

    const submissionError = (error: string) => {
        statusBorder = "#fb4934";
        statusMessage = error;
    };

    const handleSubmitCode = async (event: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement}) => {
        event.preventDefault();

        const data = new FormData(event.currentTarget);
        const formInviteCode = data.get("invite-code");

        if (formInviteCode == null) {
            submissionError("error: Code was left empty");
            return false;
        }

        inviteCode = Number.parseInt(formInviteCode as string);

        if (isNaN(inviteCode)) {
            submissionError(`error: \"${formInviteCode}\" is not a number`);
            return false;
        }

        localStorage.setItem("userId", inviteCode.toString());
        notificationsPanelOpen = true;

        return true;
    };

    const subscribeUser = async () => {
        await handleSubscribe(inviteCode);
        goto("/");
    }
</script>

<div class="flex items-center h-full welcome-dialog">
    <div>
        <h1 class="welcome-text text-center text-4xl">Welcome!</h1>
        <h2 class="text-center my-2">Thank you for trying out the grandparent meal planner!</h2>
        <div class="access-code-input rounded-md text-center">
            {#if !notificationsPanelOpen}
            <h2 class="py-8 text-lg">Please enter your invite code below</h2>
            <form class="px-8 flex flex-row gap-x-3" onsubmit={handleSubmitCode}>
                <input type="text" id="invite-code" name="invite-code" class="rounded-sm inline-block" style="border: 2px solid {statusBorder};">
                <input type="submit" value="Go" id="submit-button" class="rounded-md inline-block">
            </form>
            <div class="login-status text-center py-4 text-md">{statusMessage}</div>
            {:else}
            <div class="flex flex-col gap-y-8 py-8 px-4">
                <h2 class="text-xl">In order to deliver your menus, we need your permission to send you a notification when that happens.</h2>
                <h2 class="text-md">No rush, click the bell below to enable notifications</h2>
                <button onclick={subscribeUser} class="flex justify-center object-none text-7xl cursor-pointer rounded-md px-4 py-3"><Fa icon={faBell} id="bell-icon" /></button>
            </div>
            {/if}
        </div>
    </div>
</div>

<style>
    .access-code-input {
        background-color: #665c54;
        margin-top: 1rem;
        margin-bottom: 4rem;
    }

    @keyframes wiggle {
        0% { transform: rotate(0deg); }
        80% { transform: rotate(0deg); }
        85% { transform: rotate(5deg); }
        95% { transform: rotate(-5deg); }
        100% { transform: rotate(0deg); }
    }

    button {
        animation: wiggle 2.5s infinite;
        background-color: #458588;
        margin-left: auto;
        margin-right: auto;
    }

    #invite-code {
      background-color: #3c3836;
      height: 2rem;
      min-width: 0;
      width: 100%;
    }

    #submit-button {
        background-color: #458588;
        padding-left: 1rem;
        padding-right: 1rem;
    }
</style>