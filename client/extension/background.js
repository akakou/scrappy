chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
    chrome.runtime.sendNativeMessage(
        "com.akakou.attestation",
        { nonce: 'request.nonce', origin: 'sender.origin' },
        (response) => {
            console.log("Messaging host sais: ", response);
            console.log("ERROR: ", chrome.runtime.lastError);

            // sendResponse(response);
        });

    return true;
});