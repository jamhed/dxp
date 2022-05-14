import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'
import { createHtmlPlugin } from 'vite-plugin-html'

// https://vitejs.dev/config/
export default defineConfig({
  preview: {
    port: 8080
  },
  server: {
    port: 8080
  },
  plugins: [vue({
    template: { transformAssetUrls }
  }),
  quasar({
    sassVariables: './src/styles/quasar.variables.sass',
  }),
  createHtmlPlugin({
    minify: true,
    entry: '../src/main.ts',
    template: 'public/index.html',
    inject: {
      data: {
        title: 'AppVue',
      },
    },
  }),
  ],
})
