const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
	lintOnSave: false,
	transpileDependencies: true,
	devServer: {
		host: 'localhost',
		port: 8080,
		proxy:{
			'/api':{
				target: 'http://10.221.100.157:10086',
				changeOrigin: true,
				pathRewrite: {
					'^/api': '/'
				}
			}
		}
	}
})