const fs = require('fs')

module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  configureWebpack: config => {
    if(process.env.NODE_ENV !== 'production') {
        config.devServer = {
            https: {
              key: fs.readFileSync('./certs/agent-web.key'),
              cert: fs.readFileSync('./certs/agent-web.crt'),
            },
            port: 8070,
        }
    }
 }
}
