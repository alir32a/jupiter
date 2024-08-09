<script setup>
import {reactive, ref} from "vue";
import axios from "axios";
import {useToastStack} from "../stores/toasts.js";
import {useRouter} from "vue-router";

const page = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const username = ref("");
let connections = ref([]);

const toasts = useToastStack();

const router = useRouter();

axios.get("/api/v1/connections", {
  params: {
    page: page.value,
    page_size: pageSize.value,
    username: username.value,
  },
  withCredentials: true,
}).then((response) => {
  connections.value = response.data.result.connections;
  totalPages.value = response.data.result.total_pages <= 0 ? 1 : response.data.result.total_pages;
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
  axios.get("/api/v1/connections", {
    params: {
      page: page.value,
      page_size: pageSize.value,
      username: username.value,
    },
    withCredentials: true,
  }).then((response) => {
    connections.value = response.data.result.connections;
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

function disconnect(id) {
  axios.post(`/api/v1/connections/${id}/disconnect`, {}, {
    withCredentials: true,
  }).then((response) => {
    toasts.pushSuccess("disconnected successfully", router.go);
  }).catch((err) => {
    if (err.response) {
      toasts.pushError(err.response.data.result.error);

      return;
    }

    toasts.pushError(err.message);
  })
}
</script>

<template>
  <div class="m-4 flex flex-col gap-5">
    <h2 class="font-bold text-xl uppercase">
      Active Connections
    </h2>
    <div class="join">
      <input class="input input-bordered join-item bordered" placeholder="Username" v-model="username" />
      <button class="btn join-item" @click="search">Search</button>
    </div>
    <div class="overflow-x-auto">
      <table class="table table-zebra">
        <thead>
        <tr>
          <th>#</th>
          <th>Username</th>
          <th>Download Traffic Usage</th>
          <th>Upload Traffic Usage</th>
          <th>Connected At</th>
          <th>Remote IP</th>
          <th>Location</th>
          <th>User Agent</th>
          <th>Device</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(conn,i) in connections" :key="i">
          <th scope="row">{{i + 1}}</th>
          <td>{{conn.username}}</td>
          <td>{{conn.download_traffic_usage}}</td>
          <td>{{conn.upload_traffic_usage}}</td>
          <td>{{conn.connected_at}}</td>
          <td>{{conn.remote_ip}}</td>
          <td>{{conn.location}}</td>
          <td>{{conn.user_agent}}</td>
          <td>{{conn.hostname}}</td>
          <td>
            <div class="md:tooltip" data-tip="disconnect">
              <button @click="() => disconnect(conn.id)">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-red-400">
                  <path stroke-linecap="round" stroke-linejoin="round" d="m3 3 8.735 8.735m0 0a.374.374 0 1 1 .53.53m-.53-.53.53.53m0 0L21 21M14.652 9.348a3.75 3.75 0 0 1 0 5.304m2.121-7.425a6.75 6.75 0 0 1 0 9.546m2.121-11.667c3.808 3.807 3.808 9.98 0 13.788m-9.546-4.242a3.733 3.733 0 0 1-1.06-2.122m-1.061 4.243a6.75 6.75 0 0 1-1.625-6.929m-.496 9.05c-3.068-3.067-3.664-7.67-1.79-11.334M12 12h.008v.008H12V12Z" />
                </svg>
              </button>
            </div>
          </td>
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