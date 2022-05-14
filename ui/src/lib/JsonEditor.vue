<template>
  <div ref="el"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import JSONEditor from 'jsoneditor'
import "jsoneditor/dist/jsoneditor.min.css"

const props = defineProps(['json'])
const el = ref()
var editor: JSONEditor

onMounted(() => {
  editor = new JSONEditor(
    el.value,
    { mode: "view", search: false, navigationBar: false, statusBar: false }
  )
  editor.set(props.json)
  editor.expandAll()
})

watch(() => props.json, (json) => {
  editor.set(json)
  editor.expandAll()
})
</script>

<style>
.jsoneditor-menu {
  display: none !important;
}
.jsoneditor {
  border: 0px !important;
}
.jsoneditor-outer {
  margin-top: 0px !important;
  padding-top: 0px !important;
}
</style>
