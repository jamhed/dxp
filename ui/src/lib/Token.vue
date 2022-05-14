<template>
  <JsonEditor :json="token" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import JsonEditor from './JsonEditor.vue'

const token = ref<{
  decoded: unknown,
  encoded: string
}>();

async function getToken() {
  token.value = {
    "decoded": {},
    "encoded": ""
  }
}
onMounted(getToken)
</script>

