import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import '@fortawesome/fontawesome-free/css/all.css'
import HexPageLoader from './components/HexPageLoader.vue'
import MainNavigation from './components/MainNavigation.vue'


Vue.config.productionTip = false
Vue.component('hexpage-loader', HexPageLoader);
Vue.component('main-navigation', MainNavigation);

new Vue({
  router,
  store,
  vuetify,
  render: function (h) { return h(App) }
}).$mount('#app')
