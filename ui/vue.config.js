const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false,
  
  // 禁用TypeScript类型检查
  parallel: false,
  
  chainWebpack: config => {
    // 完全禁用TypeScript相关插件
    if (config.plugins.has('fork-ts-checker')) {
      config.plugins.delete('fork-ts-checker')
    }
    
    // 移除TypeScript相关的loader和规则
    config.module.rules.delete('ts')
    config.module.rules.delete('tsx')
    
    // 确保入口文件正确
    config.entry('app').clear().add('./src/main.js')
    
    // 禁用类型检查
    config.resolve.extensions.delete('.ts').delete('.tsx')
  },
  
  devServer: {
    port: 3001,
    historyApiFallback: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8090',
        changeOrigin: true,
      }
    }
  },
  
  css: {
    loaderOptions: {
      less: {
        lessOptions: {
          modifyVars: {
            'primary-color': '#1890ff',
            'link-color': '#1890ff',
            'border-radius-base': '4px',
          },
          javascriptEnabled: true,
        }
      }
    }
  },
  
  configureWebpack: {
    resolve: {
      alias: {
        '@': require('path').resolve(__dirname, 'src')
      }
    }
  }
})