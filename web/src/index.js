const longURLField = document.getElementById("long-url")
const shortURLField = document.getElementById("short-url")
const shortenBtn = document.getElementById("shorten-btn")
const URI = "https://goly-backend.onrender.com"
const createURL = URI + "/create"

const getShortURL = async (longURL) => {
    let data
    try {
        const res = await fetch(createURL, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            mode: "cors",
            body: JSON.stringify({ "long_url": longURL })
        })
        data = await res.json()
    } catch (e) {
        alert(e)
    } finally {
        return data
    }
}

const isValidReq = () => {
    return longURLField.value !== ""
}

shortenBtn.addEventListener("click", async (event) => {
    event.preventDefault()
    shortenBtn.disabled = true

    if (!isValidReq()) {
        shortenBtn.disabled = false
        alert("Enter valid URL")
        return
    }

    const data = await getShortURL(longURLField.value)
    if (!data) {
        shortenBtn.disabled = false
        return
    }

    shortURLField.value = URI + "/" + data["short_url"]
    shortenBtn.disabled = false
})