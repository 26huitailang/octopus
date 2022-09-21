<script setup>
import { computed, onMounted, onUpdated, ref } from 'vue'

const props = defineProps(['command'])
const command = ref(props.command)
const out = ref('')
const rows = ref(10)
const cols = ref(80)
const myinput = ref(null)

const outCols = computed(() => {
  return window.innerWidth / 20
})
const outRows = computed(() => {
  return window.innerHeight / 35
})
const calCommand = computed(() => {
  command.value = props.command
  return props.command
})

function handleClick() {
  const data = {"script": command.value}
  fetch("/agent/script", {method: "POST", body: JSON.stringify(data)})
  .then(response => response.json()).then(data => {
    console.log("success", data)
    out.value = data.out
  })
  .catch((error) => {
    console.error("error", error)
  })
}
onUpdated(() => {
  myinput.value.focus()
})
onMounted(() => {
  console.log("hello")
})
</script>

<template>
  <div class="card">
    <div class="cmdline">
      <textarea ref="myinput" v-model="command" @keyup.ctrl.enter="handleClick" :rows="rows" :cols="cols"></textarea>
      <button class="button" type="button" @click="handleClick">Run</button>
    </div>
    <div>
    <textarea v-model="out" disabled="false" :rows="outRows" :cols="outCols"></textarea>
    </div>
  </div>
</template>

<style scoped>
.card {
  height: 500px;
  padding: 2em;
}
.cmdline {
  display: flex;
  place-items: center;
  margin-bottom: 1.5em;
}
</style>
