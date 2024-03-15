import {IsNotEmpty} from "./base/utils.js"

const CREATE_GROUP_URL = "/groups"

function GoToTestPage(testID) {
    window.location.href = "/p/groups/" + testID
}

Spruce.store("data", {
    showModal: false,
    title: "",

    groups: [],

    search: ""
})

const response = await axios.get(CREATE_GROUP_URL)
$store.data.groups = response.data

Spruce.store("methods", {
    showModal() {
        $store.data.title = ""
        $store.data.showModal = true
    },

    async createGroup() {
        if (IsNotEmpty($store.data.title)) {
            const newTest = {
                title: $store.data.title
            }

            try {
                const response = await axios.post(CREATE_GROUP_URL, newTest)
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