import {IsNotEmpty} from "./base/utils.js"

const BASE_AUTH_URL = "/auth"
const REGISTER_URL = BASE_AUTH_URL + "/register"
const LOGIN_URL = BASE_AUTH_URL + "/login"

Spruce.store("data", {
    isLogin: true,

    email: "",
    password: "",
    passwordRetry: ""
})

Spruce.store("methods", {
    toggle() {
        $store.data.email = ""
        $store.data.password = ""
        $store.data.passwordRetry = ""
        $store.data.isLogin = !$store.data.isLogin
    },

    async register() {
        if ($store.data.password === $store.data.passwordRetry) {
            if (IsNotEmpty($store.data.email, $store.data.password, $store.data.passwordRetry)) {

                const user = {
                    email: $store.data.email,
                    password: $store.data.password
                }

                console.log(user)

                try {
                    const response = await axios.post(REGISTER_URL, user)
                    console.log(response)
                } catch (e) {
                    alert(e.response.data)
                    return null
                }
            }
        }
    },

    async login() {
        if (IsNotEmpty($store.data.email, $store.data.password)) {

            const user = {
                email: $store.data.email,
                password: $store.data.password
            }

            try {
                const response = await axios.post(LOGIN_URL, user)
            } catch (e) {
                alert(e.response.data)
            }

            let timestamp = new Date().getTime();

            window.location.href = "/p/tests" // + "?timestamp=" + timestamp;
        }
    }
})