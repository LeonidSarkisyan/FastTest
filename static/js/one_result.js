const RESULT_ID = Number(window.location.pathname.split("/").pop())

Spruce.store("data", {
    access: {},
    passes: [],
    students: [],
})

const response = await axios.get(`/results/${RESULT_ID}`)

$store.data.passes = response.data.passes
$store.data.students = response.data.students