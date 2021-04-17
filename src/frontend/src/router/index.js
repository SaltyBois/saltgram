import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import axios from 'axios'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
 {
   path: '/register',
   name: 'Register',
   // route level code-splitting
   // this generates a separate chunk (about.[hash].js) for this route
   // which is lazy-loaded when the route is visited.
   component: () => import(/* webpackChunkName: "register" */ '../views/Register.vue')
 },
 {
   path: '/user',
   name: 'User',
   component: () => import(/* webpackChunkName: "user" */ '../views/User.vue')
 },
 {
  path: '/email/change/:token',
  name: 'PasswordReset',
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
   component: () => import(/* webpackChunkName: "activate" */ '../views/ActivateEmail')
 }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// TODO(Jovan): Authentication
router.beforeEach((to, from, next) => {

  if(to.name === "PasswordReset") {
    let token = to.params["token"]
    axios.get("http://localhost:8081/email/change/" + token)
      .then(r => {
        console.log(r);
        next();
      })
      .catch(r => {
        console.log(r);
        next({name: "Home"});
      });
  // } else if(to.name === "Activate") {
  //   let token = to.params["token"]
  //   axios.get("http://localhost:8081/email/activate/" + token)
  //     .then(r => {
  //       console.log("ACTIVATED!");
  //       console.log(r);
  //       next({name: "Home"});
  //     })
  //     .catch(r => {
  //       console.log("NOT ACTIVATED!");
  //       console.log(r);
  //       next({name: "Home"});
  //     });
  } else {
    next()
  }
})

export default router
