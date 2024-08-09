<script setup>
import {useRouter} from "vue-router";
import {useUserStore} from "./stores/user.js";

const router = useRouter();

router.beforeEach(async (to, from, next) => {
  const user = useUserStore();

  if (!user.isLoggedIn && !to.meta.noAuth) {
    next({ path: '/login' });
  } else {
    next();
  }
});
</script>

<template>
  <RouterView />
</template>

