import { Ref, ref } from 'vue'

class Auth {
  authenticated: Ref<boolean>
  constructor(readonly url: string) {
    this.authenticated = ref(false)
  }

  checkAuth = async () => {
    const token = localStorage.getItem('token')
    const re = await fetch(this.url, {
      headers: {
        'Content-Type': 'application/json',
        authorization: `Bearer ${token}`
      }
    })
    const user = await re.json()
    if (user.user && user.user[0]) {
      this.authenticated.value = true
    } else {
      this.authenticated.value = false
    }
  }
}

let isAuth: Ref<boolean>
let checkAuth: () => any

export function createAuth(url: string) {
  ({ authenticated: isAuth, checkAuth } = new Auth(url))
}

export function useAuth() {
  return { isAuth, checkAuth }
}
