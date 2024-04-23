import { AddTestID, ChangeIndex, GetIndex } from "./base/localstorage.js";

const modalChat_ = document.getElementById("myModalChat")

const TEST_ID = Number(
    window.location.pathname.split("/")[window.location.pathname.split("/").length - 1]
);

AddTestID(TEST_ID, 0);

const QUESTION_URL = window.location.pathname.replace("/p", "") + "/questions";

function QUESTION_WITH_ID_URL(questionID) {
    return QUESTION_URL + "/" + questionID + "/" + "answers";
}

Spruce.store("data", {
    showModalChat: false,
    loading: false,
    imageLoaded: false,

    titleTheme: "",
    countQuestion: 5,

    currentIndex: GetIndex(TEST_ID),

    questions: [{type: "choose", answers: []}],
})

Spruce.store("methods", {
    ChooseIndex(index) {
        $store.data.currentIndex = index;
        ChangeIndex(TEST_ID, index);

        // try {
        //     document.getElementById("input").focus();
        // } catch (e) {}
    },

    showChatModal() {
        $store.data.titleTheme = ""
        $store.data.countQuestion = 5
        modalChat_.style.display = "block"
    },

    async CreateQuestionsFromChatGPT() {
        if ($store.data.titleTheme.trim().length === 0) {
            return null
        }

        $store.data.loading = true
        try {
            const response = await axios.post(`/tests/${TEST_ID}/questions/chat-gpt`, {
                title_theme: $store.data.titleTheme,
                count_questions: Number($store.data.countQuestion),
            })
            console.log(response)
        } catch (e) {
            alert(e.response.data)
        } finally {
            window.location.reload()
        }
    },

    async UpdateTextQuestion() {
        let questionID = $store.data.questions[$store.data.currentIndex].id;
        let text = $store.data.questions[$store.data.currentIndex].text;

        let body = {
            text: text,
        };

        try {
            const response = await axios.patch(
                QUESTION_WITH_ID_URL(questionID).replace("/answers", ""),
                body
            );
            console.log(response);
        } catch (e) {
            console.log(e);
        }
    },

    async UpdateAnswer(answerIndex) {
        let questionID = $store.data.questions[$store.data.currentIndex].id;
        let answer = $store.data.questions[$store.data.currentIndex].answers[answerIndex];

        const body = {
            text: answer.text,
            is_correct: answer.is_correct,
        };

        const response = await axios.patch(
            QUESTION_WITH_ID_URL(questionID) + "/" + answer.id,
            body
        );
        console.log(response);
    },

    async UpdateIsCorrectAnswer(answerIndex) {
        let questionID = $store.data.questions[$store.data.currentIndex].id;
        let answer = $store.data.questions[$store.data.currentIndex].answers[answerIndex];

        const body = {
            text: answer.text,
            is_correct: !answer.is_correct,
        };

        const response = await axios.patch(
            QUESTION_WITH_ID_URL(questionID) + "/" + answer.id,
            body
        );
        console.log(response.data);
    },

    async DeleteQuestion(index, questionID, event) {
        event.preventDefault()

        if (index === 0 && $store.data.questions.length === 1) {
            return;
        }

        let indexFromLocalStorage = GetIndex(TEST_ID);

        if (indexFromLocalStorage === index) {
            if (index !== 0) {
                await this.ChooseIndex(index - 1);
            }
        }

        if (indexFromLocalStorage === $store.data.questions.length - 1) {
            await this.ChooseIndex(indexFromLocalStorage - 1);
        }

        try {
            const response = await axios.delete(
                QUESTION_WITH_ID_URL(questionID).replace("/answers", "")
            );

            $store.data.questions = $store.data.questions.filter((question) => {
                return question.id !== questionID;
            });

            console.log(response);
        } catch (e) {
            console.log(e);
        }
    },

    async AddQuestionWithType(type) {
        const response = await axios.post(`/tests/${TEST_ID}/questions/type/${type}`);
        console.log(response);

        if (type === "group") {
            let newQuestion = {
                id: response.data.id,
                text: response.data.text,
                type: "group",
                data: response.data.data,
                answers: [],
            };

            $store.data.questions = [...$store.data.questions, newQuestion];
            this.ChooseIndex($store.data.questions.length - 1);
            const input = document.getElementById("input");
            input.focus();
        }
    },

    async AddQuestion() {
        const response = await axios.post(QUESTION_URL);

        let newQuestion = {
            id: response.data.id,
            text: "",
            type: "choose",
            answers: [],
        };

        for (let answerID of response.data.answers_ids) {
            newQuestion.answers.push({
                id: answerID,
                text: "",
                is_correct: false,
            });
        }

        $store.data.questions = [...$store.data.questions, newQuestion];
        this.ChooseIndex($store.data.questions.length - 1);
        const input = document.getElementById("input");
        input.focus();
    },

    async AddAnswer() {
        let questionID = $store.data.questions[$store.data.currentIndex].id;
        const response = await axios.post(QUESTION_WITH_ID_URL(questionID));

        const newAnswer = {
            id: response.data.id,
            text: "",
            is_correct: false,
        }

        $store.data.questions[$store.data.currentIndex].answers = [...$store.data.questions[$store.data.currentIndex].answers, newAnswer]
    },

    async DeleteAnswer(id) {
        if ($store.data.questions[$store.data.currentIndex].answers.length < 3) {
            return;
        }

        let questionID = $store.data.questions[$store.data.currentIndex].id;

        try {
            const response = await axios.delete(
                QUESTION_WITH_ID_URL(questionID) + "/" + id
            );

            $store.data.questions[$store.data.currentIndex].answers = $store.data.questions[$store.data.currentIndex].answers.filter((answer) => {
                return answer.id !== id;
            })
        } catch (e) {
            console.log(e);
        }
    },

    OpenFileWindow() {
        document.getElementById("file").click()
    },

    async ShowUploadImage() {
        const input = document.getElementById("file")
        const file = input.files[0]
        const ext = file.name.split(".")[file.name.split(".").length - 1].toLowerCase()

        if (ext !== "jpg" && ext !== "png" && ext !== "jpeg") {
            alert("некорректный формат файла, должен быть: jpg, jpeg или png")
            return null;
        }

        let formData = new FormData();
        formData.append('image', file);

        const response = await axios.post(`${QUESTION_URL}/${$store.data.questions[$store.data.currentIndex].id}/image`, formData)
        console.log(response.data.filename)

        $store.data.questions[$store.data.currentIndex].image_url = response.data.filename
    },

    async DeleteImage() {
        const response = await axios.delete(`${QUESTION_URL}/${$store.data.questions[$store.data.currentIndex].id}/image`)
        console.log(response)
        $store.data.questions[$store.data.currentIndex].image_url = ""
    },

    async SaveGroup() {
        const question = $store.data.questions[$store.data.currentIndex];

        const response = await axios.patch(
            `/tests/${TEST_ID}/questions/type/${question.type}/${question.id}`,
            question.data
        );
        console.log(response);
    },

    async AddAnswerToGroup() {
        const lastAnswer = $store.data.questions[$store.data.currentIndex].data.answers[
            $store.data.questions[$store.data.currentIndex].data.answers.length - 1
        ]

        $store.data.questions[$store.data.currentIndex].data.answers = [
            ...$store.data.questions[$store.data.currentIndex].data.answers,
            {
                text: "",
                group_index: lastAnswer ? lastAnswer.group_index : 0
            }
        ]

        await this.SaveGroup();

        const inputs = document.getElementsByClassName("answers__inputs")
        inputs[inputs.length - 1].focus()
    },

    async DeleteAnswerFromGroup(index) {
        $store.data.questions[$store.data.currentIndex].data.answers = $store.data.questions[$store.data.currentIndex].data.answers.filter((v, i) => {
            return index !== i
        })

        await this.SaveGroup();
    },

    async ChangeGroupIndex(index, event) {
        $store.data.questions[$store.data.currentIndex].data.answers[index].group_index = Number(event.target.value)

        await this.SaveGroup();
    },

    async AddGroup() {
        $store.data.questions[$store.data.currentIndex].data.groups = [
            ...$store.data.questions[$store.data.currentIndex].data.groups, {title: ""}
        ]

        await this.SaveGroup()

        const inputs = document.getElementsByClassName('group__input__')
        inputs[inputs.length - 1].focus()
    },

    async DeleteGroup(index) {
        if (index !== 0) {
            for (let i = 0; i < $store.data.questions[$store.data.currentIndex].data.answers.length; i++) {
                if ($store.data.questions[$store.data.currentIndex].data.answers[i].group_index === index) {
                    $store.data.questions[$store.data.currentIndex].data.answers[i].group_index = index - 1
                }
            }
        }

        $store.data.questions[$store.data.currentIndex].data.groups = $store.data.questions[$store.data.currentIndex].data.groups.filter((v, i) => {
            return index !== i
        })

        await this.SaveGroup();
    }
})

let response = await axios.get(QUESTION_URL)
console.log(response.data)
$store.data.questions = response.data