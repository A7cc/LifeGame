import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, '.'),
    },
  },
  build: {
    reportCompressedSize: false,
    rollupOptions: {
      // Vue and Element Plus ship stable browser ESM bundles. Keeping them as
      // local vendored modules lowers the production build peak from >600 MB
      // to a level suitable for the 2 GB ARM development machine.
      external: ['vue', 'element-plus'],
      output: {
        manualChunks(id) {
          if (id.includes('/node_modules/echarts/') || id.includes('/node_modules/zrender/')) {
            return 'charts-vendor'
          }
          if (
            id.includes('/node_modules/@vue/') ||
            id.includes('/node_modules/vue-router/') ||
            id.includes('/node_modules/pinia/')
          ) {
            return 'vue-vendor'
          }
        },
      },
    },
  },
})
