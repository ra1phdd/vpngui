import {defineConfig} from 'vite'
import path from 'path';
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    rollupOptions: {
      external: [],
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@wails': path.resolve(__dirname, './wailsjs'),
      '@assets': path.resolve(__dirname, './src/assets'),
      '@styles': path.resolve(__dirname, './src/assets/styles'),
      '@components': path.resolve(__dirname, './src/components'),
      '@constants': path.resolve(__dirname, './src/constants'),
      '@pages': path.resolve(__dirname, './src/pages')
    }
  },
  plugins: [react()],
})