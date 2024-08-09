<script setup>
import {ref} from "vue";
import axios from "axios";
import {useRouter} from "vue-router";
import {useToastStack} from "../stores/toasts.js";

const username = ref("");
const trafficLimit = ref(0);
const maxConnections = ref(0);
const expiry = ref(0);

const router = useRouter();

const toasts = useToastStack();

function addPackage() {
  axios.post("/api/v1/packages", {
    username: username.value,
    traffic_limit: trafficLimit.value,
    max_connections: maxConnections.value,
    expiry: expiry.value,
  }, {withCredentials: true}).then((response) => {
    if (!response.data.ok) {
      toasts.pushError(response.data.result.error);

      return;
    }

    router.go();
    toasts.pushSuccess("Package added successfully");
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
  <dialog id="packagesModal" class="modal modal-bottom sm:modal-middle">
    <div class="modal-box flex flex-col gap-6">
      <div class="flex flex-col gap-6">
        <h3 class="text-lg font-bold">Add Package</h3>
        <label class="form-control">
          <div class="label">
            <span class="label-text">Username</span>
          </div>
          <label class="input input-bordered flex items-center gap-2" :class="error ? 'input-error' : ''">
            <input type="text" class="w-full" v-model="username" @change="error = null" />
          </label>
          <div class="label">
            <span class="label-text-alt text-red-400">{{error}}</span>
          </div>
        </label>
        <label class="form-control">
          <div class="label">
            <span class="label-text">Traffic Limit</span>
          </div>
          <label class="input input-bordered flex items-center gap-2" :class="error ? 'input-error' : ''">
            <input type="number" class="w-full" v-model="trafficLimit" />
            <span class="badge badge-secondary justify-self-end">GB</span>
          </label>
        </label>
        <label class="form-control">
          <div class="label">
            <span class="label-text">Max Connections</span>
          </div>
          <label class="input input-bordered flex items-center gap-2" :class="error ? 'input-error' : ''">
            <input type="number" class="w-full" v-model="maxConnections" />
          </label>
        </label>
        <label class="form-control">
          <div class="label">
            <span class="label-text">Expiry</span>
          </div>
          <label class="input input-bordered flex items-center gap-2" :class="error ? 'input-error' : ''">
            <input type="number" class="w-full" v-model="expiry" />
            <span class="badge badge-secondary justify-self-end">Days</span>
          </label>
        </label>
      </div>
      <div class="modal-action">
          <button class="btn btn-primary" @click="addPackage">Add</button>
        <form method="dialog">
          <button class="btn">Close</button>
        </form>
      </div>
    </div>
  </dialog>
</template>

<style scoped>

</style>