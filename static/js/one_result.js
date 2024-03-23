const RESULT_ID = Number(window.location.pathname.split("/").pop())

Spruce.store("data", {
    access: {},
    passes: [],
    students: [],
    results: [],

    hideCodes: false,
})

Spruce.store("methods", {
    async Reset(passID, index) {
        const response = await axios.patch(`/results/${RESULT_ID}/reset/${passID}`)
        console.log(response)

        $store.data.passes[index].is_activated = false
        $store.data.results[index].access_id = 0
        $store.data.results[index].mark = 0
    },

    TimeProcess(seconds) {
        let hours = Math.floor(seconds / 3600);
        let minutes = Math.floor((seconds % 3600) / 60);
        let remainingSeconds = seconds % 60;

        hours = pad(hours)
        minutes = pad(minutes)
        remainingSeconds = pad(remainingSeconds)

        return hours + ':' + minutes + ':' + remainingSeconds;
    },

    CopyCode(code) {
        navigator.clipboard.writeText(code)

        let notification = document.getElementById("notification");
        notification.style.display = "block";
        setTimeout(function(){
            notification.style.display = "none";
        }, 3000);
    },

    ToggleHideCodes() {
        $store.data.hideCodes = !$store.data.hideCodes
    }
})

const response = await axios.get(`/results/${RESULT_ID}`)
console.log(response.data)

$store.data.results = response.data.results
$store.data.passes = response.data.passes
$store.data.students = response.data.students


const socket = new WebSocket(`ws://81.200.149.16:8080/results/${RESULT_ID}/ws`);

socket.onopen = async function(event) {
    console.log('WebSocket connected');
    socket.send('Hello, server!');
};

socket.onclose = function(event) {
    console.log('WebSocket disconnected');
    console.log(event)
};

socket.onerror = (event) => {
    console.log(event)
}

socket.onmessage = async function(event) {
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
