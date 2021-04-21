const fs = require('fs')

module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    https: {
      key: fs.readFileSync('./certs/localhost.key'),
      cert: fs.readFileSync('./certs/localhost.crt'),
    },
  },
}
