<script setup>
defineProps({
  stack: Array,
})
</script>

<template>
  <div class="toast toast-top toast-center z-inf">
    <div v-for="toast in stack" class="alert alert-info w-80" :class="toast.type === 'error' ? 'alert-error' : 'alert-success'">
      <svg v-if="toast.type === 'success'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
        <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
      </svg>

      <svg v-if="toast.type === 'error'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
      </svg>

      <span>{{ toast.message }}</span>
    </div>
  </div>
</template>

<style scoped>
.z-inf {
  z-index: 999999;
}
</style>