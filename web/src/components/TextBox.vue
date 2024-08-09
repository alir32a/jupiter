<script setup>
defineProps({
  title: String | Number,
  content: String | Number,
});
</script>

<template>
  <div class="join border-b-2 border-base-300 rounded-none first:border-t-2">
    <div class="bg-base-200 p-3 join-item w-3/12">
      {{ title }}
    </div>
    <div class="bg-base-100 p-2 flex items-center w-full join-item">
      {{ content }}
    </div>
  </div>
</template>

<style scoped>

</style>