<template>
  <div class="q-pa-md">
    <div class="row">
      <div class="col" v-if="isAuth"></div>
      <div class="col" v-else>
        <q-btn @click="loginWithRedirect" label="Login" icon="login" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from "../auth"
const { isAuth } = useAuth()

function loginWithRedirect() {
}
</script>
