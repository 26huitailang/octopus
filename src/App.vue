<script setup>
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
import {ref} from 'vue'
import Runner from './components/Runner.vue'
import Cmdline from './components/Cmdline.vue'
const command = ref('')

function handleCmdEvent (cmd) {
  command.value = cmd
}
</script>

<template>
  <Cmdline @cmd-event="handleCmdEvent" />
  <Runner :command="command"/>
</template>

<style scoped>
</style>
