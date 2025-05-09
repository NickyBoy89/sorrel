import { APIUrl } from "../constants";

// This function is needed because Chrome doesn't accept a base64 encoded string
// as value for applicationServerKey in pushManager.subscribe yet
// https://bugs.chromium.org/p/chromium/issues/detail?id=802280
function urlBase64ToUint8Array(base64String: string) {
    var padding = '='.repeat((4 - base64String.length % 4) % 4);
    var base64 = (base64String + padding)
      .replace(/\-/g, '+')
      .replace(/_/g, '/');
   
    var rawData = window.atob(base64);
    var outputArray = new Uint8Array(rawData.length);
   
    for (var i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
}

export const isSubscriptionValid = async (): Promise<boolean> => {
    let valid = false;

    const registration = await navigator.serviceWorker.ready;

    const subJson = await registration.pushManager.getSubscription().then((sub) => sub?.toJSON()).catch((error) => console.error(error));

    if (subJson == null) {
        return false;
    }

    await fetch(`${APIUrl}/api/push/validate`, {
        method: "POST",
        body: JSON.stringify(subJson),
    }).then((resp) => {
        if (resp.ok) {
            valid = true;
        } else if (resp.status == 404) {
            valid = false;
        } else {
            console.error(`failed to validate subcription: code ${resp.status}, message: ${resp.statusText}`);
            return false;
        }
    }).catch((error) => {
        console.error(error);
    });

    return valid;
}

export const handleSubscribe = async (inviteCode: number) => {
    const vapidKey = await fetch(`${APIUrl}/api/push/public-key`).then((resp) => resp.text());

    const registration = await navigator.serviceWorker.ready;

    await registration.pushManager.getSubscription().then((sub) => {
        sub?.unsubscribe().catch((error) => {
            console.error("Error unsubscribing", error);
            return;
        })
    })

    const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(vapidKey)
    }).catch((error) => {
        console.log(error);
    });

    console.log("Subscribed user...");

    await fetch(`${APIUrl}/api/push/subscribe`, {
        method: "POST",
        body: JSON.stringify({
            userId: inviteCode,
            sub: subscription?.toJSON()
        }),
    });
}