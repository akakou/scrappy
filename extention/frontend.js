setTimeout( ()=> {
    const c = confirm("attest to this page?")
    if (!c) return;

    chrome.runtime.sendMessage({ domain: document.domain }, function (response) {});
}, 1000)


chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    console.log("result", message);
});