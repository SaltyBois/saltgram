import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';

Vue.use(Vuetify);

export default new Vuetify({
  icons: {
    iconfont: 'fa',
  },
  theme: {
    themes: {
      light: {
        primary: '#040d21',
        secondary: '#0c162d',
        info: '#00cfc8',
        accent: '#e247a5',
        error: '#cb2431',
        success: '#34d058',
        lightsecondary: '#69708d',
      },
    },
    options: {
      customProperties: true,
    },
  },
});
