import {IsNotEmpty} from "./base/utils.js"

const GROUP_ID = Number(
    window.location.pathname.split("/")[window.location.pathname.split("/").length - 1]
)

Spruce.store("data", {
    showExcelModal: false,
    showModal: false,
    showCreateModal: false,

    inputElement: null,

    mainTitle: "",
    title: "",

    search: "",

    name: "",
    surname: "",
    patronymic: "",

    students: []
})

let response = await axios.get("/groups/" + GROUP_ID)
$store.data.title = response.data.name
$store.data.mainTitle = response.data.name

response = await axios.get(`/groups/${GROUP_ID}/students`)
console.log(response.data)
$store.data.students = response.data

Spruce.store("methods", {

    async AddStudent() {
        let nameText = name.value
        let surnameText = surname.value
        let patronymicText = patronymic.value

        if (IsNotEmpty(nameText, surnameText)) {
            try {
                const response = await axios.post(`/groups/${GROUP_ID}/students`, {
                    name: nameText,
                    surname: surnameText,
                    patronymic: patronymicText
                })

                const newStudent = {
                    id: response.data.id,
                    name: nameText,
                    surname: surnameText,
                    patronymic: patronymicText
                }

                newStudent.id = response.data.id

                $store.data.showCreateModal = false
                $store.data.students = [newStudent, ...$store.data.students]
            } catch (e) {
                console.log(e)
            }
        }
    },

    async DeleteStudent(studentID) {
        try {
            const response = await axios.delete(`/groups/${GROUP_ID}/students/${studentID}`)

            $store.data.students = $store.data.students.filter(student => {
                return student.id !== studentID
            })
        } catch (e) {

        }
    },

    ShowCreateModal() {
        name.value = ""
        surname.value = ""
        patronymic.value = ""

        $store.data.showCreateModal = true
    },

    ShowExcelModal() {
        $store.data.showExcelModal = true
        excel.value = ""
    },

    showModal() {
        $store.data.mainTitle = $store.data.title
        $store.data.showModal = true
    },

    async UpdateTitleGroup() {
        let body = {
            title: $store.data.mainTitle
        }

        try {
            const response = await axios.patch("/groups/" + GROUP_ID, body)
            console.log(response)

            $store.data.title = $store.data.mainTitle
            $store.data.showModal = false
        } catch (e) {
            console.log(e)
        }
    },

    async UploadExcelFile() {
        let formData = new FormData();
        formData.append('file', excel.files[0]);

        try {
            const response = await axios.post(`/groups/${GROUP_ID}/students/excel`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })

            $store.data.students = [...$store.data.students, ...response.data]
            $store.data.showExcelModal = false
        } catch (e) {
            alert(e.response.data)
        }
    }
})

let name = document.getElementById("name")
let surname = document.getElementById("surname")
let patronymic = document.getElementById("patronymic")

let modal = document.getElementById("myModal");
let excelModal = document.getElementById("myModalExcel")
let createModal = document.getElementById("myModalCreate")

window.onmousedown = function(event) {
    if (event.target === modal || event.target === excelModal || event.target === createModal) {
        $store.data.showModal = false
        $store.data.showExcelModal = false
        $store.data.showCreateModal = false
    }
}

const excel = document.getElementById("excel")

excel.onchange = async function(event) {
    const ext = event.target.files[0].name.split(".")[event.target.files[0].name.split(".").length - 1]

    if (ext !== "xlsx") {
        event.stopPropagation()
        this.value = '';
        alert("Формат файла должен быть xlsx")
    }
}