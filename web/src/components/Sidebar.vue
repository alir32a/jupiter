<script setup>
</script>

<template>
  <div class="drawer lg:drawer-open">
    <input id="sidebar" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content flex flex-col">
      <slot></slot>
    </div>
    <div class="drawer-side">
      <label for="sidebar" aria-label="close sidebar" class="drawer-overlay"></label>
      <ul class="menu bg-base-300 min-h-full w-64 p-4 gap-2">
          <slot name="items"></slot>
      </ul>
    </div>
  </div>
</template>