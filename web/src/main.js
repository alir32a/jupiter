import { createApp } from 'vue';
import '../dist/output.css';
import App from './App.vue';
import HomePage from "./components/HomePage.vue";
import PackagesPage from "./components/PackagesPage.vue";
import {createRouter, createWebHistory} from "vue-router";
import UsersPage from "./components/UsersPage.vue";
import OcservPage from "./components/OcservPage.vue";
import ChangePasswordPage from "./components/ChangePasswordPage.vue";
import LogOutPage from "./components/LogOutPage.vue";
import LoginPage from "./components/LoginPage.vue";
import MainPage from "./components/MainPage.vue";
import {createPinia} from "pinia";
import persist from "./plugins/piniaPersist.js";
import {themeChange} from "theme-change";
import axios from "axios";

axios.defaults.baseURL = import.meta.env.VITE_BACKEND_BASE_URL;

const routes = [
    {
        path: "/",
        component: MainPage,
        children: [
            { path: "", component: HomePage },
            { path: "packages", component: PackagesPage },
            { path: "users", component: UsersPage },
            { path: "ocserv", component: OcservPage },
            { path: "change-password", component: ChangePasswordPage },
            { path: "logout", component: LogOutPage },
        ]
    },
    {
        path: "/login",
        name: "login",
        component: LoginPage,
        meta: {
            noAuth: true,
        },
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

const pinia = createPinia();
pinia.use(persist);

themeChange(true);

createApp(App).
    use(router).
    use(pinia).
    mount('#app');
