<template>
  <div class="q-pa-md">
    <JsonEditor :json="obj" />
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useSingleStream } from '../sse'
import JsonEditor from './JsonEditor.vue'

const baseURL = 'http://localhost:8080'

const route = useRoute()
const { obj } = useSingleStream(`${baseURL}/k8s/${route.params.kind}/${route.params.namespace}/${route.params.name}`)
</script>
