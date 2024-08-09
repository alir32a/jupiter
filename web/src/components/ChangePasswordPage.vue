<script setup>
import {ref} from "vue";
import axios from "axios";
import {useRouter} from "vue-router";
import {useToastStack} from "../stores/toasts.js";

const currentPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const error = ref(null);

const router = useRouter();

const toasts = useToastStack();

function changePassword() {
  axios.post("/api/v1/change-password", {
    current_password: currentPassword.value,
    new_password: newPassword.value,
    confirm_password: confirmPassword.value,
  }, {withCredentials: true}).then((response) => {
    if (!response.data.ok) {
      toasts.pushError(response.data.result.error);

      return;
    }

    toasts.pushSuccess("Password changed successfully");
    router.push({ name: "login" });
  }).catch((err) => {
    if (err.response) {
      if (err.response.status === 401) {
        router.push("/login");

        return;
      }

      toasts.pushError(err.response.data.result.error);
      return;
    }

    toasts.pushError(err.message);
  });
}
</script>

<template>
  <div class="flex flex-col gap-5 m-4">
    <h1 class="font-bold text-xl uppercase">
      Change Password
    </h1>
    <div class="flex flex-col gap-6 w-72">
      <label class="input input-bordered flex items-center gap-2" :class="error ? 'input-error' : ''">
        <input type="password" class="grow" placeholder="Current Password" v-model="currentPassword" />
      </label>
      <label class="input input-bordered flex items-center gap-2" :class="error ? 'input-error' : ''">
        <input type="password" class="grow" placeholder="New Password" v-model="newPassword" />
      </label>
      <label class="input input-bordered flex items-center gap-2" :class="error ? 'input-error' : ''">
        <input type="password" class="grow" placeholder="Confirm Password" v-model="confirmPassword" />
      </label>
      <button class="btn btn-primary" type="submit" @click="changePassword">Submit</button>
    </div>
  </div>
</template>

<style scoped>

</style>