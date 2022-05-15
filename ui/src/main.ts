import { Quasar } from 'quasar'
import 'quasar/src/css/index.sass'
import { createApp, h } from 'vue'
import { createRouter, createWebHashHistory, NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import App from './App.vue'
import { createAuth, useAuth } from './auth'
import Login from './lib/Login.vue'
import Logout from './lib/Logout.vue'
import Profile from './lib/Profile.vue'
import Token from './lib/Token.vue'
import Pods from './lib/Pods.vue'
import KObject from './lib/KObject.vue'

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
      { path: '/', component: Login },
      { path: '/token', component: Token },
      { path: '/profile', component: Profile },
      { path: '/k8s/:kind/:namespace/:name', component: KObject },
      { path: '/pods', component: Pods },
      { path: '/logout', component: Logout }
    ],
    history: createWebHashHistory()
  })
  router.beforeEach(authGuard)
  app.use(router)
  app.mount('#app')
}

setup()