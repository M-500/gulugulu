const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8081,
    proxy: {
      '/na': {
        target: process.env.VUE_APP_PROXY_TARGET || 'http://localhost:8888',
        changeOrigin: true
      }
    }
  }
})
