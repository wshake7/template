import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'
import tailwindcss from '@tailwindcss/vite'
import { devtools } from '@tanstack/devtools-vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
// import { DevTools } from "@vitejs/devtools";
import react from '@vitejs/plugin-react'
import AutoImport from 'unplugin-auto-import/vite'
import { loadEnv } from 'vite'
// import Inject from "vite-plugin-inspect";
import { VitePWA } from 'vite-plugin-pwa'
import { defineConfig } from 'vite-plus'

export default defineConfig(({ mode }: { mode: string }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    base: env.GITHUB_ACTIONS === 'true' ? '/template/' : '/',
    staged: {
      '*': '',
    },
    resolve: {
      tsconfigPaths: true,
    },
    preview: {
      strictPort: true,
    },
    server: {
      port: Number(env.VITE_PORT),
      watch: {
        // wsl下热更新必须开
        usePolling: true,
        interval: 500,
      },
      proxy: {
        '/api': {
          target: env.VITE_API_URL,
          changeOrigin: true,
          secure: false,
          // rewrite: path => path.replace(/^\/api/, ''),
        },
      },
    },
    define: {
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: true,
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: true,
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id: string) {
            // React 核心
            if (id.includes('node_modules/react') || id.includes('node_modules/react-dom')) {
              return 'vendor-react'
            }
            // TanStack 系列
            if (id.includes('node_modules/@tanstack')) {
              return 'vendor-tanstack'
            }
            if (id.includes('zod')) {
              return 'vendor-zod'
            }
            // 其他 node_modules 统一打成 vendor
            if (id.includes('node_modules')) {
              return 'vendor'
            }
          },
        },
      },
    },
    plugins: [
      // DevTools(),
      // Inject(),
      devtools(),
      tailwindcss(),
      tanstackRouter(),
      react(),
      AutoImport({
        include: [/\.[jt]sx?$/, /\.md$/, /tsr-split/],
        imports: ['react', 'ahooks', 'react-i18next'],
        dts: 'src/auto-imports.d.ts',
        dirs: [
          'src/stores/**',
          'src/domains/**',
          'src/utils/**',
          'src/components/**',
          'src/api/**',
          'src/config/**',
        ],
      }) as any,
      VitePWA({
        registerType: 'autoUpdate', // 当有新版本时自动更新
        manifest: {
          name: '我的 React PWA 应用',
          short_name: 'ReactPWA',
          description: '一个使用 Vite 和 React 构建的 PWA 应用',
          theme_color: '#ffffff',
        },
      }),
      {
        name: 'remove-msw-in-prod',
        apply: 'build',
        closeBundle() {
          if (mode === 'prod') {
            const file = path.resolve('dist/mockServiceWorker.js')
            if (fs.existsSync(file)) {
              fs.unlinkSync(file)
              console.log('✅ removed mockServiceWorker.js')
            }
          }
        },
      },
    ],
  }
})
