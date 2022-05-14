import { Ref, ref } from 'vue'

class Auth {
  registered: Ref<boolean>
  constructor(readonly url: string) {
    this.registered = ref(false)
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
      this.registered.value = true
    } else {
      this.registered.value = false
    }
  }
}

let isAuth: Ref<boolean>
let checkAuth: () => any

export function createAuth(url: string) {
  ({ registered: isAuth, checkAuth } = new Auth(url))
}

export function useAuth() {
  return { isAuth, checkAuth }
}
