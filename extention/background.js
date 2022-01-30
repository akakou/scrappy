chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
    chrome.runtime.sendNativeMessage(
        "com.akakou.attestation",
        { message: "msg" }, 
        (response) => {
            console.log("Messaging host sais: ", response);
            console.log("ERROR: ", chrome.runtime.lastError);

            chrome.tabs.sendMessage(sender.tab.id, response);
        });

    sendResponse({});
});