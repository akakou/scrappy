var attestation = null;
var state = 'ok';

// const period = 24 * 60 * 60 * 1000
const origin = window.location.origin


document.addEventListener('onAttest', function (e) {
    const answer = confirm('Do you want to attest?')
    if (!answer) return;

    period = document.getElementById('period').value
    console.log('period:', period)

    chrome.runtime.sendMessage({ origin, period }, function (response) {
        if (!response.signature) {
            alert(response)
            return;
        }

        console.log("attestation", response)

        // attestation = {
        // signature: encodeURI(response.signature),
        // counter: encodeURI(response.counter)
        // }

        e.target.value = response.signature
        // encodeURI(response.signature)
        // JSON.stringify(attestation)
    })
})

