import { defineConfig } from 'vite'
import path from "path"
import react from '@vitejs/plugin-react-swc'
import tailwindcss from "@tailwindcss/vite"

export default defineConfig({
  base: "./",
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
    dedupe: ["react", "react-dom", "prop-types"]
  },
  optimizeDeps: {
    include: ['react', 'react-dom', 'prop-types'],
  },
  server: {
    watch: {
      usePolling: true,
    },
    host: true,
    strictPort: true,
    proxy: {
      "/api": {
        target: process.env.VITE_API_BASE_URL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, "")
      },
    },
  },
  build: {
    rollupOptions: {},
    // commonjsOptions: {
    //   include: ["/node_modules/"]
    // }
  }
})
