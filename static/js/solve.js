const URL_CHAPTERS = window.location.href.split("/")
const RESULT_ID = URL_CHAPTERS[URL_CHAPTERS.length - 3]
const PASS_ID = URL_CHAPTERS[URL_CHAPTERS.length - 1]

console.log(RESULT_ID, PASS_ID)

Spruce.store("data", {
    questions: [],
    minutes: 0,

    isPass: false,
    isComplete: false,
    currentQuestionIndex: 0,

    showModal: false,

    error: "",

    result: {}
})

Spruce.store("methods", {
    changeIndex(index) {
        if (index === $store.data.questions.length || index < 0) {
            return null
        }

        $store.data.currentQuestionIndex = index
    },

    SaveAnswer(event, index, aIndex) {
        const questions = JSON.parse(localStorage.getItem("questions"))

        console.log(event)

        if (questions[index].type === "radio") {
            for (let j = 0; j < questions[index].answers.length; j++) {
                questions[index].answers[j].is_correct = false
            }
            questions[index].answers[aIndex].is_correct = event.srcElement.checked
        } else {
            questions[index].answers[aIndex].is_correct = event.srcElement.checked
        }

        localStorage.setItem("questions", JSON.stringify(questions))
    },

    ShowModal() {
        $store.data.showModal = true
    },

    CanCompleteTest() {
        const questions = JSON.parse(localStorage.getItem("questions"))

        for (let q of questions) {
            let completeQ = false

            for (let a of q.answers) {
                if (a.is_correct === true) {
                    completeQ = true
                }
            }

            if (!completeQ) {
                return false
            }
        }

        return true
    },

    async completeTest() {
        const questions = JSON.parse(localStorage.getItem("questions"))

        const timeText = document.getElementById("timer").innerText

        let timeParts = timeText.split(':');

        let hours = parseInt(timeParts[0]);
        let minutes = parseInt(timeParts[1]);
        let seconds = parseInt(timeParts[2]);

        let totalSeconds = (hours * 3600) + (minutes * 60) + seconds;

        try {
            const response = await axios.post(`/passing/${RESULT_ID}/solving/${PASS_ID}/results`, {
                questions: questions,
                time_pass: startTotalSeconds - totalSeconds
            })

            $store.data.result = response.data.result
            $store.data.isPass = false
            $store.data.showModal = false
            $store.data.isComplete = true
            clearInterval(timer)
        } catch (e) {
            $store.data.error = e.response.data
        }
    }
})

const startButton = document.getElementById("startButton")

startButton.onclick = async () => {
    try {
        const response = await axios.get(`/passing/${RESULT_ID}/solving/${PASS_ID}/questions`)
        $store.data.questions = response.data.questions
        startTimer(response.data.access.passage_time)
        localStorage.setItem("questions", JSON.stringify(response.data.questions))
        $store.data.isPass = true
        console.log(response.data)
    } catch (e) {
        alert(e.response.data)
    }
}

let timer;
let totalSeconds;
let startTotalSeconds

function startTimer(minutes) {
    totalSeconds = minutes * 60;
    startTotalSeconds = minutes * 60;
    updateTimer()
    timer = setInterval(updateTimer, 1000);
}

async function updateTimer() {
    totalSeconds--

    if (totalSeconds < 0) {
        await $store.methods.completeTest()
    }

    let hours = Math.floor(totalSeconds / 3600);
    let minutes = Math.floor((totalSeconds % 3600) / 60);
    let remainingSeconds = totalSeconds % 60;

    document.getElementById('timer').innerText = pad(hours) + ':' + pad(minutes) + ':' + pad(remainingSeconds);
}

function pad(val) {
    return val > 9 ? val : '0' + val;
}

let modal = document.getElementById("myModal");

window.onmousedown = function(event) {
    if (event.target === modal) {
        $store.data.showModal = false
        $store.data.error = ""
    }
}

window.addEventListener('beforeunload', function (e) {
    e.preventDefault();
    e.returnValue = 'Вы уверены, что хотите покинуть эту страницу? Все несохраненные изменения будут потеряны.';
});