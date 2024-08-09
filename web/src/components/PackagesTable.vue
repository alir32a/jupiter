<script setup>
import {ref} from "vue";
import axios from "axios";
import PackagesModal from "./PackagesModal.vue";
import {useToastStack} from "../stores/toasts.js";

const page = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const username = ref("");
let packages = ref([]);

const toasts = useToastStack();

axios.get("/api/v1/packages", {
  params: {
    page: page.value,
    page_size: pageSize.value,
    username: username.value,
  },
  withCredentials: true,
}).then((response) => {
  packages.value = response.data.result.packages;
  totalPages.value = response.data.result.total_pages;
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
})

function search() {
  axios.get("/api/v1/packages", {
    params: {
      page: page.value,
      page_size: pageSize.value,
      username: username.value,
    },
    withCredentials: true,
  }).then((response) => {
    packages.value = response.data.result.packages;
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
  })
}

function nextPage() {
  page.value++;

  search();
}

function prevPage() {
  page.value--;

  search();
}
</script>

<template>
  <div class="m-4 flex flex-col gap-5">
    <h1 class="font-bold text-xl uppercase">
      Packages
    </h1>
    <div class="flex justify-between">
      <div class="join">
        <input class="input input-bordered join-item bordered" placeholder="Username" v-model="username"/>
        <button class="btn join-item" @click="search">Search</button>
      </div>
      <div class="flex-none">
        <button class="btn btn-primary" onclick="packagesModal.showModal()">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          Add Package
        </button>
      </div>
      <PackagesModal />
    </div>
    <div class="overflow-x-auto">
      <table class="table table-zebra">
        <thead>
        <tr>
          <th></th>
          <th>Username</th>
          <th>Traffic Limit</th>
          <th>Download Usage</th>
          <th>Upload Usage</th>
          <th>Max Connections</th>
          <th>Is Trial</th>
          <th>Expiry</th>
          <th>Created At</th>
          <th>Expires At</th>
        </tr>
        </thead>
        <tbody>
          <tr v-for="(pack, i) in packages" :key="i">
            <th scope="row">{{i + 1}}</th>
            <td>{{ pack.username }}</td>
            <td>{{ pack.traffic_limit }}</td>
            <td>{{ pack.download_traffic_usage }}</td>
            <td>{{ pack.upload_traffic_usage }}</td>
            <td>{{ pack.max_connections }}</td>
            <td>{{ pack.is_trial }}</td>
            <td>{{ pack.expiration_in_days }} Days</td>
            <td>{{ pack.created_at }}</td>
            <td>{{ pack.expire_at }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="join justify-center">
      <div class="join">
        <button class="join-item btn" @click="prevPage" :class="page === 1 ? 'btn-disabled' : ''">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3">
            <path fill-rule="evenodd" d="M17 10a.75.75 0 0 1-.75.75H5.612l4.158 3.96a.75.75 0 1 1-1.04 1.08l-5.5-5.25a.75.75 0 0 1 0-1.08l5.5-5.25a.75.75 0 1 1 1.04 1.08L5.612 9.25H16.25A.75.75 0 0 1 17 10Z" clip-rule="evenodd" />
          </svg>
        </button>
        <button class="join-item btn">Page {{page}} of {{totalPages}}</button>
        <button class="join-item btn" @click="nextPage" :class="page === totalPages ? 'btn-disabled' : ''">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3">
            <path fill-rule="evenodd" d="M3 10a.75.75 0 0 1 .75-.75h10.638L10.23 5.29a.75.75 0 1 1 1.04-1.08l5.5 5.25a.75.75 0 0 1 0 1.08l-5.5 5.25a.75.75 0 1 1-1.04-1.08l4.158-3.96H3.75A.75.75 0 0 1 3 10Z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>