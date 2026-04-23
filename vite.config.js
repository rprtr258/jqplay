import { defineConfig } from 'vite'
import { createHtmlPlugin } from 'vite-plugin-html'

export default defineConfig({
  base: "",
  plugins: [createHtmlPlugin({ minify: true })],
  build: {
    outDir: 'dist',
    sourcemap: true
  },
  server: {
    port: 8080,
    open: true
  },
  publicDir: 'public'
})
