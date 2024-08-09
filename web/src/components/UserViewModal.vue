<script setup>
import TextBox from "./TextBox.vue";
import axios from "axios";
import {useToastStack} from "../stores/toasts.js";
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";

const router = useRouter();

const user = defineModel();
const emits = defineEmits(["closeModal"]);

const toasts = useToastStack();

const packages = ref(null);
const connections = ref(null);

onMounted(() => {
  document.getElementById("userViewModal").showModal();

  axios.get("/api/v1/active-packages", {
    params: {
      user_id: user.value.id,
    },
    withCredentials: true,
  }).then((response) => {
    if (!response.data.ok) {
      toasts.pushError(response.data.result.error);

      return;
    }

    packages.value = response.data.result;
  }).catch((err) => {
    if (err.response) {
      if (err.response.status === 401) {
        router.push("/login");

        return;
      }

      toasts.pushError(err.response.result.error);

      return;
    }

    toasts.pushError(err.message);
  });

  axios.get("/api/v1/user-connections", {
    params: {
      username: user.value.username,
    },
    withCredentials: true,
  }).then((response) => {
    if (!response.data.ok) {
      toasts.pushError(response.data.result.error);

      return;
    }

    connections.value = response.data.result;
  }).catch((err) => {
    if (err.response) {
      if (err.response.status === 401) {
        router.push("/login");

        return;
      }

      toasts.pushError(err.response.result.error);

      return;
    }

    toasts.pushError(err.message);
  });
});

function closeModal() {
  emits("closeModal");
}
</script>

<template>
  <dialog id="userViewModal" class="modal w-screen sm:modal-middle" @close="closeModal">
    <div class="modal-box flex flex-col gap-6">
      <form method="dialog" class="mb-4">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
          </svg>
        </button>
      </form>
      <div class="flex flex-col">
        <TextBox title="ID" :content="user.id"/>
        <TextBox title="Username" :content="user.username"/>
        <TextBox title="User Type" :content="user.user_type"/>
        <TextBox title="Referral Code" :content="user.referral_code"/>
        <TextBox title="Referral" :content="user.referral"/>
        <TextBox title="Created At" :content="user.created_at"/>
      </div>
      <div class="flex flex-col gap-4">
        <h1 class="uppercase font-bold my-4">Packages</h1>
        <h3 v-if="packages.length === 0" class="text-center">{{ user.username }} doesn't have an active package</h3>
        <div v-if="packages.length > 0" class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
            <tr>
              <th>#</th>
              <th>Status</th>
              <th>Traffic Limit</th>
              <th>Used Traffic</th>
              <th>Expires At</th>
            </tr>
            </thead>
            <tbody>
              <tr v-for="(pack,i) in packages">
                <td>{{ i + 1 }}</td>
                <td>
                  <div class="badge gap-2" :class="pack.status === 'active' ? 'badge-primary' : 'badge-secondary'">
                    {{pack.status}}
                  </div>
                </td>
                <td>{{ pack.traffic_limit }}</td>
                <td>{{ pack.used_traffic }}</td>
                <td>{{ pack.expires_at }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <h1 class="uppercase font-bold my-4">Active Connections</h1>
        <h3 v-if="packages.length === 0" class="text-center">{{ user.username }} doesn't have an active connection</h3>
        <div v-if="packages.length > 0" class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
            <tr>
              <th>#</th>
              <th>Used Traffic</th>
              <th>IP</th>
              <th>User Agent</th>
              <th>Connected At</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="(conn,i) in connections">
              <td>{{ i + 1 }}</td>
              <td>{{ conn.used_traffic }}</td>
              <td>{{ conn.ip }}</td>
              <td>{{ conn.user_agent }}</td>
              <td>{{ conn.connected_at }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="modal-action">
        <form method="dialog">
          <button class="btn">Close</button>
        </form>
      </div>
    </div>
  </dialog>
</template>

<style scoped>

</style>