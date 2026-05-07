import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    port: 9001,
    proxy: {
      '/api': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
    },
  },
});
