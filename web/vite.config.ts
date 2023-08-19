import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    rollupOptions: {
      output: {
        dir: '../static/'
      }
    }
  },
  server: {
    proxy: {
      '/api': 'http://192.168.50.236:6060'
    }
  }
})
