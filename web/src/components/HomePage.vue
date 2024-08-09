<script setup>
import ConnectionsTable from "./ConnectionsTable.vue";
import SystemStat from "./SystemStat.vue";

</script>

<template>
  <div class="flex flex-col gap-8 p-6">
    <SystemStat />
    <ConnectionsTable />
  </div>
</template>