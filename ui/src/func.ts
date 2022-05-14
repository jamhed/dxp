export function prepare(value: { __typename?: unknown }): object {
  delete value["__typename"]
  return value
}

export function clone(value: object): any {
  return JSON.parse(JSON.stringify(value))
}

export function first(value: object): any {
  const [head] = Object.values(value)
  return head
}
