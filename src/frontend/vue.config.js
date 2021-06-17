const fs = require('fs')

module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
      https: {
        key: fs.readFileSync('./certs/saltgram-web-server.key'),
        cert: fs.readFileSync('./certs/saltgram-web-server.crt'),
      },
  },
//  devServer: {
//    proxy: {
//      '^/api': {
//        target: "http://localhost:8081/",
//        changeOrigin: true,
//        logLevel: 'debug',
//        pathRewrite: {
//          '^/api': ''
//        }
//      }
//    }
//  }
}
