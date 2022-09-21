<script setup>
import { computed, ref, defineEmits } from 'vue'

const emit = defineEmits(['cmd-event'])
</script>

<template>
  <h1>Scirpt</h1>
  <button @click="emit('cmd-event', 'whoami')">whoami</button>
  <button @click="emit('cmd-event', 'ls -lh')">open-22</button>
  <button>open-20022</button>
</template>

<style scoped>
</style>
