<template>
  <div class="q-pa-md" style="max-width: 400px">
    <q-form class="q-gutter-sm">
      <q-input filled v-model="client.firstName" label="First name" />
      <q-input filled v-model="client.lastName" label="Last name" />
      <q-input filled v-model="client.emailAddress" label="Email" />
      <q-input filled v-model="client.mobilePhone" label="Mobile phone" />
      <q-input filled v-model="client.birthDate" label="Birthday">
        <template v-slot:append>
          <q-icon name="event" class="cursor-pointer">
            <q-popup-proxy ref="qDateProxy" cover transition-show="scale" transition-hide="scale">
              <q-date v-model="client.birthDate" mask="YYYY-MM-DD">
                <div class="row items-center justify-end">
                  <q-btn v-close-popup label="Close" color="primary" flat />
                </div>
              </q-date>
            </q-popup-proxy>
          </q-icon>
        </template>
      </q-input>
      <div class="q-pa-md q-gutter-md">
        <q-btn v-if="isAuth" label="Update" @click="update" color="primary" />
        <q-btn v-else label="Register" @click="register" color="primary" />
        <q-btn v-if="isAuth" label="Unregister" @click="unregister" color="deep-orange" />
      </div>
    </q-form>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '../auth'

const client = {
  firstName: '',
  lastName: '',
  emailAddress: '',
  mobilePhone: '',
  birthDate: ''
}

const { isAuth, checkAuth } = useAuth()

function update() {
}

async function register() {
  checkAuth()
}

async function unregister() {
  checkAuth()
}
</script>
