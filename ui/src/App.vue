<template>
  <div class="q-pa-md">
    <q-layout view="hHh lpR fFf">
      <q-header elevated class="bg-black">
        <q-toolbar>
          <q-btn flat @click="drawer = !drawer" round dense icon="menu" />
          <q-toolbar-title>dxp</q-toolbar-title>
        </q-toolbar>
      </q-header>
      <q-drawer v-if="isAuthenticated" v-model="drawer" show-if-above :width="200" :breakpoint="500">
        <q-scroll-area class="fit">
          <q-list padding class="menu-list">
            <q-item clickable to="/profile">
              <q-item-section avatar>
                <q-icon name="info" />
              </q-item-section>
              <q-item-section>Profile</q-item-section>
            </q-item>
            <q-item clickable to="/token">
              <q-item-section avatar>
                <q-icon name="token" />
              </q-item-section>
              <q-item-section>Token</q-item-section>
            </q-item>
            <q-item clickable to="/pods">
              <q-item-section avatar>
                <q-icon name="bookmarks" />
              </q-item-section>
              <q-item-section>Pods</q-item-section>
            </q-item>
            <q-item clickable to="/services">
              <q-item-section avatar>
                <q-icon name="electrical_services" />
              </q-item-section>
              <q-item-section>Services</q-item-section>
            </q-item>
            <q-item clickable to="/logout">
              <q-item-section avatar>
                <q-icon name="logout" />
              </q-item-section>
              <q-item-section>Logout</q-item-section>
            </q-item>
          </q-list>
        </q-scroll-area>
      </q-drawer>
      <q-page-container>
        <q-page padding>
          <router-view />
        </q-page>
      </q-page-container>
    </q-layout>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from "./auth"
import { ref } from "vue"

const drawer = ref(false)

const { isAuthenticated } = useAuth()
</script>
