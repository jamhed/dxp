<template>
  <div class="q-pa-md">
    <q-table title="Services" :rows="rows" :columns="columns" :pagination="pagination" row-key="name">
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td key="namespace" :props="props">
            {{ props.row.metadata.namespace }}
          </q-td>
          <q-td key="name" :props="props">
            <router-link :to="object_uri(props.row)">{{ props.row.metadata.name }}</router-link>
          </q-td>
          <q-td key="links" :props="props">
            <a v-for="port in props.row.spec.ports" :key="port.port" target="_blank"
              :href="proxy_uri(port.port, props.row)">
              {{ port.port }}
            </a>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </div>
</template>

<script setup lang="ts">
import { useStream } from '../sse'

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
  { align: 'left', name: 'name', sortable: true, label: 'name', field: (row: any) => row.metadata.name },
  { align: 'left', name: 'links', sortable: true, label: 'links', field: (row: any) => row.metadata.name },
]

function object_uri(row: any) {
  return `/k8s/service/${row.metadata.namespace}/${row.metadata.name}`
}

function proxy_uri(port: string, row: any): string {
  return `${baseURL}/proxy/${row.metadata.namespace}/${row.metadata.name}/${port}`
}

const { rows } = useStream(`${baseURL}/watch/services`)
</script>
