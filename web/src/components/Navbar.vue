<script setup>

import ThemeController from "./ThemeController.vue";
import {useUserStore} from "../stores/user.js";

const user = useUserStore();
</script>

<template>
  <div class="navbar bg-base-200">
    <div class="flex-none">
      <button class="btn btn-square btn-ghost lg:hidden show">
        <label for="sidebar" aria-label="open sidebar" class="btn btn-square btn-ghost">
        <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            class="inline-block h-5 w-5 stroke-current">
          <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"></path>
        </svg>
        </label>
      </button>
    </div>
    <div class="flex-1">
      <a class="btn btn-ghost text-xl">Jupiter</a>
    </div>
    <div class="flex-none p-4">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-6 h-6">
        <path d="M10 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM3.465 14.493a1.23 1.23 0 0 0 .41 1.412A9.957 9.957 0 0 0 10 18c2.31 0 4.438-.784 6.131-2.1.43-.333.604-.903.408-1.41a7.002 7.002 0 0 0-13.074.003Z" />
      </svg>
      <h6 class="p-1">{{ user.username }}</h6>
    </div>
    <div class="flex-none">
      <ThemeController />
    </div>
  </div>
</template>
