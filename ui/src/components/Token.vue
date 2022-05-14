<template>
  <JsonEditor :json="token" />
</template>

<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue'
import { ref, onMounted } from 'vue'
import jwt_decode from "jwt-decode"
import JsonEditor from './JsonEditor.vue'

const { getAccessTokenSilently } = useAuth0()

const token = ref<{
  decoded: unknown,
  encoded: string
}>();

async function getToken() {
  const accessToken = await getAccessTokenSilently()
  token.value = {
    "decoded": jwt_decode(accessToken),
    "encoded": accessToken
  }
}
onMounted(getToken)
</script>

