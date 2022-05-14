<template>
  <div></div>
</template>

<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue'
const { logout } = useAuth0()
logout()
</script>
