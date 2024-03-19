const URL_CHAPTERS = window.location.href.split("/")
const RESULT_ID = URL_CHAPTERS[URL_CHAPTERS.length - 3]
const PASS_ID = URL_CHAPTERS[URL_CHAPTERS.length - 1]

console.log(RESULT_ID, PASS_ID)

Spruce.store("data", {
    questions: [],
    minutes: 0,
})

const startButton = document.getElementById("startButton")

startButton.onclick = async () => {
    try {
        const response = await axios.get(`/passing/${RESULT_ID}/solving/${PASS_ID}/questions`)
        $store.data.questions = response.data.questions
        console.log(response.data)
    } catch (e) {
        alert(e.response.data)
    }
}
