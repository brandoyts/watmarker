import { defineConfig, loadEnv } from 'vite'

import path from 'path'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "")

  return {
    base: './',
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
      dedupe: ['react', 'react-dom'],
    },
    optimizeDeps: {
      include: ['react', 'react-dom'],
    },
    server: {
      watch: {
        usePolling: true,
      },
      host: true,
      strictPort: true,
      proxy: {
        "/api": {
          target: env.API_BASE_URL,
          changeOrigin: true,
          rewrite: (p) => p.replace(/^\/api/, ""),
        },
      },
    },
    build: {
      sourcemap: mode === 'production' ? false : true,
      minify: 'esbuild',
      rollupOptions: {
        output: {
          manualChunks: {
            react: ['react', 'react-dom'],
          },
        },
      },
    },
  }
})
