const TEST_ID = Number(
    window.location.pathname.split("/")[window.location.pathname.split("/").length - 2]
)

let today = new Date();

let formattedDate = today.getFullYear() + '-' + (
    today.getMonth() + 1).toString().padStart(2, '0') + '-' + today.getDate().toString().padStart(2, '0'
);

document.getElementById('start').value = formattedDate;
document.getElementById('end').value = formattedDate;

const input5 = document.getElementById("5")
const input4 = document.getElementById("4")
const input3 = document.getElementById("3")

input5.addEventListener('input', function() {
    let count = Number(document.getElementById("count").innerText)

    let numberValue = Number(this.value)

    if (numberValue <= 0 || numberValue > count) {
        this.value = '';
    }
});

input4.addEventListener("input", function () {
    let count = Number(input5.value)

    console.log(count)

    if (this.value <= 0 || this.value >= count) {
        this.value = '';
    }
})

input3.addEventListener("input", function () {
    let count = Number(input4.value)

    let numberValue = Number(this.value)

    if (numberValue >= count || numberValue <= 0) {
        this.value = ""
    }
})

Spruce.store("data", {
    showModal: false,
    search: "",

    minutes: "",
    dateStart: formattedDate,
    dateEnd: formattedDate,
    shuffle: true,

    groups: [],
    selectedGroup: {
        id: 0,
        name: "Группа не выбрана"
    }
})

const response = await axios.get("/groups")
$store.data.groups = response.data

Spruce.store("methods", {
    ShowModal() {
        const div = document.getElementById("groupList")
        div.scrollTop = 0;

        $store.data.search = ""
        $store.data.showModal = true
    },

    ChooseGroup(index) {
        $store.data.selectedGroup = $store.data.groups[index]
        $store.data.showModal = false
    },

    async CreateAccess() {
        const shuffle = $store.data.shuffle
        const groupID = $store.data.selectedGroup.id
        const minutes = Number($store.data.minutes)
        const start = $store.data.dateStart
        const end = $store.data.dateEnd
        const five = Number(input5.value)
        const four = Number(input4.value)
        const three = Number(input3.value)

        const today = new Date(formattedDate)
        const dateStart = new Date(start)
        const dateEnd = new Date(end)

        if (groupID === 0) {
            alert("выберите группу, которую будете тестировать")
            return
        }

        if (minutes === 0) {
            alert("укажите время прохождения теста")
            return
        }

        if (dateStart > dateEnd) {
            alert("дата конца не может быть раньше даты начала")
            return
        }

        if (today > dateStart || today > dateEnd) {
            alert(`нельзя установить дату, раньше сегодняшней: ${formattedDate.replaceAll("-", ".")}`)
            return
        }

        if (five === 0 || four === 0 || three === 0) {
            alert("введите все критерии для выставления оценки")
            return
        }

        if (five <= four || five <= three || four <= three) {
            alert("некорректно выставлены значения в критерии")
            return
        }

        const access = {
            shuffle: shuffle,
            passage_time: minutes,
            date_start: start,
            date_end: end,
            criteria: {
                five,
                four,
                three
            }
        }

        try {
            const response = await axios.post(`/tests/${TEST_ID}/access/${groupID}`, access)
            console.log(response)
        } catch (e) {
            alert(e.response.data)   
        }
    }
})

let modal = document.getElementById("modal")

window.onmousedown = function(event) {
    if (event.target === modal) {
        $store.data.showModal = false
    }
}