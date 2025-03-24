import { backendRootURL } from "../constants";

export const isSubscriptionValid = async (): Promise<boolean> => {
    const registration = await navigator.serviceWorker.ready;

    const subJson = await registration.pushManager.getSubscription().then((sub) => sub?.toJSON()).catch((error) => console.error(error));

    if (subJson == null) {
        return false;
    }

    const resp = await fetch(`${backendRootURL}/api/push/validate`, {
        method: "POST",
        body: JSON.stringify(subJson),
    })

    if (resp.status != 200) {
        console.error(`failed to validate subcription: code ${resp.status}, message: ${resp.statusText}`);
        return false;
    }

    return true;
}

export const handleSubscribe = async (inviteCode: number) => {
    const vapidKey = await fetch(`${backendRootURL}/api/push/public-key`).then((resp) => resp.text());

    const registration = await navigator.serviceWorker.ready;

    await registration.pushManager.getSubscription().then((sub) => {
        sub?.unsubscribe().catch((error) => {
            console.error("Error unsubscribing", error);
            return;
        })
    })

    const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: vapidKey
    }).catch((error) => {
        console.log(error);
    });

    console.log("Subscribed user...");

    await fetch(`${backendRootURL}/api/push/subscribe`, {
        method: "POST",
        body: JSON.stringify({
            userId: inviteCode,
            sub: subscription?.toJSON()
        }),
    });
}