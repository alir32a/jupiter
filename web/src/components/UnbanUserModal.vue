<script setup>
import {useRouter} from "vue-router";
import axios from "axios";
import {useToastStack} from "../stores/toasts.js";

const props = defineProps({
  id: Number,
  username: String
});

const toasts = useToastStack();

const router = useRouter();

function banUser() {
  axios.post(`/api/v1/users/${props.id}/unban`, {}, {
    withCredentials: true,
  }).then(() => {
    router.go();
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
</script>

<template>
  <dialog id="unbanUserModal" class="modal modal-bottom sm:modal-middle">
    <div class="modal-box flex flex-col gap-6">
      <div class="flex flex-col gap-6">
        <h3 class="text-lg font-bold">Unban user</h3>
        <p>
          Are you sure you want to unban {{ username }}?
        </p>
      </div>
      <div class="modal-action">
        <button class="btn btn-primary" @click="banUser">Unban</button>
        <form method="dialog">
          <button class="btn">Close</button>
        </form>
      </div>
    </div>
  </dialog>
</template>

<style scoped>

</style>