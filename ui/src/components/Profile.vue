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
import { watch, ref } from 'vue'
import { prepare, first, clone } from '../func'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { useAuth0 } from '@auth0/auth0-vue'
import { useAuth } from '../auth'
import gql from 'graphql-tag'

const GET_CLIENT = gql`
  query getClient {
    Client {
      id encodedKey firstName lastName emailAddress mobilePhone birthDate
    }
  }`

const UPDATE_CLIENT = gql`
  mutation Mutate($client: ClientInput) {
    updateClient(client: $client) {
      id encodedKey firstName lastName emailAddress mobilePhone birthDate
    }
  }
`
const CREATE_CLIENT = gql`
  mutation Mutate($client: ClientInput) {
    createClient(client: $client) {
      id encodedKey firstName lastName emailAddress mobilePhone birthDate
    }
  }
`

const DELETE_CLIENT = gql`
  mutation deleteClient($oidc_id: String!) {
    delete_user_by_pk(oidc_id: $oidc_id) {
      oidc_id
    }
  }
`

type Client = {
  __typename?: unknown,
  firstName?: string,
  lastName?: string,
  emailAddress?: string,
  mobilePhone?: string,
  birthDate?: string,
}

const client = ref<Client>({})
const { user } = useAuth0()

const { isAuth, checkAuth } = useAuth()

if (!isAuth.value) {
  client.value = {
    firstName: user.value.given_name,
    lastName: user.value.family_name,
    emailAddress: user.value.email
  }
}

const { mutate: updateClient } = useMutation(UPDATE_CLIENT)
const { mutate: createClient } = useMutation(CREATE_CLIENT)
const { mutate: deleteClient } = useMutation(DELETE_CLIENT)

const { result } = useQuery(GET_CLIENT)

if (result && result.value) {
  client.value = clone(first(result.value))
}

watch(result, result => { client.value = clone(first(result)) })

function update() {
  updateClient({ client: prepare(client.value) })
}

async function register() {
  await createClient({ client: prepare(client.value) })
  checkAuth()
}

async function unregister() {
  await deleteClient({ oidc_id: user.value.sub })
  checkAuth()
}
</script>
