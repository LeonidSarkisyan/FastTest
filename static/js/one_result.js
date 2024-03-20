const RESULT_ID = Number(window.location.pathname.split("/").pop())

Spruce.store("data", {
    access: {},
    passes: [],
    students: [],
    results: []
})

Spruce.store("methods", {
    TimeProcess(seconds) {
        let hours = Math.floor(seconds / 3600);
        let minutes = Math.floor((seconds % 3600) / 60);
        let remainingSeconds = seconds % 60;

        hours = pad(hours)
        minutes = pad(minutes)
        remainingSeconds = pad(remainingSeconds)

        return hours + ':' + minutes + ':' + remainingSeconds;
    }
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

const socket = new WebSocket(`ws://localhost:8080/results/${RESULT_ID}/ws`);

socket.onopen = function(event) {
    console.log('WebSocket connected');
    socket.send('Hello, server!');
};

socket.onclose = function(event) {
    console.log('WebSocket disconnected');
};

socket.onmessage = function(event) {
    const newResult = JSON.parse(event.data)

    let index = 0

    for (let p of $store.data.passes) {
        if (p.id === newResult.pass_id) {
            if (newResult.mark === -1) {
                $store.data.passes[index].is_activated = true
                $store.data.results[index].time_pass = newResult.time_pass
                break
            } else {
                $store.data.passes[index].is_activated = true
                $store.data.results[index] = newResult
                break
            }
        }
        index++
    }

    console.log('Message received:', newResult);
};

function pad(val) {
    return val > 9 ? val : '0' + val;
}
