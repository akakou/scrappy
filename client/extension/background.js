chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
    console.log(request.domain)
    chrome.runtime.sendNativeMessage(
        "com.akakou.attestation",
        { nonce: request.nonce, domain: request.domain },
        (response) => {
            console.log("Messaging host sais: ", response);
            console.log("ERROR: ", chrome.runtime.lastError);

            chrome.tabs.sendMessage(sender.tab.id, response);
        });

    sendResponse({});
});