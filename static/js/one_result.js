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

                const timerDiv = document.getElementById("timer_" + index)

                startTimer(timerDiv)

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


const timers = []

function startTimer(elem) {
    updateTimer()
    timer = setInterval(function () {
        updateTimer(elem)
    }, 1000);
}

function updateTimer(elem) {
    totalSeconds++

    let hours = 0
    let minutes = 0
    let remainingSeconds = 0

    elem.innerText = pad(hours) + ':' + pad(minutes) + ':' + pad(remainingSeconds);
}

function pad(val) {
    return val > 9 ? val : '0' + val;
}
