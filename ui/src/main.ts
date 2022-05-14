import { authGuard } from '@auth0/auth0-vue'
import { Quasar } from 'quasar'
import 'quasar/src/css/index.sass'
import { createApp, h } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import Login from './components/Login.vue'
import Logout from './components/Logout.vue'
import Profile from './components/Profile.vue'
import quasarUserOptions from './quasar-user-options'

function setup() {
  const app = createApp({
    setup() {
    },
    render: () => h(App)
  })
  app.use(Quasar, quasarUserOptions)

  app.use(createRouter({
    routes: [
      { path: '/', name: 'login', component: Login },
      { path: '/profile', name: 'profile', component: Profile, beforeEnter: authGuard },
      { path: '/logout', name: 'logout', component: Logout, beforeEnter: authGuard }
    ],
    history: createWebHashHistory()
  }))
  app.mount('#app')
}

setup()