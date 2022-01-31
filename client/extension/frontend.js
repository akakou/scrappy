var attestation = null;
var state = 'ok'; 

const period = 24 * 60 * 60 * 1000


window.onload = async (_) => {
    await cookieStore.delete('attestation')
    chrome.runtime.sendMessage({ domain: document.domain }, function (response) { });
}

window.onsubmit = async () => {
    const answer = confirm('Attest not attacking?')

    if (answer)
        await cookieStore.set({
            name: "attestation",
            value: attestation,
            expires: Date.now() + period,
        })
}

chrome.runtime.onMessage.addListener(async (message, sender, sendResponse) => {
    console.log(sender)
    console.log("attestation", message)
    attestation = encodeURI(message)
});