import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { vitePluginForArco } from '@arco-plugins/vite-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ArcoResolver } from 'unplugin-vue-components/resolvers'
import { resolve } from 'path'
import VueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig(({ command }) => {
  const isProduction = command === 'build'
  
  return {
    // 支持反向代理的基础路径配置
    base: isProduction ? './' : '/',
    
    plugins: [
      vue(),
      VueDevTools(),
      vitePluginForArco({
        style: 'css'
      }),
      AutoImport({
        imports: ['vue', 'vue-router', 'pinia'],
        resolvers: [ArcoResolver()],
      }),
      Components({
        resolvers: [
          ArcoResolver({
            sideEffect: true
          })
        ]
      })
    ],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
      },
    },
    build: {
      outDir: '../static/dist',
      emptyOutDir: true,
      // 确保资源使用相对路径
      assetsDir: 'assets',
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: ['vue', 'vue-router', 'pinia'],
            arco: ['@arco-design/web-vue']
          },
          // 确保文件名哈希一致性
          assetFileNames: 'assets/[name].[hash][extname]',
          chunkFileNames: 'assets/[name].[hash].js',
          entryFileNames: 'assets/[name].[hash].js'
        }
      }
    },
    server: {
      port: 3000,
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true
        }
      }
    }
  }
})