import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'
import Home from '../views/Home.vue'
import axios from 'axios'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    beforeEnter: (to, from, next) => {
      if(store.state.frontjws) {
        axios.get('/agent/isagent', {headers: {'Authorization': 'Bearer ' + store.state.frontjws}})
          .then(() => {
            next({path: '/agent'});
          })
          .catch(() => {
            next({path: '/home'});
          })
      } else {
        next({path: '/home/'});
      }
    }
  },
  {
    path: '/home',
    name: 'Home',
    component: Home,
  },
  {
    path: '/logout',
    beforeEnter: () => {
      store.state.frontjws = '';
      router.go(0);
    },
  },
  {
    path: '/signin',
    name: 'Signin',
    component: () => import('../views/Signin.vue'),
  },
  {
    path: '/signup',
    name: 'Signup',
    component: () => import('../views/Signup.vue'),
  },
  {
    path: '/agent',
    name: 'Agent',
    component: () => import('../views/Agent.vue'),
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
