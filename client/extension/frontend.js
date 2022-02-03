var attestation = null;
var state = 'ok'; 

const period = 24 * 60 * 60 * 1000


document.addEventListener('onAttest', function (e) {
    const answer = confirm('Attest not attacking?')
    if (!answer) return;

    nonce = document.getElementById('nonce').value
    console.log('nonce:', nonce)

    chrome.runtime.sendMessage({ domain: document.domain, nonce: nonce }, function (response) {
        if (!response.signature || !response.counter)
            alert(response)

        console.log("attestation", response)
        attestation = {
            signature: encodeURI(response.signature),
            counter: encodeURI(response.counter)
        }

        e.target.value = JSON.stringify(attestation)
    })
})

