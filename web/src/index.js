const longURLField = document.getElementById("long-url")
const shortURLField = document.getElementById("short-url")
const shortenBtn = document.getElementById("shorten-btn")

const getShortURL = async (longURL) => {
    let data
    try {
        const res = await fetch("https://goly-backend.onrender.com/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Access-Control-Allow-Origin": "*",
            },
            // mode: "no-cors",
            body: JSON.stringify({ "long_url": longURL })
        })
        data = res.json()
    } catch (e) {
        console.log(e);
    }
    console.log(data)
}

shortenBtn.addEventListener("click", async (event) => {
    event.preventDefault()
    shortenBtn.disabled = true

    const longURL = longURLField.value
    await getShortURL(longURL)
    shortenBtn.disabled = false
})