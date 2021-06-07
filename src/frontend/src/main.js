import Vue from 'vue'
import Vuex from 'vuex'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueRecaptcha from 'vue-recaptcha'
import Login from './components/Login.vue'
import Register from './components/Register.vue'
import vuetify from './plugins/vuetify'
import 'material-design-icons-iconfont/dist/material-design-icons.css'
import PortalVue from 'portal-vue';


Vue.config.productionTip = false

Vue.use(VueAxios, axios);
Vue.use(Vuex);
Vue.use(PortalVue);

Vue.component("vue-recaptcha", VueRecaptcha);
Vue.component("sg-login", Login);
Vue.component("sg-register", Register);

axios.defaults.withCredentials = true;
axios.defaults.baseURL = process.env.VUE_APP_API_ENDPOINT;

const store = new Vuex.Store({
  state: {
    jws: "",
  }
});
export default store;

Vue.mixin({
  methods: {
    getAHeader: function () {
       console.log("aaaa", store.state.jws);
      return {'Authorization': 'Bearer ' + store.state.jws};
    },
    refreshToken: async function (aHeader) {
       console.log("aaa", store.state.jws);
      return axios.get("auth/refresh", {headers: aHeader});
    },
  }
})

new Vue({
  router,
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')
