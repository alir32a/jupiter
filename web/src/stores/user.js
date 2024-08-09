import {defineStore} from "pinia";
import {ref} from "vue";

export const useUserStore = defineStore("user", () => {
    const username = ref("");
    const isLoggedIn = ref(false);

    function login(user) {
        username.value = user;
        isLoggedIn.value = true;
    }

    function logout() {
        isLoggedIn.value = false;
        username.value = "";
    }

    return {
        username,
        isLoggedIn,
        login,
        logout,
    }
},{
    persist: true,
})