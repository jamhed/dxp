<template>
  <div class="q-pa-md">
    <q-card class="fixed-center" v-if="!isAuthenticated">
      <q-card-section>
        <div class="text-h6">Please authenticate</div>
      </q-card-section>

      <q-separator />

      <q-card-actions vertical>
        <q-btn @click="login" label="Login" icon="login" />
      </q-card-actions>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from "../auth"
const { isAuthenticated, login } = useAuth()
</script>
