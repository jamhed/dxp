import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'
import { createHtmlPlugin } from 'vite-plugin-html'

export default defineConfig({
  preview: {
    port: 8080
  },
  server: {
    port: 8080
  },
  base: "/ui/",
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
        title: 'dxp',
      },
    },
  }),
  ],
})
