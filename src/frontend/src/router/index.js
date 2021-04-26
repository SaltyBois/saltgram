import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import axios from 'axios'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    // NOTE(Jovan): Logged out if switching to home, just as protonmail
    // beforeEnter: (to, from, next) => {
    //   if(this.$store.state.jws) 
    //     next("/user");
    //   else
    //     next();
    // },
  },
 {
   path: '/user',
   name: 'User',
   component: () => import(/* webpackChunkName: "user" */ '../views/User.vue')
 },
 {
   // NOTE(Jovan): Confirm reset
  path: '/email/change/:token',
  name: 'PasswordReset',
  beforeEnter: (to, from, next) => {
    let token = to.params["token"]
    axios.put("/email/change/" + token)
      .then(r => {
        console.log(r);
        next();
      })
      .catch(r => {
        console.log(r);
        next({name: "Home"});
      });
  },
  component: () => import(/* webpackChunkName: "passwordReset" */ '../views/PasswordReset.vue')
 },
 {
   path: '/forgotpassword',
   name: 'ForgotPassword',
   component: () => import(/* webpackChunkName: "forgotPassword" */ '../views/ForgotPassword.vue')
 },
 {
   path: '/email/activate/:token',
   name: 'ActivateEmail',
   beforeEnter: (to, from, next) => {
    let token = to.params["token"]
    axios.put("email/activate/" + token)
      .finally(function(){
        next({ name: "Home"});
      })
   },
 }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// TODO(Jovan): Authentication
// router.beforeEach((to, from, next) => {

// });

export default router
