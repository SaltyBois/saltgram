import Vue from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueRecaptcha from 'vue-recaptcha';
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false

Vue.use(VueAxios, axios)
Vue.component("vue-recaptcha", VueRecaptcha)

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
