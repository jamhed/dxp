import { computed, ComputedRef, reactive, ref, Ref } from 'vue'

type Cache = {
  [key: string]: any
}

type UseStreamRet = {
  cache: Cache
  rows: ComputedRef<unknown[]>
}

type UseSingleStreamRet = {
  obj: Ref<{}>
}

async function sse(url: string, onMessage: (event: MessageEvent<any>) => any): Promise<EventSource> {
  let es = new EventSource(url, { withCredentials: true })
  es.onerror = (err) => console.log("SSE", err)
  es.onmessage = onMessage
  return es
}

function objKey(obj: any) {
  let key = obj.kind + obj.metadata.namespace + obj.metadata.name
  return key
}

export function useStream(url: string): UseStreamRet {
  const cache: Cache = reactive({})
  const rows = computed(() => Object.values(cache).sort((a, b) => objKey(a).localeCompare(objKey(b))))

  sse(url, (event: MessageEvent<any>) => {
    var msg = JSON.parse(event.data)
    var obj = msg.Content
    switch (msg.Action) {
      case "delete":
        delete cache[obj.metadata.uid]
        break
      case "add":
        cache[obj.metadata.uid] = obj
        break
      case "update":
        cache[obj.metadata.uid] = obj
        break
    }
  })
  return { cache, rows }
}

export function useSingleStream(url: string): UseSingleStreamRet {
  const obj = ref({})
  sse(url, (event: MessageEvent<any>) => {
    var msg = JSON.parse(event.data)
    var k8s = msg.Content
    switch (msg.Action) {
      case "delete":
        obj.value = {}
        break
      case "add":
        obj.value = k8s
        break
      case "update":
        obj.value = k8s
        break
    }
  })
  return { obj }
}