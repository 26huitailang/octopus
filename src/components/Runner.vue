<script setup>
import { computed, onMounted, onUpdated, ref, toRefs, watch } from "vue";

const scripts = ref([]);
const command = ref("");
const out = ref("");
const rows = ref(10);
const cols = ref(80);
const myinput = ref(null);

const outCols = computed(() => {
  return window.innerWidth / 20;
});
const outRows = computed(() => {
  return window.innerHeight / 35;
});

function fetchScipts() {
  fetch("/agent/scripts", { method: "GET" })
    .then((response) => response.json())
    .then((data) => {
      console.log("scripts", data);
      scripts.value = [...data.scripts];
    })
    .catch((error) => {
      console.log("error", error);
    });
}
function handleCmdEvent(cmd) {
  fetch(`/agent/scripts/${cmd}`, { method: "GET" })
    .then((response) => response.json())
    .then((data) => {
      console.log("data", data);
      command.value = data.content;
    })
    .catch((error) => console.log("error", error));
}

function handleClick() {
  const data = { script: command.value };
  fetch("/agent/script/run", { method: "POST", body: JSON.stringify(data) })
    .then((response) => response.json())
    .then((data) => {
      console.log("success", data);
      out.value = data.out;
    })
    .catch((error) => {
      console.error("error", error);
    });
}

function handleRefresh() {
  fetchScipts();
}
onUpdated(() => {
  myinput.value.focus();
});
onMounted(() => {
  console.log("onMounted");
  fetchScipts();
});
</script>

<template>
  <div class="header">
    <h2 style="padding-right: 2em">Scirpt</h2>
    <button @click="handleRefresh">refresh</button>
  </div>
  <div>
    <button v-for="item in scripts" :key="item" @click="handleCmdEvent(item)">
      {{ item }}
    </button>
  </div>
  <div class="card">
    <div class="cmdline">
      <textarea
        ref="myinput"
        v-model="command"
        @keyup.ctrl.enter="handleClick"
        :rows="rows"
        :cols="cols"
      ></textarea>
      <button class="button" type="button" @click="handleClick">Run</button>
    </div>
    <div>
      <textarea
        v-model="out"
        disabled="false"
        :rows="outRows"
        :cols="outCols"
      ></textarea>
    </div>
  </div>
</template>

<style scoped>
.header {
  display: flex;
  place-items: center;
  text-align: center;
  justify-content: center;
}
.card {
  height: 500px;
  padding: 2em;
}
.cmdline {
  display: flex;
  place-items: center;
  margin: 0 0 1.5em 3em;
}
</style>
