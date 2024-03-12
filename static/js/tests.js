import {IsNotEmpty} from "./base/utils.js"

const CREATE_TEST_URL = "/tests"

function GoToTestPage(testID) {
    window.location.href = "/p/tests/" + testID
}

Spruce.store("data", {
    showModal: false,
    title: "",

    tests: []
})

const response = await axios.get(CREATE_TEST_URL)
$store.data.tests = response.data

Spruce.store("methods", {
    showModal() {
        $store.data.title = ""
        $store.data.showModal = true
    },

    async createTest() {
        if (IsNotEmpty($store.data.title)) {
            console.log("Создаём тест! title = " + $store.data.title)

            const newTest = {
                title: $store.data.title
            }

            try {
                const response = await axios.post(CREATE_TEST_URL, newTest)
                GoToTestPage(response.data.id)
            } catch (e) {
                console.log(e)
                return null
            }
        }
    }
})

let modal = document.getElementById("myModal");

window.onmousedown = function(event) {
    if (event.target === modal) {
        $store.data.showModal = false
    }
}