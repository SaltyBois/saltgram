const fs = require('fs')

module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  configureWebpack: config => {
    if(process.env.NODE_ENV !== 'production') {
        config.devServer = {
            https: {
              key: fs.readFileSync('./certs/saltgram-web-server.key'),
              cert: fs.readFileSync('./certs/saltgram-web-server.crt'),
            },
        }
    }
 }
}
