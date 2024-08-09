<script setup>
import {ref, watch} from "vue";
import axios from "axios";
import {useToastStack} from "../stores/toasts.js";
import BanUserModal from "./BanUserModal.vue";
import UnbanUserModal from "./UnbanUserModal.vue";
import {useRouter} from "vue-router";
import UserViewModal from "./UserViewModal.vue";

const router = useRouter();
const toasts = useToastStack();

const page = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const username = ref("");
let users = ref([]);

const selectedUser = ref(null);
const userModalIsOpen = ref(false);

axios.get("/api/v1/users", {
  params: {
    page: page.value,
    page_size: pageSize.value,
    username: username.value,
  },
  withCredentials: true,
}).then((response) => {
  users.value = response.data.result.users;
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
  axios.get("/api/v1/users", {
    params: {
      page: page.value,
      page_size: pageSize.value,
      username: username.value,
    },
    withCredentials: true,
  }).then((response) => {
    users.value = response.data.result.users;
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

function showBanModal(banned) {
  if (banned) {
    document.getElementById("unbanUserModal").showModal();

    return;
  }

  document.getElementById("banUserModal").showModal();
}

function showUserModal(user) {
  selectedUser.value = user;
  userModalIsOpen.value = true;
}

function closeUserModal() {
  userModalIsOpen.value = false;
}
</script>

<template>
  <div class="m-4 flex flex-col gap-5">
    <h2 class="font-bold text-xl uppercase">
      Users
    </h2>
    <div class="join">
      <input class="input input-bordered join-item bordered" placeholder="Username" v-model="username" />
      <button class="btn join-item" @click="search">Search</button>
    </div>
    <div class="overflow-x-auto">
      <table class="table table-zebra">
        <thead>
        <tr>
          <th></th>
          <th>Username</th>
          <th>User Type</th>
          <th>Banned At</th>
          <th>Referral Code</th>
          <th>Referral</th>
          <th>Created At</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(user, i) in users" :key="i">
          <th scope="row">{{i + 1}}</th>
          <td>{{user.username}}</td>
          <td>{{user.user_type}}</td>
          <td>
            <div class="badge gap-2" :class="user.banned_at ? 'badge-error' : 'badge-info'">
              {{user.banned_at ? 'banned' : 'active'}}
            </div>
          </td>
          <td>{{user.referral_code}}</td>
          <td>{{user.referral}}</td>
          <td>{{user.created_at}}</td>
          <td>
            <div class="flex gap-2">
              <div class="md:tooltip" data-tip="view user">
                <button @click="() => showUserModal(user)">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-sky-600">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                  </svg>
                </button>
              </div>
              <div class="md:tooltip" data-tip="ban/unban user">
                <button @click="() => showBanModal(user.banned_at)">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-red-600">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 0 0 5.636 5.636m12.728 12.728A9 9 0 0 1 5.636 5.636m12.728 12.728L5.636 5.636" />
                  </svg>
                </button>
              </div>
            </div>
          </td>
          <BanUserModal :id="user.id" :username="user.username" />
          <UnbanUserModal :id="user.id" :username="user.username" />
        </tr>
        </tbody>
      </table>
      <UserViewModal v-if="userModalIsOpen" v-model="selectedUser" @close-modal="closeUserModal" />
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