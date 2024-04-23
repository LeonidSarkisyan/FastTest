import { AddTestID, ChangeIndex, GetIndex } from "./base/localstorage.js";

const TEST_ID = Number(
    window.location.pathname.split("/")[window.location.pathname.split("/").length - 1]
);

AddTestID(TEST_ID, 0);

const QUESTION_URL = window.location.pathname.replace("/p", "") + "/questions";

function QUESTION_WITH_ID_URL(questionID) {
    return QUESTION_URL + "/" + questionID + "/" + "answers";
}

let data = {
    showModal: false,
    inputElement: null,
    mainTitle: "",
    title: "",
    currentIndex: GetIndex(TEST_ID),
    titleTheme: "",
    showModalChat: false,
    countQuestion: 5,
    loading: false,
    questions: [],
    showModalDelete: false,
    render: true,
};

export const methods = {
    async CreateQuestionsFromChatGPT() {
        if (data.titleTheme === "") {
            return null;
        }

        const params = {
            title_theme: data.titleTheme,
            count_questions: Number(data.countQuestion),
        };

        console.log(params);

        let lastIndex = data.questions.length - 1;

        data.loading = true;

        try {
            const response = await axios.post(`/tests/${TEST_ID}/questions/chat-gpt`, params);
            data.questions = [...data.questions, ...response.data];
            data.showModalChat = false;

            this.chooseIndex(lastIndex + 1);
        } catch (e) {
            console.log(e);
        } finally {
            data.loading = false;
        }
    },

    showModalDelete() {
        document.getElementById("test__title__delete").value = "";
        data.showModalDelete = true;
    },

    async DeleteTest() {
        if (inputs.deleteTest.value.trim() === data.title) {
            const response = await axios.delete(`/tests/${TEST_ID}`);
            console.log(response);
            window.location.href = "/p/tests";
        }
    },

    showChatModal() {
        data.titleTheme = "";
        data.countQuestion = 5;
        data.showModalChat = true;
    },

    chooseIndex(index) {
        data.currentIndex = index;
        ChangeIndex(TEST_ID, index);

        try {
            document.getElementById("input").focus();
        } catch (e) {}
    },

    async addQuestionWithType(type) {
    },

    async addQuestion() {
    },

    async deleteQuestion(index, questionID) {
        if (index === 0 && data.questions.length === 1) {
            return;
        }

        let indexFromLocalStorage = GetIndex(TEST_ID);

        if (indexFromLocalStorage === index) {
            if (index !== 0) {
                await this.chooseIndex(index - 1);
            }
        }

        if (indexFromLocalStorage === data.questions.length - 1) {
            await this.chooseIndex(indexFromLocalStorage - 1);
        }

        try {
            const response = await axios.delete(
                QUESTION_WITH_ID_URL(questionID).replace("/answers", "")
            );

            data.questions = data.questions.filter((question) => {
                return question.id !== questionID;
            });

            console.log(response);
        } catch (e) {
            console.log(e);
        }
    },

    async updateTextQuestion() {
        let questionID = data.questions[data.currentIndex].id;
        let text = data.questions[data.currentIndex].text;

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

    async addAnswer() {
        let questionID = data.questions[data.currentIndex].id;
        const response = await axios.post(QUESTION_WITH_ID_URL(questionID));
        console.log(data.questions[data.currentIndex].answers);

        const newAnswer = {
            id: response.data.id,
            text: "",
            is_correct: false,
        }

        let q = data.questions

        q[data.currentIndex].answers = [...q[data.currentIndex].answers, newAnswer]

        data.questions = q
    },

    async updateAnswer(answerIndex) {
        let questionID = data.questions[data.currentIndex].id;
        let answer = data.questions[data.currentIndex].answers[answerIndex];

        const body = {
            text: answer.text,
            is_correct: answer.is_correct,
        };

        const response = await axios.patch(
            QUESTION_WITH_ID_URL(questionID) + "/" + answer.id,
            body
        );
        console.log(response.data);
    },

    async UpdateIsCorrectAnswer(answerIndex) {
        let questionID = data.questions[data.currentIndex].id;
        let answer = data.questions[data.currentIndex].answers[answerIndex];

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

    async focusAnswer(index, event) {
        if (event) {
            if (event.ctrlKey) {
                console.log(index);
                data.questions[data.currentIndex].answers[index].is_correct = !data.questions[
                    data.currentIndex
                    ].answers[index].is_correct;
                await this.updateAnswer(index);
                return;
            }
        }

        let inputs = document.getElementsByClassName("answer__input");

        if (index < inputs.length - 1) {
            inputs[index + 1].focus();
        } else {
            await this.addAnswer();
            let inputs = document.getElementsByClassName("answer__input");
            inputs[inputs.length - 1].focus();
        }
    },

    async deleteAnswer(id) {
        if (data.questions[data.currentIndex].answers.length < 3) {
            return;
        }

        let questionID = data.questions[data.currentIndex].id;

        try {
            const response = await axios.delete(
                QUESTION_WITH_ID_URL(questionID) + "/" + id
            );

            let q = data.questions

            q[data.currentIndex].answers = q[data.currentIndex].answers.filter((answer) => {
                return answer.id !== id;
            })

            data.questions = q
        } catch (e) {
            console.log(e);
        }
    },

    showModal() {
        data.mainTitle = data.title;
        elements.modal.style.display = "block"
    },

    async UpdateTitleTest(newTitle) {
        if (!newTitle.length) {
            return null
        }

        const body = {
            title: newTitle,
        };

        try {
            const response = await axios.patch("/tests/" + TEST_ID, body);

            elements.title.innerText = inputs.testTitle.value
            elements.modal.style.display = "none"
            data.title = inputs.testTitle.value
        } catch (e) {
            console.log(e);
        }
    },

    async SaveGroup() {
        const question = data.questions[data.currentIndex];

        const response = await axios.patch(
            `/tests/${TEST_ID}/questions/type/${question.type}/${question.id}`,
            question.data
        );
        console.log(response);
    },

    async AddAnswerToGroup(index) {
        console.log(index);

        let q = data.questions

        q[data.currentIndex].data.groups[index].answers = [
            ...data.questions[data.currentIndex].data.groups[index].answers,
            "",
        ]

        data.questions = q

        let s = window.scrollY;

        console.log(s);

        data.render = false;

        setTimeout(function () {
            data.render = true;
            window.scrollTo(0, s);
        }, 0);

        await this.SaveGroup();
    },

    async RemoveAnswerFromGroup(groupIndex, answerIndex) {
        let qq = data.questions

        console.log(groupIndex, answerIndex)
        console.log(qq[data.currentIndex].data.groups[groupIndex].answers)

        qq[data.currentIndex].data.groups[groupIndex].answers.splice(answerIndex, 1)

        data.questions = qq

        await this.SaveGroup();
    },
};

const modal = document.getElementById("myModal");
const modalChat = document.getElementById("myModalChat");
const myModalDelete = document.getElementById("myModalDelete");

window.onmousedown = function (event) {
    if (event.target === modal || event.target === modalChat || event.target === myModalDelete) {
        data.showModal = false;
        data.showModalChat = false;
        data.showModalDelete = false;
    }
};

window.onkeydown = function (event) {
    if (event.key === "ArrowRight" && event.ctrlKey) {
        let index = data.currentIndex + 1;

        if (index >= data.questions.length) {
            methods.chooseIndex(0);
        } else {
            methods.chooseIndex(index);
        }
    }

    if (event.key === "ArrowLeft" && event.ctrlKey) {
        event.preventDefault();
        let index = data.currentIndex - 1;

        if (index < 0) {
            methods.chooseIndex(data.questions.length - 1);
        } else {
            methods.chooseIndex(index);
        }
    }
};

async function init () {
    let response = await axios.get("/tests/" + TEST_ID)
    data.title = response.data.title
    data.mainTitle = response.data.title

    inputs.testTitle.value = data.mainTitle

    response = await axios.get(QUESTION_URL)
    data.questions = response.data

    elements.app.style.display = "block"
    elements.modal.style.display = "none"
    elements.deleteModal.style.display = "none"
    elements.myModalChat.style.display = "none"

    elements.title.innerText = data.title
}

function renderQuestionsList() {
    // while (lists.questions.firstChild) {
    //     lists.questions.removeChild(lists.questions.firstChild);
    // }
    //
    // for (let i = data.questions.length - 1; i >= 0; i--) {
    //     let question = data.questions[i];
    //     let index = i;
    //
    //     let div = document.createElement('div');
    //
    //     div.title = question.text
    //     div.innerText = String(index + 1)
    //     div.classList.add("question__list__point")
    //
    //     div.addEventListener("click", function () {
    //         methods.chooseIndex(index)
    //     })
    //
    //     div.addEventListener("contextmenu", async function (event) {
    //         event.preventDefault()
    //         await methods.deleteQuestion(index, question.id)
    //     })
    //
    //     if (data.currentIndex === index) {
    //         div.classList.add('question__list__point__active');
    //     }
    //
    //     lists.questions.prepend(div);
    // }
    //
    // //             <div class="question__list__point plus" id="plus">
    // //                 +
    // //             </div>
    //
    // const addQuestionButton = document.createElement("div")
    // addQuestionButton.classList.add("question__list__point", "plus")
    // addQuestionButton
    // addQuestionButton.innerText = "+"
}

function renderCurrentQuestion(index) {
    // try {
    //     document.getElementById("question").remove()
    //     document.getElementById("answer__list").remove()
    //     document.getElementById("addAnswerButton").remove()
    // } catch (e) {}
    //
    // try {
    //     document.getElementById("group__name").remove()
    // } catch (e) {}
    //
    // const question = data.questions[index]
    //
    // switch (question.type) {
    //     case "choose":
    //         renderChooseQuestions(index, question)
    //         break
    //     case "group":
    //         renderGroupQuestions(index, question)
    //         break
    // }
}

function renderAnswer(index, answer, i) {
    const answerPoint = document.createElement('div');
    answerPoint.classList.add('answer__point');
    answerPoint.id = 'answer_' + answer.id;

    const checkbox = document.createElement('input');
    checkbox.classList.add('answer__point__checkbox');
    checkbox.type = 'checkbox';
    checkbox.checked = answer.is_correct;

    checkbox.addEventListener('change', async function(event) {
        data.questions[data.currentIndex].answers[i].is_correct = checkbox.checked
        await methods.updateAnswer(i);
    });

    const input = document.createElement('input');
    input.classList.add('base__input', 'text', 'answer__input');
    input.value = answer.text;

    input.addEventListener('keydown', async function(event) {
        if (event.key === 'Enter') {
            await methods.focusAnswer(i, event);
        }
    });

    input.addEventListener('change', async function(event) {
        data.questions[data.currentIndex].answers[i].text = event.target.value
        await methods.updateAnswer(i);
    });

    const deleteButton = document.createElement('div');
    deleteButton.classList.add('answer__point__delete');
    deleteButton.innerHTML = '&#10006;';

    deleteButton.addEventListener('click', async function() {
        await methods.deleteAnswer(answer.id);
    });

    if (data.questions[index].answers.length < 3) {
        deleteButton.classList.add('answer__point__delete__disable');
    }

    answerPoint.appendChild(checkbox);
    answerPoint.appendChild(input);
    answerPoint.appendChild(deleteButton);

    return answerPoint
}

function renderChooseQuestions(index, question) {
    const questionForm = document.createElement('div')
    questionForm.classList.add('question__form')
    questionForm.id = "question"

    const questionTitle = document.createElement('div')
    questionTitle.classList.add('question__form__title')
    questionTitle.textContent = "Текст вопроса:"
    questionForm.appendChild(questionTitle)

    const questionInput = document.createElement('input');
    questionInput.id = 'input';
    questionInput.value = question.text
    questionInput.classList.add('base__input', 'text');
    questionInput.addEventListener('change', async function(event) {
        data.questions[data.currentIndex].text = event.target.value
        await methods.updateTextQuestion()
    });
    questionForm.appendChild(questionInput);

    const answerList = document.createElement('div');
    answerList.classList.add('question__form');
    answerList.id = 'answer__list';

    const answersTitle = document.createElement('div');
    answersTitle.classList.add('question__form__title');
    answersTitle.textContent = 'Варианты ответов:';
    answerList.appendChild(answersTitle);

    for (let i = 0; i < question.answers.length; i++) {
        const answer = question.answers[i]
        answerList.appendChild(renderAnswer(index, answer, i))
    }

    let addAnswerButton = document.createElement("div")
    addAnswerButton.classList.add("base__button")
    addAnswerButton.innerText = "Добавить вариант ответа"
    addAnswerButton.id = "addAnswerButton"

    addAnswerButton.addEventListener("click", async function() {
        await methods.addAnswer()
    })

    elements.app.appendChild(questionForm);
    elements.app.appendChild(answerList);
    elements.app.appendChild(addAnswerButton)
}

function renderGroupQuestions(index, question) {
    const questionForm = document.createElement('div')
    questionForm.classList.add('question__form')
    questionForm.id = "question"

    const questionTitle = document.createElement('div')
    questionTitle.classList.add('question__form__title')
    questionTitle.textContent = "Текст вопроса:"
    questionForm.appendChild(questionTitle)

    const questionInput = document.createElement('input');
    questionInput.id = 'input';
    questionInput.value = question.text
    questionInput.classList.add('base__input', 'text');
    questionInput.addEventListener('change', async function(event) {
        data.questions[data.currentIndex].text = event.target.value
        await methods.updateTextQuestion()
    });
    questionForm.appendChild(questionInput);
    elements.app.appendChild(questionForm)

    const qForm = document.createElement("div")
    qForm.classList.add("question__form")
    qForm.id = "group__name"

    const flex = document.createElement("div")
    flex.classList.add("display-flex", "justify-between")

    const variants = document.createElement("div")
    variants.id = "variants"

    const variantsTitle = document.createElement("div")
    variantsTitle.classList.add("question__form__title")
    variantsTitle.innerText = "Варианты:"

    variants.appendChild(variantsTitle)

    flex.appendChild(variants)
    qForm.appendChild(flex)
    elements.app.append(qForm)

    for (let i = 0; i < data.questions[index].data.groups.length; i++) {
        for (let j = 0; j < data.questions[index].data.groups[i].answers.length; j++) {
            const answerPoint = document.createElement("div");
            answerPoint.classList.add("answer__point", "variants", "groups__names");

            const label1 = document.createElement("label");
            const input = document.createElement("input");
            input.classList.add("base__input", "answers__inputs");

            input.onchange = async event => {
                data.questions[index].data.groups[i].answers[j] = event.target.value
                await methods.SaveGroup()
            }

            input.onkeydown = async event => {
                if (event.key === "Enter") {
                    let q = data.questions
                    q[index].data.groups[q[index].data.groups.length - 1].answers.push("")
                    data.questions = q
                    await methods.SaveGroup()

                    const inputs = document.getElementsByClassName("answers__inputs")
                    inputs[inputs.length - 1].focus()
                }
            }

            input.value = data.questions[index].data.groups[i].answers[j]

            label1.appendChild(input);
            answerPoint.appendChild(label1);

            const selectLabel = document.createElement("label");
            const select = document.createElement("select");
            select.classList.add("base__input", "choose__group");

            for (let w = 0; w < data.questions[index].data.groups.length; w++) {
                const option = document.createElement("option");
                option.innerText = String(w + 1);
                select.appendChild(option);
            }

            select.onchange = async event => {
                let q = data.questions

                let indexDelete = q[index].data.groups[i].answers.indexOf(input.value)

                q[index].data.groups[i].answers.splice(indexDelete, 1)
                q[index].data.groups[event.target.value - 1].answers.push(input.value)

                data.questions = q
                await methods.SaveGroup()
            }

            select.value = String(i + 1);
            selectLabel.appendChild(select);
            answerPoint.appendChild(selectLabel);

            let deleteAnswer = document.createElement("div");
            deleteAnswer.className = "answer__point__delete group__delete";
            deleteAnswer.innerHTML = "&#10006;";

            deleteAnswer.addEventListener("click", async function() {
                let q = data.questions
                q[index].data.groups[i].answers.splice(j, 1)
                data.questions = q
                await methods.SaveGroup()
            });

            answerPoint.appendChild(deleteAnswer)

            variants.appendChild(answerPoint)
        }
    }

    const addVariant = document.createElement("div")
    addVariant.classList.add("base__button", "edit", "margin__button16")
    addVariant.innerText = "Добавить вариант"
    addVariant.addEventListener("click", async function () {
        let q = data.questions
        q[index].data.groups[q[index].data.groups.length - 1].answers.push("")
        data.questions = q
        await methods.SaveGroup()
    })

    variants.appendChild(addVariant)


    const groupList = document.createElement("div")
    groupList.classList.add("groups__names")

    const groupsTitle = document.createElement("div");
    groupsTitle.classList.add("question__form__title", "groups");
    groupsTitle.textContent = "Группы:";
    groupList.appendChild(groupsTitle);

    data.questions[index].data.groups.forEach((group, index_) => {
        const groupDiv = document.createElement("div");
        groupDiv.classList.add("answer__point", "align-items", "justify-between");

        const numberGroupDiv = document.createElement("div");
        numberGroupDiv.classList.add("number__group");
        numberGroupDiv.textContent = index_ + 1 + ".";

        const label = document.createElement("label");
        const input = document.createElement("input");
        input.type = "text";
        input.classList.add("base__input", "group__input__");
        input.value = group.name;

        input.addEventListener("change", async function(event) {
            data.questions[index].data.groups[index_].name = event.target.value
            await methods.SaveGroup()
        })

        input.addEventListener("keydown", async function(event){
            if (event.key === "Enter") {
                let q = data.questions
                q[index].data.groups.push({
                    name: "",
                    answers: []
                })
                data.questions = q
                await methods.SaveGroup()
                const inputs = document.getElementsByClassName("group__input__")
                inputs[inputs.length - 1].focus()
            }
        })

        let deleteGroup = document.createElement("div");
        deleteGroup.className = "answer__point__delete group__delete";
        deleteGroup.innerHTML = "&#10006;";

        if (data.questions[index].data.groups.length > 1) {
            deleteGroup.addEventListener("click", async function() {
                let q = data.questions

                if (index_ > 0) {
                    q[index].data.groups[index_ - 1].answers = [
                        ...q[index].data.groups[index_ - 1].answers, ...q[index].data.groups[index_].answers,
                    ]
                } else {
                    q[index].data.groups[index_ + 1].answers = [
                        ...q[index].data.groups[index_ + 1].answers, ...q[index].data.groups[index_].answers,
                    ]
                }

                q[index].data.groups.splice(index_, 1)
                data.questions = q
                await methods.SaveGroup()
            });
        } else {
            deleteGroup.classList.add("answer__point__delete__disable")
        }

        label.appendChild(input);
        groupDiv.appendChild(numberGroupDiv);
        groupDiv.appendChild(label);
        groupDiv.appendChild(deleteGroup)
        groupList.appendChild(groupDiv);
    });

    const addButton = document.createElement("div");
    addButton.classList.add("base__button", "edit");
    addButton.textContent = "Добавить группу";

    addButton.addEventListener("click", async function(event) {
        let q = data.questions
        q[index].data.groups.push({
            name: "",
            answers: []
        })
        data.questions = q
        await methods.SaveGroup()
    })

    groupList.appendChild(addButton);

    flex.appendChild(groupList)
}

const elements = {
    app: document.getElementById("app"),
    title: document.getElementById("title"),
    modal: document.getElementById("myModal"),
    change: document.getElementById("change__title"),
    changeTitleButton: document.getElementById("change__title__button"),

    deleteModal: document.getElementById("myModalDelete"),
    deleteButton: document.getElementById("test__delete__button"),
    deleteOpen: document.getElementById("test__delete__open"),

    myModalChat: document.getElementById("myModalChat"),

    choose: document.getElementById("choose"),
    groupCreate: document.getElementById("groupCreate"),
}

const lists = {
    questions: document.getElementById("questionList")
}

const inputs = {
    testTitle: document.getElementById("test__title"),
    deleteTest: document.getElementById("test__title__delete")
}

await init()
renderQuestionsList()
renderCurrentQuestion(data.currentIndex)

// Change Title Test Modal //////////////////////////////////////////////////////////////////////////////////

elements.change.addEventListener("click", function () {
    elements.modal.style.display = "block"
})

window.onmousedown = function(event) {
    console.log("qwe")

    if (event.target === elements.modal || event.target === elements.deleteModal || event.target === modalChat) {
        elements.modal.style.display = "none"
        elements.deleteModal.style.display = "none"
        modalChat.style.display = "none"

        inputs.testTitle.value = data.title
    }
}

elements.changeTitleButton.addEventListener("click", async function () {
    await methods.UpdateTitleTest(inputs.testTitle.value)
})

inputs.testTitle.addEventListener("keydown", async function(event) {
    if (event.key === "Enter") {
        await methods.UpdateTitleTest(inputs.testTitle.value)
    }
})

// Delete Test Modal //////////////////////////////////////////////////////////////////////////////////

elements.deleteOpen.addEventListener("click", function () {
    elements.deleteModal.style.display = "block"
})

elements.deleteButton.addEventListener("click", async function() {
    await methods.DeleteTest()
})

// Choose active questions

let target = {
    name: 'Target Object',
    value: 42
};

data = new Proxy(data, {
    set(target, prop, val) {
        console.log(prop)
        if (prop === "currentIndex") {
            const q = document.getElementsByClassName("question__list__point")

            for (let i = 0; i < q.length; i++) {
                if (val === i) {
                    q[i].classList.add("question__list__point__active")
                } else {
                    q[i].classList.remove("question__list__point__active")
                }
            }

            renderCurrentQuestion(val)
            target[prop] = val
        } else if (prop === "questions") {
            target[prop] = val
            renderCurrentQuestion(data.currentIndex)
            renderQuestionsList()
        }
        return true
    }
});

// Create Question //////////////////////////////////////////////////////////////////////////////////

elements.choose.onclick = async (event) => {
    await methods.addQuestion()
}

elements.groupCreate.onclick = async (event) => {
    await methods.addQuestionWithType("group")
}