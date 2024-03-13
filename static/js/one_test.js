import {AddTestID, ChangeIndex, GetIndex} from "./base/localstorage.js"

const TEST_ID = Number(
    window.location.pathname.split("/")[window.location.pathname.split("/").length - 1]
)

AddTestID(TEST_ID, 0)

const QUESTION_URL = window.location.pathname.replace("/p", "") + "/questions";

function QUESTION_WITH_ID_URL(questionID) {
    return QUESTION_URL + "/" + questionID + "/" + "answers"
}

Spruce.store("data", {
    showModal: false,
    inputElement: null,

    mainTitle: "",
    title: "",
    currentIndex: GetIndex(TEST_ID),

    questions: [],
    answers: []
})

let response = await axios.get("/tests/" + TEST_ID)
$store.data.title = response.data.title
$store.data.mainTitle = response.data.title

response = await axios.get(QUESTION_URL)
$store.data.questions = response.data

response = await axios.get(QUESTION_WITH_ID_URL($store.data.questions[$store.data.currentIndex].id))
$store.data.answers = response.data

Spruce.store("methods", {
    async chooseIndex(index) {
        const response = await axios.get(QUESTION_WITH_ID_URL($store.data.questions[index].id))
        $store.data.answers = response.data
        console.log(response.data)
        $store.data.currentIndex = index
        ChangeIndex(TEST_ID, index)
    },

    async addQuestion() {
        const response = await axios.post(QUESTION_URL)
        $store.data.questions.push({
            id: response.data.id,
            text: ""
        })
        await this.chooseIndex($store.data.questions.length - 1)
        const input = document.getElementById("input")
        input.focus()
    },

    async deleteQuestion(index, questionID) {
        if (index === 0 && $store.data.questions.length === 1) {
            return
        }

        let indexFromLocalStorage = GetIndex(TEST_ID)

        if (indexFromLocalStorage === index) {
            if (index !== 0) {
                await this.chooseIndex(index - 1)
            }
        }

        if (indexFromLocalStorage === $store.data.questions.length - 1) {
            await this.chooseIndex(indexFromLocalStorage - 1)
        }

        try {
            const response = await axios.delete(
                QUESTION_WITH_ID_URL(questionID).replace("/answers", "")
            )

            $store.data.questions = $store.data.questions.filter(question => {
                return question.id !== questionID
            })

            console.log(response)
        } catch (e) {
            console.log(e)
        }

        if (index === 0) {
            await this.chooseIndex(0)
        }
    },

    async updateTextQuestion() {
        let questionID = $store.data.questions[$store.data.currentIndex].id
        let text = $store.data.questions[$store.data.currentIndex].text

        let body = {
            text: text
        }

        try {
            const response = await axios.patch(
                QUESTION_WITH_ID_URL(questionID).replace("/answers", ""), body
            )
            console.log(response)
        } catch (e) {
            console.log(e)
        }
    },

    async addAnswer() {
        let questionID = $store.data.questions[$store.data.currentIndex].id
        const response = await axios.post(QUESTION_WITH_ID_URL(questionID))
        $store.data.answers.push({
            id: response.data.id,
            text: ""
        })
    },

    async updateAnswer(answerIndex) {
        let questionID = $store.data.questions[$store.data.currentIndex].id
        let answer = $store.data.answers[answerIndex]

        const body = {
            text: answer.text,
            is_correct: answer.is_correct
        }

        console.log(body)

        const response = await axios.patch(QUESTION_WITH_ID_URL(questionID) + "/" + answer.id, body)
        console.log(response.data)
    },

    async focusAnswer(index) {
        let inputs = document.getElementsByClassName("answer__input")

        if (index < inputs.length - 1) {
            inputs[index + 1].focus()
        } else {
            await this.addAnswer()
            let inputs = document.getElementsByClassName("answer__input")
            inputs[inputs.length - 1].focus()
        }
    },

    async deleteAnswer(id) {
        if ($store.data.answers.length < 3) {
            return
        }

        let questionID = $store.data.questions[$store.data.currentIndex].id

        try {
            const response = await axios.delete(QUESTION_WITH_ID_URL(questionID) + "/" + id)

            $store.data.answers = $store.data.answers.filter(answer => {
                return answer.id !== id
            })
        } catch (e) {
            console.log(e)
        }
    },

    showModal() {
        $store.data.mainTitle = $store.data.title
        $store.data.showModal = true
    },

    async UpdateTitleTest() {
        const body = {
            title: $store.data.mainTitle
        }

        try {
            const response = await axios.patch("/tests/" + TEST_ID, body)
            console.log(response)

            $store.data.title = $store.data.mainTitle
            $store.data.showModal = false
        } catch (e) {
            console.log(e)
        }
    }
})

let modal = document.getElementById("myModal");

window.onmousedown = function(event) {
    if (event.target === modal) {
        $store.data.showModal = false
    }
}