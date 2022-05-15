<template>
  <div class="q-pa-md">
    <q-table title="Pods" :rows="rows" :columns="columns" :pagination="pagination" row-key="name" />
  </div>
</template>

<script setup lang="ts">
import { reactive, computed } from 'vue'
const baseURL = 'http://localhost:8080'

const pagination = { rowsPerPage: 50, sortBy: 'creationDate' }

const columns = [
  { align: 'left', name: 'namespace', sortable: true, label: 'namespace', field: (row: any) => row.metadata.namespace },
  { align: 'left', name: 'name', sortable: true, label: 'name', field: (row: any) => row.metadata.name }
]

function objKey(obj: any) {
  let key = obj.kind + obj.metadata.namespace + obj.metadata.name
  return key
}

const rows = computed(() => Object.values(cache).sort((a, b) => objKey(a).localeCompare(objKey(b))))

type Cache = {
  [key: string]: any
}

const cache: Cache = reactive({})

function sse(url: string, onMessage: (event: MessageEvent<any>) => any): EventSource {
  let es = new EventSource(baseURL + url, { withCredentials: true })
  es.onerror = (err) => console.log("SSE", err)
  es.onmessage = onMessage
  return es
}

function setupStream() {
  sse('/watch/pods', (event: MessageEvent<any>) => {
    var msg = JSON.parse(event.data)
    var obj = msg.Content
    switch (msg.Action) {
      case "delete":
        delete cache[obj.metadata.uid]
        break
      case "add":
        console.log(obj)
        cache[obj.metadata.uid] = obj
        break
      case "update":
        cache[obj.metadata.uid] = obj
        break
    }
  })
}

setupStream()
</script>
