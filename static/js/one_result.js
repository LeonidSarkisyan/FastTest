const RESULT_ID = Number(window.location.pathname.split("/").pop())

Spruce.store("data", {
    access: {},
    passes: [],
    students: [],
    results: []
})

const response = await axios.get(`/results/${RESULT_ID}`)
console.log(response.data)

$store.data.results = response.data.results
$store.data.passes = response.data.passes
$store.data.students = response.data.students

const codes = document.getElementsByClassName("code")

for (let code of codes) {
    code.addEventListener("click", event => {
        navigator.clipboard.writeText(event.target.textContent)
    })
}