<template>
  <div class="q-pa-md">
    <JsonEditor :json="object" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import JsonEditor from './JsonEditor.vue'

const baseURL = 'http://localhost:8080'

const object = ref({})

function sse(url: string, onMessage: (event: MessageEvent<any>) => any): EventSource {
  let es = new EventSource(baseURL + url, { withCredentials: true })
  es.onerror = (err) => console.log("SSE", err)
  es.onmessage = onMessage
  return es
}

function setupStream(kind: string, namespace: string, name: string) {
  sse(`/k8s/${kind}/${namespace}/${name}`, (event: MessageEvent<any>) => {
    var msg = JSON.parse(event.data)
    var obj = msg.Content
    switch (msg.Action) {
      case "delete":
        object.value = {}
        break
      case "add":
        object.value = obj
        break
      case "update":
        object.value = obj
        break
    }
  })
}
const route = useRoute()
setupStream(route.params.kind, route.params.namespace, route.params.name)
</script>
