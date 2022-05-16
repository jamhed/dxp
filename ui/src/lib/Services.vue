<template>
  <div class="q-pa-md">
    <q-table title="Services" :rows="rows" :columns="columns" :pagination="pagination" row-key="metadata.uid">
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

const columns: Array<ColumnType> = [
  { align: 'left', name: 'namespace', sortable: true, label: 'namespace', field: 'namespace' },
  { align: 'left', name: 'name', sortable: true, label: 'name', field: 'name' },
  { align: 'left', name: 'links', sortable: true, label: 'links', field: 'links' }
]

function object_uri(row: any) {
  return `/k8s/service/${row.metadata.namespace}/${row.metadata.name}`
}

function proxy_uri(port: string, row: any): string {
  return makeURL(`/proxy/${row.metadata.namespace}/${row.metadata.name}/${port}`)
}

const { rows } = useStream(makeURL('/watch/services'))
</script>
