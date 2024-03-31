const logout = document.getElementById("logout")
const change = document.getElementById("change")
const reset = document.getElementById("reset")
const changeGo = document.getElementById("changeGo")
const success = document.getElementById("success")

reset.style.display = "none"

logout.onclick = async function (event) {
    console.log("/auth/logout")
    try {
        const response = await axios.get("/auth/logout")
        window.location.href = "/auth"
    } catch (e) {
        console.log(e)
    }
}

change.onclick = async function (event) {
    try {
        const response = await axios.post("/auth/reset")
        console.log(response)
    } catch (e) {
        console.log(e)
    }
    reset.style.display = "block"
}

changeGo.onclick = async function (event) {
    const password = document.getElementById("newPassword").value
    const passwordRepeat = document.getElementById("newPasswordRepeat").value

    if (password !== passwordRepeat) {
        alert("пароли должны совпадать")
        return
    }

    try {

        const response = await axios.post(`/auth/change/${document.getElementById("code").value}`, {
            password: password
        })

        console.log(response)

        reset.style.display = "none"
        success.style.display = "block"

        setTimeout(function () {
            success.style.opacity = 0
        }, 3000)

        setTimeout(function () {
            success.style.display = "none"
        }, 3250)

    } catch (e) {
        alert(e.response.data)
    }
}