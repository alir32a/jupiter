<script setup>
import {ref} from "vue";
import axios from "axios";
import {useUserStore} from "../stores/user.js";
import {useRouter} from "vue-router";

const username = ref("");
const password = ref("");
const error = ref("");

const router = useRouter();

function login() {
  const store = useUserStore();

  axios.post("/api/v1/login", {
    username: username.value,
    password: password.value,
  }, {withCredentials: true}).then((response) => {
    if (!response.data.ok) {
      error.value = response.data.result.error;

      return;
    }

    store.login(username.value);
    router.push({path: "/"});
  }).catch((err) => {
    if (err.response) {
      error.value = err.response.data.result.error;
    }

    error.value = err.message;
  });
}
</script>

<template>
  <div class="container mx-auto bg-base-200 flex flex-col gap-8 shadow-md my-36 w-5/12 items-center rounded-md">
    <h1 class="font-bold text-3xl text-center uppercase mt-16">
      Jupiter
    </h1>
    <div class="flex flex-col w-5/12 gap-6">
      <label class="input input-bordered flex items-center gap-2">
        <input v-model="username" type="text" class="grow" placeholder="Username" />
      </label>
      <label class="input input-bordered flex items-center gap-2">
        <input v-model="password" type="password" class="grow" placeholder="Password" />
      </label>
      <p class="text-red-400 text-sm text-center">{{error}}</p>
      <button type="submit" class="btn btn-primary justify-center gap-4 mb-32 w-6/12 items-center" @click="login">
        Login
      </button>
    </div>
  </div>
</template>

<style scoped>

</style>