<template>
  <JsonEditor :json="token" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import jwt_decode from "jwt-decode"
import JsonEditor from './JsonEditor.vue'

const token = ref<{
  decoded: unknown,
  encoded: string
}>();

async function getToken() {
  const re = await fetch("/token")
  const value = await re.json()
  token.value = {
    "decoded": jwt_decode(value),
    "encoded": value
  }
}
onMounted(getToken)
</script>

