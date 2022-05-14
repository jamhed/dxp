<template>
  <div class="q-pa-md">
    <div class="row">
      <div class="col" v-if="isAuthenticated"></div>
      <div class="col" v-else>
        <q-btn @click="loginWithRedirect" label="Login" icon="login" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue'
const { loginWithRedirect, isAuthenticated } = useAuth0()
</script>
