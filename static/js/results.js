Spruce.store("data", {
    results: [],

    tests: [],
    groups: [],

    search: "",

    selectedTestID: 0,
    selectedGroupID: 0,
})

Spruce.store("methods", {
    FilterResults(event, type) {
        switch (type) {
            case "test":
                console.log(`test id = ${event.target.value}`)
                console.log(`group id = ${$store.data.selectedGroupID}`)
                break
            case "group":
                console.log(`test id = ${$store.data.selectedTestID}`)
                console.log(`group id = ${event.target.value}`)
                break
        }
    }
})

let response = await axios.get("/results")
$store.data.results = response.data

response = await axios.get("/tests")
$store.data.tests = response.data

response = await axios.get("/groups")
$store.data.groups = response.data

console.log($store.data)