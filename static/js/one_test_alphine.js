const modalChat_ = document.getElementById("myModalChat")

const TEST_ID = Number(
    window.location.pathname.split("/")[window.location.pathname.split("/").length - 1]
);

Spruce.store("data", {
    showModalChat: false,
    loading: false,

    titleTheme: "",
    countQuestion: 5,
})

Spruce.store("methods", {
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
    }
})