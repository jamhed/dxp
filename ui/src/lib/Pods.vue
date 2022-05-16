<template>
  <div class="q-pa-md">
    <q-table title="Pods" :rows="rows" :columns="columns" :pagination="pagination" row-key="metadata.uid">
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td key="namespace" :props="props">
            {{ props.row.metadata.namespace }}
          </q-td>
          <q-td key="name" :props="props">
            <router-link :to="object_uri(props.row)">{{ props.row.metadata.name }}</router-link>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </div>
</template>

<script setup lang="ts">
import { useStream } from '../sse'
import { useConfig } from '../config'

const { makeURL } = useConfig()

const pagination = { rowsPerPage: 50, sortBy: 'creationDate' }

interface ColumnType {
  name: string
  label: string
  field: string
  sortable?: boolean
  align: "left" | "right" | "center"
}

function object_uri(row: any) {
  return `/k8s/pod/${row.metadata.namespace}/${row.metadata.name}`
}

const columns: Array<ColumnType> = [
  { align: 'left', name: 'namespace', sortable: true, label: 'namespace', field: 'namespace' },
  { align: 'left', name: 'name', sortable: true, label: 'name', field: 'name' }
]

const { rows } = useStream(makeURL('/watch/pods'))
</script>
