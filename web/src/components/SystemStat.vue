<script setup>
import {ref} from "vue";
import axios from "axios";
import {useToastStack} from "../stores/toasts.js";

const totalConnections = ref(0);
const onlineUsers = ref(0);
const totalUsers = ref(0);
const totalDownloadUsage = ref("");
const totalUploadUsage = ref("");

const toasts = useToastStack();

axios.get("/api/v1/system-statuses", {withCredentials: true}).then((response) => {
  totalConnections.value = response.data.result.total_active_connections;
  onlineUsers.value = response.data.result.online_users;
  totalUsers.value = response.data.result.total_users;
  totalDownloadUsage.value = response.data.result.total_download_usage;
  totalUploadUsage.value = response.data.result.total_upload_usage;
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
</script>

<template>
  <div class="m-4 flex flex-col gap-5">
    <h2 class="font-bold text-xl uppercase">
      System Status
    </h2>
    <div class="stats shadow overflow-x-auto">
      <div class="stat">
        <div class="stat-figure text-secondary">
          <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="inline-block h-8 w-8 stroke-current">
            <path
                stroke-linecap="round"
                stroke-linejoin="round" d="M8.288 15.038a5.25 5.25 0 0 1 7.424 0M5.106 11.856c3.807-3.808 9.98-3.808 13.788 0M1.924 8.674c5.565-5.565 14.587-5.565 20.152 0M12.53 18.22l-.53.53-.53-.53a.75.75 0 0 1 1.06 0Z"
            />
          </svg>
        </div>
        <div class="stat-title">Connections</div>
        <div class="stat-value">{{totalConnections}}</div>
      </div>

      <div class="stat">
        <div class="stat-figure text-secondary">
          <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="inline-block h-8 w-8 stroke-current">
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z"
            />
          </svg>

        </div>
        <div class="stat-title">Online Users</div>
        <div class="stat-value">{{onlineUsers}}</div>
        <div class="stat-desc">of {{totalUsers}}</div>
      </div>

      <div class="stat">
        <div class="stat-figure text-secondary">
          <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="inline-block h-8 w-8 stroke-current">
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
            />
          </svg>

        </div>
        <div class="stat-title">Total Download Usage</div>
        <div class="stat-value">{{totalDownloadUsage}}</div>
      </div>

      <div class="stat">
        <div class="stat-figure text-secondary">
          <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="inline-block h-8 w-8 stroke-current">
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5"
            />
          </svg>
        </div>
        <div class="stat-title">Total Upload Usage</div>
        <div class="stat-value">{{totalUploadUsage}}</div>
      </div>
    </div>
  </div>
</template>