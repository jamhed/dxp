declare global {
  interface Window {
    dxp: {
      baseURL: string
    }
  }
}

class Config {
  baseURL: string
  constructor() {
    if (window.dxp && window.dxp.baseURL) {
      this.baseURL = window.dxp.baseURL
    } else {
      this.baseURL = ''
    }
  }
  makeURL(url: string): string {
    return `${this.baseURL}${url}`
  }
}

let makeURL: (url: string) => string

export function createConfig() {
  ({ makeURL } = new Config())
}

export function useConfig() {
  return { makeURL }
}
