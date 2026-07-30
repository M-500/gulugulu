const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8081,
    proxy: {
      '/na': {
        target: process.env.VUE_APP_PROXY_TARGET || 'http://localhost:8888',
        changeOrigin: true
      },
      '/api': {
        target: process.env.VUE_APP_PROXY_TARGET || 'http://localhost:8888',
        changeOrigin: true
      },
      '/gulugulu-media': {
        target: process.env.VUE_APP_MINIO_PROXY_TARGET || 'http://localhost:9000',
        changeOrigin: true
      },
      '/gulugulu-public': {
        target: process.env.VUE_APP_MINIO_PROXY_TARGET || 'http://localhost:9000',
        changeOrigin: true
      }
    }
  }
})
