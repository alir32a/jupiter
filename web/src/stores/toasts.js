import {defineStore} from "pinia";
import {ref} from "vue";

export const useToastStack = defineStore("toasts", () => {
    const stack = ref([]);

    function pushError(err) {
        stack.value.push({
            message: err,
            type: "error"
        });

        setTimeout(() => {
            stack.value.pop();
        }, 5000);
    }

    function pushSuccess(msg, afterFn = null) {
        stack.value.push({
            message: msg,
            type: "success"
        });

        setTimeout(() => {
            stack.value.pop();

            if (afterFn) {
                afterFn();
            }
        }, 5000);
    }

    return {
        stack,
        pushError,
        pushSuccess,
    }
});