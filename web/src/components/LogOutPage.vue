<script setup>
  import {useUserStore} from "../stores/user.js";
  import {useRouter} from "vue-router";
  import axios from "axios";
  import {ref} from "vue";
  import {useToastStack} from "../stores/toasts.js";

  const router = useRouter();

  const toasts = useToastStack();

  axios.post("/api/v1/logout", {}, {withCredentials: true}).then((response) => {
    if (!response.data.ok) {
      error.value = response.data.result.error;

      return;
    }

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

  const user = useUserStore();
  user.logout();
  router.push({ name: "login" });
</script>

<template>

</template>