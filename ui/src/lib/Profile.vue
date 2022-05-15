<template>
  <div class="q-pa-md" style="max-width: 400px">
    <q-form class="q-gutter-sm">
      <q-input filled v-model="client.given_name" label="First name" readonly />
      <q-input filled v-model="client.family_name" label="Last name" readonly />
      <q-input filled v-model="client.email" label="Email" readonly />
    </q-form>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '../auth'
const { profile } = useAuth()
const client = profile()
</script>
