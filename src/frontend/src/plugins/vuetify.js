import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    themes: {
      light: {
        // primary: '#f5efe3',
        // secondary: '#211010',
        // accent: '#d11515',
        // error: '#ef5858',
        // info: '#2196F3',
        // success: '#4CAF50',
        // warning: '#FFC107'
      },
      dark: {
        accent: 'accent'
      }
    },
  },
  icons: {
    iconfont: 'fa4',
  },
});
