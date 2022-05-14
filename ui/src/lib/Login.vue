<template>
  <div class="q-pa-md">
    <div class="row">
      <div class="col" v-if="isAuthenticated()"></div>
      <div class="col" v-else>
        <q-btn @click="loginWithRedirect" label="Login" icon="login" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
function isAuthenticated(): boolean {
  return true
}

function loginWithRedirect() {
}
</script>
