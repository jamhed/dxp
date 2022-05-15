import { Ref, ref } from 'vue'

class Auth {
  isAuthenticated: Ref<boolean>

  constructor(readonly url: string) {
    this.isAuthenticated = ref(false)
  }

  checkAuth = async () => {
    var profile = localStorage.getItem('profile')
    if (profile) {
      this.isAuthenticated.value = true
      return
    }
    const re = await fetch(this.url, {
      headers: {
        'Content-Type': 'application/json'
      }
    })
    profile = await re.json()
    if (profile) {
      localStorage.setItem('profile', JSON.stringify(profile))
      this.isAuthenticated.value = true
      return
    }
    this.isAuthenticated.value = false
  }

  profile() {
    return JSON.parse(localStorage.getItem('profile') || "{}")
  }

  login() {
    localStorage.clear()
    window.location.href = '/auth'
  }

  logout() {
    localStorage.clear()
    window.location.href = '/logout'
  }
}

let isAuthenticated: Ref<boolean>
let checkAuth: () => any
let profile: () => any
let logout: () => any
let login: () => any


export function createAuth(url: string) {
  ({ isAuthenticated, checkAuth, profile, logout, login } = new Auth(url))
}

export function useAuth() {
  return { isAuthenticated, checkAuth, profile, logout, login }
}
