const RESULT_ID = Number(window.location.pathname.split("/").pop())

Spruce.store("data", {
    passes: [],
    students: [],
})

const response = await axios.get()