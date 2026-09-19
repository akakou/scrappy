chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
    chrome.runtime.sendNativeMessage(
        "com.akakou.scrappy",
        { origin: sender.origin, period: Number(request.period) },
        (response) => {
            console.log("origin: ", sender.origin)
            console.log("period: ", sender.period)
            console.log("Messaging host sais: ", response);
            console.log("ERROR: ", chrome.runtime.lastError);

            sendResponse(response);
        });

    return true;
});