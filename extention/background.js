setTimeout(async function () {
    console.log("Sending:  ping");
    
    chrome.runtime.sendNativeMessage(
        "com.akakou.attestation",
        { message: "msg" }, 
        (response)=>{
            console.log("Messaging host sais: ", response);
            console.log("ERROR: ", chrome.runtime.lastError);

        });

    
}, 5000)

// });