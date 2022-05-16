<template>
  <div class="q-pa-md">
    <q-table title="Pods" :rows="rows" :columns="columns" :pagination="pagination" row-key="name"
      @row-click="onClick" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router"
import { useStream } from '../sse'

const router = useRouter()

const baseURL = 'http://localhost:8080'

const pagination = { rowsPerPage: 50, sortBy: 'creationDate' }

interface ColumnType {
  name: string
  label: string
  field: string | ((row: any) => any);
  sortable?: boolean
  align: "left" | "right" | "center"
}

const columns: Array<ColumnType> = [
  { align: 'left', name: 'namespace', sortable: true, label: 'namespace', field: (row: any) => row.metadata.namespace },
  { align: 'left', name: 'name', sortable: true, label: 'name', field: (row: any) => row.metadata.name }
]

function onClick(_ev: any, row: any) {
  router.push(`/k8s/pod/${row.metadata.namespace}/${row.metadata.name}`)
}

const { rows } = useStream(`${baseURL}/watch/pods`)
</script>
