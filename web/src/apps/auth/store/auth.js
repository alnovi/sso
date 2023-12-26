import { defineStore } from "pinia";
import { ref } from "vue";
import { useToastStore } from "./toast.js";
import api from "../../../service/api.js";

export const useAuthStore = defineStore('auth', () => {
    const toastStore = useToastStore()
    const credential = ref({login: "admin", password: "admin", remember: false})
    const errors = ref({})

    function signIn() {
        errors.value = {}
        return api.post("/api/oauth/signin", credential.value)
            .then(response => {
                window.location.replace(response.data.location)
            })
            .catch(error => {
                toastStore.error(error.response.data.message)
                if (error.response.status === 422) {
                    errors.value = error.response.data.errors || {}
                }
            })
    }

    function forgotPassword() {
        errors.value = {}

        let postData = {
            login: credential.value.login,
            client_id: (new URL(location)).searchParams.get("client_id")
        }

        return api.post("/forgot-password", postData)
            .then(response => {
                toastStore.info(response.data.message)
            })
            .catch(error => {
                toastStore.error(error.response.data.message)
                if (error.response.status === 422) {
                    errors.value = error.response.data.errors || {}
                }
            })
    }

    return { credential, errors, signIn, forgotPassword }
})
