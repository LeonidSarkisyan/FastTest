Spruce.store("data", {
    results: [],

    search: ""
})

const response = await axios.get("/results")
$store.data.results = response.data