import { Quasar } from 'quasar'
import 'quasar/src/css/index.sass'
import { createApp, h } from 'vue'
import { createRouter, createWebHashHistory, NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import App from './App.vue'
import { createAuth, useAuth } from './auth'
import Login from './lib/Login.vue'
import Logout from './lib/Logout.vue'
import Profile from './lib/Profile.vue'
import quasarUserOptions from './quasar-user-options'

createAuth("/profile")

const { checkAuth } = useAuth()

async function authGuard(_to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) {
  await checkAuth()
  next()
}

function setup() {
  const app = createApp({
    setup() {
    },
    render: () => h(App)
  })
  app.use(Quasar, quasarUserOptions)
  const router = createRouter({
    routes: [
      { path: '/', name: 'login', component: Login },
      { path: '/profile', name: 'profile', component: Profile },
      { path: '/logout', name: 'logout', component: Logout }
    ],
    history: createWebHashHistory()
  })
  router.beforeEach(authGuard)
  app.use(router)
  app.mount('#app')
}

setup()