import { Quasar } from 'quasar'
import 'quasar/src/css/index.sass'
import { createApp, h } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import Login from './lib/Login.vue'
import Logout from './lib/Logout.vue'
import Profile from './lib/Profile.vue'
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
      { path: '/profile', name: 'profile', component: Profile },
      { path: '/logout', name: 'logout', component: Logout }
    ],
    history: createWebHashHistory()
  }))
  app.mount('#app')
}

setup()