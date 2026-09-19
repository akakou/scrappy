const origin = window.location.origin



const sign = (origin, period) => {
    return new Promise((resolve, reject) => {
        // setTimeout(resolve, 3000)
        chrome.runtime.sendMessage({ origin, period }, function (response) {
            if (!response.signature) {
                reject("scrappy-error: " + response.error)
                return
            }
            resolve(response.signature);
        })
    })

}

const getElemetsByNamesOnTheElement = (srcElement, name) => {
    const elements = []

    const children = srcElement.children

    for (const child of children) {
        const childName = child.getAttribute("name")

        if (childName == name) {
            elements.push(child)
        }
    }

    return elements
}

document.addEventListener('submit', function (e) {
    const targetForm = e.target

    const tags = getElemetsByNamesOnTheElement(targetForm, 'scrappy')
    if (tags.length == 0) { return }
    const tag = tags[0]

    e.preventDefault()

    const period = tag.getAttribute('period');

    (async () => {

        const signature = await sign(origin, period)
        tag.value = signature

        e.target.submit()
    })()
})

