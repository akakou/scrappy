var attestation = null;
var state = 'ok'; 

const period = 24 * 60 * 60 * 1000


window.onload = async () => {
    nonce = await cookieStore.get('nonce')
    chrome.runtime.sendMessage({ domain: document.domain, nonce: nonce.value }, function (response) { });
}

window.onsubmit = async () => {
    const answer = confirm('Attest not attacking?')

    if (answer){
        await cookieStore.set({
            name: "signature",
            value: attestation.signature,
            expires: Date.now() + period,
        })

        await cookieStore.set({
            name: "counter",
            value: attestation.counter,
            expires: Date.now() + period,
        })
    }     

    console.log('cookie: ', document.cookie)
    alert(1)
}

chrome.runtime.onMessage.addListener(async (message, sender, sendResponse) => {
    console.log(sender)
    console.log("attestation", message)
    attestation = { 
        signature: encodeURI(message.signature), 
        counter: encodeURI(message.counter)
    }
});