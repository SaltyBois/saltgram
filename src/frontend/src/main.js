import Vue from 'vue'
import Vuex from 'vuex'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueRecaptcha from 'vue-recaptcha'
import Login from './components/Login.vue'
import Register from './components/Register.vue'
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false

Vue.use(VueAxios, axios);
Vue.use(Vuex);

Vue.component("vue-recaptcha", VueRecaptcha);
Vue.component("sg-login", Login);
Vue.component("sg-register", Register);

axios.defaults.withCredentials = true;

const store = new Vuex.Store({
  state: {
    jws: "",
  },
});

new Vue({
  router,
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')
