import { defineStore } from "pinia";
import { ref } from "vue";

export const useToastStore = defineStore('toast', () => {
    const messages = ref([])

    function success(message) {
        if (message) {
            messages.value.push({
                severity: 'success', detail: message, life: 5000
            })
        }
    }

    function info(message) {
        if (message) {
            messages.value.push({
                severity: 'info', detail: message, life: 5000
            })
        }
    }

    function warn(message) {
        if (message) {
            messages.value.push({
                severity: 'warn', detail: message, life: 5000
            })
        }
    }

    function error(message) {
        if (message) {
            messages.value.push({
                severity: 'error', detail: message, life: 5000
            })
        }
    }

    function getMessage() {
        return messages.value.shift()
    }

    return { messages, success, info, warn, error, getMessage }
})
