<template>
  <div></div>
</template>

<script setup lang="ts">
import { useAuth } from "../auth"
const { logout } = useAuth()
logout()
</script>
